document.getElementById("registerForm").addEventListener("submit", async function(event) {
    event.preventDefault();

    const username = document.getElementById("username").value;
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, email, password })
        });

        if (response.ok) {
            const data = await response.json();
            window.location.href = data.redirect; // Redirige vers /home
        } else {
            const errorText = await response.text();
            alert("Erreur : " + errorText);
        }
    } catch (error) {
        console.error("Erreur réseau :", error);
        alert("Une erreur réseau est survenue. Veuillez réessayer.");
    }
});