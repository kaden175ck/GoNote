package v24

type Set[T comparable] map[T]struct{}

// 构造函数
func NewSet[T comparable](n int) Set[T] {
	m := make(map[T]struct{}, n)
	return Set[T](m)
}

// 往Set里面添加元素
func (set Set[T]) Add(ele T) {
	set[ele] = struct{}{}
}

// 获取Set的长度
func (set Set[T]) Len() int {
	return len(set)
}

// 删除元素
func (set Set[T]) Remove(ele T) {
	delete(set, ele)
}

// 判断某个元素是否存在
func (set Set[T]) Exists(ele T) bool {
	_, exists := set[ele]
	return exists
}

// 遍历Set
func (set Set[T]) Range(f func(ele T)) {
	for key := range set {
		f(key)
	}
}
