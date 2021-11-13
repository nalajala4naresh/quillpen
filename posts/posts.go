package posts

import (
	"fmt"
	"net/http"
	"quillpen/models"
	"quillpen/storage"

	"html/template"
	"strings"
	"errors"

	"github.com/gorilla/mux"
	md "github.com/russross/blackfriday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/html"
)



var templates *template.Template
func init(){

	templates = template.Must(template.ParseFiles("templates/posts.html"))



}

func extract_header_node(html_string string) (*html.Node, error) {
    
	var header *html.Node
	node, err := html.Parse(strings.NewReader(html_string))
	if err!= nil {
		panic("Invalid HTML String")
	}
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
	
        if node.Type == html.ElementNode && (node.Data == "h1" || node.Data == "h2" || node.Data == "h3") {
            header = node
            return
        }
    for child := node.FirstChild; child != nil; child = child.NextSibling {
            crawler(child)
        }
    }

	crawler(node)
	if header != nil {
        return header, nil
    }
    return nil, errors.New("Missing <body> in the node tree")



}




func Write_post(resp http.ResponseWriter, req *http.Request) {


	req.ParseForm()
    extensions := 0
	var post models.Post
	fmt.Println(req.Form)
	obj_id := primitive.NewObjectID()
	post.PostId = obj_id.Hex()
	post.MD_Content = []byte(req.Form["md_content"][0])
    
	html_render := md.HtmlRenderer(md.HTML_SKIP_HTML,"","")
	html_bytes := md.Markdown(post.MD_Content,html_render,extensions)
	
    node, header_extract_error := extract_header_node(string(html_bytes))

	if header_extract_error != nil {
		fmt.Println("header missing ")
	}

	post.Title = node.FirstChild.Data

	
	// now save the html bytes to Storage
	_ ,err := storage.InsertOne(&post,storage.POSTS_COLLECTIONS)
	if err != nil {
		panic("Unable to write the post")

	}


	// redirect user to the same post he wrote, to make any modifications
	Read_Posts(resp,req )
	




}

func Read_Posts(resp http.ResponseWriter, req *http.Request)  {


	query := bson.D{}
	extensions := 0
	result_set := storage.FindMany(query, storage.POSTS_COLLECTIONS,10)
	html_render := md.HtmlRenderer(md.HTML_SKIP_HTML,"","")
    
	if result_set.Posts == nil {


	
	}
	
	for _, post := range result_set.Posts {

		if post.MD_Content == nil {
			continue
		}
        
		html := md.Markdown(post.MD_Content,html_render,extensions)
		post.HTML_Content = template.HTML(html)

	}
	templates.ExecuteTemplate(resp,"posts/list",result_set.Posts)



}

func Read_A_Post(resp http.ResponseWriter, req *http.Request) {
    
	extensions := 0
	uri_params  := mux.Vars(req)
	
	var result models.Result
	result = storage.FindOne(bson.M{"post_id":uri_params["postid"]},storage.POSTS_COLLECTIONS)
	
	html_render := md.HtmlRenderer(md.HTML_SKIP_HTML,"","")

	if result.Post == nil {
		http.NotFound(resp, req)
		return
	}
	
	html := md.Markdown(result.Post.MD_Content,html_render,extensions)
	result.Post.HTML_Content = template.HTML(html)

	

	templates.ExecuteTemplate(resp,"posts/one", result.Post)



}
