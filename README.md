# URLShortener

URL Shortener permet de raccourcir une URL longue en une URL courte unique.

## Fonctionnalités

- Inscription et connexion pour les utilisateurs.
- Raccourcir une URL.

## Commencer

git clone
cd votre-projet
cd serveur

### Prérequis

Assurez-vous d'avoir installé Go et MongoDB sur votre système.

```
[Installation de Go](https://golang.org/doc/install)
[Installation de MongoDB](https://docs.mongodb.com/manual/installation/)
```

### Installation

- go run .

### Utilisation

Entrez une URL longue dans le formulaire "URL à raccourcir", puis appuyez sur "Raccourcir l'URL". Copiez l'URL dans le pop-up.

## Construit Avec

- [MongoDB](https://www.mongodb.com/) - La base de données utilisée
- [Go](https://golang.org/) - Le langage de programmation utilisé

## Auteurs

- **Arnaud Gibelli**
- **Hugo Cleret**
- **Alex Corceiro**

### Tester de l'application

- go test -v

### Packages

main:

- "context" : définir des contextes
- "log" : Enregistrer les messages de logs
- "math/rand" : Génération de données aléatoires
- "net/http" : Gérer des requêtes HTTP
- "time" : fournir le temps actuel
- "go.mongodb.org/mongo-driver/mongo" : Utilisé pour les interactions avec la BDD

handler:

- "encoding/json": Encoder et décoder des données JSON
- "errors": Utilisé pour créer des erreurs personnalisées
- "net/url": Manipuler des URLs
- "strings": Offre des fonctions pour manipuler des chaînes de caractères

- "github.com/dgrijalva/jwt-go": Travailler avec JSON Web Tokens
- "go.mongodb.org/mongo-driver/bson": Travailler avec des données BSON
- "golang.org/x/crypto/bcrypt": Utilisé pour le hachage de mots de passe

router:

- "github.com/gorilla/mux": Permet de créer des routeurs

App_test:

- "bytes": Fournit des fonctions pour manipuler des slices d'octets
- "testing": Utilisé pour écrire des cas de test unitaires
