package app

import (
	"github.com/ezeoleaf/tblogs/api"
	"github.com/ezeoleaf/tblogs/helpers"
	"github.com/ezeoleaf/tblogs/models"
	"github.com/gdamore/tcell"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
)

var listBlogs *tview.List
var listPosts *tview.List
var blogs models.Blogs
var blogPage *tview.Flex

func (a *App) generateBlogsList() {
	blogPage = tview.NewFlex()

	blogs = api.GetBlogs(a.Config)

	listBlogs.ShowSecondaryText(false)
	listBlogs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			a.followBlogs()
			return nil
		} else if event.Key() == tcell.KeyCtrlF {
			pages.ShowPage(blogsModalName)
			//TODO: Search blogs
		}
		return event
	})

	for _, blog := range blogs.Blogs {
		r := emptyRune
		isIn, _ := helpers.IsIn(blog.ID, a.Config.APP.FollowingBlogs)
		if isIn {
			r = followRune
		}
		listBlogs.AddItem(blog.Name, blog.Company, r, emptyFunc)
	}

	listPosts = getList()
	listPosts.SetDoneFunc(func() {
		a.App.SetFocus(listBlogs)
	})
	listBlogs.SetSelectedFunc(func(x int, s string, s1 string, r rune) {
		listPosts.Clear()
		blogID := blogs.Blogs[x].ID
		posts := api.GetPostsByBlog(blogID, a.Config)
		for _, post := range posts.Posts {
			r := emptyRune
			isIn, _ := helpers.IsHash(post.Hash, a.Config.APP.SavedPosts)
			if isIn {
				r = savedRune
			}
			listPosts.AddItem(post.Title, post.Published, r, emptyFunc)
		}

		listPosts.SetSelectedFunc(func(x int, s string, s1 string, r rune) {
			post := posts.Posts[x]
			browser.OpenURL(post.Link)
		})

		listPosts.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlS {

				x := listPosts.GetCurrentItem()

				post := posts.Posts[x]

				r := emptyRune
				isIn, ix := helpers.IsHash(post.Hash, a.Config.APP.SavedPosts)
				if !isIn {
					r = savedRune
					a.Config.APP.SavedPosts = append(a.Config.APP.SavedPosts, post)
				} else {
					a.Config.APP.SavedPosts = append(a.Config.APP.SavedPosts[:ix], a.Config.APP.SavedPosts[ix+1:]...)
				}
				a.Config.UpdateConfig()
				updateItemList(listPosts, x, post.Title, post.Published, r, emptyFunc)
				generateSavedPosts()
				return nil
			}
			return event
		})
		a.App.SetFocus(listPosts)
	})

	blogPage.AddItem(listBlogs, 0, 1, true).
		AddItem(listPosts, 0, 1, false)
}

func (a *App) followBlogs() {

	x := listBlogs.GetCurrentItem()

	blog := blogs.Blogs[x]

	r := emptyRune
	isIn, ix := helpers.IsIn(blog.ID, a.Config.APP.FollowingBlogs)
	if !isIn {
		r = followRune
		a.Config.APP.FollowingBlogs = append(a.Config.APP.FollowingBlogs, blog.ID)
	} else {
		a.Config.APP.FollowingBlogs = append(a.Config.APP.FollowingBlogs[:ix], a.Config.APP.FollowingBlogs[ix+1:]...)
	}
	a.Config.UpdateConfig()

	updateItemList(listBlogs, x, blog.Name, blog.Company, r, emptyFunc)
	generateHomeList()
}

func (a *App) blogsPage(nextSlide func()) (title string, content tview.Primitive) {
	pages = tview.NewPages()

	listBlogs = getList()

	a.generateBlogsList()

	pages.AddPage("blogs", blogPage, true, true)

	return blogsSection, pages
}

func searchBlogs() {
	//TODO: Add ability to search
}

func cancelSearchBlogs() {
	//TODO: Add ability to search
}
