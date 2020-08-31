package main

import (
	"github.com/ezeoleaf/tblogs/api"
	"github.com/rivo/tview"
)

func Blogs(nextSlide func()) (title string, content tview.Primitive) {

	code := tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true)
	code.SetBorderPadding(1, 1, 2, 0)

	listBlogs := tview.NewList()

	// basic := func() {
	// 	table.SetBorders(false).
	// 		SetSelectable(false, false).
	// 		SetSeparator(' ')
	// 	code.Clear()
	// 	fmt.Fprint(code, tableBasic)
	// }

	// separator := func() {
	// 	table.SetBorders(false).
	// 		SetSelectable(false, false).
	// 		SetSeparator(tview.Borders.Vertical)
	// 	code.Clear()
	// 	fmt.Fprint(code, tableSeparator)
	// }

	// borders := func() {
	// 	table.SetBorders(true).
	// 		SetSelectable(false, false)
	// 	code.Clear()
	// 	fmt.Fprint(code, tableBorders)
	// }

	// selectRow := func() {
	// 	table.SetBorders(false).
	// 		SetSelectable(true, false).
	// 		SetSeparator(' ')
	// 	code.Clear()
	// 	fmt.Fprint(code, tableSelectRow)
	// }

	// selectColumn := func() {
	// 	table.SetBorders(false).
	// 		SetSelectable(false, true).
	// 		SetSeparator(' ')
	// 	code.Clear()
	// 	fmt.Fprint(code, tableSelectColumn)
	// }

	// selectCell := func() {
	// 	table.SetBorders(false).
	// 		SetSelectable(true, true).
	// 		SetSeparator(' ')
	// 	code.Clear()
	// 	fmt.Fprint(code, tableSelectCell)
	// }

	// navigate := func() {
	// 	app.SetFocus(table)
	// 	table.SetDoneFunc(func(key tcell.Key) {
	// 		app.SetFocus(list)
	// 	}).SetSelectedFunc(func(row int, column int) {
	// 		app.SetFocus(list)
	// 	})
	// }

	b := api.GetBlogs()
	listBlogs.SetBorderPadding(1, 1, 2, 2)
	listBlogs.ShowSecondaryText(false)
	for _, blog := range b.Blogs {
		listBlogs.AddItem(blog.Name, blog.Company, ' ', func() {
			return
		})
	}

	listPosts := tview.NewList()
	listPosts.SetBorderPadding(1, 1, 2, 2)
	listPosts.SetDoneFunc(func() {
		app.SetFocus(listBlogs)
	})

	listBlogs.SetSelectedFunc(func(x int, s string, s1 string, r rune) {
		listPosts.Clear()
		blogID := b.Blogs[x].ID
		posts := api.GetPostsByBlog(blogID)
		// code.Clear()
		for _, post := range posts.Posts {
			listPosts.AddItem(post.Title, post.Published, '-', func() {
				return
			})
		}
		app.SetFocus(listPosts)
		// fmt.Fprint(code, posts)
	})

	return "Blogs", tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listBlogs, 0, 1, true), 0, 1, true).
		AddItem(listPosts, 100, 1, false)
}
