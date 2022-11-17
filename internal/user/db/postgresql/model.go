package postgresql

type SortOptions struct {
	Field, Order string
}

//func (sort *sortOptions) EnrichQuery(query squirrel.SelectBuilder) string {
//	query.OrderBy(fmt.Sprintf("%s %s", sort.Field, sort.Order))
//}
