# Quiz de Pâques

Bienvenue sur le Quiz de Pâques, une application interactive conçue pour évaluer vos connaissances sur le thème de Pâques de manière ludique et éducative. Développé en Go pour le traitement côté serveur et en HTML/CSS pour une interface utilisateur élégante, ce quiz garantit une expérience fluide et captivante.

## Caractéristiques

- **Niveaux de Difficulté**: L'application propose trois niveaux de difficulté (Facile, Moyen, Difficile), chacun offrant 5 questions uniques par session pour enrichir l'expérience utilisateur à chaque nouvelle tentative.
- **Retour Instantané**: Obtenez une réponse immédiate après chaque question pour faciliter un apprentissage interactif et efficace.
- **Bilan des Résultats**: À la conclusion du quiz, consultez un récapitulatif de vos réponses pour évaluer votre performance et identifier les domaines à améliorer.

## Technologies Employées

- **Backend** : Go (Golang)
- **Frontend** : HTML, CSS

## Installation et utilisation avec Makefile

Pour utiliser cette application, vous devez d'abord configurer votre environnement :

1. **Installer Go** :
   Avant de cloner le dépôt et de construire le projet, assurez-vous que Go est installé sur votre ordinateur. Si ce n'est pas déjà fait, téléchargez et installez Go à partir de [la documentation officielle de Go](https://golang.org/doc/install). Ceci est nécessaire pour compiler et exécuter l'application.

2. **Clonez le dépôt** :
   Utilisez la commande suivante pour cloner le dépôt et naviguer dans le répertoire du projet :

   ```bash
   git clone https://github.com/EnzoTurpin/EasterQuiz
   cd chemin_vers_le_projet
   ```

3. **Construisez le projet** :
   À la racine du projet, exécutez la commande suivante pour construire l'application à l'aide du Makefile :

   ```bash
   make build
   ```

   Cela compilera les sources et générera un exécutable nommé `QuizDePaques`.

4. **Exécutez l'application** :
   Pour démarrer l'application, exécutez :

   ```bash
   make run
   ```

   Puis, ouvrez votre navigateur et accédez à `http://localhost:8080` pour commencer à utiliser l'application.

5. **Nettoyer le projet** :
   Pour supprimer l'exécutable et nettoyer les fichiers générés lors de la construction, exécutez :

   ```bash
   make clean
   ```

## Contribuer

Les contributions sont vivement encouragées ! Si vous souhaitez apporter des améliorations au Quiz de Pâques, vous pouvez forker le dépôt et soumettre une pull request.

## Licence

Ce projet est open source. Vous pouvez le redistribuer et/ou le modifier selon les conditions de votre choix.

## Contact

Si vous avez des questions, veuillez me contacter à [enzoturpin3531@gmail.com].
