package object

type Cookbook struct {
	page          map[string]Object
	extended_from *Cookbook
}

func NewCookbook() *Cookbook {
	s := make(map[string]Object)
	return &Cookbook{page: s, extended_from: nil}
}

func NewExtendedCookbook(extended_from *Cookbook) *Cookbook {
	book := NewCookbook()
	book.extended_from = extended_from
	return book
}

func (c *Cookbook) Get(name string) (Object, bool) {
	obj, ok := c.page[name]
	if !ok && c.extended_from != nil {
		obj, ok = c.extended_from.Get(name)
	}
	return obj, ok
}

func (c *Cookbook) Set(name string, val Object) Object {
	c.page[name] = val
	return val
}
