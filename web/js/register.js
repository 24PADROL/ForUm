console.log("Script register.js chargé avec succès !");
document.getElementById("registerForm").addEventListener("submit", async function(event) {
    event.preventDefault(); // Empêche l'envoi classique du formulaire

    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, email, password }) // Envoie les données en JSON
        });

        if (response.ok) {
            const data = await response.json();
            alert(data.message); // Affiche un message de succès
            window.location.href = "/home"; // Redirige vers la page d'accueil
        } else {
            const errorData = await response.json();
            alert(errorData.message || "Échec de l'inscription");
        }
    } catch (error) {
        console.error("Erreur lors de l'inscription :", error);
        alert("Une erreur est survenue. Veuillez réessayer.");
    }
});