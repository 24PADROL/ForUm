document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM fully loaded and parsed');
    document.getElementById('registerForm').addEventListener('submit', function(event) {
        event.preventDefault();
        console.log('Form submit intercepted');
        
        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username: username, email: email, password: password })
        })
        .then(response => {
            console.log('Response status:', response.status);
            if (response.status === 201) {
                window.location.href = '/web/html/home.html';
            } else {
                response.text().then(text => {
                    console.error('Registration failed:', text);
                    alert(text);
                });
            }
        })        
        })
        .catch(error => {
            console.error('Erreur:', error);
        });
});
