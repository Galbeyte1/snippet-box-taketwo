package transport

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
)

/*
File Server and Template Parsing example

Step	Browser action												Go server reaction
1 		Browser requests http://localhost:8080/static/style.css		FileServer looks inside ./static/style.css and streams it
2		Browser requests http://localhost:8080/						HomeHandler runs, tmpl.Execute renders index.html with Name = "Alice"
3		index.html links to static/style.css						Browser automatically requests style.css, served again by FileServer

Building a server cache for templates WHILE file server caching handled by os/browser
*/

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
	data := templates.NewTemplateData(r)
	data.Snippets = snippets
	app.Render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}
	snippet, err := app.Snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotRecord) {
			http.NotFound(w, r)
		} else {
			app.ServerError(w, r, err)
		}
		return
	}

	data := templates.NewTemplateData(r)
	data.Snippet = snippet

	app.Render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := templates.NewTemplateData(r)

	app.Render(w, r, http.StatusOK, "create.tmpl", data)

}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := `الإمام الشافعي`
	//content := `دع الأيام تفعل ما تشاءُ وطب نفساً إذا حكم القضاءُ ولا تجزع لحادثة الليالي فما لحوادث الدنيا بقاءُ`
	expires := 20
	//content := `فدَعهُ ولا تكثرْ عليه التأسُّفا ففي الناسِ أبدالٌ وفي الترك راحةٌ وفي القلبِ صبرٌ للحبيبِ ولو جفا فما كلُّ من تهواهُ يهواكَ قلبُهُ ولا كلُّ مَن صافيته لك قد صفا إذا لم يكن صفوُ الودادِ طبيعةً فلا خيرَ في ودٍّ يجيءُ تكلُّفا ولا خيرَ في خلٍّ يخونُ خليلهُ ويلقاهُ من بعدِ المودّةِ بالجفا ويُنكِرُ عيشًا قد تقادمَ عهدهُ ويُظهِرُ سرًّا كان بالأمسِ قد خفا سلامٌ على الدنيا إذا لم يكن بها صديقٌ صدوقٌ صادقُ الوعدِ منصفا`
	//content := `سهرت أعينٌ ونامت عيونُ في أمورٍ تكون أو لا تكونُ فادرأ الهمّ ما استطعت عن النفس فحملانك الهموم جنونُ إن رباً كفاك بالأمس ما كان سيكفيك في غدٍ ما يكونُ`
	content := `ارحل بنفسك من أرضٍ تضامُ بها ولا تكنْ من فراقِ الأهلِ في حرقِ فالعنبر الخام روثٌ في مواطنه وفي التغربِ محمولٌ على العنقِ والكحلُ نوعٌ من الأحجارِ تنظره في أرضه وهو مرميٌ على الطرقِ لما تغربَ حاز الفضل أجمعه فصار يُحملُ بين الجفنِ والحدقِ`
	id, err := app.Snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
	log.Println("Successfully inserted snippet with ID", id) // <-- And this
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
