package utopia

import (
	//"log"
	//"image"
	//"math"
	//"sort"

	//"github.com/utopiagio/gio/text"
	//"golang.org/x/image/math/fixed"
)

import (
	"bufio"
	"image"
	"io"
	//"log"
	"math"
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/utopiagio/gio/f32"
	"github.com/utopiagio/gio/io/system"
	"github.com/utopiagio/gio/layout"
	"github.com/utopiagio/gio/op"
	"github.com/utopiagio/gio/op/clip"
	"github.com/utopiagio/gio/op/paint"
	"github.com/utopiagio/gio/text"
	"github.com/utopiagio/gio/unit"
	//widget_gio "github.com/utopiagio/gio/widget"
	"golang.org/x/image/math/fixed"

)

type offEntry struct {
	runes int
	bytes int
}

type maskReader struct {
	// rr is the underlying reader.
	rr      io.RuneReader
	maskBuf [utf8.UTFMax]byte
	// mask is the utf-8 encoded mask rune.
	mask []byte
	// overflow contains excess mask bytes left over after the last Read call.
	overflow []byte
}

func (m *maskReader) Reset(r io.Reader, mr rune) {
	m.rr = bufio.NewReader(r)
	n := utf8.EncodeRune(m.maskBuf[:], mr)
	m.mask = m.maskBuf[:n]
}

// Read reads from the underlying reader and replaces every
// rune with the mask rune.
func (m *maskReader) Read(b []byte) (n int, err error) {
	for len(b) > 0 {
		var replacement []byte
		if len(m.overflow) > 0 {
			replacement = m.overflow
		} else {
			var r rune
			r, _, err = m.rr.ReadRune()
			if err != nil {
				break
			}
			if r == '\n' {
				replacement = []byte{'\n'}
			} else {
				replacement = m.mask
			}
		}
		nn := copy(b, replacement)
		m.overflow = replacement[nn:]
		n += nn
		b = b[nn:]
	}
	return n, err
}

// textSource provides text data for use in widgets. If the underlying data type
// can fail due to I/O errors, it is the responsibility of that type to provide
// its own mechanism to surface and handle those errors. They will not always
// be returned by widgets using these functions.
type textSource interface {
	io.ReaderAt
	// Size returns the total length of the data in bytes.
	Size() int64
	// Changed returns whether the contents have changed since the last call
	// to Changed.
	Changed() bool
	// ReplaceRunes replaces runeCount runes starting at byteOffset within the
	// data with the provided string. Implementations of read-only text sources
	// are free to make this a no-op.
	ReplaceRunes(byteOffset int64, runeCount int64, replacement string)
}

// GioTextView provides efficient shaping and indexing of interactive text. When provided
// with a TextSource, GioTextView will shape and cache the runes within that source.
// It provides methods for configuring a viewport onto the shaped text which can
// be scrolled, and for configuring and drawing text selection boxes.
type GioTextView struct {
	Alignment text.Alignment
	// SingleLine forces the text to stay on a single line.
	// SingleLine also sets the scrolling direction to
	// horizontal.
	SingleLine bool
	// MaxLines limits the shaped text to a specific quantity of shaped lines.
	MaxLines int
	// Mask replaces the visual display of each rune in the contents with the given rune.
	// Newline characters are not masked. When non-zero, the unmasked contents
	// are accessed by Len, Text, and SetText.
	Mask rune

	font               text.Font
	shaper             *text.Shaper
	textSize           fixed.Int26_6
	seekCursor         int64
	rr                 textSource
	maskReader         maskReader
	lastMask           rune
	maxWidth, minWidth int
	viewSize           image.Point
	valid              bool
	regions            []Region
	dims               layout.Dimensions

	// offIndex is an index of rune index to byte offsets.
	offIndex []offEntry

	index glyphIndex

	caret struct {
		// xoff is the offset to the current position when moving between lines.
		xoff fixed.Int26_6
		// start is the current caret position in runes, and also the start position of
		// selected text. end is the end position of selected text. If start
		// == end, then there's no selection. Note that it's possible (and
		// common) that the caret (start) is after the end, e.g. after
		// Shift-DownArrow.
		start int
		end   int
	}

	scrollOff image.Point

	locale system.Locale
}

func (e *GioTextView) Changed() bool {
	return e.rr.Changed()
}

// Dimensions returns the dimensions of the visible text.
func (e *GioTextView) Dimensions() layout.Dimensions {
	basePos := e.dims.Size.Y - e.dims.Baseline
	return layout.Dimensions{Size: e.viewSize, Baseline: e.viewSize.Y - basePos}
}

// FullDimensions returns the dimensions of all shaped text, including
// text that isn't visible within the current viewport.
func (e *GioTextView) FullDimensions() layout.Dimensions {
	return e.dims
}

// SetSource initializes the underlying data source for the Text. This
// must be done before invoking any other methods on Text.
func (e *GioTextView) SetSource(source textSource) {
	e.rr = source
	e.invalidate()
	e.seekCursor = 0
}

// ReadRuneAt reads the rune starting at the given byte offset, if any.
func (e *GioTextView) ReadRuneAt(off int64) (rune, int, error) {
	var buf [utf8.UTFMax]byte
	b := buf[:]
	n, err := e.rr.ReadAt(b, off)
	b = b[:n]
	r, s := utf8.DecodeRune(b)
	return r, s, err
}

// ReadRuneAt reads the run prior to the given byte offset, if any.
func (e *GioTextView) ReadRuneBefore(off int64) (rune, int, error) {
	var buf [utf8.UTFMax]byte
	b := buf[:]
	if off < utf8.UTFMax {
		b = b[:off]
		off = 0
	} else {
		off -= utf8.UTFMax
	}
	n, err := e.rr.ReadAt(b, off)
	b = b[:n]
	r, s := utf8.DecodeLastRune(b)
	return r, s, err
}

func (e *GioTextView) makeValid() {
	if e.valid {
		return
	}
	e.layoutText(e.shaper)
	e.valid = true
}

func (e *GioTextView) closestToRune(runeIdx int) combinedPos {
	e.makeValid()
	pos, _ := e.index.closestToRune(runeIdx)
	return pos
}

func (e *GioTextView) closestToLineCol(line, col int) combinedPos {
	e.makeValid()
	return e.index.closestToLineCol(screenPos{line: line, col: col})
}

func (e *GioTextView) closestToXY(x fixed.Int26_6, y int) combinedPos {
	e.makeValid()
	return e.index.closestToXY(x, y)
}

func (e *GioTextView) MoveLines(distance int, selAct selectionAction) {
	caretStart := e.closestToRune(e.caret.start)
	x := caretStart.x + e.caret.xoff
	// Seek to line.
	pos := e.closestToLineCol(caretStart.lineCol.line+distance, 0)
	pos = e.closestToXY(x, pos.y)
	e.caret.start = pos.runes
	e.caret.xoff = x - pos.x

	e.updateSelection(selAct)

}

// calculateViewSize determines the size of the current visible content,
// ensuring that even if there is no text content, some space is reserved
// for the caret.
func (e *GioTextView) calculateViewSize(gtx layout.Context) image.Point {
	base := e.dims.Size
	if caretWidth := e.caretWidth(gtx); base.X < caretWidth {
		base.X = caretWidth
	}
	return gtx.Constraints.Constrain(base)
}

// Update the text, reshaping it as necessary. If not nil, eventHandling will be invoked after reshaping the text to
// allow parent widgets to adapt to any changes in text content or positioning. If eventHandling modifies the contents
// of the GioTextView, it is guaranteed to be reshaped (and ready for painting) before Update returns.
func (e *GioTextView) Update(gtx layout.Context, lt *text.Shaper, font text.Font, size unit.Sp, eventHandling func(gtx layout.Context)) {
	if e.locale != gtx.Locale {
		e.locale = gtx.Locale
		e.invalidate()
	}
	textSize := fixed.I(gtx.Sp(size))
	if e.font != font || e.textSize != textSize {
		e.invalidate()
		e.font = font
		e.textSize = textSize
	}
	maxWidth := gtx.Constraints.Max.X
	if e.SingleLine {
		maxWidth = math.MaxInt
	}
	minWidth := gtx.Constraints.Min.X
	if maxWidth != e.maxWidth {
		e.maxWidth = maxWidth
		e.invalidate()
	}
	if minWidth != e.minWidth {
		e.minWidth = minWidth
		e.invalidate()
	}
	if lt != e.shaper {
		e.shaper = lt
		e.invalidate()
	}
	if e.Mask != e.lastMask {
		e.lastMask = e.Mask
		e.invalidate()
	}

	e.makeValid()
	if eventHandling != nil {
		eventHandling(gtx)
		e.makeValid()
	}

	if viewSize := e.calculateViewSize(gtx); viewSize != e.viewSize {
		e.viewSize = viewSize
		e.invalidate()
	}
	e.makeValid()
}

// PaintSelection clips and paints the visible text selection rectangles. Callers
// are expected to apply an appropriate paint material with a paint.ColorOp or
// similar prior to calling PaintSelection.
func (e *GioTextView) PaintSelection(gtx layout.Context) {
	localViewport := image.Rectangle{Max: e.viewSize}
	docViewport := image.Rectangle{Max: e.viewSize}.Add(e.scrollOff)
	defer clip.Rect(localViewport).Push(gtx.Ops).Pop()
	e.regions = e.index.locate(docViewport, e.caret.start, e.caret.end, e.regions)
	for _, region := range e.regions {
		area := clip.Rect(region.Bounds).Push(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		area.Pop()
	}
}

// PaintText clips and paints the visible text glyph outlines. Callers
// are expected to apply an appropriate paint material with a paint.ColorOp or
// similar prior to calling PaintSelection.
func (e *GioTextView) PaintText(gtx layout.Context) {
	//log.Println("*GioTextView::PaintText()")
	m := op.Record(gtx.Ops)
	viewport := image.Rectangle{
		Min: e.scrollOff,
		Max: e.viewSize.Add(e.scrollOff),
	}
	it := textIterator{viewport: viewport}

	startGlyph := 0
	for _, line := range e.index.lines {
		if line.descent.Ceil()+line.yOff >= viewport.Min.Y {
			break
		}
		startGlyph += line.glyphs
	}
	var glyphs [32]text.Glyph
	line := glyphs[:0]
	for _, g := range e.index.glyphs[startGlyph:] {
		var ok bool
		if line, ok = it.paintGlyph(gtx, e.shaper, g, line); !ok {
			break
		}
	}

	call := m.Stop()
	viewport.Min = viewport.Min.Add(it.padding.Min)
	viewport.Max = viewport.Max.Add(it.padding.Max)
	defer clip.Rect(viewport.Sub(e.scrollOff)).Push(gtx.Ops).Pop()
	call.Add(gtx.Ops)
}

// caretWidth returns the width occupied by the caret for the current
// gtx.
func (e *GioTextView) caretWidth(gtx layout.Context) int {
	carWidth2 := gtx.Dp(1) / 2
	if carWidth2 < 1 {
		carWidth2 = 1
	}
	return carWidth2
}

// PaintCaret clips and paints the caret rectangle. Callers
// are expected to apply an appropriate paint material with a paint.ColorOp or
// similar prior to calling PaintSelection.
func (e *GioTextView) PaintCaret(gtx layout.Context) {
	carWidth2 := e.caretWidth(gtx)
	caretPos, carAsc, carDesc := e.CaretInfo()

	carRect := image.Rectangle{
		Min: caretPos.Sub(image.Pt(carWidth2, carAsc)),
		Max: caretPos.Add(image.Pt(carWidth2, carDesc)),
	}
	cl := image.Rectangle{Max: e.viewSize}
	carRect = cl.Intersect(carRect)
	if !carRect.Empty() {
		defer clip.Rect(carRect).Push(gtx.Ops).Pop()
		paint.PaintOp{}.Add(gtx.Ops)
	}
}

func (e *GioTextView) CaretInfo() (pos image.Point, ascent, descent int) {
	caretStart := e.closestToRune(e.caret.start)

	ascent = caretStart.ascent.Ceil()
	descent = caretStart.descent.Ceil()

	pos = image.Point{
		X: caretStart.x.Round(),
		Y: caretStart.y,
	}
	pos = pos.Sub(e.scrollOff)
	return
}

// ByteOffset returns the start byte of the rune at the given
// rune offset, clamped to the size of the text.
func (e *GioTextView) ByteOffset(runeOffset int) int64 {
	return int64(e.runeOffset(e.closestToRune(runeOffset).runes))
}

// Len is the length of the editor contents, in runes.
func (e *GioTextView) Len() int {
	e.makeValid()
	return e.closestToRune(math.MaxInt).runes
}

// Text returns the contents of the editor. If the provided buf is large enough, it will
// be filled and returned. Otherwise a new buffer will be allocated.
// Callers can guarantee that buf is large enough by giving it capacity e.Len()*utf8.UTFMax.
func (e *GioTextView) Text(buf []byte) []byte {
	size := e.rr.Size()
	if cap(buf) < int(size) {
		buf = make([]byte, size)
	}
	buf = buf[:size]
	e.Seek(0, io.SeekStart)
	n, _ := io.ReadFull(e, buf)
	buf = buf[:n]
	return buf
}

func (e *GioTextView) ScrollBounds() image.Rectangle {
	var b image.Rectangle
	if e.SingleLine {
		if len(e.index.lines) > 0 {
			line := e.index.lines[0]
			b.Min.X = line.xOff.Floor()
			if b.Min.X > 0 {
				b.Min.X = 0
			}
		}
		b.Max.X = e.dims.Size.X + b.Min.X - e.viewSize.X
	} else {
		b.Max.Y = e.dims.Size.Y - e.viewSize.Y
	}
	return b
}

func (e *GioTextView) ScrollRel(dx, dy int) {
	e.scrollAbs(e.scrollOff.X+dx, e.scrollOff.Y+dy)
}

// ScrollOff returns the scroll offset of the text viewport.
func (e *GioTextView) ScrollOff() image.Point {
	return e.scrollOff
}

func (e *GioTextView) scrollAbs(x, y int) {
	e.scrollOff.X = x
	e.scrollOff.Y = y
	b := e.ScrollBounds()
	if e.scrollOff.X > b.Max.X {
		e.scrollOff.X = b.Max.X
	}
	if e.scrollOff.X < b.Min.X {
		e.scrollOff.X = b.Min.X
	}
	if e.scrollOff.Y > b.Max.Y {
		e.scrollOff.Y = b.Max.Y
	}
	if e.scrollOff.Y < b.Min.Y {
		e.scrollOff.Y = b.Min.Y
	}
}

func (e *GioTextView) MoveCoord(pos image.Point) {
	x := fixed.I(pos.X + e.scrollOff.X)
	y := pos.Y + e.scrollOff.Y
	e.caret.start = e.closestToXY(x, y).runes
	e.caret.xoff = 0
}

func (e *GioTextView) layoutText(lt *text.Shaper) {
	e.Seek(0, io.SeekStart)
	var r io.Reader = e
	if e.Mask != 0 {
		e.maskReader.Reset(e, e.Mask)
		r = &e.maskReader
	}
	e.index = glyphIndex{}
	it := textIterator{viewport: image.Rectangle{Max: image.Point{X: math.MaxInt, Y: math.MaxInt}}}
	if lt != nil {
		lt.Layout(text.Parameters{
			Font:      e.font,
			PxPerEm:   e.textSize,
			Alignment: e.Alignment,
			MaxLines:  e.MaxLines,
		}, e.minWidth, e.maxWidth, e.locale, r)
		for glyph, ok := it.processGlyph(lt.NextGlyph()); ok; glyph, ok = it.processGlyph(lt.NextGlyph()) {
			e.index.Glyph(glyph)
		}
	} else {
		// Make a fake glyph for every rune in the reader.
		b := bufio.NewReader(r)
		for _, _, err := b.ReadRune(); err != io.EOF; _, _, err = b.ReadRune() {
			g, _ := it.processGlyph(text.Glyph{Runes: 1, Flags: text.FlagClusterBreak}, true)
			e.index.Glyph(g)

		}
	}
	dims := layout.Dimensions{Size: it.bounds.Size()}
	dims.Baseline = dims.Size.Y - it.baseline
	e.dims = dims
}

// CaretPos returns the line & column numbers of the caret.
func (e *GioTextView) CaretPos() (line, col int) {
	pos := e.closestToRune(e.caret.start)
	return pos.lineCol.line, pos.lineCol.col
}

// CaretCoords returns the coordinates of the caret, relative to the
// editor itself.
func (e *GioTextView) CaretCoords() f32.Point {
	pos := e.closestToRune(e.caret.start)
	return f32.Pt(float32(pos.x)/64-float32(e.scrollOff.X), float32(pos.y-e.scrollOff.Y))
}

// indexRune returns the latest rune index and byte offset no later than r.
func (e *GioTextView) indexRune(r int) offEntry {
	// Initialize index.
	if len(e.offIndex) == 0 {
		e.offIndex = append(e.offIndex, offEntry{})
	}
	i := sort.Search(len(e.offIndex), func(i int) bool {
		entry := e.offIndex[i]
		return entry.runes >= r
	})
	// Return the entry guaranteed to be less than or equal to r.
	if i > 0 {
		i--
	}
	return e.offIndex[i]
}

// runeOffset returns the byte offset into e.rr of the r'th rune.
// r must be a valid rune index, usually returned by closestPosition.
func (e *GioTextView) runeOffset(r int) int {
	const runesPerIndexEntry = 50
	entry := e.indexRune(r)
	lastEntry := e.offIndex[len(e.offIndex)-1].runes
	for entry.runes < r {
		if entry.runes > lastEntry && entry.runes%runesPerIndexEntry == runesPerIndexEntry-1 {
			e.offIndex = append(e.offIndex, entry)
		}
		_, s, _ := e.ReadRuneAt(int64(entry.bytes))
		entry.bytes += s
		entry.runes++
	}
	return entry.bytes
}

func (e *GioTextView) invalidate() {
	e.offIndex = e.offIndex[:0]
	e.valid = false
}

// Replace the text between start and end with s. Indices are in runes.
// It returns the number of runes inserted.
func (e *GioTextView) Replace(start, end int, s string) int {
	if start > end {
		start, end = end, start
	}
	startPos := e.closestToRune(start)
	endPos := e.closestToRune(end)
	startOff := e.runeOffset(startPos.runes)
	replaceSize := endPos.runes - startPos.runes
	sc := utf8.RuneCountInString(s)
	newEnd := startPos.runes + sc

	e.rr.ReplaceRunes(int64(startOff), int64(replaceSize), s)
	adjust := func(pos int) int {
		switch {
		case newEnd < pos && pos <= endPos.runes:
			pos = newEnd
		case endPos.runes < pos:
			diff := newEnd - endPos.runes
			pos = pos + diff
		}
		return pos
	}
	e.caret.start = adjust(e.caret.start)
	e.caret.end = adjust(e.caret.end)
	e.invalidate()
	return sc
}

func (e *GioTextView) MovePages(pages int, selAct selectionAction) {
	caret := e.closestToRune(e.caret.start)
	x := caret.x + e.caret.xoff
	y := caret.y + pages*e.viewSize.Y
	pos := e.closestToXY(x, y)
	e.caret.start = pos.runes
	e.caret.xoff = x - pos.x
	e.updateSelection(selAct)
}

// MoveCaret moves the caret (aka selection start) and the selection end
// relative to their current positions. Positive distances moves forward,
// negative distances moves backward. Distances are in runes.
func (e *GioTextView) MoveCaret(startDelta, endDelta int) {
	e.caret.xoff = 0
	e.caret.start = e.closestToRune(e.caret.start + startDelta).runes
	e.caret.end = e.closestToRune(e.caret.end + endDelta).runes
}

func (e *GioTextView) MoveStart(selAct selectionAction) {
	caret := e.closestToRune(e.caret.start)
	caret = e.closestToLineCol(caret.lineCol.line, 0)
	e.caret.start = caret.runes
	e.caret.xoff = -caret.x
	e.updateSelection(selAct)
}

func (e *GioTextView) MoveEnd(selAct selectionAction) {
	caret := e.closestToRune(e.caret.start)
	caret = e.closestToLineCol(caret.lineCol.line, math.MaxInt)
	e.caret.start = caret.runes
	e.caret.xoff = fixed.I(e.maxWidth) - caret.x
	e.updateSelection(selAct)
}

// MoveWord moves the caret to the next word in the specified direction.
// Positive is forward, negative is backward.
// Absolute values greater than one will skip that many words.
// BUG(whereswaldon): this method's definition of a "word" is currently
// whitespace-delimited. Languages that do not use whitespace to delimit
// words will experience counter-intuitive behavior when navigating by
// word.
func (e *GioTextView) MoveWord(distance int, selAct selectionAction) {
	// split the distance information into constituent parts to be
	// used independently.
	words, direction := distance, 1
	if distance < 0 {
		words, direction = distance*-1, -1
	}
	// atEnd if caret is at either side of the buffer.
	caret := e.closestToRune(e.caret.start)
	atEnd := func() bool {
		return caret.runes == 0 || caret.runes == e.Len()
	}
	// next returns the appropriate rune given the direction.
	next := func() (r rune) {
		off := e.runeOffset(caret.runes)
		if direction < 0 {
			r, _, _ = e.ReadRuneBefore(int64(off))
		} else {
			r, _, _ = e.ReadRuneAt(int64(off))
		}
		return r
	}
	for ii := 0; ii < words; ii++ {
		for r := next(); unicode.IsSpace(r) && !atEnd(); r = next() {
			e.MoveCaret(direction, 0)
			caret = e.closestToRune(e.caret.start)
		}
		e.MoveCaret(direction, 0)
		caret = e.closestToRune(e.caret.start)
		for r := next(); !unicode.IsSpace(r) && !atEnd(); r = next() {
			e.MoveCaret(direction, 0)
			caret = e.closestToRune(e.caret.start)
		}
	}
	e.updateSelection(selAct)
}

func (e *GioTextView) ScrollToCaret() {
	//log.Println("GioTextView::ScrollToCaret()")
	caret := e.closestToRune(e.caret.start)
	if e.SingleLine {
		var dist int
		if d := caret.x.Floor() - e.scrollOff.X; d < 0 {
			dist = d
		} else if d := caret.x.Ceil() - (e.scrollOff.X + e.viewSize.X); d > 0 {
			dist = d
		}
		e.ScrollRel(dist, 0)
	} else {
		
		miny := caret.y - caret.ascent.Ceil()
		maxy := caret.y + caret.descent.Ceil()
		var dist int
		if d := miny - e.scrollOff.Y; d < 0 {
			dist = d
		} else if d := maxy - (e.scrollOff.Y + e.viewSize.Y); d > 0 {
			dist = d
		}
		//log.Println("ScrollRel - dist =", dist)
		e.ScrollRel(0, dist)
	}
}

// SelectionLen returns the length of the selection, in runes; it is
// equivalent to utf8.RuneCountInString(e.SelectedText()).
func (e *GioTextView) SelectionLen() int {
	return e.absValue(e.caret.start - e.caret.end)
}

// Selection returns the start and end of the selection, as rune offsets.
// start can be > end.
func (e *GioTextView) Selection() (start, end int) {
	return e.caret.start, e.caret.end
}

// SetCaret moves the caret to start, and sets the selection end to end. start
// and end are in runes, and represent offsets into the editor text.
func (e *GioTextView) SetCaret(start, end int) {
	e.caret.start = e.closestToRune(start).runes
	e.caret.end = e.closestToRune(end).runes
}

// SelectedText returns the currently selected text (if any) from the editor,
// filling the provided byte slice if it is large enough or allocating and
// returning a new byte slice if the provided one is insufficient.
// Callers can guarantee that the buf is large enough by providing a buffer
// with capacity e.SelectionLen()*utf8.UTFMax.
func (e *GioTextView) SelectedText(buf []byte) []byte {
	startOff := e.runeOffset(e.caret.start)
	endOff := e.runeOffset(e.caret.end)
	start := e.minValue(startOff, endOff)
	end := e.maxValue(startOff, endOff)
	if cap(buf) < end-start {
		buf = make([]byte, end-start)
	}
	buf = buf[:end-start]
	n, _ := e.rr.ReadAt(buf, int64(start))
	// There is no way to reasonably handle a read error here. We rely upon
	// implementations of textSource to provide other ways to signal errors
	// if the user cares about that, and here we use whatever data we were
	// able to read.
	return buf[:n]
}

func (e *GioTextView) updateSelection(selAct selectionAction) {
	if selAct == selectionClear {
		e.ClearSelection()
	}
}

// ClearSelection clears the selection, by setting the selection end equal to
// the selection start.
func (e *GioTextView) ClearSelection() {
	e.caret.end = e.caret.start
}

// WriteTo implements io.WriterTo.
func (e *GioTextView) WriteTo(w io.Writer) (int64, error) {
	e.Seek(0, io.SeekStart)
	return io.Copy(w, struct{ io.Reader }{e})
}

// Seek implements io.Seeker.
func (e *GioTextView) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		e.seekCursor = offset
	case io.SeekCurrent:
		e.seekCursor += offset
	case io.SeekEnd:
		e.seekCursor = e.rr.Size() + offset
	}
	return e.seekCursor, nil
}

// Read implements io.Reader.
func (e *GioTextView) Read(p []byte) (int, error) {
	n, err := e.rr.ReadAt(p, e.seekCursor)
	e.seekCursor += int64(n)
	return n, err
}

// ReadAt implements io.ReaderAt.
func (e *GioTextView) ReadAt(p []byte, offset int64) (int, error) {
	return e.rr.ReadAt(p, offset)
}

// Regions returns visible regions covering the rune range [start,end).
func (e *GioTextView) Regions(start, end int, regions []Region) []Region {
	viewport := image.Rectangle{
		Min: e.scrollOff,
		Max: e.viewSize.Add(e.scrollOff),
	}
	return e.index.locate(viewport, start, end, regions)
}

func (e *GioTextView) maxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (e *GioTextView) minValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (e *GioTextView) absValue(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// SPDX-License-Identifier: Unlicense OR MIT

type lineInfo struct {
	xOff            fixed.Int26_6
	yOff            int
	width           fixed.Int26_6
	ascent, descent fixed.Int26_6
	glyphs          int
}

type glyphIndex struct {
	glyphs []text.Glyph
	// positions contain all possible caret positions, sorted by rune index.
	positions []combinedPos
	// lines contains metadata about the size and position of each line of
	// text.
	lines []lineInfo

	// currentLineMin and currentLineMax track the dimensions of the line
	// that is being indexed.
	currentLineMin, currentLineMax fixed.Int26_6
	// currentLineGlyphs tracks how many glyphs are contained within the
	// line that is being indexed.
	currentLineGlyphs int
	// pos tracks attributes of the next valid cursor position within the indexed
	// text.
	pos combinedPos
	// prog tracks the current glyph text progression to detect bidi changes.
	prog text.Flags
	// clusterAdvance accumulates the advances of glyphs in a glyph cluster.
	clusterAdvance fixed.Int26_6
	// skipPrior controls whether a text position is inserted "before" the
	// next glyph. Usually this should not happen, but the boundaries of
	// lines and bidi runs require it.
	skipPrior bool
}

// screenPos represents a character position in text line and column numbers,
// not pixels.
type screenPos struct {
	// col is the column, measured in runes.
	// FIXME: we only ever use col for start or end of lines.
	// We don't need accurate accounting, so can we get rid of it?
	col  int
	line int
}

// combinedPos is a point in the editor.
type combinedPos struct {
	// runes is the offset in runes.
	runes int

	lineCol screenPos

	// Pixel coordinates
	x fixed.Int26_6
	y int

	ascent, descent fixed.Int26_6

	// runIndex tracks which run this position is within, counted each time
	// the index processes an end of run marker.
	runIndex int
	// towardOrigin tracks whether this glyph's run is progressing toward the
	// origin or away from it.
	towardOrigin bool
}

// incrementPosition returns the next position after pos (if any). Pos _must_ be
// an unmodified position acquired from one of the closest* methods. If eof is
// true, there was no next position.
func (g *glyphIndex) incrementPosition(pos combinedPos) (next combinedPos, eof bool) {
	candidate, index := g.closestToRune(pos.runes)
	for candidate != pos && index+1 < len(g.positions) {
		index++
		candidate = g.positions[index]
	}
	if index+1 < len(g.positions) {
		return g.positions[index+1], false
	}
	return candidate, true

}

// Glyph indexes the provided glyph, generating text cursor positions for it.
func (g *glyphIndex) Glyph(gl text.Glyph) {
	g.glyphs = append(g.glyphs, gl)
	g.currentLineGlyphs++
	if len(g.positions) == 0 {
		// First-iteration setup.
		g.currentLineMin = math.MaxInt32
		g.currentLineMax = 0
	}
	if gl.X < g.currentLineMin {
		g.currentLineMin = gl.X
	}
	if end := gl.X + gl.Advance; end > g.currentLineMax {
		g.currentLineMax = end
	}
	if !g.skipPrior || gl.Flags&text.FlagTowardOrigin != g.prog || gl.Flags&text.FlagParagraphStart != 0 {
		// Set the new text progression based on that of the first glyph.
		g.prog = gl.Flags & text.FlagTowardOrigin
		g.pos.towardOrigin = g.prog == text.FlagTowardOrigin
		// Create the text position prior to the first glyph.
		pos := g.pos
		pos.x = gl.X
		pos.y = int(gl.Y)
		pos.ascent = gl.Ascent
		pos.descent = gl.Descent
		if pos.towardOrigin {
			pos.x += gl.Advance
		}
		g.pos = pos
		g.positions = append(g.positions, pos)
		g.skipPrior = true
	}
	needsNewLine := gl.Flags&text.FlagLineBreak != 0
	needsNewRun := gl.Flags&text.FlagRunBreak != 0
	breaksParagraph := gl.Flags&text.FlagParagraphBreak != 0

	// We should insert new positions if the glyph we're processing terminates
	// a glyph cluster.
	insertPositionAfter := gl.Flags&text.FlagClusterBreak != 0 && !breaksParagraph && gl.Runes > 0
	if breaksParagraph {
		// Paragraph breaking clusters shouldn't have positions generated for both
		// sides of them. They're always zero-width, so doing so would
		// create two visually identical cursor positions. Just reset
		// cluster state, increment by their runes, and move on to the
		// next glyph.
		g.clusterAdvance = 0
		g.pos.runes += int(gl.Runes)
	}
	// Always track the cumulative advance added by the glyph, even if it
	// doesn't terminate a cluster itself.
	g.clusterAdvance += gl.Advance
	if insertPositionAfter {
		// Construct the text position _after_ gl.
		pos := g.pos
		pos.y = int(gl.Y)
		pos.ascent = gl.Ascent
		pos.descent = gl.Descent
		width := g.clusterAdvance
		perRune := width / fixed.Int26_6(gl.Runes)
		adjust := fixed.Int26_6(0)
		if pos.towardOrigin {
			// If RTL, subtract increments from the width of the cluster
			// instead of adding.
			adjust = width
			perRune = -perRune
		}
		for i := 1; i <= int(gl.Runes); i++ {
			pos.x = gl.X + adjust + perRune*fixed.Int26_6(i)
			pos.runes++
			pos.lineCol.col++
			g.positions = append(g.positions, pos)
		}
		g.pos = pos
		g.clusterAdvance = 0
	}
	if needsNewRun {
		g.pos.runIndex++
	}
	if needsNewLine {
		g.lines = append(g.lines, lineInfo{
			xOff:    g.currentLineMin,
			yOff:    int(gl.Y),
			width:   g.currentLineMax - g.currentLineMin,
			ascent:  g.positions[len(g.positions)-1].ascent,
			descent: g.positions[len(g.positions)-1].descent,
			glyphs:  g.currentLineGlyphs,
		})
		g.pos.lineCol.line++
		g.pos.lineCol.col = 0
		g.pos.runIndex = 0
		g.currentLineMin = math.MaxInt32
		g.currentLineMax = 0
		g.currentLineGlyphs = 0
		g.skipPrior = false
	}
}

func (g *glyphIndex) closestToRune(runeIdx int) (combinedPos, int) {
	if len(g.positions) == 0 {
		return combinedPos{}, 0
	}
	i := sort.Search(len(g.positions), func(i int) bool {
		pos := g.positions[i]
		return pos.runes >= runeIdx
	})
	if i > 0 {
		i--
	}
	closest := g.positions[i]
	closestI := i
	for ; i < len(g.positions); i++ {
		if g.positions[i].runes == runeIdx {
			return g.positions[i], i
		}
	}
	return closest, closestI
}

func (g *glyphIndex) closestToLineCol(lineCol screenPos) combinedPos {
	if len(g.positions) == 0 {
		return combinedPos{}
	}
	i := sort.Search(len(g.positions), func(i int) bool {
		pos := g.positions[i]
		return pos.lineCol.line > lineCol.line || (pos.lineCol.line == lineCol.line && pos.lineCol.col >= lineCol.col)
	})
	if i > 0 {
		i--
	}
	prior := g.positions[i]
	if i+1 >= len(g.positions) {
		return prior
	}
	next := g.positions[i+1]
	if next.lineCol != lineCol {
		return prior
	}
	return next
}

func dist(a, b fixed.Int26_6) fixed.Int26_6 {
	if a > b {
		return a - b
	}
	return b - a
}

func (g *glyphIndex) closestToXY(x fixed.Int26_6, y int) combinedPos {
	if len(g.positions) == 0 {
		return combinedPos{}
	}
	i := sort.Search(len(g.positions), func(i int) bool {
		pos := g.positions[i]
		return pos.y+pos.descent.Round() >= y
	})
	// If no position was greater than the provided Y, the text is too
	// short. Return either the last position or (if there are no
	// positions) the zero position.
	if i == len(g.positions) {
		return g.positions[i-1]
	}
	first := g.positions[i]
	// Find the best X coordinate.
	closest := i
	closestDist := dist(first.x, x)
	line := first.lineCol.line
	// NOTE(whereswaldon): there isn't a simple way to accelerate this. Bidi text means that the x coordinates
	// for positions have no fixed relationship. In the future, we can consider sorting the positions
	// on a line by their x coordinate and caching that. It'll be a one-time O(nlogn) per line, but
	// subsequent uses of this function for that line become O(logn). Right now it's always O(n).
	for i := i + 1; i < len(g.positions) && g.positions[i].lineCol.line == line; i++ {
		candidate := g.positions[i]
		distance := dist(candidate.x, x)
		// If we are *really* close to the current position candidate, just choose it.
		if distance.Round() == 0 {
			return g.positions[i]
		}
		if distance < closestDist {
			closestDist = distance
			closest = i
		}
	}
	return g.positions[closest]
}

// makeRegion creates a text-aligned rectangle from start to end. The vertical
// dimensions of the rectangle are derived from the provided line's ascent and
// descent, and the y offset of the line's baseline is provided as y.
func makeRegion(line lineInfo, y int, start, end fixed.Int26_6) Region {
	if start > end {
		start, end = end, start
	}
	dotStart := image.Pt(start.Round(), y)
	dotEnd := image.Pt(end.Round(), y)
	return Region{
		Bounds: image.Rectangle{
			Min: dotStart.Sub(image.Point{Y: line.ascent.Ceil()}),
			Max: dotEnd.Add(image.Point{Y: line.descent.Floor()}),
		},
		Baseline: line.descent.Floor(),
	}
}

// Region describes the position and baseline of an area of interest within
// shaped text.
type Region struct {
	// Bounds is the coordinates of the bounding box relative to the containing
	// widget.
	Bounds image.Rectangle
	// Baseline is the quantity of vertical pixels between the baseline and
	// the bottom of bounds.
	Baseline int
}

// region is identical to Region except that its coordinates are in document
// space instead of a widget coordinate space.
type region = Region

// locate returns highlight regions covering the glyphs that represent the runes in
// [startRune,endRune). If the rects parameter is non-nil, locate will use it to
// return results instead of allocating, provided that there is enough capacity.
// The returned regions have their Bounds specified relative to the provided
// viewport.
func (g *glyphIndex) locate(viewport image.Rectangle, startRune, endRune int, rects []Region) []Region {
	if startRune > endRune {
		startRune, endRune = endRune, startRune
	}
	rects = rects[:0]
	caretStart, _ := g.closestToRune(startRune)
	caretEnd, _ := g.closestToRune(endRune)

	for lineIdx := caretStart.lineCol.line; lineIdx < len(g.lines); lineIdx++ {
		if lineIdx > caretEnd.lineCol.line {
			break
		}
		pos := g.closestToLineCol(screenPos{line: lineIdx})
		if int(pos.y)+pos.descent.Ceil() < viewport.Min.Y {
			continue
		}
		if int(pos.y)-pos.ascent.Ceil() > viewport.Max.Y {
			break
		}
		line := g.lines[lineIdx]
		if lineIdx > caretStart.lineCol.line && lineIdx < caretEnd.lineCol.line {
			startX := line.xOff
			endX := startX + line.width
			// The entire line is selected.
			rects = append(rects, makeRegion(line, pos.y, startX, endX))
			continue
		}
		selectionStart := caretStart
		selectionEnd := caretEnd
		if lineIdx != caretStart.lineCol.line {
			// This line does not contain the beginning of the selection.
			selectionStart = g.closestToLineCol(screenPos{line: lineIdx})
		}
		if lineIdx != caretEnd.lineCol.line {
			// This line does not contain the end of the selection.
			selectionEnd = g.closestToLineCol(screenPos{line: lineIdx, col: math.MaxInt})
		}

		var (
			startX, endX fixed.Int26_6
			eof          bool
		)
	lineLoop:
		for !eof {
			startX = selectionStart.x
			if selectionStart.runIndex == selectionEnd.runIndex {
				// Commit selection.
				endX = selectionEnd.x
				rects = append(rects, makeRegion(line, pos.y, startX, endX))
				break
			} else {
				currentDirection := selectionStart.towardOrigin
				previous := selectionStart
			runLoop:
				for !eof {
					// Increment the start position until the next logical run.
					for startRun := selectionStart.runIndex; selectionStart.runIndex == startRun; {
						previous = selectionStart
						selectionStart, eof = g.incrementPosition(selectionStart)
						if eof {
							endX = selectionStart.x
							rects = append(rects, makeRegion(line, pos.y, startX, endX))
							break runLoop
						}
					}
					if selectionStart.towardOrigin != currentDirection {
						endX = previous.x
						rects = append(rects, makeRegion(line, pos.y, startX, endX))
						break
					}
					if selectionStart.runIndex == selectionEnd.runIndex {
						// Commit selection.
						endX = selectionEnd.x
						rects = append(rects, makeRegion(line, pos.y, startX, endX))
						break lineLoop
					}
				}
			}
		}
	}
	for i := range rects {
		rects[i].Bounds = rects[i].Bounds.Sub(viewport.Min)
	}
	return rects
}