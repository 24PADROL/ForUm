document.getElementById("loginForm").addEventListener("submit", async function(event) {
    event.preventDefault(); // Empêche l'envoi classique du formulaire

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ email, password }) // Envoie les données en JSON
        });

        if (response.ok) {
            console.log("Connexion réussie !");
            window.location.href = "/web/html/home.html"; // Redirige vers la page d'accueil
        } else {
            const errorText = await response.text();
            console.error("Erreur de connexion :", errorText);
            alert("Erreur : " + errorText);
        }
    } catch (error) {
        console.error("Erreur réseau :", error);
        alert("Une erreur réseau est survenue. Veuillez réessayer.");
    }
        window.location.href = "/auth/google";
});

document.getElementById("googleLogin").addEventListener("click", function () {

});