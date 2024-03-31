package object

func NewCookbook() *Cookbook {
	s := make(map[string]Object)
	return &Cookbook{page: s}
}

type Cookbook struct {
	page map[string]Object
}

func (c *Cookbook) Get(name string) (Object, bool) {
	obj, ok := c.page[name]
	return obj, ok
}

func (c *Cookbook) Set(name string, val Object) Object {
	c.page[name] = val
	return val
}
