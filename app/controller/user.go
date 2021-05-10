package controller

import (
    "github.com/sar5430/crackmes.one/app/model"
    "log"
    "net/http"
    "sort"
    //"app/shared/session"
    "github.com/sar5430/crackmes.one/app/shared/view"

    "github.com/gorilla/context"
    "github.com/julienschmidt/httprouter"
)

type By func(p1, p2 *model.User) bool

func (by By) Sort(users []model.User) {
    ps := &userSorter{
        users: users,
        by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
    }
    sort.Sort(ps)
}

type userSorter struct {
    users []model.User
    by    func(p1, p2 *model.User) bool // Closure used in the Less method.
}

func (s *userSorter) Len() int {
    return len(s.users)
}

// Swap is part of sort.Interface.
func (s *userSorter) Swap(i, j int) {
    s.users[i], s.users[j] = s.users[j], s.users[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *userSorter) Less(i, j int) bool {
    return s.by(&s.users[i], &s.users[j])
}

// NotepadReadGET displays the notes in the notepad
func UserGET(w http.ResponseWriter, r *http.Request) {
    var params httprouter.Params
    params = context.Get(r, "params").(httprouter.Params)
    name := params.ByName("name")

    user, err := model.UserByName(name)
    if err != nil {
        log.Println(err)
        Error404(w, r)
        return
    }

    crackmes, err := model.CrackmesByUser(name)
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    nbCrackmes, err := model.CountCrackmesByUser(name)
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    nbSolutions, err := model.CountSolutionsByUser(name)
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    nbComments, err := model.CountCommentsByUser(name)
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    solutions, err := model.SolutionsByUser(name)
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    comments, err := model.CommentsByUser(name)
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    solutionsext := make([]model.SolutionExtended, len(solutions))

    for i := range solutions {
        solutionsext[i].Solution = &solutions[i]
        solutionsext[i].Crackmeshexid = (&solutions[i]).CrackmeId.Hex()
        tmpcrackme, err := model.CrackmeByHexId(solutionsext[i].Crackmeshexid)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }
        solutionsext[i].Crackmename = tmpcrackme.Name
    }

    for i, c := range crackmes {
        crackmes[i].NbComments, err = model.CountCommentsByCrackme(c.HexId)

        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }

        crackmes[i].NbSolutions, err = model.CountSolutionsByCrackme(c.HexId)

        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }
    }

    user.NbCrackmes = nbCrackmes
    user.NbSolutions = nbSolutions
    user.NbComments = nbComments

    // Display the view
    v := view.New(r)
    v.Name = "user/read"
    v.Vars["username"] = user.Name
    v.Vars["NbCrackmes"] = user.NbCrackmes
    v.Vars["NbSolutions"] = user.NbSolutions
    v.Vars["NbComments"] = user.NbComments
    v.Vars["crackmes"] = crackmes
    v.Vars["solutions"] = solutionsext
    v.Vars["comments"] = comments
    v.Render(w)
}

func UsersGET(w http.ResponseWriter, r *http.Request) {

    users, err := model.AllUsersVisible()
    name := func(p1, p2 *model.User) bool {
        return p1.Name < p2.Name
    }
    By(name).Sort(users)

    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    for _, user := range users {
        nbSolutions, err := model.CountSolutionsByUser(user.Name)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }

        nbComments, err := model.CountCommentsByUser(user.Name)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }
        nbCrackmes, err := model.CountCrackmesByUser(user.Name)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }
        user.NbSolutions = nbSolutions
        user.NbComments = nbComments
        user.NbCrackmes = nbCrackmes
    }
    // Display the view
    v := view.New(r)
    v.Name = "users/read"
    v.Vars["users"] = users
    v.Render(w)
}
