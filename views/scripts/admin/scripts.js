/*!
    * Start Bootstrap - SB Admin v7.0.7 (https://startbootstrap.com/template/sb-admin)
    * Copyright 2013-2023 Start Bootstrap
    * Licensed under MIT (https://github.com/StartBootstrap/startbootstrap-sb-admin/blob/master/LICENSE)
    */
    // 
// Scripts
//
const baseURL = "http://localhost:3000";
window.addEventListener('DOMContentLoaded', event => {
    // Toggle the side navigation
    const sidebarToggle = document.body.querySelector('#sidebarToggle');
    if (sidebarToggle) {
        // Uncomment Below to persist sidebar toggle between refreshes
        // if (localStorage.getItem('sb|sidebar-toggle') === 'true') {
        //     document.body.classList.toggle('sb-sidenav-toggled');
        // }
        sidebarToggle.addEventListener('click', event => {
            event.preventDefault();
            document.body.classList.toggle('sb-sidenav-toggled');
            localStorage.setItem('sb|sidebar-toggle', document.body.classList.contains('sb-sidenav-toggled'));
        });
    }


    const datatablesSimple = document.getElementById('datatablesSimple');
    if (datatablesSimple) {
        new simpleDatatables.DataTable(datatablesSimple);
    }
});
////////

////////
function selectUserRole(user) {
    const selectElement = document.getElementById('role');
    for (let i = 0; i < selectElement.options.length; i++) {
        selectElement.options[i].selected = false;
        if (user.role_id.toString() === selectElement.options[i].value){
            selectElement.options[i].selected = true
        }
    }
}

function sendUser(id, method) {
    const name = document.getElementById('name');
    const password = document.getElementById('password');
    const role_id = document.getElementById('role');
    const formData = new FormData();


    formData.append('name', name.value)
    formData.append('password', password.value)
    formData.append('role_id', role_id.value)


    const urlEncodedData = new URLSearchParams(formData).toString();
    fetch(baseURL+'/api/user/' + id, {
        method: method,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
        },
        body: urlEncodedData
    }).then((resposne) => {
        window.location.replace(baseURL+'/admin/users/' + id)
    })
}

function deleteUser(id){
    fetch(baseURL + '/api/user/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.reload()
    })
}
////////


////////
function selectRolePermission(role) {
    role.Permissions.forEach((permission) => {
        console.log(permission);

        document.getElementById(permission.name).checked = true;
    })

    for (let i = 0; i < selectElement.options.length; i++) {
        selectElement.options[i].selected = false;
        if (user.role_id.toString() === selectElement.options[i].value){
            selectElement.options[i].selected = true
        }
    }
}

function sendRole(id, method) {
    const name = document.getElementById('name');
    const cb = document.querySelectorAll('input[name="permission"]:checked');
    const permissions = []
    cb.forEach((checkbox) => {permissions.push(Number(checkbox.value));});
    console.log(permissions);

    /*  formData.append('name', name.value)
        formData.append('permissions', permissions.value)
        const urlEncodedData = new URLSearchParams(formData).toString(); */

    const data = JSON.stringify({name: name.value, permissions: permissions})
    console.log(data)

    fetch(baseURL+'/api/roles/' + id, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    }).then((resposne) => {
       window.location.replace(baseURL+'/admin/roles/' + id)
    })
}

function deleteRole(id){
    fetch(baseURL + '/api/role/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.reload()
    })
}
////////

function sendPermission(id, method) {
    const name = document.getElementById('name');
    const data = JSON.stringify({name: name.value})
    fetch(baseURL+'/api/permissions/' + id, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    }).then((response) =>{
        if(response.status === 200){
            window.location.replace(baseURL+'/admin/permissions/' + id)
        }
        else {
            response.json().then(data =>{

                document.getElementById('error').innerText = data;
                document.getElementById('breakBeforeErr').hidden = false;
            });
        }
    })

}

function deletePermission(id){
    fetch(baseURL + '/api/permissions/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.reload()
    })
}
////////

////////
function deleteImage(id){
    fetch(baseURL + '/api/files/images/blog/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.reload()
    })
}

function sendBlog(id, method) {
    const title = document.getElementById('title');
    const body = document.getElementById('editor');
    const user_id = document.getElementById('user_id');
    const formData = new FormData();

    const sum = document.getElementById('summary');
    console.log(user_id.value)
    if (sum.value === ""){
        sum.value = document.getElementById('editor').firstElementChild.textContent.slice(0, 50)
    if (sum.length === 50) {
        sum.value += '...'
    }
    }


    formData.append('user_id', user_id.value)
    formData.append('title', title.value)
    formData.append('body', body.firstElementChild.innerHTML)
    formData.append('summary', sum.value)

    const urlEncodedData = new URLSearchParams(formData).toString();

    fetch(baseURL+'/api/blog/' + id, {
        method: method,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        },
        body: urlEncodedData
    }).then((resposne) => resposne.json()).then(data => {
        const image = document.getElementById('image').files[0]
        if (image !== null) {
            const imgForm = new FormData()
            imgForm.append('image', image)
            fetch(baseURL + '/api/files/images/blog/' + data.data.id, {
                method: 'POST',
                body: imgForm
            }).then((response) => {
                window.location.assign(baseURL + '/admin/blogs/' + id)
            })
        }else {
            window.location.assign(baseURL + '/admin/blogs/' + id)
        }
    })
}

function deleteBlog(id){
    fetch(baseURL + '/api/blog/' + id,{
        method: 'DELETE',
    }).then((response)=>{
        window.location.assign(baseURL + '/admin/blogs/')
    })
}
////////


function logout(){
    fetch(baseURL+"/api/logout",{
        method: 'POST'
    }).then(response => {
        window.location.assign(baseURL);
    })
}

function controlAuthorizationAndHide(isHavePermission){
    console.log(!isHavePermission)
    let elements;
    if (!isHavePermission) {
        elements = document.getElementsByClassName('hideable_control')
        console.log(elements)
        for (let i = 0; i < elements.length; i++) {
            elements[i].style.display = 'none'
            console.log(elements[i])
        }

    }
}

