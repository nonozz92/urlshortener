document.addEventListener("DOMContentLoaded", function () {
  const loginButton = document.getElementById("loginButton");
  const logoutButton = document.getElementById("logoutButton");
  const loginStatus = document.getElementById("loginStatus");
  const form = document.querySelector("form");

  function updateLoginUI(isLoggedIn) {
    if (isLoggedIn) {
      if (loginStatus) loginStatus.textContent = "Connecté";
      if (loginButton) loginButton.style.display = "none";
      if (logoutButton) logoutButton.style.display = "block";
    } else {
      if (loginStatus) loginStatus.textContent = "Déconnecté";
      if (loginButton) loginButton.style.display = "block";
      if (logoutButton) logoutButton.style.display = "none";
    }
  }

  if (logoutButton) {
    logoutButton.addEventListener("click", function () {
      localStorage.removeItem("authToken");
      updateLoginUI(false);
      window.location.href = "home.html";
    });
  }

  if (form) {
    form.addEventListener("submit", function (event) {
      event.preventDefault();
      const submitButton = form.querySelector('button[type="submit"]');
      submitButton.disabled = true;

      const formData = {
        username: form.querySelector('input[type="text"]').value,
        password: form.querySelector('input[type="password"]').value,
      };

      const endpoint =
        form.id === "loginForm"
          ? "http://localhost:8000/login"
          : "http://localhost:8000/register";

      fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error(
              "Erreur lors de la connexion ou de l'enregistrement."
            );
          }
          return response.json();
        })
        .then((data) => {
          if (form.id === "registerForm") {
            alert("Compte créé avec succès!");
            setTimeout(() => {
              window.location.href = "login.html";
            }, 2000);
          } else {
            localStorage.setItem("authToken", data.token);
            updateLoginUI(true);
            window.location.href = "home.html";
          }
        })
        .catch((error) => {
          console.error("Error:", error);
          alert("Erreur lors du processus: " + error.message);
          submitButton.disabled = false;
        });
    });
  }

  updateLoginUI(!!localStorage.getItem("authToken"));
});

document.getElementById("shortenUrl").addEventListener("click", function () {
  const longUrl = document.getElementById("longUrl").value;
  const authToken = localStorage.getItem("authToken");

  const headers = {
    "Content-Type": "application/json",
  };
  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  if (longUrl) {
    fetch("http://localhost:8000/shorten", {
      method: "POST",
      headers: headers,
      body: JSON.stringify({ longUrl: longUrl, userToken: authToken }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Échec de la requête: " + response.statusText);
        }
        return response.json();
      })
      .then((data) => {
        if (data.shortUrl) {
          alert("URL courte générée : " + data.shortUrl);
        } else {
          alert("Erreur lors de la génération de l'URL courte.");
        }
      })
      .catch((error) => {
        console.error("Erreur:", error);
        alert("Erreur lors de la communication avec le serveur.");
      });
  } else {
    alert("Veuillez entrer une URL à raccourcir.");
  }
});

document
  .getElementById("redirectToLongUrl")
  .addEventListener("click", function () {
    const shortUrl = document.getElementById("shortUrl").value;
    if (shortUrl) {
      const encodedShortUrl = encodeURIComponent(shortUrl);
      const authToken = localStorage.getItem("authToken");

      fetch(`http://localhost:8000/resolve?url=${encodedShortUrl}`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      })
        .then((response) => {
          if (!response.ok) {
            console.error("Response status:", response.status);
            if (response.status === 401) {
              throw new Error("Unauthorized: Invalid or missing token");
            } else if (response.status === 403) {
              throw new Error(
                "Cette URL est privée et vous ne pouvez pas y accéder."
              );
            } else if (response.status === 404) {
              throw new Error("Aucun URL trouvé");
            } else {
              throw new Error("Cannot access URL");
            }
          }
          return response.json();
        })
        .then((data) => {
          window.open(data.longUrl, "_blank");
        });
    }
  });
function fetchAndDisplayStats() {
  fetch("http://localhost:8000/link-stats")
    .then((response) => {
      if (!response.ok) {
        throw new Error("Problème lors de la récupération des statistiques");
      }
      return response.json();
    })
    .then((data) => {
      const totalUrlsElement = document.getElementById("totalUrls");
      const urlListElement = document.getElementById("urlList");

      totalUrlsElement.textContent = `Nombre d'URL raccourci : ${data.length}`;
      urlListElement.innerHTML = "";

      data.forEach((url) => {
        const listItem = document.createElement("li");
        listItem.classList.add(
          "list-group-item",
          "d-flex",
          "justify-content-between",
          "align-items-center"
        );

        const textSpan = document.createElement("span");
        textSpan.textContent = `${url.ShortUrl} - Nombre de visites: ${url.VisitCount}`;

        const copyButton = document.createElement("button");
        copyButton.classList.add("btn", "btn-primary");
        copyButton.textContent = "Copier";
        copyButton.onclick = function () {
          navigator.clipboard.writeText(url.ShortUrl).then(
            () => {
              alert("URL copiée dans le presse-papiers !");
            },
            () => {
              alert("Erreur lors de la copie");
            }
          );
        };

        listItem.appendChild(textSpan);
        listItem.appendChild(copyButton);
        urlListElement.appendChild(listItem);
      });
    })
    .catch((error) => {
      console.error("Erreur lors de la récupération des statistiques:", error);
    });
}

fetchAndDisplayStats();

function redirectToLogin() {
  window.location.href = "login.html";
}
