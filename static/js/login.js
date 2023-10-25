

const form = document.querySelector('#form');
const containerForm = document.querySelector('.form-container');

addListeners();
function addListeners() {
   form.addEventListener("submit", handlerDataJSON);
};

function handlerDataJSON(e) {
   e.preventDefault();

   const email = form.querySelector('#email').value;
   const password = form.querySelector('#password').value;

   dataJSON(email, password);
};

async function dataJSON(email, password) {
   try {
      const response = await fetch("/login", {
         method: 'POST',
         headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
         },
         body: JSON.stringify({email, password})
      });

      if (!response.ok) {
         throw new Error("La solicitud no se pudo completar correctamente.");
      };

      const data = await response.json();
      if (data.success) {
         setTimeout(() => {
            window.location.href = '/';
         }, 1000);
      } else {
         messageErrorHTML('Datos incorrectos');
         
         setTimeout(() => {
            clearHTML();
         }, 3000);
      };
   } catch (error) {
      console.error("Error");
   };
};

function messageErrorHTML(message) {
   clearHTML();

   const errorHTML = document.createElement('P');
   errorHTML.textContent = message;
   errorHTML.classList.add('error-message');

   containerForm.appendChild(errorHTML);
};

function clearHTML() {
   const errorMessageHTML = containerForm.querySelector('.error-message');

   if (errorMessageHTML) {
      errorMessageHTML.remove();
   };
};
