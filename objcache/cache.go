package objcache

// ТРЕБУЕТСЯ: написать реализацию пула кешируемых элементов - по сути КЭШ. Кол-во требуемых элементов в кэше неизвестно.
// Работа кэша происходит по следующему алгоритму: по мере необходимости из кеша запрашиваются элементы методом Get,
// каждый из полученных элементов имеет собственный адрес в памяти. После того как работа с объектами выполнена,
// вызывается метод Clear и все запрошенные ранее элементы считаются свободными и могут быть заново выданы в следующий раз.

// Оценивается скорость работы и количество аллокаций. При повторном запросе элементов, аллокаций в памяти быть не должно.

// Cache примерный принцип реализации пула. Не проходит всех тестов.
type Cache[V any] struct {
	d [1000]V
	i int
}

// Get выдает адрес к переиспользуемому участку памяти, этот адрес будет зарезервирован до следующего вызова Clear и следующий вызов Get выдаст другой участок памяти.
func (c *Cache[V]) Get() *V {
	c.i++
	if c.i < 1001 {
		return &c.d[c.i-1]
	}
	var v V
	return &v
}

// Clear снимает резервирование и помечает все адреса памяти, как свободные. После вызова этого метода Get будет выдавать адреса с самого первого.
func (c *Cache[V]) Clear() {
	c.i = 0
}