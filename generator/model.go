package main

type smw struct {
	Characters []string
	Actions    []string
	Palette    map[byte]string
	ASCII      map[string][]string
}

type color struct {
	chars   []byte
	index   int
	mapping map[string]byte
}

func (c *color) char(rgb string) byte {
	if b, ok := c.mapping[rgb]; ok {
		return b
	}

	if c.index >= len(c.chars) {
		return chars[0]
	}

	c.mapping[rgb] = c.chars[c.index]
	c.index++
	return c.mapping[rgb]
}
