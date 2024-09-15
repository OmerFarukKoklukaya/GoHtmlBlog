/*!
* Start Bootstrap - Blog Home v5.0.9 (https://startbootstrap.com/template/blog-home)
* Copyright 2013-2023 Start Bootstrap
* Licensed under MIT (https://github.com/StartBootstrap/startbootstrap-blog-home/blob/master/LICENSE)
*/
// This file is intentionally blank
// Use this file to add JavaScript to your project

const baseURL = "http://localhost:3000";

function logout(){
    fetch(baseURL+"/api/logout",{
        method: 'POST'
    }).then(response => {
        window.location.assign(baseURL);
    })
}

function sendBlog(id, method) {
    const title = document.getElementById('title');
    const body = document.getElementById('editor');
    let sum = body.firstElementChild.textContent.slice(0, 50)

    const formData = new FormData();
    if (sum.length === 50) {
        sum += '...'
    }

    formData.append('user_id', '{{.Meta.AuthedUser.ID}}')
    formData.append('title', title.value)
    formData.append('body', body.firstElementChild.innerHTML)
    formData.append('summary', sum)

    const urlEncodedData = new URLSearchParams(formData).toString();
    fetch(baseURL+'/api/blog/' + id, {
        method: method,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
        },
        body: urlEncodedData
    }).then((resposne) => resposne.json()).then(data => {
        const image = document.getElementById('image').files[0]
        const imgForm = new FormData()
        imgForm.append('image', image)
        fetch(baseURL+'/api/files/images/blog/' + data.data.id, {
            method: 'POST',
            body: imgForm
        }).then( (response) => {
            window.location.assign(baseURL+'/blogs/' + id)
        })

    })
}

function deleteImage(id){
    fetch(baseURL + '/api/files/images/blog/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.assign(baseURL + '/editor/' + id)
    })
}

function deleteBlog(id){
    fetch(baseURL + '/api/blog/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.assign(baseURL + '/blogs/')
    })
}

function updateUserName(id) {
    const name = document.getElementById('name');
    const data = JSON.stringify({name: name.value, id: id})


    fetch(baseURL+'/api/user' , {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    }).then(data => {
         window.location.assign(baseURL+'/user/'+id)
    })
}


function updatePassword(id) {
    const password = document.getElementById('password');
    const new_password = document.getElementById('new_password');
    const new_password_verification = document.getElementById('new_password_verification');
    const formData = new FormData();

    formData.append('password', password.value)
    formData.append('new_password', new_password.value)
    formData.append('new_password_verification', new_password_verification.value)

    const urlEncodedData = new URLSearchParams(formData).toString();

    fetch(baseURL+'/api/user/password' , {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
        },
        body: urlEncodedData
    }).then((resposne) => {
        fetch(baseURL+'/api/logout', {
            method: 'POST'
        }).then(response => {
           window.location.assign(baseURL+'/login')
        })
    })
}