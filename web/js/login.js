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
            const data = await response.json();
            alert(data.message); // Affiche un message de succès
            window.location.href = "/accueil"; // Redirige vers la page d'accueil
        } else {
            const errorData = await response.json();
            alert(errorData.message || "Échec de la connexion");
        }
    } catch (error) {
        console.error("Erreur lors de la connexion :", error);
        alert("Une erreur est survenue. Veuillez réessayer.");
    }
});