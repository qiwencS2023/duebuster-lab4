package main

import "fmt"

type CoordinatorCache interface {
	InvalidateCache(table string, pkValue string)
	GetCache(request *GetLineRequest) (*Line, bool)
	PutCache(request *GetLineRequest, line *Line)
}

type CoordinatorCacheImpl struct {
	// key: table name, getLine query
	// value: getLine result
	cache map[string]map[string]*Line
}

func (c *CoordinatorCacheImpl) InvalidateCache(table string, pkValue string) {
	delete(c.cache[table], pkValue)
}

func (c *CoordinatorCacheImpl) GetCache(request *GetLineRequest) (*Line, bool) {
	tableName := request.Table.Name
	line, ok := c.cache[tableName][request.PrimaryKeyValue]
	return line, ok
}

func (c *CoordinatorCacheImpl) PutCache(request *GetLineRequest, line *Line) {
	tableName := request.Table.Name
	getLineQuery := fmt.Sprintf("%v", request)
	c.cache[tableName][getLineQuery] = line
}

func NewCoordinatorCache() *CoordinatorCacheImpl {
	return &CoordinatorCacheImpl{
		cache: make(map[string]map[string]*Line),
	}
}
