

const form2 = document.querySelector('#form');
const formContainer = document.querySelector('.form-container');

addListeners2();
function addListeners2() {
   form.addEventListener("submit", handlerDataJSON2);
};

function handlerDataJSON2(e) {
   e.preventDefault();

   const email = document.querySelector('#email').value;
   const password = document.querySelector('#password').value;
   const confirmPassword = document.querySelector('#confirm-password').value;

   if (password !== confirmPassword) {
      createMessageHTML2("Contraseñas no coinciden");
   } else {
      dataJSON2(email, password);
   };
};

async function dataJSON2(email, password) {
   try {
      const response = await fetch('/register', {
         method: 'POST',
         headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
         },
         body: JSON.stringify({email, password}),
      });

      if (!response.ok) {
         throw new Error("La solicitud no se pudo completar correctamente.");
      };

      const data = await response.json();
      if (data.success) {
         createMessageHTML2("Registrado correctamente");
         setTimeout(() => {
            window.location.href = '/';
         }, 2000);
      };
   } catch (error) {
      console.error("Error");
   }
}

function createMessageHTML2(message) {
   clearHTML2()

   const errorHTML = document.createElement('P');
   errorHTML.textContent = message;
   errorHTML.classList.add('error-message');
   formContainer.appendChild(errorHTML);
};

function clearHTML2() {
   const result = formContainer.querySelector('.error-message');

   if (result) {
      result.remove();
   };
};