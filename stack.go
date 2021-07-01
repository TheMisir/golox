package main

type ResolverStack struct {
	list []map[string]bool
}

func (s *ResolverStack) Push(m map[string]bool) {
	s.list = append(s.list, m)
}

func (s *ResolverStack) Pop() {
	s.list = s.list[:len(s.list)-1]
}

func (s *ResolverStack) Peek() map[string]bool {
	return s.list[len(s.list)-1]
}

func (s *ResolverStack) PeekAll() []map[string]bool {
	return s.list
}

func (s *ResolverStack) Clear() {
	s.list = make([]map[string]bool, 0)
}

func (s *ResolverStack) IsEmpty() bool {
	return len(s.list) == 0
}

func (s *ResolverStack) IsNotEmpty() bool {
	return !s.IsEmpty()
}

func (s *ResolverStack) Size() int {
	return len(s.list)
}

func (s *ResolverStack) Get(i int) map[string]bool {
	return s.list[i]
}
