# Qu'est-ce ?

Ce script fonctionne en duo avec le HTTP Files Server et sert à :

- récupérer la liste des fichiers disponibles au téléchargement sur le serveur
- télécharger une copie de chaque fichier en local
- supprimer la copie distante après le téléchargement

## Comment utiliser le script

Il est possible d'utiliser `go install` pour générer un executable. Si `go` est installé sur la machine, il est également possible de lancer le script avec la commande `go run main.go`.

Le script attend les arguments suivants (dans l'ordre) :

- répertoire de destination des fichiers téléchargés
- le domaine + le port du script serveur
- le token de sécurité

## Exemple d'utilisation

`go run main.go ./files localhost:8888 azerty`
