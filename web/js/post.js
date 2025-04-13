document.getElementById("postForm").addEventListener("submit", async function(event) {
    event.preventDefault(); // Empêche le rechargement de la page

    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;
    const image = document.getElementById("image").value;
    const category = document.getElementById("category").value;

    try {
        const response = await fetch("/post", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ title, content, image, category }) // Envoie les données en JSON
        });

        if (response.ok) {
            alert("Post créé avec succès !");
            document.getElementById("postForm").reset(); // Réinitialise le formulaire
        } else {
            const errorText = await response.text();
            alert("Erreur : " + errorText);
        }
    } catch (error) {
        console.error("Erreur réseau :", error);
        alert("Une erreur réseau est survenue. Veuillez réessayer.");
    }
});