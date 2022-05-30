// function myFunction() {
//     document.getElementById("myDropdown").classList.toggle("show");
// }

//   // Close the dropdown if the user clicks outside of it
// window.onclick = function(event) {
// if (!event.target.matches('.dropbtn')) {
//     var dropdowns = document.getElementsByClassName("dropdown-content");
//     var i;
//     for (i = 0; i < dropdowns.length; i++) {
//     var openDropdown = dropdowns[i];
//     if (openDropdown.classList.contains('show')) {
//         openDropdown.classList.remove('show');
//     }
//     }
// }
// }
let validation_length = document.getElementById("validation_length")
let validation_upper =  document.getElementById("validation_upper")
let validation_special = document.getElementById("validation_special")
let validation_length_30 = document.getElementById("validation_length_30")
let password = document.getElementById("password_check")
let main_password = document.getElementById("password")




$(".password").on("input",function(ev){
    convert($(ev.target).val())
});

$(".password_check").on("input",function(ev){
    check($(ev.target))
});

$(".password").on("input",function(ev){
    check($(ev.target))
});

function check(){
    
    if (password.value == "" && main_password.value == "") {
        password.style.border = "1px solid #cdd0d1"
        main_password.style.border = "1px solid #cdd0d1"
    } else if (password.value === main_password.value) {
        password.style.border = "thin solid green"
        main_password.style.border = "thin solid green"
    } else {
        password.style.border = "thin solid red"
        main_password.style.border = "thin solid red"
    }
  
    
}

function convert(getNum){
    let text = getNum
    const regex = /[A-Z]/g;
    const regex_special = /[!@#\$%\^\&*\)\(+=._-]/g;
    
    if (text.length < 8) {
        validation_length.style.color = "red";
    } else if (text.length > 30 ) {
        validation_length_30.textContent = "Пароль должен содержать не более 30 символов"
        validation_length_30.style.color = "red"
    }else {
        validation_length.style.color = "green";
        validation_length_30.style.color = "green";
        validation_length_30.textContent = ""
    }
    
    if (!text.match(regex)){
        validation_upper.style.color = "red";
    } else {
        validation_upper.style.color = "green";
    }

    if (!text.match(regex_special)){
        validation_special.style.color = "red";
    } else {
        validation_special.style.color = "green";
    }
}



if (performance.navigation.type == performance.navigation.TYPE_RELOAD) {
    let login = document.getElementById("login")
    login.value = ""
} 


function ValidateForm() {
    
    if(validation_length.style.color == "red" || validation_upper.style.color == "red" || validation_special.style.color == "red" || validation_length_30.style.color == "red" || password.style.border == "thin solid red" || main_password.style.border == "thin solid red"){ 
      alert("Невалидно")
     
      return false;
    }
  
    return true;
}




function Like(ev) {
    let like = document.getElementById("like")
    let dislike = document.getElementById("dislike")
    like.style.backgroundColor = "#e06149"
    dislike.style.backgroundColor = "black"
}

function Dislike(ev) {
    let like = document.getElementById("like")
    let dislike = document.getElementById("dislike")
    like.style.backgroundColor = "black"
    dislike.style.backgroundColor = "#e06149"
}
