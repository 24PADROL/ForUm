document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    fetch('/login', {  // Vérifie cette URL, elle doit pointer vers un serveur, pas un fichier HTML
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email: email, password: password })
    })
    .then(response => {
        console.log("Statut de la réponse:", response.status); // Debug du statut HTTP
        console.log("Headers de la réponse:", response.headers); // Debug des headers

        return response.text().then(text => {
            console.log("Réponse brute:", text); // Affiche la réponse brute

            if (response.ok) {
                window.location.href = '/accueil.html';
            } else {
                alert(text);
            }
        });
    })
    .catch(error => console.error('Erreur:', error));
});
