# gorille



Ce projet vise à permettre à un client d'envoyer un fichier audio à un serveur qui va distortionner ce fichier puis le lui renvoyer.
Il y a deux fichiers, un contenant le code du serveur, et un autre contenant le code du client ainsi que les fichiers audio en format .wav a modifier.

Utilisation du projet:
-lancer le serveur
-lancer le client
-entrer le nom du fichier audio à modifier
-et voilà ! le fichier de sortie a le même nom que celui d'entrée avec _out en plus

La modification du fichier est parallélisée ainsi qu'une partie du parsing du fichier audio (la fonction parseRawData).
