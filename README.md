# gorille

Ce projet vise à permettre à un client d'envoyer un fichier audio à un serveur qui va distortionner ce fichier puis le lui renvoyer.
Il y a deux fichiers, un contenant le code du serveur, et un autre contenant le code du client ainsi que les fichiers audio en format .wav a modifier.

## Utilisation du projet:

- lancer le serveur
- lancer le client
- entrer le nom du fichier audio à modifier
- et voilà ! le fichier de sortie a le même nom que celui d'entrée avec _out en plus

## Parallélisation

L'algorithme de distortion est parallélisé, et les fonctions les plus lentes de la librairie utilisée ont été parallélisées aussi, le tout pour un gain moyen de 83% de temps d'exécution.a

Fonctions de la lib parallélisées :

- deleteJunk dans reader.go
- parseRawData dans reader.go
- samplesToRawData dans writer.go

## Regard critique

Nous n'avons utilisé aucun channel. Chaque routine écrit dans un tableau à des indexes que seule elle accède (on a vérifié que ça ne crée pas de race condition avec `go run -race <prog>`).

Nous créons ~16 routines à chaque parallélisation (nombre de coeur logique de nos PCs), nous ne nous sommes pas interéssés aux pools (ce serait une grande amélioration puisque pour l'instant, ce n'est optimisé que pour un client à la fois...)
