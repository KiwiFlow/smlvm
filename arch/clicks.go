package arch

import (
	"container/list"
	"fmt"

	"e8vm.io/e8vm/arch/vpc"
)

type click struct {
	line uint8
	col  uint8
}

type clicks struct {
	q      *list.List
	p      *pageOffset
	send   vpc.Sender
	intBus intBus
}

func newClicks(p *page, i intBus, s vpc.Sender) *clicks {
	return &clicks{
		q:      list.New(),
		p:      &pageOffset{p, clicksBase},
		send:   s,
		intBus: i,
	}
}

func (c *clicks) addClick(line, col uint8) error {
	if line > 24 {
		return fmt.Errorf("line too big: %d", line)
	}
	if col > 80 {
		return fmt.Errorf("col too big: %d", col)
	}

	if c.q.Len() >= 16 {
		return fmt.Errorf("click event queue full")
	}

	c.q.PushBack(&click{line: line, col: col})
	c.send.Send([]byte{byte(line), byte(col)})

	return nil
}

func (c *clicks) Tick() {
	if c.q.Len() == 0 {
		return
	}

	if c.p.readByte(0) != 0 {
		return
	}

	front := c.q.Front()
	pos := front.Value.(*click)
	c.q.Remove(front)

	buf := []byte{1, 0, pos.line, pos.col}
	c.p.writeWord(0, Endian.Uint32(buf))
}
