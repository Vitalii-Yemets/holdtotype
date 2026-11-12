package main

func init() {
	msgs["de"] = map[string]string{
		"menu.sub.start": "Untertitel starten",
		"menu.sub.stop": "Untertitel beenden",
		"sub.title": "Untertitel",
		"sub.title.tr": "Übersetzung",
		"sub.follow": "zum Neuesten ↓",
		"sub.loading": "Modell wird geladen…",
		"sub.listening": "Höre auf den Ton des Computers…",
		"sub.restarting": "Die Erkennung startet neu…",
		"sub.err.nomodel": "Untertitel: für ihre Sprache ist kein Modell installiert",
		"sub.err.engine": "Untertitel: die Erkennungs-Engine startete nicht",
		"sub.err.device": "Untertitel: kein Wiedergabegerät zum Mithören",
		"sub.err.lost": "Untertitel beendet: die Erkennungs-Engine fiel immer wieder aus",
		"app.name": "{app}", "already.running": "Die Anwendung läuft bereits (Tray-Symbol).",
		"err.title": "{app} — Fehler", "cfg.err.title": "{app} — Konfigurationsfehler",
		"err.details": "\n\nDetails in {log}", "err.hook": "Tastatur-Hook fehlgeschlagen: %s",
		"err.mic": "Mikrofon: %s", "err.hotkey.cfg": "Tastenkürzel in config.json: %s",
		"err.model.notfound": "Modelldatei nicht gefunden: %s\nProjekt neu bauen (.\\build.ps1) oder \"model\" in config.json korrigieren",
		"err.server.repeat":  "whisper-server stürzt wiederholt ab — siehe {log}",
		"err.server.dead":    "Erkennungsserver %s antwortet nicht (server_url/server_autostart in config.json)",
		"err.server.timeout": "whisper-server antwortete nicht innerhalb von %s",
		"err.server.start":   "whisper-server beim Start beendet (siehe Log)",
		"err.srv.noaddr":    "der externe Erkennungsserver ist nicht eingerichtet — Adresse in den Einstellungen setzen",
		"err.webview":        "Das Einstellungsfenster benötigt die Microsoft WebView2 Runtime (in Windows 11 enthalten).\nInstallieren Sie sie oder bearbeiten Sie config.json manuell.",
		"status.loading":     "Modell wird geladen…", "status.ready": "Bereit — %s halten und sprechen",
		"status.recording": "Aufnahme läuft…", "status.transcribing": "Erkenne…", "status.disabled": "Deaktiviert",
		"status.server.restart": "Erkennungsserver abgestürzt, Neustart…", "status.cfg.err": "Fehler in config.json (siehe Log)",
		"menu.settings": "Einstellungen…", "menu.enable": "Aktivieren", "menu.disable": "Deaktivieren",
		"menu.open.config": "config.json öffnen", "menu.open.log": "Log öffnen",
		"menu.about": "Über", "menu.quit": "Beenden",
		"err.port.busy": "Port %d ist belegt: ein anderes Programm hört darauf. Ändern Sie den Port in den Einstellungen.",
		"ov.speak": "Sprechen…", "ov.transcribing": "Erkenne", "ov.asking": "Warte auf Ihre Antwort", "ov.inserted": "Eingefügt: %d Zeichen", "ov.left": "noch %d s", "ov.esc": "1…9 · Enter · Esc bricht ab", "err.net.dns": "Keine Verbindung zu %s — Internet prüfen", "err.net.timeout": "Der Server hat nicht rechtzeitig geantwortet — erneut versuchen", "err.net.down": "Verbindung nicht möglich — Internet prüfen", "err.net.cert": "Die Verbindung ist nicht sicher — Datum und Virenschutz prüfen", "err.answer": "Der Server antwortete unverständlich — später versuchen", "err.file.missing": "Datei nicht gefunden", "err.file.denied": "Kein Zugriff auf die Datei — das Programm schließen, das sie hält", "err.disk.full": "Kein Speicherplatz auf der Festplatte", "err.cancelled": "Abgebrochen", "err.generic": "Hat nicht geklappt — Einzelheiten im Protokoll", "err.server.launch": "%s konnte nicht gestartet werden — Serverpfad unter System prüfen",
		"ov.err.recognize": "Erkennungsfehler (siehe Log)", "ov.err.paste": "Nicht eingefügt — der Text steht im letzten Ergebnis",
		"ov.moved":  "Das Fenster hat gewechselt — der Text liegt in der Zwischenablage",
		"copy.ok":   "Kopiert",
		"copy.none": "Nichts zu kopieren",
		"copy.fail": "Kopieren fehlgeschlagen: %s",
		"mic.busy":  "Ein Diktat läuft, jetzt geht das nicht", "mic.check.ok": "Klingt gut: Spitze %.0f dB, Sprache in %.0f%% der Aufnahme",
		"mic.check.quiet": "Zu leise: Spitze %.0f dB — Mikrofonpegel in Windows anheben oder näher sitzen", "mic.check.clipped": "Übersteuert: %.1f%% der Abtastwerte abgeschnitten — Mikrofonpegel senken", "mic.check.silent": "Keine Sprache gehört — prüfen Sie, ob das richtige Mikrofon gewählt und nicht stumm ist",
		"ov.quiet": "Zu leise, es war fast nichts zu hören", "ov.clipped": "Übersteuert — der Ton wurde abgeschnitten",
		"ov.cmd.cancelled": "Per Sprache abgebrochen",
		"ov.undo.done": "Gelöscht: %d Zeichen", "ov.undo.none": "Nichts zu löschen", "ov.undo.typed": "Nicht gelöscht: seitdem wurde getippt oder geklickt",
		"ov.undo.window": "Nicht gelöscht: ein anderes Fenster", "ov.undo.sent": "Nicht gelöscht: schon mit Enter gesendet", "ov.undo.old": "Nicht gelöscht: die Einfügung ist zu alt",
		"ov.undo.breaks": "Nicht gelöscht: der Text hatte Zeilenumbrüche", "ov.undo.long": "Nicht gelöscht: der Text ist zu lang", "ov.undo.moved": "Löschen abgebrochen: das Fenster hat gewechselt",
		"ov.shape.next": "Nächste Phrase: %s", "cmd.heading": "Überschrift", "cmd.list": "Liste", "cmd.capital": "großgeschrieben", "cmd.quotes": "in Anführungszeichen", "cmd.nostop": "ohne Punkt",
		"ov.silence":       "Stille — nichts erkannt", "ov.server.loading": "Server lädt noch",
		"ov.stop.silence":  "Bei Stille beendet",
		"ov.tooshort":  "Zu kurz — Tasten länger halten",
		"ov.cancelled": "Abgebrochen", "ov.editing": "Bearbeite: %s", "ov.translating": "Übersetze",
		"ov.llm.needed": "Diese Sprache benötigt das LLM-Modul", "td.title": "Übersetzen nach:", "td.plain": "In die Ausgangssprache übersetzen",
		"cap.title": "{app} — Tastenkürzel", "cap.prompt": "Neue Tastenkombination drücken\n\njetzt: %s   ·   Esc — Abbrechen",
		"cap.selected": "Gewählt: %s", "cap.cancelled": "Abgebrochen", "hk.taken": "%s ist von Windows belegt: %s. Das Diktat startet womöglich nie", "hk.lock": "Computer sperren", "hk.desktop": "Desktop anzeigen", "hk.explorer": "Explorer", "hk.run": "Ausführen-Fenster", "hk.settings": "Windows-Einstellungen", "hk.search": "Suche", "hk.center": "Infocenter", "hk.menu": "Power-User-Menü", "hk.clipboard": "Zwischenablageverlauf", "hk.gamebar": "Gamebar", "hk.voice": "Windows-Spracheingabe", "hk.project": "Projizieren", "hk.tasks": "Taskansicht", "hk.layout": "Tastaturlayout wechseln", "hk.newdesktop": "neuer Desktop", "hk.closedesktop": "Desktop schließen", "hk.snip": "Screenshot-Werkzeug", "hk.switch": "Fenster wechseln", "hk.close": "Fenster schließen", "hk.cycle": "Fenster durchblättern", "hk.start": "Startmenü", "hk.taskmgr": "Task-Manager", "hk.secure": "Sicherheitsbildschirm",
		"model.switching": "Modellwechsel — der Erkenner startet neu…", "model.del.active": "Aktives Modell kann nicht gelöscht werden",
		"model.del.ok":   "Modell gelöscht",
		"about.text":     "{app} %s\n\nSprache → Text an der Cursorposition.\nCursor in ein Eingabefeld setzen, %s halten, sprechen, loslassen — der Text wird eingefügt.\n\nErkennung: whisper.cpp, vollständig lokal und offline.\nModell: %s (Sprache: %s)\n\nEinstellungen: Klick auf das Tray-Symbol oder config.json.\nLogs: {log} (max. ~2 MB).",
		"status.nomodel": "Kein Erkennungsmodell geladen — wählen Sie eines in den Einstellungen",
		"tf.title": "%s übersetzt nicht. Trotzdem übersetzen?",
		"tf.yes": "Ja, über %s",
		"tf.no": "Wie erkannt",
		"hk.lone": "%s ist eine einzelne Taste: solange sie das Kürzel ist, lässt sie sich nicht mehr normal tippen",
		"state.loaded.none": "nichts geladen",
		"state.week": "%d Diktate · %d Zeichen",
		"snd.ok": "Pegel in Ordnung",
		"snd.quiet": "leise — sprechen Sie näher am Mikrofon",
		"snd.clipped": "zu laut, der Ton übersteuert",
		"snd.silent": "Stille in der Aufnahme",
		"status.parked": "Erkennung entladen — das Tastenkürzel weckt sie wieder",
		"status.nomodel.lang": "Für %s ist das Modell %s nicht installiert — öffnen Sie „Sprachen & Modelle“",
		"menu.lastcopy":  "Letztes Ergebnis kopieren",
		"ov.copied":      "In die Zwischenablage kopiert", "ov.kept": "Abgebrochen — Text bleibt im letzten Ergebnis",
		"ov.llm.skipped": "Eingefügt ohne das Profil „%s“",
		"fd.title":       "Fokus gewechselt — einfügen?", "fd.here": "Hier einfügen", "fd.copy": "Kopieren",
		"fd.keep":            "Behalten",
		"ov.err.mic":         "Mikrofon nicht verfügbar — prüfen Sie das Gerät in den Einstellungen",
		"ov.notranslate":     "Das aktive Modell kann nicht übersetzen — wie erkannt eingefügt",
		"ov.engine.fallback": "Die andere Engine startete nicht — es bleibt bei der aktuellen",
		"route.speech":       "Sprache auf %s", "route.other": "Andere Sprachen", "route.translate": "Übersetzung",
		"route.lang.auto":    "beliebige Sprache",
		"route.why.language": "hier genauer, mit Satzzeichen", "route.why.otherlang": "99 Sprachen",
		"route.why.translate": "nur Whisper übersetzt", "route.why.notinstalled": "das russische Modell ist nicht installiert",
		"route.why.unknownlang": "keine Sprache gesetzt — nur Whisper erkennt sie", "route.why.forced": "in config.json erzwungen",
		"status.line": "Bereit · %s · %.1f GB frei", "state.ram.free": "%d MB frei",
		"ago.now": "gerade eben", "ago.min": "vor %d Min.", "ago.hour": "vor %d Std.",
		"chars": "%d Zeichen", "inserted.into": "eingefügt in %s",
		"punct.prompt":         "Setze Satzzeichen und Großschreibung. Ändere die Wörter nicht, übersetze nicht, füge nichts hinzu. Gib nur den korrigierten Text zurück.",
		"err.sherpa.notfound":  "sherpa-Erkennung nicht gefunden: %s",
		"err.sherpa.start":     "sherpa-server hat sich beim Start beendet (siehe Protokoll)",
		"err.sherpa.translate": "dieses Modell kann nicht übersetzen",
		"err.sherpa.model":     "Modelldatei nicht gefunden: %s — laden Sie sie in den Einstellungen oder korrigieren Sie sherpa_model in config.json",
		"err.hotkey.dup":       "Das Kürzel „%s“ ist doppelt vergeben — Kürzel müssen eindeutig sein",
		"cfg.err.recovered":    "config.json ist beschädigt (%s).\nDie Datei wurde als %s gesichert, die Einstellungen wurden zurückgesetzt.",
		"err.disk.space":       "wenig Speicherplatz: %d MB frei, ~%d MB nötig",
		"err.model.dir":       "der Ordner %s konnte nicht angelegt werden — prüfen Sie den Zugriff",
		"err.save":             "Einstellungen nicht gespeichert: %s — die alten bleiben",
		"err.port":             "Port %d passt nicht: eine Nummer zwischen 1024 und 65535 wählen",
		"err.nolangs":          "lassen Sie mindestens eine Sprache für die Übersetzungsfrage stehen",
		"ov.mic.lost":          "Mikrofon weg — Aufnahme abgebrochen",
		"err.hash":             "die heruntergeladene Datei ist beschädigt — bitte erneut versuchen",
		"models.check.ok":      "Geprüfte Modelle: %d — alle Dateien in Ordnung",
		"models.check.none":    "Nichts zu prüfen — kein installiertes Modell hat einen Referenz-Hash",
		"models.check.bad":     "Beschädigte Dateien: %s — Modell erneut herunterladen",
		"hist.insert.gone":     "Eintrag nicht gefunden",
		"ov.aim": "Klicken Sie das Feld zum Einfügen an · Esc bricht ab",
		"hist.aim.armed": "klicken Sie das Feld zum Einfügen an",
		"hist.aim.busy": "warte bereits auf einen Klick",
		"hist.aim.off": "Einfügen abgebrochen",
		"hist.insert.nowin":    "kein Ziel zum Einfügen — der Text liegt in der Zwischenablage",
		"hist.insert.ok":       "in „%s“ eingefügt",
		"lists.bad":            "diese Datei passt nicht",
		"lists.saved":          "gespeichert in %s",
		"lists.added":          "hinzugefügt: %d, übersprungen: %d",
		"lists.save.title":     "Wohin die Listen speichern",
		"lists.open.title":     "Welche Datei laden",
		"un.title":             "{app} — Deinstallation", "un.confirm": "{app} von diesem Rechner entfernen?",
		"un.data": "Auch Einstellungen, geladene Modelle und Untertitel-Mitschriften löschen?", "un.done": "{app} wurde entfernt.",
		"srv.restarting": "Erkennung wird mit den neuen Einstellungen neu gestartet…",
	}
	msgs["fr"] = map[string]string{
		"menu.sub.start": "Démarrer les sous-titres",
		"menu.sub.stop": "Arrêter les sous-titres",
		"sub.title": "Sous-titres",
		"sub.title.tr": "Traduction",
		"sub.follow": "au plus récent ↓",
		"sub.loading": "Chargement du modèle…",
		"sub.listening": "À l'écoute du son de l'ordinateur…",
		"sub.restarting": "La reconnaissance redémarre…",
		"sub.err.nomodel": "Sous-titres : aucun modèle n'est installé pour leur langue",
		"sub.err.engine": "Sous-titres : le moteur de reconnaissance n'a pas démarré",
		"sub.err.device": "Sous-titres : aucun périphérique de lecture à écouter",
		"sub.err.lost": "Sous-titres arrêtés : le moteur de reconnaissance a planté à répétition",
		"app.name": "{app}", "already.running": "L'application est déjà lancée (icône de la barre).",
		"err.title": "{app} — erreur", "cfg.err.title": "{app} — erreur de configuration",
		"err.details": "\n\nDétails dans {log}", "err.hook": "Échec du hook clavier : %s",
		"err.mic": "Microphone : %s", "err.hotkey.cfg": "Raccourci dans config.json : %s",
		"err.model.notfound": "Modèle introuvable : %s\nRecompilez le projet (.\\build.ps1) ou corrigez \"model\" dans config.json",
		"err.server.repeat":  "whisper-server plante en boucle — voir {log}",
		"err.server.dead":    "Le serveur %s ne répond pas (server_url/server_autostart dans config.json)",
		"err.server.timeout": "whisper-server n'a pas répondu en %s",
		"err.server.start":   "whisper-server s'est arrêté au démarrage (voir log)",
		"err.srv.noaddr":    "le serveur de reconnaissance distant n’est pas configuré — renseignez son adresse dans les réglages",
		"err.webview":        "La fenêtre des réglages nécessite Microsoft WebView2 Runtime (inclus dans Windows 11).\nInstallez-le ou modifiez config.json manuellement.",
		"status.loading":     "Chargement du modèle…", "status.ready": "Prêt — maintenez %s et parlez",
		"status.recording": "Enregistrement…", "status.transcribing": "Reconnaissance…", "status.disabled": "Désactivé",
		"status.server.restart": "Serveur planté, redémarrage…", "status.cfg.err": "Erreur dans config.json (voir log)",
		"menu.settings": "Réglages…", "menu.enable": "Activer", "menu.disable": "Désactiver",
		"menu.open.config": "Ouvrir config.json", "menu.open.log": "Ouvrir le log",
		"menu.about": "À propos", "menu.quit": "Quitter",
		"err.port.busy": "Le port %d est occupé par un autre programme. Changez le port dans les réglages.",
		"ov.speak": "Parlez…", "ov.transcribing": "Reconnaissance", "ov.asking": "En attente de votre réponse", "ov.inserted": "Inséré : %d caractères", "ov.left": "%d s restantes", "ov.esc": "1…9 · Entrée · Échap annule", "err.net.dns": "Pas de connexion à %s — vérifiez internet", "err.net.timeout": "Le serveur n'a pas répondu à temps — réessayez", "err.net.down": "Connexion impossible — vérifiez internet", "err.net.cert": "La connexion n'est pas sécurisée — vérifiez la date et l'antivirus", "err.answer": "Le serveur a répondu de façon incompréhensible — réessayez plus tard", "err.file.missing": "Fichier introuvable", "err.file.denied": "Pas d'accès au fichier — fermez le programme qui le retient", "err.disk.full": "Le disque est plein", "err.cancelled": "Annulé", "err.generic": "Cela n'a pas fonctionné — détails dans le journal", "err.server.launch": "Impossible de lancer %s — vérifiez le chemin du serveur dans Système",
		"ov.err.recognize": "Erreur de reconnaissance (voir log)", "ov.err.paste": "Non collé — le texte est dans le dernier résultat",
		"ov.moved":  "La fenêtre a changé — le texte est dans le presse-papiers",
		"copy.ok":   "Copié",
		"copy.none": "Rien à copier",
		"copy.fail": "Copie impossible : %s",
		"mic.busy":  "Une dictée est en cours, impossible de vérifier", "mic.check.ok": "Bon signal : crête %.0f dB, parole sur %.0f%% de l'enregistrement",
		"mic.check.quiet": "Trop faible : crête %.0f dB — montez le niveau du micro dans Windows ou rapprochez-vous", "mic.check.clipped": "Saturation : %.1f%% des échantillons écrêtés — baissez le niveau du micro", "mic.check.silent": "Aucune parole entendue — vérifiez que le bon micro est choisi et qu'il n'est pas coupé",
		"ov.quiet": "Trop faible, presque rien n'a été entendu", "ov.clipped": "Saturation — le son a été écrêté",
		"ov.cmd.cancelled": "Annulé à la voix",
		"ov.undo.done": "Effacé : %d caractères", "ov.undo.none": "Rien à effacer", "ov.undo.typed": "Non effacé : vous avez tapé ou cliqué depuis",
		"ov.undo.window": "Non effacé : une autre fenêtre", "ov.undo.sent": "Non effacé : déjà envoyé avec Entrée", "ov.undo.old": "Non effacé : inséré il y a trop longtemps",
		"ov.undo.breaks": "Non effacé : le texte contenait des sauts de ligne", "ov.undo.long": "Non effacé : le texte est trop long", "ov.undo.moved": "Effacement interrompu : la fenêtre a changé",
		"ov.shape.next": "Phrase suivante : %s", "cmd.heading": "titre", "cmd.list": "liste", "cmd.capital": "majuscule", "cmd.quotes": "entre guillemets", "cmd.nostop": "sans point",
		"ov.silence":       "Silence — rien reconnu", "ov.server.loading": "Le serveur charge encore",
		"ov.stop.silence":  "Arrêté sur un silence",
		"ov.tooshort":  "Trop court — maintenez les touches plus longtemps",
		"ov.cancelled": "Annulé", "ov.editing": "Édition : %s", "ov.translating": "Traduction",
		"ov.llm.needed": "Cette langue nécessite le module LLM", "td.title": "Traduire vers :", "td.plain": "Traduire dans la langue source",
		"cap.title": "{app} — raccourci", "cap.prompt": "Appuyez une nouvelle combinaison\n\nactuel : %s   ·   Échap — annuler",
		"cap.selected": "Choisi : %s", "cap.cancelled": "Annulé", "hk.taken": "%s est pris par Windows : %s. La dictée risque de ne jamais démarrer", "hk.lock": "verrouiller l'ordinateur", "hk.desktop": "afficher le bureau", "hk.explorer": "explorateur de fichiers", "hk.run": "fenêtre Exécuter", "hk.settings": "paramètres Windows", "hk.search": "recherche", "hk.center": "centre de notifications", "hk.menu": "menu avancé", "hk.clipboard": "historique du presse-papiers", "hk.gamebar": "barre de jeu", "hk.voice": "saisie vocale Windows", "hk.project": "projeter sur un écran", "hk.tasks": "affichage des tâches", "hk.layout": "changer de disposition", "hk.newdesktop": "nouveau bureau", "hk.closedesktop": "fermer le bureau", "hk.snip": "outil capture d'écran", "hk.switch": "changer de fenêtre", "hk.close": "fermer la fenêtre", "hk.cycle": "parcourir les fenêtres", "hk.start": "menu Démarrer", "hk.taskmgr": "gestionnaire des tâches", "hk.secure": "écran de sécurité",
		"model.switching": "Changement de modèle — redémarrage…", "model.del.active": "Impossible de supprimer le modèle actif",
		"model.del.ok":   "Modèle supprimé",
		"about.text":     "{app} %s\n\nVoix → texte à la position du curseur.\nPlacez le curseur, maintenez %s, parlez, relâchez — le texte s'insère.\n\nReconnaissance : whisper.cpp, entièrement locale et hors ligne.\nModèle : %s (langue : %s)\n\nRéglages : clic sur l'icône ou config.json.\nLogs : {log} (max ~2 Mo).",
		"status.nomodel": "Aucun modèle de reconnaissance téléchargé — choisissez-en un dans les réglages",
		"tf.title": "%s ne traduit pas. Traduire quand même ?",
		"tf.yes": "Oui, via %s",
		"tf.no": "Tel quel",
		"hk.lone": "%s est une touche seule : tant qu'elle sert de raccourci, on ne peut plus la taper normalement",
		"state.loaded.none": "rien de chargé",
		"state.week": "%d dictées · %d caractères",
		"snd.ok": "niveau correct",
		"snd.quiet": "faible — parlez plus près du micro",
		"snd.clipped": "trop fort, le son a saturé",
		"snd.silent": "silence dans l’enregistrement",
		"status.parked": "Moteur déchargé — le raccourci le réveillera",
		"status.nomodel.lang": "Pour %s, le modèle %s n’est pas installé — ouvrez « Langues et modèles »",
		"menu.lastcopy":  "Copier le dernier résultat",
		"ov.copied":      "Copié dans le presse-papiers", "ov.kept": "Annulé — le texte reste dans le dernier résultat",
		"ov.llm.skipped": "Inséré sans le profil « %s »",
		"fd.title":       "Le focus a changé — insérer ?", "fd.here": "Insérer ici", "fd.copy": "Copier",
		"fd.keep":            "Garder",
		"ov.err.mic":         "Microphone indisponible — vérifiez le périphérique dans les réglages",
		"ov.notranslate":     "Le modèle actif ne sait pas traduire — inséré tel que reconnu",
		"ov.engine.fallback": "L'autre moteur n'a pas démarré — on garde l'actuel",
		"route.speech":       "Parole en %s", "route.other": "Autres langues", "route.translate": "Traduction",
		"route.lang.auto":    "n'importe quelle langue",
		"route.why.language": "plus précis ici, avec la ponctuation", "route.why.otherlang": "99 langues",
		"route.why.translate": "seul Whisper traduit", "route.why.notinstalled": "le modèle russe n'est pas installé",
		"route.why.unknownlang": "aucune langue définie — seul Whisper la détecte", "route.why.forced": "forcé dans config.json",
		"status.line": "Prêt · %s · %.1f Go libres", "state.ram.free": "%d Mo libres",
		"ago.now": "à l'instant", "ago.min": "il y a %d min", "ago.hour": "il y a %d h",
		"chars": "%d caractères", "inserted.into": "inséré dans %s",
		"punct.prompt":         "Ajoute la ponctuation et les majuscules. Ne change pas les mots, ne traduis pas, n'ajoute rien. Renvoie uniquement le texte corrigé.",
		"err.sherpa.notfound":  "reconnaissance sherpa introuvable : %s",
		"err.sherpa.start":     "sherpa-server s'est arrêté au démarrage (voir le journal)",
		"err.sherpa.translate": "ce modèle ne sait pas traduire",
		"err.sherpa.model":     "Fichier de modèle introuvable : %s — téléchargez-le dans les réglages ou corrigez sherpa_model dans config.json",
		"err.hotkey.dup":       "Le raccourci « %s » est attribué deux fois — les raccourcis doivent être uniques",
		"cfg.err.recovered":    "config.json est corrompu (%s).\nLe fichier a été enregistré sous %s et les réglages ont été réinitialisés.",
		"err.disk.space":       "espace disque faible : %d Mo libres, ~%d Mo nécessaires",
		"err.model.dir":       "impossible de créer le dossier %s — vérifiez les droits d'accès",
		"err.save":             "réglages non enregistrés : %s — les anciens sont conservés",
		"err.port":             "le port %d ne convient pas : choisissez un nombre entre 1024 et 65535",
		"err.nolangs":          "laissez au moins une langue pour la question de traduction",
		"ov.mic.lost":          "Micro débranché — enregistrement interrompu",
		"err.hash":             "le fichier téléchargé est endommagé — réessayez",
		"models.check.ok":      "Modèles vérifiés : %d — tous les fichiers sont intacts",
		"models.check.none":    "Rien à vérifier — aucun modèle installé n'a d'empreinte de référence",
		"models.check.bad":     "Fichiers endommagés : %s — téléchargez le modèle à nouveau",
		"hist.insert.gone":     "entrée introuvable",
		"ov.aim": "Cliquez le champ où coller · Échap annule",
		"hist.aim.armed": "cliquez le champ où coller",
		"hist.aim.busy": "j'attends déjà un clic",
		"hist.aim.off": "collage annulé",
		"hist.insert.nowin":    "rien où coller — le texte est dans le presse-papiers",
		"hist.insert.ok":       "collé dans « %s »",
		"lists.bad":            "ce fichier ne convient pas",
		"lists.saved":          "enregistré dans %s",
		"lists.added":          "ajoutés : %d, ignorés : %d",
		"lists.save.title":     "Où enregistrer les listes",
		"lists.open.title":     "Quel fichier charger",
		"un.title":             "{app} — Désinstallation", "un.confirm": "Supprimer {app} de cet ordinateur ?",
		"un.data": "Supprimer aussi les réglages, les modèles téléchargés et les transcriptions des sous-titres ?", "un.done": "{app} a été supprimé.",
		"srv.restarting": "Redémarrage de la reconnaissance avec les nouveaux réglages…",
	}
	msgs["es"] = map[string]string{
		"menu.sub.start": "Iniciar subtítulos",
		"menu.sub.stop": "Detener subtítulos",
		"sub.title": "Subtítulos",
		"sub.title.tr": "Traducción",
		"sub.follow": "a lo último ↓",
		"sub.loading": "Cargando el modelo…",
		"sub.listening": "Escuchando el sonido del ordenador…",
		"sub.restarting": "El reconocimiento se está reiniciando…",
		"sub.err.nomodel": "Subtítulos: no hay ningún modelo instalado para su idioma",
		"sub.err.engine": "Subtítulos: el motor de reconocimiento no arrancó",
		"sub.err.device": "Subtítulos: no hay ningún dispositivo de reproducción que escuchar",
		"sub.err.lost": "Subtítulos detenidos: el motor de reconocimiento fallaba una y otra vez",
		"app.name": "{app}", "already.running": "La aplicación ya está en ejecución (icono de bandeja).",
		"err.title": "{app} — error", "cfg.err.title": "{app} — error de configuración",
		"err.details": "\n\nDetalles en {log}", "err.hook": "Fallo del hook de teclado: %s",
		"err.mic": "Micrófono: %s", "err.hotkey.cfg": "Atajo en config.json: %s",
		"err.model.notfound": "Modelo no encontrado: %s\nRecompile el proyecto (.\\build.ps1) o corrija \"model\" en config.json",
		"err.server.repeat":  "whisper-server falla repetidamente — vea {log}",
		"err.server.dead":    "El servidor %s no responde (server_url/server_autostart en config.json)",
		"err.server.timeout": "whisper-server no respondió en %s",
		"err.server.start":   "whisper-server terminó al iniciar (vea el log)",
		"err.srv.noaddr":    "el servidor de reconocimiento remoto no está configurado: indique su dirección en los ajustes",
		"err.webview":        "La ventana de ajustes requiere Microsoft WebView2 Runtime (incluido en Windows 11).\nInstálelo o edite config.json manualmente.",
		"status.loading":     "Cargando modelo…", "status.ready": "Listo — mantenga %s y hable",
		"status.recording": "Grabando…", "status.transcribing": "Reconociendo…", "status.disabled": "Desactivado",
		"status.server.restart": "Servidor caído, reiniciando…", "status.cfg.err": "Error en config.json (vea el log)",
		"menu.settings": "Ajustes…", "menu.enable": "Activar", "menu.disable": "Desactivar",
		"menu.open.config": "Abrir config.json", "menu.open.log": "Abrir log",
		"menu.about": "Acerca de", "menu.quit": "Salir",
		"err.port.busy": "El puerto %d está ocupado por otro programa. Cambia el puerto en los ajustes.",
		"ov.speak": "Hable…", "ov.transcribing": "Reconociendo", "ov.asking": "Esperando tu respuesta", "ov.inserted": "Insertado: %d caracteres", "ov.left": "quedan %d s", "ov.esc": "1…9 · Enter · Esc cancela", "err.net.dns": "Sin conexión con %s — revise internet", "err.net.timeout": "El servidor no respondió a tiempo — inténtelo otra vez", "err.net.down": "No se pudo conectar — revise internet", "err.net.cert": "La conexión no es segura — revise la fecha y el antivirus", "err.answer": "El servidor respondió de forma incomprensible — inténtelo más tarde", "err.file.missing": "Archivo no encontrado", "err.file.denied": "Sin acceso al archivo — cierre el programa que lo retiene", "err.disk.full": "El disco está lleno", "err.cancelled": "Cancelado", "err.generic": "No funcionó — los detalles están en el registro", "err.server.launch": "No se pudo iniciar %s — revise la ruta del servidor en Sistema",
		"ov.err.recognize": "Error de reconocimiento (vea el log)", "ov.err.paste": "No se pegó: el texto está en el último resultado",
		"ov.moved":  "La ventana cambió: el texto está en el portapapeles",
		"copy.ok":   "Copiado",
		"copy.none": "Nada que copiar",
		"copy.fail": "No se pudo copiar: %s",
		"mic.busy":  "Hay un dictado en curso, ahora no se puede comprobar", "mic.check.ok": "Se oye bien: pico %.0f dB, voz en el %.0f%% de la grabación",
		"mic.check.quiet": "Demasiado bajo: pico %.0f dB — sube el nivel del micrófono en Windows o acércate", "mic.check.clipped": "Saturación: %.1f%% de muestras recortadas — baja el nivel del micrófono", "mic.check.silent": "No se oye voz — comprueba que el micrófono elegido es el correcto y no está silenciado",
		"ov.quiet": "Demasiado bajo, casi no se oyó nada", "ov.clipped": "Saturación: el sonido se recortó",
		"ov.cmd.cancelled": "Cancelado por voz",
		"ov.undo.done": "Borrado: %d caracteres", "ov.undo.none": "Nada que borrar", "ov.undo.typed": "No se borró: escribiste o hiciste clic después",
		"ov.undo.window": "No se borró: otra ventana", "ov.undo.sent": "No se borró: ya se envió con Enter", "ov.undo.old": "No se borró: se insertó hace demasiado",
		"ov.undo.breaks": "No se borró: el texto tenía saltos de línea", "ov.undo.long": "No se borró: el texto es demasiado largo", "ov.undo.moved": "Borrado interrumpido: la ventana cambió",
		"ov.shape.next": "Siguiente frase: %s", "cmd.heading": "título", "cmd.list": "lista", "cmd.capital": "mayúscula", "cmd.quotes": "entre comillas", "cmd.nostop": "sin punto",
		"ov.silence":       "Silencio — nada reconocido", "ov.server.loading": "El servidor aún carga",
		"ov.stop.silence":  "Detenido por silencio",
		"ov.tooshort":  "Demasiado corto: mantén las teclas más tiempo",
		"ov.cancelled": "Cancelado", "ov.editing": "Editando: %s", "ov.translating": "Traduciendo",
		"ov.llm.needed": "Este idioma requiere el módulo LLM", "td.title": "Traducir a:", "td.plain": "Traducir al idioma de origen",
		"cap.title": "{app} — atajo", "cap.prompt": "Pulse una nueva combinación\n\nactual: %s   ·   Esc — cancelar",
		"cap.selected": "Elegido: %s", "cap.cancelled": "Cancelado", "hk.taken": "%s lo usa Windows: %s. Puede que el dictado no empiece nunca", "hk.lock": "bloquear el equipo", "hk.desktop": "mostrar el escritorio", "hk.explorer": "explorador de archivos", "hk.run": "ventana Ejecutar", "hk.settings": "configuración de Windows", "hk.search": "búsqueda", "hk.center": "centro de notificaciones", "hk.menu": "menú avanzado", "hk.clipboard": "historial del portapapeles", "hk.gamebar": "barra de juego", "hk.voice": "escritura por voz de Windows", "hk.project": "proyectar en pantalla", "hk.tasks": "vista de tareas", "hk.layout": "cambiar la distribución", "hk.newdesktop": "escritorio nuevo", "hk.closedesktop": "cerrar el escritorio", "hk.snip": "herramienta de recortes", "hk.switch": "cambiar de ventana", "hk.close": "cerrar la ventana", "hk.cycle": "recorrer ventanas", "hk.start": "menú Inicio", "hk.taskmgr": "administrador de tareas", "hk.secure": "pantalla de seguridad",
		"model.switching": "Cambiando modelo — reiniciando…", "model.del.active": "No se puede borrar el modelo activo",
		"model.del.ok":   "Modelo borrado",
		"about.text":     "{app} %s\n\nVoz → texto en la posición del cursor.\nColoque el cursor, mantenga %s, hable, suelte — el texto se inserta.\n\nReconocimiento: whisper.cpp, totalmente local y sin conexión.\nModelo: %s (idioma: %s)\n\nAjustes: clic en el icono o config.json.\nLogs: {log} (máx ~2 MB).",
		"status.nomodel": "No hay ningún modelo descargado — elige uno en los ajustes",
		"tf.title": "%s no traduce. ¿Traducir de todos modos?",
		"tf.yes": "Sí, con %s",
		"tf.no": "Tal cual",
		"hk.lone": "%s es una sola tecla: mientras sea el atajo, no se podrá teclear con normalidad",
		"state.loaded.none": "nada cargado",
		"state.week": "%d dictados · %d caracteres",
		"snd.ok": "nivel correcto",
		"snd.quiet": "bajo — hable más cerca del micrófono",
		"snd.clipped": "demasiado alto, el sonido se recortó",
		"snd.silent": "silencio en la grabación",
		"status.parked": "Motor descargado — el atajo lo despertará",
		"status.nomodel.lang": "Para %s no está instalado el modelo %s — abra «Idiomas y modelos»",
		"menu.lastcopy":  "Copiar el último resultado",
		"ov.copied":      "Copiado al portapapeles", "ov.kept": "Cancelado — el texto queda en el último resultado",
		"ov.llm.skipped": "Insertado sin el perfil «%s»",
		"fd.title":       "Cambió el foco, ¿insertar?", "fd.here": "Insertar aquí", "fd.copy": "Copiar",
		"fd.keep":            "Conservar",
		"ov.err.mic":         "Micrófono no disponible — revisa el dispositivo en los ajustes",
		"ov.notranslate":     "El modelo activo no traduce — insertado tal como se reconoció",
		"ov.engine.fallback": "El otro motor no arrancó — se sigue con el actual",
		"route.speech":       "Habla en %s", "route.other": "Otros idiomas", "route.translate": "Traducción",
		"route.lang.auto":    "cualquier idioma",
		"route.why.language": "más preciso aquí, con puntuación", "route.why.otherlang": "99 idiomas",
		"route.why.translate": "solo Whisper traduce", "route.why.notinstalled": "el modelo ruso no está instalado",
		"route.why.unknownlang": "sin idioma definido — solo Whisper lo detecta", "route.why.forced": "forzado en config.json",
		"status.line": "Listo · %s · %.1f GB libres", "state.ram.free": "%d MB libres",
		"ago.now": "ahora mismo", "ago.min": "hace %d min", "ago.hour": "hace %d h",
		"chars": "%d caracteres", "inserted.into": "insertado en %s",
		"punct.prompt":         "Añade puntuación y mayúsculas. No cambies las palabras, no traduzcas, no añadas nada. Devuelve solo el texto corregido.",
		"err.sherpa.notfound":  "no se encuentra el reconocedor sherpa: %s",
		"err.sherpa.start":     "sherpa-server se cerró al arrancar (mira el registro)",
		"err.sherpa.translate": "este modelo no puede traducir",
		"err.sherpa.model":     "No se encuentra el archivo del modelo: %s — descárgalo en los ajustes o corrige sherpa_model en config.json",
		"err.hotkey.dup":       "El atajo «%s» está asignado dos veces — los atajos deben ser únicos",
		"cfg.err.recovered":    "config.json está dañado (%s).\nEl archivo se guardó como %s y los ajustes volvieron a los valores por defecto.",
		"err.disk.space":       "poco espacio en disco: %d MB libres, ~%d MB necesarios",
		"err.model.dir":       "no se pudo crear la carpeta %s — comprueba el acceso",
		"err.save":             "ajustes no guardados: %s — se mantienen los anteriores",
		"err.port":             "el puerto %d no sirve: elige un número entre 1024 y 65535",
		"err.nolangs":          "deja al menos un idioma para la pregunta de traducción",
		"ov.mic.lost":          "Micrófono desconectado: grabación interrumpida",
		"err.hash":             "el archivo descargado está dañado: inténtalo de nuevo",
		"models.check.ok":      "Modelos comprobados: %d, todos los archivos están intactos",
		"models.check.none":    "Nada que comprobar: ningún modelo instalado tiene hash de referencia",
		"models.check.bad":     "Archivos dañados: %s — descarga el modelo de nuevo",
		"hist.insert.gone":     "no se encuentra la entrada",
		"ov.aim": "Haga clic en el campo donde pegar · Esc cancela",
		"hist.aim.armed": "haga clic en el campo donde pegar",
		"hist.aim.busy": "ya estoy esperando un clic",
		"hist.aim.off": "pegado cancelado",
		"hist.insert.nowin":    "no hay dónde pegar: el texto está en el portapapeles",
		"hist.insert.ok":       "pegado en «%s»",
		"lists.bad":            "este archivo no sirve",
		"lists.saved":          "guardado en %s",
		"lists.added":          "añadidos: %d, omitidos: %d",
		"lists.save.title":     "Dónde guardar las listas",
		"lists.open.title":     "Qué archivo cargar",
		"un.title":             "{app} — Desinstalación", "un.confirm": "¿Quitar {app} de este equipo?",
		"un.data": "¿Borrar también los ajustes, los modelos descargados y las transcripciones de subtítulos?", "un.done": "{app} se ha quitado.",
		"srv.restarting": "Reiniciando el reconocedor con los nuevos ajustes…",
	}
	msgs["it"] = map[string]string{
		"menu.sub.start": "Avvia i sottotitoli",
		"menu.sub.stop": "Ferma i sottotitoli",
		"sub.title": "Sottotitoli",
		"sub.title.tr": "Traduzione",
		"sub.follow": "all'ultimo ↓",
		"sub.loading": "Caricamento del modello…",
		"sub.listening": "In ascolto del suono del computer…",
		"sub.restarting": "Il riconoscimento si sta riavviando…",
		"sub.err.nomodel": "Sottotitoli: nessun modello installato per la loro lingua",
		"sub.err.engine": "Sottotitoli: il motore di riconoscimento non è partito",
		"sub.err.device": "Sottotitoli: nessun dispositivo di riproduzione da ascoltare",
		"sub.err.lost": "Sottotitoli fermati: il motore di riconoscimento continuava a cadere",
		"app.name": "{app}", "already.running": "L'applicazione è già in esecuzione (icona nella tray).",
		"err.title": "{app} — errore", "cfg.err.title": "{app} — errore di configurazione",
		"err.details": "\n\nDettagli in {log}", "err.hook": "Hook tastiera non riuscito: %s",
		"err.mic": "Microfono: %s", "err.hotkey.cfg": "Scorciatoia in config.json: %s",
		"err.model.notfound": "Modello non trovato: %s\nRicompilare il progetto (.\\build.ps1) o correggere \"model\" in config.json",
		"err.server.repeat":  "whisper-server continua a bloccarsi — vedi {log}",
		"err.server.dead":    "Il server %s non risponde (server_url/server_autostart in config.json)",
		"err.server.timeout": "whisper-server non ha risposto entro %s",
		"err.server.start":   "whisper-server terminato all'avvio (vedi log)",
		"err.srv.noaddr":    "il server di riconoscimento remoto non è configurato — imposta il suo indirizzo nelle impostazioni",
		"err.webview":        "La finestra impostazioni richiede Microsoft WebView2 Runtime (incluso in Windows 11).\nInstallalo o modifica config.json manualmente.",
		"status.loading":     "Caricamento modello…", "status.ready": "Pronto — tieni %s e parla",
		"status.recording": "Registrazione…", "status.transcribing": "Riconoscimento…", "status.disabled": "Disattivato",
		"status.server.restart": "Server bloccato, riavvio…", "status.cfg.err": "Errore in config.json (vedi log)",
		"menu.settings": "Impostazioni…", "menu.enable": "Attiva", "menu.disable": "Disattiva",
		"menu.open.config": "Apri config.json", "menu.open.log": "Apri log",
		"menu.about": "Informazioni", "menu.quit": "Esci",
		"err.port.busy": "La porta %d è occupata da un altro programma. Cambia la porta nelle impostazioni.",
		"ov.speak": "Parla…", "ov.transcribing": "Riconoscimento", "ov.asking": "Attendo la tua risposta", "ov.inserted": "Inserito: %d caratteri", "ov.left": "%d s rimasti", "ov.esc": "1…9 · Invio · Esc annulla", "err.net.dns": "Nessuna connessione a %s — controlla internet", "err.net.timeout": "Il server non ha risposto in tempo — riprova", "err.net.down": "Impossibile connettersi — controlla internet", "err.net.cert": "La connessione non è sicura — controlla la data e l'antivirus", "err.answer": "Il server ha risposto in modo incomprensibile — riprova più tardi", "err.file.missing": "File non trovato", "err.file.denied": "Nessun accesso al file — chiudi il programma che lo tiene", "err.disk.full": "Il disco è pieno", "err.cancelled": "Annullato", "err.generic": "Non ha funzionato — i dettagli sono nel registro", "err.server.launch": "Impossibile avviare %s — controlla il percorso del server in Sistema",
		"ov.err.recognize": "Errore di riconoscimento (vedi log)", "ov.err.paste": "Non incollato: il testo è nell'ultimo risultato",
		"ov.moved":  "La finestra è cambiata: il testo è negli appunti",
		"copy.ok":   "Copiato",
		"copy.none": "Niente da copiare",
		"copy.fail": "Impossibile copiare: %s",
		"mic.busy":  "C'è una dettatura in corso, ora non si può controllare", "mic.check.ok": "Si sente bene: picco %.0f dB, voce nel %.0f%% della registrazione",
		"mic.check.quiet": "Troppo basso: picco %.0f dB — alza il livello del microfono in Windows o avvicinati", "mic.check.clipped": "Distorsione: %.1f%% dei campioni tagliati — abbassa il livello del microfono", "mic.check.silent": "Nessuna voce sentita — controlla che sia scelto il microfono giusto e non sia muto",
		"ov.quiet": "Troppo basso, non si è sentito quasi nulla", "ov.clipped": "Distorsione: il suono è stato tagliato",
		"ov.cmd.cancelled": "Annullato a voce",
		"ov.undo.done": "Cancellati: %d caratteri", "ov.undo.none": "Niente da cancellare", "ov.undo.typed": "Non cancellato: hai scritto o cliccato nel frattempo",
		"ov.undo.window": "Non cancellato: un'altra finestra", "ov.undo.sent": "Non cancellato: già inviato con Invio", "ov.undo.old": "Non cancellato: inserito troppo tempo fa",
		"ov.undo.breaks": "Non cancellato: il testo andava a capo", "ov.undo.long": "Non cancellato: il testo è troppo lungo", "ov.undo.moved": "Cancellazione interrotta: la finestra è cambiata",
		"ov.shape.next": "Prossima frase: %s", "cmd.heading": "titolo", "cmd.list": "elenco", "cmd.capital": "maiuscola", "cmd.quotes": "tra virgolette", "cmd.nostop": "senza punto",
		"ov.silence":       "Silenzio — nulla riconosciuto", "ov.server.loading": "Il server sta ancora caricando",
		"ov.stop.silence":  "Fermato dal silenzio",
		"ov.tooshort":  "Troppo breve — tieni premuti i tasti più a lungo",
		"ov.cancelled": "Annullato", "ov.editing": "Modifica: %s", "ov.translating": "Traduzione",
		"ov.llm.needed": "Questa lingua richiede il modulo LLM", "td.title": "Traduci in:", "td.plain": "Tradurre nella lingua di partenza",
		"cap.title": "{app} — scorciatoia", "cap.prompt": "Premi una nuova combinazione\n\nattuale: %s   ·   Esc — annulla",
		"cap.selected": "Scelto: %s", "cap.cancelled": "Annullato", "hk.taken": "%s è occupata da Windows: %s. La dettatura potrebbe non partire", "hk.lock": "bloccare il computer", "hk.desktop": "mostrare il desktop", "hk.explorer": "esplora file", "hk.run": "finestra Esegui", "hk.settings": "impostazioni di Windows", "hk.search": "ricerca", "hk.center": "centro notifiche", "hk.menu": "menu avanzato", "hk.clipboard": "cronologia appunti", "hk.gamebar": "barra di gioco", "hk.voice": "dettatura di Windows", "hk.project": "proiettare su schermo", "hk.tasks": "visualizzazione attività", "hk.layout": "cambio layout", "hk.newdesktop": "nuovo desktop", "hk.closedesktop": "chiudere il desktop", "hk.snip": "strumento di cattura", "hk.switch": "cambiare finestra", "hk.close": "chiudere la finestra", "hk.cycle": "scorrere le finestre", "hk.start": "menu Start", "hk.taskmgr": "gestione attività", "hk.secure": "schermata di sicurezza",
		"model.switching": "Cambio modello — riavvio…", "model.del.active": "Impossibile eliminare il modello attivo",
		"model.del.ok":   "Modello eliminato",
		"about.text":     "{app} %s\n\nVoce → testo alla posizione del cursore.\nPosiziona il cursore, tieni %s, parla, rilascia — il testo viene inserito.\n\nRiconoscimento: whisper.cpp, completamente locale e offline.\nModello: %s (lingua: %s)\n\nImpostazioni: clic sull'icona o config.json.\nLog: {log} (max ~2 MB).",
		"status.nomodel": "Nessun modello di riconoscimento scaricato — scegline uno nelle impostazioni",
		"tf.title": "%s non traduce. Tradurre comunque?",
		"tf.yes": "Sì, con %s",
		"tf.no": "Così com’è",
		"hk.lone": "%s è un tasto singolo: finché è la scorciatoia, non si può più digitare normalmente",
		"state.loaded.none": "niente in memoria",
		"state.week": "%d dettature · %d caratteri",
		"snd.ok": "livello a posto",
		"snd.quiet": "piano — parlate più vicino al microfono",
		"snd.clipped": "troppo forte, il suono è saturato",
		"snd.silent": "silenzio nella registrazione",
		"status.parked": "Motore scaricato — la scorciatoia lo risveglia",
		"status.nomodel.lang": "Per %s il modello %s non è installato — aprite «Lingue e modelli»",
		"menu.lastcopy":  "Copia l'ultimo risultato",
		"ov.copied":      "Copiato negli appunti", "ov.kept": "Annullato — il testo resta nell'ultimo risultato",
		"ov.llm.skipped": "Inserito senza il profilo «%s»",
		"fd.title":       "Il focus è cambiato — inserire?", "fd.here": "Inserisci qui", "fd.copy": "Copia",
		"fd.keep":            "Tieni",
		"ov.err.mic":         "Microfono non disponibile — controlla il dispositivo nelle impostazioni",
		"ov.notranslate":     "Il modello attivo non traduce — inserito così com'è stato riconosciuto",
		"ov.engine.fallback": "L'altro motore non è partito — si resta su quello attuale",
		"route.speech":       "Parlato in %s", "route.other": "Altre lingue", "route.translate": "Traduzione",
		"route.lang.auto":    "qualsiasi lingua",
		"route.why.language": "qui è più preciso, con la punteggiatura", "route.why.otherlang": "99 lingue",
		"route.why.translate": "solo Whisper traduce", "route.why.notinstalled": "il modello russo non è installato",
		"route.why.unknownlang": "nessuna lingua impostata — solo Whisper la riconosce", "route.why.forced": "forzato in config.json",
		"status.line": "Pronto · %s · %.1f GB liberi", "state.ram.free": "%d MB liberi",
		"ago.now": "proprio ora", "ago.min": "%d min fa", "ago.hour": "%d h fa",
		"chars": "%d caratteri", "inserted.into": "inserito in %s",
		"punct.prompt":         "Aggiungi punteggiatura e maiuscole. Non cambiare le parole, non tradurre, non aggiungere nulla. Restituisci solo il testo corretto.",
		"err.sherpa.notfound":  "riconoscitore sherpa non trovato: %s",
		"err.sherpa.start":     "sherpa-server si è chiuso durante l'avvio (vedi il registro)",
		"err.sherpa.translate": "questo modello non sa tradurre",
		"err.sherpa.model":     "File del modello non trovato: %s — scaricalo nelle impostazioni o correggi sherpa_model in config.json",
		"err.hotkey.dup":       "La scorciatoia «%s» è assegnata due volte — le scorciatoie devono essere uniche",
		"cfg.err.recovered":    "config.json è danneggiato (%s).\nIl file è stato salvato come %s e le impostazioni sono tornate ai valori predefiniti.",
		"err.disk.space":       "poco spazio su disco: %d MB liberi, ~%d MB necessari",
		"err.model.dir":       "non è stato possibile creare la cartella %s — controlla i permessi",
		"err.save":             "impostazioni non salvate: %s — restano quelle di prima",
		"err.port":             "la porta %d non va bene: scegli un numero tra 1024 e 65535",
		"err.nolangs":          "lascia almeno una lingua per la domanda sulla traduzione",
		"ov.mic.lost":          "Microfono scollegato: registrazione interrotta",
		"err.hash":             "il file scaricato è danneggiato — riprova",
		"models.check.ok":      "Modelli controllati: %d — tutti i file sono integri",
		"models.check.none":    "Niente da controllare: nessun modello installato ha un hash di riferimento",
		"models.check.bad":     "File danneggiati: %s — scarica di nuovo il modello",
		"hist.insert.gone":     "voce non trovata",
		"ov.aim": "Clicca il campo dove incollare · Esc annulla",
		"hist.aim.armed": "clicca il campo dove incollare",
		"hist.aim.busy": "sto già aspettando un clic",
		"hist.aim.off": "incollaggio annullato",
		"hist.insert.nowin":    "non c'è dove incollare: il testo è negli appunti",
		"hist.insert.ok":       "incollato in «%s»",
		"lists.bad":            "questo file non va bene",
		"lists.saved":          "salvato in %s",
		"lists.added":          "aggiunti: %d, saltati: %d",
		"lists.save.title":     "Dove salvare le liste",
		"lists.open.title":     "Quale file caricare",
		"un.title":             "{app} — Disinstallazione", "un.confirm": "Rimuovere {app} da questo computer?",
		"un.data": "Eliminare anche impostazioni, modelli scaricati e trascrizioni dei sottotitoli?", "un.done": "{app} è stato rimosso.",
		"srv.restarting": "Riavvio del riconoscitore con le nuove impostazioni…",
	}
	msgs["pl"] = map[string]string{
		"menu.sub.start": "Włącz napisy",
		"menu.sub.stop": "Wyłącz napisy",
		"sub.title": "Napisy",
		"sub.title.tr": "Tłumaczenie",
		"sub.follow": "do najnowszego ↓",
		"sub.loading": "Wczytuję model…",
		"sub.listening": "Słucham dźwięku komputera…",
		"sub.restarting": "Rozpoznawanie uruchamia się ponownie…",
		"sub.err.nomodel": "Napisy: dla ich języka nie zainstalowano modelu",
		"sub.err.engine": "Napisy: silnik rozpoznawania nie wystartował",
		"sub.err.device": "Napisy: brak urządzenia odtwarzania do nasłuchu",
		"sub.err.lost": "Napisy zatrzymane: silnik rozpoznawania raz po raz padał",
		"app.name": "{app}", "already.running": "Aplikacja już działa (ikona w zasobniku).",
		"err.title": "{app} — błąd", "cfg.err.title": "{app} — błąd konfiguracji",
		"err.details": "\n\nSzczegóły w {log}", "err.hook": "Błąd hooka klawiatury: %s",
		"err.mic": "Mikrofon: %s", "err.hotkey.cfg": "Skrót w config.json: %s",
		"err.model.notfound": "Nie znaleziono modelu: %s\nZbuduj projekt (.\\build.ps1) lub popraw \"model\" w config.json",
		"err.server.repeat":  "whisper-server ciągle się zawiesza — zobacz {log}",
		"err.server.dead":    "Serwer %s nie odpowiada (server_url/server_autostart w config.json)",
		"err.server.timeout": "whisper-server nie odpowiedział w ciągu %s",
		"err.server.start":   "whisper-server zakończył się przy starcie (zobacz log)",
		"err.srv.noaddr":    "zdalny serwer rozpoznawania nie jest skonfigurowany — podaj adres w ustawieniach",
		"err.webview":        "Okno ustawień wymaga Microsoft WebView2 Runtime (zawarty w Windows 11).\nZainstaluj go lub edytuj config.json ręcznie.",
		"status.loading":     "Ładowanie modelu…", "status.ready": "Gotowy — przytrzymaj %s i mów",
		"status.recording": "Nagrywanie…", "status.transcribing": "Rozpoznawanie…", "status.disabled": "Wyłączone",
		"status.server.restart": "Serwer padł, restart…", "status.cfg.err": "Błąd w config.json (zobacz log)",
		"menu.settings": "Ustawienia…", "menu.enable": "Włącz", "menu.disable": "Wyłącz",
		"menu.open.config": "Otwórz config.json", "menu.open.log": "Otwórz log",
		"menu.about": "O programie", "menu.quit": "Zakończ",
		"err.port.busy": "Port %d jest zajęty przez inny program. Zmień port w ustawieniach.",
		"ov.speak": "Mów…", "ov.transcribing": "Rozpoznawanie", "ov.asking": "Czekam na odpowiedź", "ov.inserted": "Wstawiono: %d znaków", "ov.left": "zostało %d s", "ov.esc": "1…9 · Enter · Esc anuluje", "err.net.dns": "Brak połączenia z %s — sprawdź internet", "err.net.timeout": "Serwer nie odpowiedział na czas — spróbuj ponownie", "err.net.down": "Nie udało się połączyć — sprawdź internet", "err.net.cert": "Połączenie nie jest bezpieczne — sprawdź datę i antywirus", "err.answer": "Serwer odpowiedział niezrozumiale — spróbuj później", "err.file.missing": "Nie znaleziono pliku", "err.file.denied": "Brak dostępu do pliku — zamknij program, który go trzyma", "err.disk.full": "Brak miejsca na dysku", "err.cancelled": "Anulowano", "err.generic": "Nie udało się — szczegóły w dzienniku", "err.server.launch": "Nie udało się uruchomić %s — sprawdź ścieżkę serwera w sekcji System",
		"ov.err.recognize": "Błąd rozpoznawania (zobacz log)", "ov.err.paste": "Nie wklejono — tekst jest w ostatnim wyniku",
		"ov.moved":  "Okno się zmieniło — tekst jest w schowku",
		"copy.ok":   "Skopiowano",
		"copy.none": "Nie ma czego kopiować",
		"copy.fail": "Nie udało się skopiować: %s",
		"mic.busy":  "Trwa dyktowanie, teraz nie można sprawdzić", "mic.check.ok": "Brzmi dobrze: szczyt %.0f dB, mowa w %.0f%% nagrania",
		"mic.check.quiet": "Za cicho: szczyt %.0f dB — podnieś poziom mikrofonu w Windows albo usiądź bliżej", "mic.check.clipped": "Przesterowanie: obcięto %.1f%% próbek — zmniejsz poziom mikrofonu", "mic.check.silent": "Nie słychać mowy — sprawdź, czy wybrany jest właściwy mikrofon i czy nie jest wyciszony",
		"ov.quiet": "Za cicho, prawie nic nie było słychać", "ov.clipped": "Przesterowanie — dźwięk został obcięty",
		"ov.cmd.cancelled": "Anulowano głosem",
		"ov.undo.done": "Usunięto: %d znaków", "ov.undo.none": "Nie ma czego usunąć", "ov.undo.typed": "Nie usunięto: od tamtej pory coś wpisano lub kliknięto",
		"ov.undo.window": "Nie usunięto: inne okno", "ov.undo.sent": "Nie usunięto: już wysłano Enterem", "ov.undo.old": "Nie usunięto: wstawiono zbyt dawno",
		"ov.undo.breaks": "Nie usunięto: tekst miał złamania wiersza", "ov.undo.long": "Nie usunięto: tekst jest za długi", "ov.undo.moved": "Usuwanie przerwane: okno się zmieniło",
		"ov.shape.next": "Następna fraza: %s", "cmd.heading": "nagłówek", "cmd.list": "lista", "cmd.capital": "wielka litera", "cmd.quotes": "w cudzysłowie", "cmd.nostop": "bez kropki",
		"ov.silence":       "Cisza — nic nie rozpoznano", "ov.server.loading": "Serwer wciąż się ładuje",
		"ov.stop.silence":  "Zatrzymano po ciszy",
		"ov.tooshort":  "Za krótko — przytrzymaj klawisze dłużej",
		"ov.cancelled": "Anulowano", "ov.editing": "Edycja: %s", "ov.translating": "Tłumaczenie",
		"ov.llm.needed": "Ten język wymaga modułu LLM", "td.title": "Tłumacz na:", "td.plain": "Przetłumacz na język źródłowy",
		"cap.title": "{app} — skrót", "cap.prompt": "Naciśnij nową kombinację\n\nobecna: %s   ·   Esc — anuluj",
		"cap.selected": "Wybrano: %s", "cap.cancelled": "Anulowano", "hk.taken": "%s jest zajęty przez Windows: %s. Dyktowanie może się nie zacząć", "hk.lock": "blokada komputera", "hk.desktop": "pokaż pulpit", "hk.explorer": "eksplorator plików", "hk.run": "okno Uruchom", "hk.settings": "ustawienia Windows", "hk.search": "wyszukiwanie", "hk.center": "centrum powiadomień", "hk.menu": "menu zaawansowane", "hk.clipboard": "historia schowka", "hk.gamebar": "pasek gier", "hk.voice": "pisanie głosem Windows", "hk.project": "projektowanie obrazu", "hk.tasks": "widok zadań", "hk.layout": "zmiana układu klawiatury", "hk.newdesktop": "nowy pulpit", "hk.closedesktop": "zamknięcie pulpitu", "hk.snip": "narzędzie wycinania", "hk.switch": "przełączanie okien", "hk.close": "zamknięcie okna", "hk.cycle": "przeglądanie okien", "hk.start": "menu Start", "hk.taskmgr": "menedżer zadań", "hk.secure": "ekran zabezpieczeń",
		"model.switching": "Zmiana modelu — restart…", "model.del.active": "Nie można usunąć aktywnego modelu",
		"model.del.ok":   "Model usunięty",
		"about.text":     "{app} %s\n\nGłos → tekst w pozycji kursora.\nUstaw kursor, przytrzymaj %s, mów, puść — tekst zostanie wstawiony.\n\nRozpoznawanie: whisper.cpp, w pełni lokalnie i offline.\nModel: %s (język: %s)\n\nUstawienia: kliknij ikonę lub config.json.\nLogi: {log} (maks. ~2 MB).",
		"status.nomodel": "Nie pobrano żadnego modelu — wybierz go w ustawieniach",
		"tf.title": "%s nie tłumaczy. Przetłumaczyć mimo to?",
		"tf.yes": "Tak, przez %s",
		"tf.no": "Jak rozpoznano",
		"hk.lone": "%s to pojedynczy klawisz: dopóki jest skrótem, nie da się go normalnie wpisać",
		"state.loaded.none": "nic nie wczytano",
		"state.week": "%d dyktowań · %d znaków",
		"snd.ok": "poziom w porządku",
		"snd.quiet": "cicho — mów bliżej mikrofonu",
		"snd.clipped": "za głośno, dźwięk się przesterował",
		"snd.silent": "cisza w nagraniu",
		"status.parked": "Silnik wyładowany — skrót go obudzi",
		"status.nomodel.lang": "Dla języka %s model %s nie jest zainstalowany — otwórz „Języki i modele”",
		"menu.lastcopy":  "Kopiuj ostatni wynik",
		"ov.copied":      "Skopiowano do schowka", "ov.kept": "Anulowano — tekst został w ostatnim wyniku",
		"ov.llm.skipped": "Wstawiono bez profilu „%s”",
		"fd.title":       "Zmieniło się okno — wstawić?", "fd.here": "Wstaw tutaj", "fd.copy": "Kopiuj",
		"fd.keep":            "Zachowaj",
		"ov.err.mic":         "Mikrofon niedostępny — sprawdź urządzenie w ustawieniach",
		"ov.notranslate":     "Aktywny model nie tłumaczy — wstawiono tak, jak rozpoznano",
		"ov.engine.fallback": "Drugi silnik nie wystartował — zostaje bieżący",
		"route.speech":       "Mowa w %s", "route.other": "Inne języki", "route.translate": "Tłumaczenie",
		"route.lang.auto":    "dowolny język",
		"route.why.language": "tu dokładniej, ze znakami", "route.why.otherlang": "99 języków",
		"route.why.translate": "tylko Whisper tłumaczy", "route.why.notinstalled": "model rosyjski nie jest zainstalowany",
		"route.why.unknownlang": "nie ustawiono języka — rozpozna go tylko Whisper", "route.why.forced": "wymuszone w config.json",
		"status.line": "Gotowe · %s · %.1f GB wolnych", "state.ram.free": "%d MB wolnych",
		"ago.now": "przed chwilą", "ago.min": "%d min temu", "ago.hour": "%d godz. temu",
		"chars": "%d znaków", "inserted.into": "wstawiono do %s",
		"punct.prompt":         "Dodaj znaki interpunkcyjne i wielkie litery. Nie zmieniaj słów, nie tłumacz, nic nie dopisuj. Zwróć tylko poprawiony tekst.",
		"err.sherpa.notfound":  "nie znaleziono rozpoznawania sherpa: %s",
		"err.sherpa.start":     "sherpa-server zakończył się przy starcie (zobacz dziennik)",
		"err.sherpa.translate": "ten model nie umie tłumaczyć",
		"err.sherpa.model":     "Nie znaleziono pliku modelu: %s — pobierz go w ustawieniach albo popraw sherpa_model w config.json",
		"err.hotkey.dup":       "Skrót „%s” jest przypisany dwa razy — skróty muszą być unikalne",
		"cfg.err.recovered":    "config.json jest uszkodzony (%s).\nPlik zapisano jako %s, a ustawienia wróciły do domyślnych.",
		"err.disk.space":       "mało miejsca na dysku: %d MB wolnych, potrzeba ~%d MB",
		"err.model.dir":       "nie udało się utworzyć folderu %s — sprawdź uprawnienia",
		"err.save":             "ustawienia nie zapisane: %s — zostają poprzednie",
		"err.port":             "port %d się nie nadaje: wybierz numer od 1024 do 65535",
		"err.nolangs":          "zostaw przynajmniej jeden język do pytania o tłumaczenie",
		"ov.mic.lost":          "Mikrofon odłączony — nagranie przerwane",
		"err.hash":             "pobrany plik jest uszkodzony — spróbuj ponownie",
		"models.check.ok":      "Sprawdzone modele: %d — wszystkie pliki są całe",
		"models.check.none":    "Nie ma czego sprawdzać — żaden zainstalowany model nie ma wzorcowego skrótu",
		"models.check.bad":     "Uszkodzone pliki: %s — pobierz model ponownie",
		"hist.insert.gone":     "nie znaleziono wpisu",
		"ov.aim": "Kliknij pole, w które wkleić · Esc anuluje",
		"hist.aim.armed": "kliknij pole, w które wkleić",
		"hist.aim.busy": "już czekam na kliknięcie",
		"hist.aim.off": "wklejanie anulowane",
		"hist.insert.nowin":    "nie ma gdzie wkleić — tekst jest w schowku",
		"hist.insert.ok":       "wklejono w „%s”",
		"lists.bad":            "ten plik nie pasuje",
		"lists.saved":          "zapisano w %s",
		"lists.added":          "dodano: %d, pominięto: %d",
		"lists.save.title":     "Gdzie zapisać listy",
		"lists.open.title":     "Który plik wczytać",
		"un.title":             "{app} — Odinstalowanie", "un.confirm": "Usunąć {app} z tego komputera?",
		"un.data": "Usunąć także ustawienia, pobrane modele i transkrypcje napisów?", "un.done": "{app} został usunięty.",
		"srv.restarting": "Ponowne uruchamianie rozpoznawania z nowymi ustawieniami…",
	}

	settingsStrings["de"] = map[string]string{
		"S_SUBS_ON": "Untertitel laufen",
		"S_SUBS_HERO_HINT": "starten — Gespräch oder Video wird zu Text auf dem Bildschirm",
		"S_SUBS_HERO_START": "Manuell starten",
		"S_SUBS_HERO_STOP": "Beenden",
		"S_UPDATED_TO": "Auf Version %s aktualisiert",
		"S_SUBS_HERO_LIVE": "%s · Wörter schon beim Sprechen",
		"S_SUBS_HERO_PHRASES": "%s · Text nach jedem Satz",
		"S_NAV_SUBS": "Untertitel",
		"S_SUBS_HINT": "Die App schreibt mit, was der Computer abspielt – ein Gespräch, ein Video – und zeigt die Worte in einem Fenster über allen anderen.",
		"S_SUBS_NOTICE": "Ihre Gesprächspartner sollten wissen, dass ihre Worte mitgeschrieben werden.",
		"S_SUBS_START": "Untertitel starten",
		"S_SUBS_STOP": "Untertitel beenden",
		"S_SUBS_OFF": "Untertitel sind aus",
		"S_SUBS_LOADING": "%s wird geladen…",
		"S_SUBS_LISTEN_LIVE": "Höre auf den Ton des Computers · %s, Wörter schon beim Sprechen",
		"S_SUBS_LISTEN_PHRASES": "Höre auf den Ton des Computers · %s, Text nach jedem Satz",
		"S_SUBS_RESTARTING": "Die Erkennung startet neu…",
		"S_SUBS_LANG": "Sprache der Rede",
		"S_SUBS_LANG_SUB": "Die Sprache dessen, was Sie hören, getrennt von der Diktiersprache",
		"S_SUBS_M_LIVE": "%s schreibt mit: Die Wörter erscheinen schon beim Sprechen.",
		"S_SUBS_M_PHRASES": "%s schreibt mit: Der Text erscheint, nachdem ein Satz zu Ende ist.",
		"S_SUBS_M_ABSENT": "Für diese Sprache ist %s gewählt, aber nicht installiert.",
		"S_SUBS_M_WHISPER": "%s ist ein Whisper-Modell: Auf dem Prozessor braucht es viele Sekunden pro Satz, der Text hinkt der Rede deutlich hinterher. Für Untertitel passt ein Streaming-Modell oder ein schnelles Modell besser.",
		"S_SUBS_GO_MODELS": "„Sprachen & Modelle“ öffnen",
		"S_SUBS_REMOTE": "Die Erkennung läuft auf einem anderen Rechner: Jeder Satz wird an diesen Server gesendet.",
		"S_SUBS_TR": "Übersetzung",
		"S_SUBS_TR_SUB": "Das Fenster zeigt die Übersetzung statt des erkannten Textes",
		"S_SUBS_TR_OFF": "Aus",
		"S_SUBS_TR_LOCAL": "%s übersetzt: Die Übersetzung kommt Satzteil für Satzteil, einen Moment nach der Rede.",
		"S_SUBS_TR_API": "Ihr Server für die Nachbearbeitung (%s) übersetzt: Jeder Satzteil wird an ihn gesendet, und die Übersetzung kommt Satzteil für Satzteil, einen Moment nach der Rede.",
		"S_SUBS_TR_NONE": "Es gibt weder ein Übersetzungsmodell noch einen Server für die Nachbearbeitung noch ein Bearbeitungsmodell, daher gibt es keine Übersetzung.",
		"S_SUBS_TR_GO": "„Nachbearbeitung“ öffnen",
		"S_SUBS_TR_SAME": "Die Übersetzungssprache ist dieselbe wie die Sprache der Rede: Es gibt nichts zu übersetzen.",
		"S_SUBS_WHO": "Wer übersetzt",
		"S_SUBS_WHO_SUB": "Das Programm nimmt den besten verfügbaren, oder Sie legen einen fest",
		"S_SUBS_WHO_AUTO": "Das Programm wählt",
		"S_SUBS_WHO_QUALITY": "Hochwertig: ein Modell auf der Grafikkarte",
		"S_SUBS_WHO_SERVER": "Ihr eigener Server",
		"S_SUBS_TRM": "Übersetzungsmodell",
		"S_SUBS_TRM_SUB": "%s, läuft auf einer Grafikkarte mit mindestens 4 GB eigenem Speicher",
		"S_SUBS_TRM_ABSENT": "Nicht heruntergeladen",
		"S_SUBS_TRM_GETTING": "Wird heruntergeladen… %s %",
		"S_SUBS_TRM_CHECKING": "Die Datei wird geprüft…",
		"S_SUBS_TRM_READY": "Bereit",
		"S_SUBS_TRM_RUNNING": "Bereit, auf der Grafikkarte geladen",
		"S_SUBS_TRM_BAD": "Die Datei im Ordner models ist beschädigt oder unvollständig",
		"S_SUBS_TRM_ERR": "Nicht heruntergeladen: %s",
		"S_SUBS_TRM_SMALL": "Auf diesem Computer läuft es nicht: Die Grafikkarte hat %s MB eigenen Speicher, nötig sind 4 GB.",
		"S_SUBS_TRM_INTEGRATED": "Auf diesem Computer läuft es nicht: Es gibt nur integrierte Grafik, die sich den Speicher mit dem Computer teilt.",
		"S_SUBS_TRM_NOCARD": "Auf diesem Computer läuft es nicht: Windows meldet keine Grafikkarte.",
		"S_SUBS_TRM_UNKNOWN": "Auf diesem Computer läuft es nicht: Diese Windows-Version beschreibt ihre Grafikkarte nicht.",
		"S_SUBS_TRM_DOWN": "Nicht gestartet: %s",
		"S_SUBS_TRM_DOWN_NEXT": "Bis zum nächsten Start der Untertitel übersetzt der nächste in der Reihe.",
		"S_SUBS_TRM_DOWN_HEARD": "Bis zum nächsten Start der Untertitel wird der Text so gezeigt, wie er gehört wurde.",
		"S_SUBS_TRM_DEL_SURE": "Löschen?",
		"S_SUBS_TRM_W_VULKAN": "der Grafiktreiber bietet kein Vulkan",
		"S_SUBS_TRM_W_DEVICE": "die Engine sieht diese Grafikkarte nicht",
		"S_SUBS_TRM_W_MEMORY": "auf der Grafikkarte ist nicht genug Speicher frei",
		"S_SUBS_TRM_W_CRASH": "die Engine hat beim Laden aufgehört",
		"S_SUBS_TRM_W_TIMEOUT": "die Engine hat nicht rechtzeitig geantwortet",
		"S_SUBS_TRM_W_PORT": "ihr Port ist von einem anderen Programm belegt",
		"S_SUBS_TRM_W_EXE": "llama-server-gpu.exe fehlt neben dem Programm",
		"S_SUBS_TRM_DL": "Herunterladen",
		"S_SUBS_TRM_CANCEL": "Abbrechen",
		"S_SUBS_TRM_DEL": "Löschen",
		"S_SUBS_TRM_TERMS": "Nutzungsbedingungen des Modells",
		"S_SUBS_TR_QUALITY": "Das Übersetzungsmodell auf der Grafikkarte (%s) übersetzt: Die Übersetzung kommt Satzteil für Satzteil, einen Moment nach der Rede.",
		"S_SUBS_TR_MISS_CARD": "„Hochwertig“ ist gewählt, aber die Grafikkarte dieses Computers passt nicht, daher wird der Text so gezeigt, wie er gehört wurde.",
		"S_SUBS_TR_MISS_MODEL": "„Hochwertig“ ist gewählt, aber das Übersetzungsmodell ist nicht heruntergeladen, daher wird der Text so gezeigt, wie er gehört wurde.",
		"S_SUBS_TR_MISS_START": "„Hochwertig“ ist gewählt, aber das Übersetzungsmodell ist auf der Grafikkarte nicht gestartet, daher wird der Text bis zum nächsten Start der Untertitel so gezeigt, wie er gehört wurde.",
		"S_SUBS_TR_MISS_SERVER": "Ihr eigener Server ist gewählt, aber es ist kein Server für die Nachbearbeitung eingerichtet, daher wird der Text so gezeigt, wie er gehört wurde.",
		"S_SUBS_ORIG": "Auch den Originaltext zeigen",
		"S_SUBS_ORIG_SUB": "das Fenster zeigt über der Übersetzung, was gesagt wurde; die Linie dazwischen lässt sich mit der Maus ziehen",
		"S_SUBS_KEY": "Tastenkürzel",
		"S_SUBS_KEY_SUB": "Startet und beendet die Untertitel aus jedem Fenster",
		"S_SUBS_SAVE": "Mitschrift speichern",
		"S_SUBS_SAVE_SUB": "Eine Textdatei pro Durchlauf im Ordner transcripts neben dem Programm, bei eingeschalteter Übersetzung eine zweite mit der Übersetzung",
		"S_SUBS_FOLDER": "Ordner der Mitschriften",
		"S_SUBS_WINDOW": "Untertitelfenster",
		"S_SUBS_ROWS": "Textzeilen",
		"S_SUBS_SIZE": "Textgröße",
		"S_SUBS_SIZE_S": "Klein",
		"S_SUBS_SIZE_M": "Mittel",
		"S_SUBS_SIZE_L": "Groß",
		"S_SUBS_BOX": "Fenstergröße",
		"S_SUBS_BOX_SUB": "ziehen Sie Ränder und Ecken des Fensters mit der Maus; die Größe wird für jede Bildschirmauflösung gemerkt, „Textzeilen“ gilt, solange es keine gibt",
		"S_SUBS_BOX_BTN": "Größe zurücksetzen",
		"S_SUBS_BOX_DONE": "Das Fenster kehrt zur Standardgröße zurück",
		"S_SUBS_FREE": "Eigener Platz",
		"S_SUBS_FREE_SUB": "das Fenster lässt sich überallhin ziehen, der Platz wird für jede Bildschirmauflösung gemerkt",
		"S_SUBS_POS": "Wo das Fenster erscheint",
		"S_SUBS_POS_SCHEME_SUB": "klicken Sie auf den Bildschirm — das Fenster rückt dorthin",
		"S_SUBS_POS_DRAG_SUB": "ziehen Sie die Markierung oder das Fenster selbst — es landet überall",
		"S_SUBS_MON": "Bildschirm",
		"S_SUBS_MON_SUB": "auf welchem Monitor das Untertitelfenster erscheint",
		"S_TITLE": "{app} — Einstellungen", "S_DICT_HINT": "Begriffe, Namen und Abkürzungen, durch Kommas getrennt — ein Hinweis fürs Gehör, keine Befehle. Gilt für Whisper; russische Sprache über GigaAM ignoriert das. Der Standardsatz folgt der Erkennungssprache, bis Sie eigenes eintragen.",
		"S_TR_DEFAULT": "Sprache der Textausgabe ändern", "S_TR_TARGET": "Standardsprache der Textausgabe", "S_TR_ASK": "Nach der Ausgabesprache fragen", "S_TR_ASK_NEVER": "Nicht fragen — sofort übersetzen",
		"S_SRCLANG_SUB": "Sie sprechen sie; sie bestimmt das Erkennungsmodell",
		"S_TR_LANGS_SUB": "diese Sprachen werden beim Einfügen zu Schaltflächen auf der Leiste",
		"S_TR_UNAVAIL": "%s kann nicht übersetzen — der Text bleibt in der Ausgangssprache",
		"S_TR_GET_WHISPER": "%s für die Übersetzung laden (%d MB)",
		"S_TR_STANDIN": "%s kann nicht übersetzen — die Leiste bietet stattdessen %s an",
		"S_TR_LOCK": "%s kann nicht aus der Liste entfernt werden — sie ist die Standardsprache der Textausgabe. Wählen Sie eine andere Standardsprache, dann lässt sich %s ausschließen.",
		"S_TR_LOCK_OK": "Verstanden",
		"S_TR_ONE": "Mehrere Sprachen sind angehakt, aber ohne Nachfrage geht die Übersetzung immer in eine — %s (Standardsprache der Textausgabe). Die übrigen bleiben angehakt, werden aber deaktiviert.",
		"S_TR_NOMODEL": "%s kann nicht übersetzen. Wenn Sie fortfahren, wird die Übersetzung ausgeschaltet und ist nicht verfügbar, solange dieses Modell arbeitet.",
		"S_TR_CONFIRM": "Bestätigen",
		"S_TR_ASK_ALWAYS": "Jedes Mal fragen", "S_TR_ASK_TIMEOUT": "Fragen, mit Timeout", "S_TR_SECONDS": "Timeout, Sek.",
		"S_TR_LANGS": "Sprachen im Dialog",
		"S_LLM_HINT": "Angehakte Profile greifen nacheinander, von oben nach unten, beim normalen Diktat. Nichts angehakt — der Text wird unverändert eingefügt.",
		"S_PROF_ADD": "Hinzufügen", "S_PROF_NAME": "Name", "S_PROF_PROMPT": "Prompt", "S_PROF_TEST": "Test",
		"S_HOTKEY": "Tastenkürzel", "S_CHANGE": "Ändern…", "S_UILANG": "Sprache der Oberfläche", "S_AUTO": "Wie im System",
		"S_SEC_SOUND": "Ton", "S_BEEP": "Tonsignale der Aufnahme", "S_SOUND": "Signalton",
		"S_SND_SPEECH": "Windows-Stimme", "S_SND_CHIME": "Glöckchen", "S_SND_SOFT": "Sanft", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Enter nach dem Einfügen drücken (Auto-Senden)", "S_RESTORE": "Zwischenablage nach Einfügen wiederherstellen",
		"S_NAV_HISTORY": "Verlauf", "S_HIST_ON": "Verlauf der Diktate führen", "S_HIST_ON_SUB": "nur Text, auf diesem Rechner; Ton wird nie gespeichert",
		"S_HIST_DAYS": "Wie viele Tage aufbewahren", "S_HIST_MAX": "Wie viele Einträge aufbewahren",
		"S_HIST_SKIP": "Aus diesen Programmen nie aufzeichnen", "S_HIST_SKIP_SUB": "durch Komma getrennt: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Programm hinzufügen", "S_SKIP_EDIT_DLG": "Programm bearbeiten", "S_SKIP_NAME": "Programmname", "S_SKIP_NAME_SUB": "Dateiname ohne Pfad: keepass.exe. Ein Stern am Ende erfasst alle Versionen: 1password*", "S_SKIP_OPEN": "Derzeit geöffnete Programme", "S_SKIP_REFRESH": "Liste aktualisieren", "S_SKIP_PICKED": "%d von %d gewählt", "S_SKIP_NONE": "Nichts gewählt", "S_SKIP_EMPTY": "Die Liste ist leer — der Verlauf wird aus allen Programmen geführt", "S_SKIP_ADD_BTN": "Programm hinzufügen", "S_SKIP_HINT": "Was Sie in diese Programme diktieren, landet nicht im Verlauf. Das Einfügen selbst funktioniert wie gewohnt.",
		"S_HIST_LIST": "Einträge", "S_HIST_CLEAR": "Leeren", "S_HIST_TILL": "bis %s", "S_HIST_TILL1": "bis morgen", "S_HIST_TILL_FULL": "Wird am %s gelöscht — Aufbewahrung %s", "S_HIST_LIST_HINT": "Was diktiert wurde: kopieren, in ein beliebiges Fenster einfügen oder löschen.", "S_HIST_COPY": "Kopieren",
		"S_HIST_KEEP": "Wie lange aufbewahren",
		"S_UNIT_MIN": "Minuten",
		"S_UNIT_HOUR": "Stunden",
		"S_UNIT_DAY": "Tage",
		"S_HIST_FIND": "Im Verlauf suchen…", "S_HIST_EMPTY": "Noch kein Verlauf", "S_HIST_ASK": "Den gesamten Diktatverlauf löschen?",
		"S_SEC_CMD": "Sprachbefehle", "S_CMD_HINT": "Gesagtes wird zu einem Zeilenumbruch, einem Zeichen, einem Abbruch, dem Löschen der letzten Einfügung oder einer Form für die Phrase, statt im Text zu landen. Als ganze Wörter erkannt, von oben nach unten angewendet, nach den Ersetzungen; Formbefehle wirken am Ende des Diktats.",
		"S_CMD_ADD": "Befehl hinzufügen", "S_CMD_PRESET": "Übliche hinzufügen", "S_CMD_PH": "die Phrase, die Sie sagen werden",
		"S_CMD_NEWLINE": "Zeilenumbruch", "S_CMD_PARAGRAPH": "neuer Absatz", "S_CMD_TEXT": "Text einfügen", "S_CMD_CANCEL": "Diktat abbrechen",
		"S_CMD_TEXT_PH": "was einfügen", "S_CMD_EMPTY": "Noch keine Befehle", "S_CMD_DEL": "Befehl löschen",
		"S_CMD_P_NEWLINE": "neue Zeile", "S_CMD_P_PARAGRAPH": "neuer Absatz", "S_CMD_P_CANCEL": "abbrechen",
		"S_CMD_UNDO": "letzte Einfügung löschen", "S_CMD_HEADING": "als Überschrift setzen", "S_CMD_LIST": "als Liste setzen",
		"S_CMD_CAPITAL": "großer Anfangsbuchstabe", "S_CMD_QUOTES": "in Anführungszeichen setzen", "S_CMD_NOSTOP": "Schlusspunkt weglassen",
		"S_CMD_HEADING_PH": "Zeichen vor der Überschrift, z. B. #", "S_CMD_LIST_PH": "Wörter vor dem letzten Eintrag, z. B. und", "S_CMD_QUOTES_PH": "öffnendes und schließendes Zeichen, z. B. „“",
		"S_CMD_END_NOTE": "Am Ende des Diktats sagen; allein gesagt, formt der Befehl das nächste.",
		"S_CMD_UNDO_NOTE": "Allein sagen. Löscht die letzte Einfügung, wenn Sie seitdem nichts getippt oder geklickt haben und das Fenster dasselbe ist; Text, der schon dort stand, wird nie angefasst.",
		"S_CMD_P_UNDO": "lösch das", "S_CMD_P_HEADING": "als Überschrift", "S_CMD_P_LIST": "als Liste", "S_CMD_P_CAPITAL": "großgeschrieben",
		"S_CMD_P_QUOTES": "in Anführungszeichen", "S_CMD_P_NOSTOP": "ohne Punkt", "S_CMD_P_AND": "und", "S_CMD_P_MARKS": "„“",
		"S_SEC_REPLACE": "Ersetzungen nach der Erkennung", "S_REPLACE_HINT": "Falsch Gehörtes wird zu dem, was gemeint war — direkt nach der Erkennung, vor den Prompts. Von oben nach unten angewendet.",
		"S_REPL_WHOLE_FULL": "Nur ganze Wörter", "S_REPL_CASE_FULL": "Groß-/Kleinschreibung beachten", "S_CMD_ACTION": "Aktion",
		"S_FM_ADD": "Hinzufügen",
		"S_TIP_REPL_LANG": "Die Regel greift nur, wenn Sie in der gewählten Sprache diktieren. „alle Sprachen“ — sie greift immer.",
		"S_TIP_REPL_CASE": "Groß- und Kleinbuchstaben zählen: „git“ und „Git“ sind verschiedene Wörter. Aus — die Schreibung ist egal.",
		"S_TIP_REPL_WHOLE": "Die Ersetzung greift nur, wenn der Text als eigenes Wort steht. Aus — sie trifft auch innerhalb anderer Wörter.",
		"S_TIP_CMD_ACTION": "Was passiert, wenn Sie die Phrase sagen: Zeilenumbruch, neuer Absatz, eigener Text, Abbruch des Diktats, Löschen der letzten Einfügung oder eine Form für die Phrase — Überschrift, Liste, großer Anfangsbuchstabe, Anführungszeichen oder ohne Schlusspunkt.",
		"S_LIST_FILTER_PH": "suchen…",
		"S_REPL_DEL": "Ersetzung löschen",
		"S_LIST_NOTHING": "Nichts gefunden: „%s“",
		"S_FM_T_REPL_ADD": "Ersetzung hinzufügen", "S_FM_T_REPL_EDIT": "Ersetzung bearbeiten",
		"S_FM_T_CMD_ADD": "Befehl hinzufügen", "S_FM_T_CMD_EDIT": "Befehl bearbeiten",
		"S_MT_DEL": "Modell löschen", "S_MT_DEL_PROMPT": "Prompt löschen", "S_MT_DL": "Modell herunterladen",
		"S_MT_TR_OFF": "Übersetzung ausschalten", "S_MT_TR_ONE": "Übersetzen ohne Nachfrage", "S_MT_TR_LOCK": "Standard-Ausgabesprache",
		"S_MT_REMOTE": "Entfernter Server", "S_MT_POST": "Externer Server", "S_MT_HIST": "Verlauf leeren",
		"S_MT_RESET": "Einstellungen zurücksetzen", "S_MT_EXE": "Server-Pfad",
		"S_DICT_ADD": "Wort hinzufügen", "S_FM_T_DICT_ADD": "Wort hinzufügen", "S_DICT_EMPTY": "Noch keine Wörter",
		"S_DICT_ADD_PH": "ein Wort oder mehrere durch Kommas",
		"S_DICT_NOMODEL": "Das aktuelle Modell %s unterstützt das Wörterbuch nicht — nur Whisper-Modelle lesen es.",
		"S_OV_FREE": "Eigener Platz", "S_OV_FREE_SUB": "die Leiste lässt sich überallhin ziehen",
		"S_OVPOS_DRAG_SUB": "ziehen Sie die Leiste mit der Maus — sie landet überall",
		"S_OVMON_N": "Bildschirm %d",
		"S_POST_ENABLE": "Nachbearbeitung einschalten",
		"S_API_SUM_URL": "Adresse", "S_API_SUM_MODEL": "Modell", "S_API_SUM_KEY": "Schlüssel", "S_API_SUM_TIMEOUT": "Wartezeit",
		"S_API_SUM_STATE": "Status", "S_API_NO_MODEL": "nicht angegeben",
		"S_API_NONE": "nicht eingerichtet — die Nachbearbeitung läuft lokal",
		"S_POSTAPI_SETUP": "Einrichten", "S_API_EDIT": "Ändern", "S_API_KEY_DEL": "Schlüssel löschen", "S_API_DLG": "Externer Server",
		"S_LLM_CATALOG": "Modellkatalog", "S_LLM_BLOCK": "Installierte Modelle", "S_LLM_NONE_HINT": "Noch kein Modell installiert — laden Sie ein gefundenes mit dem Pfeil herunter, dann erscheint es hier", "S_LLM_IN_MEM": "im Speicher", "S_LLM_ON_DISK": "auf der Festplatte", "S_LLM_EJECT": "Aus dem Speicher entladen", "S_LLM_FOUND": "%d gefunden", "S_LLM_NOSEARCH": "noch keine Suche", "S_LLM_SEARCH_HINT": "Modellnamen eingeben und „Suchen“ drücken", "S_LLM_PICK_WAIT": "Verfügbar, sobald das Modell geladen ist", "S_LLM_INSTALLED": "installiert",
		"S_LLM_SUM_MODEL": "Modell", "S_LLM_SUM_SIZE": "Größe", "S_LLM_SUM_COUNT": "installiert", "S_LLM_SUM_RAM": "Speicher",
		"S_DLG_CLOSE": "Schließen", "S_LLM_NOPICK": "nicht gewählt", "S_NO_PROMPTS": "Noch keine Prompts", "S_PROF_DRAG": "ziehen, um die Reihenfolge zu ändern",
		"S_PROF_NAME_PH": "wie der Prompt heißen soll", "S_PROF_TEST_PH": "Satz zum Ausprobieren eingeben",
		"S_PF_NEW": "Neuer Prompt", "S_PF_EDIT": "Prompt bearbeiten",
		"S_POST_NO_MODEL": "an, aber kein Modell gewählt", "S_POST_NO_API": "an, aber kein Server eingerichtet", "S_POST_BAD": "Server antwortete nicht: %s", "S_POST_NO_PROMPT": "an, aber kein Prompt ausgewählt", "S_API_TEST": "Verbindung testen", "S_API_TEST_RUN": "Prüfe…", "S_API_TEST_OK": "Der Server hat geantwortet", "S_API_CLEAR": "Löschen", "S_API_CLEAR_ASK": "Adresse, Modell und Schlüssel des externen Servers löschen? Die Nachbearbeitung kehrt zum lokalen Modell zurück.", "S_RAM_AVAIL": "Verfügbarer Speicher: %s GB von %s GB", "S_RAM_OF": "%s GB von %s GB",
		"S_REPL_ADD": "Ersetzung hinzufügen", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "ganze Wörter", "S_REPL_CASE": "Groß-/Kleinschreibung", "S_REPL_EMPTY": "Noch keine Ersetzungen",
		"S_PASTE_DELAY": "Verzögerung vor dem Einfügen", "S_PASTE_DELAY_SUB": "wenn das Programm den Text noch nicht annimmt",
		"S_OVPOS": "Wo die Leiste erscheint", "S_OVPOS_SUB": "am Cursor — neben der Eingabestelle; zeigt die App sie nicht, dann neben dem Mauszeiger",
		"S_OVPOS_CARET": "Am Cursor",
		"S_OVTEXT": "Erkannten Text anzeigen", "S_OVTEXT_SUB": "auf der Leiste nach dem Einfügen, statt der Zeichenzahl",
		"S_OVERLAY": "Leiste anzeigen", "S_OVERLAY_SUB": "während des Diktats zeigt der Bildschirm, dass aufgenommen wird", "S_TYPEMODE": "Zeichenweise Eingabe (für Felder ohne Einfügen)",
		"S_RECLANG": "Sprache der Rede", "S_RECAUTO": "Automatisch",
		"S_LANG_KEY": "Eigenes Kürzel", "S_LANG_KEY_SET": "Zuweisen…", "S_LANG_KEY_CLEAR": "Kürzel entfernen", "S_LANG_KEY_HINT": "Halten Sie es gedrückt, um in dieser Sprache mit ihrem Modell zu diktieren. Das Hauptkürzel verwendet weiterhin die oben gewählte Sprache der Rede.",
		"S_DL": "Laden", "S_DEL": "Löschen",
		"S_M_BASE": "schnell, für schwache PCs", "S_M_SMALL": "ausgewogen", "S_M_MED": "Whisper für jede Sprache; das Modell, das übersetzt", "S_M_TURBO": "beste Genauigkeit auf CPU", "S_M_GIGAAM": "genauer auf Russisch, setzt Satzzeichen selbst", "S_M_PARAKEET": "25 europäische Sprachen, setzt Satzzeichen selbst",
		"S_THREADS":  "CPU-Threads", "S_MINMS": "Min. Aufnahme, ms", "S_MAXSEC": "Max. Aufnahme, s",
		"S_AUTOSTART": "whisper-server automatisch starten", "S_PORT": "Port", "S_SERVEREXE": "Pfad zu whisper-server", "S_SERVEREXE_SUB": "wird automatisch ausgefüllt; nur ändern, wenn der Server woanders liegt", "S_EXE_RESET": "Zurücksetzen", "S_EXE_WARN": "Die App findet whisper-server neben sich. Mit einem handgeschriebenen Pfad startet die Erkennung nach dem Verschieben des Ordners nicht mehr. Ändern?", "S_RESET_ALL": "Einstellungen zurücksetzen", "S_RESET_ALL_SUB": "alles außer Modellen und Verlauf auf Werkszustand", "S_RESET_ALL_BTN": "Zurücksetzen", "S_RESET_ALL_ASK": "Alle Einstellungen auf Werkszustand? Modelle, Verlauf und Prompts bleiben erhalten.", "S_RELOAD_CFG": "config.json neu einlesen", "S_RELOAD_CFG_SUB": "falls Sie die Datei von Hand bearbeitet haben", "S_RELOAD_CFG_BTN": "Neu einlesen", "S_UPD_FOUND": "Version %s ist da", "S_THEME": "Farbe", "S_THEME_SUB": "Farbe von Fenster, Leiste und Symbol im Infobereich", "S_THEME_GREEN": "Grün", "S_THEME_AMBER": "Bernstein", "S_THEME_BLUE": "Blau", "S_THEME_PINK": "Rosa", "S_THEME_EDITOR": "Editor", "S_THEME_NEON": "Neon", "S_WND_MAX": "Bildschirm füllen", "S_WND_RESTORE": "Vorherige Größe", "S_WND_MIN": "In den Infobereich", "S_WND_CLOSE": "Fenster schließen", "S_SKIN": "Design", "S_SKIN_SUB": "Schrift, Form, Effekte und Bewegung", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Sanft", "S_SKIN_PAPER": "Dokument", "S_SKIN_FLUENT": "Fluent", "S_SKIN_STUDIO": "Studio",
		"S_SERVERURL": "Externer Server (URL)", "S_URLHINT": "Falls gesetzt, wird kein lokaler Server gestartet",
		"S_STT_SRV": "Erkennungsserver",
		"S_STT_SRV_HINT": "Whisper-Modelle laufen in einem eigenen Programm. Es kann auf diesem oder auf einem anderen Rechner laufen — wählen Sie, welches genutzt wird.",
		"S_SRV_LOCAL": "Auf diesem Rechner",
		"S_SRV_REMOTE": "Auf einem anderen Rechner",
		"S_SRV_REMOTE_HINT": "Derselbe whisper-server, nur woanders gestartet: Heimserver, Rechner mit Grafikkarte, der Rechner nebenan.",
		"S_SRV_K_AUTO": "Autostart",
		"S_SRV_K_FILE": "Datei",
		"S_SRV_K_ADDR": "Adresse",
		"S_SRV_K_CHECK": "letzte Prüfung",
		"S_SRV_NEAR": "whisper-server.exe neben der App",
		"S_SRV_NOADDR": "nicht gesetzt",
		"S_SRV_NOCHECK": "nie geprüft",
		"S_SRV_LOCAL_DLG": "Lokaler Erkennungsserver",
		"S_SRV_ADDR": "Serveradresse",
		"S_SRV_ADDR_SUB": "Adresse des whisper-server auf dem anderen Rechner, mit Port",
		"S_SRV_ON": "ein",
		"S_SRV_OFF": "aus",
		"S_SRV_K_THREADS": "CPU-Threads",
		"S_SRV_K_PORT": "Port",
		"S_SRV_DOWN": "Erkennung nicht verfügbar",
		"S_SRV_DOWN_WHY": "der externe Erkennungsserver ist nicht eingerichtet — Adresse in den Einstellungen setzen",
		"S_SRV_DOWN_GO": "Servereinstellungen öffnen",
		"S_SRV_WARN_NOW": "Das Diktat funktioniert gerade nicht: Der externe Server ist gewählt, seine Adresse fehlt.",
		"S_SRV_WARN_LATER": "Sobald ein Whisper-Modell gewählt wird, funktioniert die Erkennung nicht: Die Adresse des externen Servers fehlt.",
		"S_SAVED":      "Gespeichert",
		"S_ABOUT_HTML": "<p><b>Stimme → Text an der Cursorposition.</b></p><p>Cursor in ein Eingabefeld setzen, Shortcut halten, sprechen, loslassen — der Text wird eingefügt.</p><p>Vollständig lokal und offline. Technik: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; Modelle von Hugging Face.</p><p>Logs überschreiten nie ~2 MB.</p>",
		  "S_SEARCH": "Einstellung finden…",
		"S_GRP_GENERAL": "Allgemein", "S_GRP_SPEECH": "Sprachverarbeitung", "S_GRP_INFO": "Infos", "S_NAV_POST": "Nachbearbeitung", "S_NAV_HELP": "Hilfe", "S_NAV_CONTACTS": "Kontakte", "S_HIST_ADD": "Hinzufügen", "S_CONTACT_MAIL": "E-Mail", "S_DICT_MODEL": "Erkennungsmodell", "S_LIB_ACC": "Genauigkeit", "S_LIB_SPD": "Tempo",
		"S_HELP_TOC": "Auf dieser Seite",
		"S_HELP_TOC_SHOW": "Inhalt anzeigen — das Fenster wird breiter",
		"S_HELP_TOC_HIDE": "Inhalt ausblenden und Fensterbreite zurücksetzen",
		"S_CONTACT_TITLE": "Kontakt aufnehmen",
		"S_ABOUT_DEPS": "Externe Module",
		"S_ABOUT_DEPS_HINT": "Fremder Code im Programm und seine Lizenzen. Ein Klick auf den Namen öffnet die Projektseite.",
		"S_DEP_WHISPER": "führt Whisper-Modelle aus",
		"S_DEP_LLAMA": "Nachbearbeitung, GGUF-Modelle, Übersetzung der Untertitel auf der Grafikkarte (llama-server-gpu.exe)",
		"S_DEP_VULKAN": "Header und Importbibliothek der Engine für die Grafikkarte; die Vulkan-Laufzeit selbst kommt mit dem Grafiktreiber",
		"S_DEP_SHERPA": "Engine für GigaAM, Parakeet, Canary, Qwen, Moonshine",
		"S_DEP_GGML": "Tensor-Bibliothek in whisper.cpp und llama.cpp",
		"S_DEP_ONNX": "führt die Modelle in sherpa-onnx aus",
		"S_DEP_WEBVIEW": "Einstellungsfenster auf WebView2",
		"S_DEP_WV2RT": "Windows-Komponente, die dieses Fenster zeichnet",
		"S_DEP_MALGO": "Mikrofonaufnahme",
		"S_DEP_MINIAUDIO": "Audioschicht in malgo",
		"S_DEP_WS": "Verbindung zum sherpa-server",
		"S_DEP_XSYS": "WinAPI-Aufrufe aus Go",
		"S_DEP_WINLOADER": "DLL-Laden in go-webview2",
		"S_DEP_PLEX": "Schrift der Oberfläche",
		"S_DEP_HF": "Katalog, aus dem Modelle geladen werden",
		"S_CONTACT_HINT": "Fehler, Idee, Frage zu einer Einstellung — schreiben Sie eine E-Mail, wenn es persönlich ist, oder öffnen Sie ein Issue, wenn es ein Fehler ist.",
		"S_CONTACT_REPO": "Repository",
		"S_CONTACT_ISSUES": "Fehler und Ideen",
		"S_CONTACT_WRITE": "E-Mail schreiben",
		"S_CONTACT_OPEN": "Öffnen",
		"S_STATE_ACTIVE": "Erkennt",
		"S_STATE_USED": "Eingesetzte Modelle",
		"S_STATE_INST": "Lokal installiert",
		"S_STATE_INST_SUB": "Modelle auf der Festplatte, bereit zur Zuweisung",
		"S_PRESETS": "Welches Modell für welche Sprache",
		"S_PRESETS_HINT": "Das Modell der Zeile «Automatisch» bedient alle Sprachen ohne eigenes Modell — sie stehen grau in der Liste. Klicken Sie eine Sprache an, um ihr ein eigenes zu geben. Dasselbe Modell erkennt auch, welche Sprache gesprochen wird, wenn keine Ausgangssprache gesetzt ist.",
		"S_MFOLDER": "Eigenes Modell",
		"S_DICT_SAVE": "Speichern",
		"S_OWNM_SUB": "Fügen Sie ein lokales Spracherkennungsmodell hinzu",
		"S_OWNM_ONEFILE": "Eine Datei",
		"S_OWNM_FOLDERF": "Ordner mit den Modelldateien",
		"S_OWNM_S1": "Öffnen Sie den Modellordner",
		"S_OWNM_S1S": "Zielordner:",
		"S_OWNM_S2": "Kopieren Sie das Modell",
		"S_OWNM_S2S": "Wählen Sie eine der unterstützten Strukturen",
		"S_OWNM_S3": "Starten Sie die Anwendung neu",
		"S_OWNM_S3S": "Das Modell erscheint für die Sprachen, die es unterstützt",
		"S_AS_AUTO": "wie automatische Erkennung",
		"S_REC_CHIP": "empfohlen",
		"S_BACK_AUTO": "Zurück zur automatischen Erkennung",
		"S_LANGS_COUNT": "Sprachen: %d",
		"S_LANGS_UNKNOWN": "Sprachen: unbekannt",
		"S_TR_EN": "übersetzt in %d Sprachen",
		"S_TR_LIST": "übersetzt: %s",
		"S_DL_GOING": "wird geladen:",
		"S_OPEN_FOLDER": "Ordner öffnen",
		"S_UNLOAD": "Aus dem Speicher entladen",
		"S_UNLOAD_SUB": "gibt den Speicher frei; das nächste Diktat lädt das Modell erneut",
		"S_UNLOAD_GO": "Entladen",
		"S_UNLOADED": "Entladen",
		"S_NOT_FOR_LANG": "%s erkennt diese Sprache nicht",
		"S_MANUAL_NOTE": "Kann nicht aus der App geladen werden — die Lizenz verbietet die Weitergabe. Laden Sie das Archiv selbst und entpacken Sie es nach models/moonshine-uk.",
		"S_MANUAL_LINK": "Selbst laden",
		"S_HF_FIT": "nur was auf diesen Rechner passt",
		"S_HF_HIDDEN": "ausgeblendet: %s",
		"S_WIZ_SKIP_DL": "Später laden",
		"S_WIZ_SKIP_NOTE": "Ohne Modell funktioniert das Diktieren nicht. Laden können Sie es unter „Sprachen & Modelle“.",
		"S_M_GIGAAM2": "die Vorgängergeneration des russischen Modells: gleiches Tempo, aber ohne Satzzeichen",
		"S_M_MOONUK": "ukrainisches Moonshine-Modell: schnell und leicht, ohne Satzzeichen",
		"S_M_LOCAL": "im Ordner models gefunden; Eigenschaften unbekannt, darum keine Balken",
		"S_ALL_LANGS": "alle Sprachen",
		"S_OVPOS_SCHEME_SUB": "klicken Sie auf den Bildschirm — die Leiste landet dort",
		"S_OVDRAG": "Ziehen Sie sie, wohin Sie wollen",
		"S_OVMON": "Bildschirm",
		"S_OVMON_SUB": "auf welchem Monitor die Leiste erscheint",
		"S_OVMON_CURSOR": "Der Bildschirm mit dem Cursor",
		"S_M_NEMOTRON": "schreibt beim Sprechen: der Text erscheint auf der Leiste, während Sie reden; 40 Sprachen, setzt Satzzeichen selbst",
		"S_M_TINY": "das kleinste und schnellste, für sehr schwache Rechner; spürbar ungenauer",
		"S_STATE_LOADED": "Gerade im Speicher",
		"S_STATE_LOADED_SUB": "Modelle entladen sich nach Leerlauf selbst",
		"S_STATE_WEEK": "Diese Woche",
		"S_ST_SUMMARY": "Übersicht", "S_ST_OVERLAY": "Leiste auf dem Bildschirm", "S_ST_BEEP": "Tonsignal", "S_ST_AUTORUN": "Start mit Windows", "S_ST_POST": "Nachbearbeitung", "S_ST_LOCAL": "lokal", "S_ST_CHECKED": "geprüft %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "an", "S_ST_OFF_M": "aus", "S_ST_MIC_OK": "Signal ist in Ordnung", "S_ST_MIC_BAD": "das Mikrofon schweigt", "S_ST_CHECK": "Prüfen", "S_ST_RECOG": "erkannt von %s", "S_ST_VER": "Version %s", "S_ST_LATEST": "aktuell", "S_ST_OUTDATED": "veraltet", "S_ST_UPD_OK": "Sie haben die neueste Version", "S_ST_UPD_DL": "Update wird geladen…",
		"S_ST_QUICK": "Schnelleinstellungen",
		"S_ST_MODELS": "Modelle",
		"S_ST_USAGE": "Diese Woche",
		"S_ST_READY": "Bereit zum Diktieren",
		"S_ST_OFF": "Im Infobereich ausgeschaltet",
		"S_ST_OFF_SUB": "das Kürzel tut nichts, bis Sie es wieder einschalten",
		"S_ST_ENABLE": "Einschalten",
		"S_ST_GOTO": "Diese Einstellung auf ihrem Tab öffnen",
		"S_ST_HOTKEY_GO": "Kürzel ändern",
		"S_ST_UPD_LAST": "Version %s — die neueste",
		"S_ST_UPD_HAVE": "Version %s ist verfügbar",
		"S_ST_MEM": "%s GB frei von %s",
		"S_ST_MEM_SUB": "im Speicher: %s · auf der Platte: %d Modelle, %s GB",
		"S_ST_MEM_NONE": "nichts",
		"S_ST_LANG": "Sprache",
		"S_ST_ASR": "Erkennung",
		"S_ST_ON": "an",
		"S_ST_OFF_W": "aus",
		"S_ST_ON_F": "an",
		"S_ST_OFF_F": "aus",
		"S_ST_ACTIVE": "aktiv",
		"S_ST_IDLE": "nicht gestartet",
		"S_ST_DISK": "auf der Platte, %s",
		"S_ST_USAGE_SUB": "%d Zeichen · heute %d · im Schnitt %d Zeichen",
		"S_WEEK_OTHER": "andere",
		"S_ST_NO_WEEK": "diese Woche keine Diktate",
		"S_ST_AUTORUN_SUB": "die App startet nicht von selbst",
		"S_ST_OVERLAY_SUB": "beim Aufnehmen sichtbar",
		"S_REPL_LANG": "Sprache der Regel",
		"S_REPL_LANG_ALL": "alle Sprachen",
		"S_M_CANARY": "Englisch, Deutsch, Spanisch, Französisch — und übersetzt selbst zwischen ihnen",
		"S_M_QWEN3": "rund 30 Sprachen, setzt Satzzeichen selbst; das schwerste und genaueste im Katalog",
		"S_POSTAPI": "Externer Server",
		"S_POST_HINT": "Bearbeitet den erkannten Text nach Ihren Prompts: entfernt Füllwörter, richtet die Zeichensetzung, ändert den Stil. Aus — der Text wird so eingefügt, wie er erkannt wurde.",
		"S_POST_MODEL": "Modell",
		"S_SRC_LOCAL": "Lokal",
		"S_SRC_USED": "in Gebrauch",
		"S_HF_GO": "Suchen",
		"S_POSTAPI_HINT": "Standardmäßig leer — die Nachbearbeitung läuft lokal. Tragen Sie eine Adresse ein, und die Prompts laufen auf einem externen Server: OpenAI, Groq, eigenes vLLM — alles mit kompatibler API.",
		"S_POSTAPI_URL": "Adresse",
		"S_POSTAPI_URL_SUB": "leer = das lokale Modell; Beispiel: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Modell",
		"S_POSTAPI_KEY": "API-Schlüssel",
		"S_POSTAPI_KEY_SET": "Schlüssel gespeichert (mit Windows DPAPI verschlüsselt)",
		"S_POSTAPI_KEY_NONE": "kein Schlüssel",
		"S_POSTAPI_SAVE": "Schlüssel speichern",
		"S_POSTAPI_TIMEOUT": "Antwortzeit", "S_SEC_SHORT": "s",
		"S_POSTAPI_WARN": "⚠ Der erkannte Text Ihrer Diktate geht an diese Adresse. Ton verlässt den Rechner nie. Der Schlüssel liegt verschlüsselt.",
		"S_POSTAPI_ASK": "Erkannten Text an %s senden? Der Ton bleibt auf diesem Rechner, der Text verlässt ihn aber.",
		"S_POSTAPI_BADGE": "externer Server",
		"S_NOT_INSTALLED": "nicht installiert",
		"S_NAV_STATE": "Status", "S_NAV_DICT": "Steuerung & Verhalten", "S_NAV_MIC": "Mikrofon", "S_NAV_MODELS": "Sprachen & Modelle",
		"S_NAV_TEXT": "Regeln", "S_NAV_TR": "Übersetzung", "S_NAV_SYSTEM": "System", "S_NAV_ABOUT": "Über",
		"S_STATE_HINT": "halten und sprechen — der Text landet dort, wo der Cursor steht",
		"S_STATE_PROC": "Nachbearbeitung",
		"S_CHANGE_MODEL": "Wechseln", "S_PICK_MODEL": "Auswählen", "S_STATE_GET": "Laden",
		"S_RETRY": "Erneut versuchen", "S_BERR_OPEN": "Servereinstellungen öffnen",
		"S_STATE_LAST": "Letztes Diktat", "S_STATE_COPY": "Kopieren", "S_STATE_MEM": "Speicher",
		"S_STATE_MEM_SUB": "Modelle bleiben geladen, der erste Satz kommt ohne Verzögerung",
		"S_HOTMODE":       "Modus", "S_HOTMODE_HOLD": "halten", "S_HOTMODE_TOGGLE": "umschalten",
		"S_SUB_HOTMODE": "Tasten halten oder einmal drücken zum Starten und noch einmal zum Beenden",
		"S_SILENCE":     "Bei Stille beenden", "S_SILENCE_MS": "Pause vor dem Beenden, ms",
		"S_SUB_SILENCE": "Modus „umschalten“: sobald Sie etwas gesagt haben, beendet eine Pause die Aufnahme, als hätten Sie die Tasten erneut gedrückt",
		"S_SUB_SILENCE_MS": "wie lange Sie schweigen dürfen; vor Ihren ersten Worten endet die Aufnahme nie",
		"S_SUB_MINMS":   "ignoriert versehentliche Tastendrücke",
		"S_SUB_ENTER":   "schickt die Nachricht sofort ab",
		"S_SUB_CLIP":    "Bilder und Dateien kommen unverändert zurück",
		"S_SUB_TYPE":    "hilft dort, wo ein Feld das Einfügen verweigert",
		"S_SEC_OVERLAY": "Bildschirmanzeige",
		"S_MIC_CHECK":   "Mikrofon prüfen", "S_MIC_CHECK_SUB": "drei Sekunden Aufnahme und ein Urteil: Pegel, Übersteuerung, ob Sprache da ist", "S_MIC_CHECKING": "Prüfe…",
		"S_MCHECK": "Installierte Modelle prüfen", "S_MCHECK_SUB": "vergleicht die Modelldateien mit den Referenz-Hashes", "S_MCHECK_GO": "Prüfen", "S_MCHECK_RUN": "Prüfe…",
		"S_HIST_INSERT": "Einfügen",
		"S_MIC": "Mikrofon", "S_MIC_DEFAULT": "Systemstandard", "S_MIC_REFRESH": "Liste aktualisieren",
		"S_MIC_LEVEL": "Eingangspegel", "S_MIC_QUIET": "still",
		"S_SUB_THREADS": "mehr Threads sind nicht immer schneller — messen Sie es auf Ihrem Rechner",
		"S_SEC_LLM":     "Bearbeitungsmodell",
		"S_PUNCT":       "Satzzeichen und Großschreibung", "S_SUB_PUNCT": "woher Satzzeichen und Großschreibung kommen",
		"S_PUNCT_MODEL": "vom Modell", "S_PUNCT_LLM": "vom Bearbeitungsmodell", "S_PUNCT_OFF": "entfernen",
		"S_SUB_DICT": "Wörterbuch", "S_SUB_PROMPTS": "Prompts",
		"S_SUB_TRTARGET": "in sie wird der Text übersetzt; der Dialog auf der Leiste bietet sie zuerst an",
		"S_REMOTE_ABOUT": "Ein entfernter Server ist eingestellt: Audio wird dorthin gesendet, und das Versprechen oben gilt so lange nicht.",
		"S_UPD":          "Aktualisierungen", "S_UPD_CHECK": "Nach Updates suchen", "S_UPD_AUTO": "Beim Start prüfen",
		"S_SUB_UPD":  "die einzige Netzwerkanfrage außer Modell-Downloads",
		"S_UPD_NONE": "Sie haben die neueste Version", "S_BADGE_MODELS": "Installierte Modelle", "S_BADGE_MISS": "Ein Modell ist nicht geladen", "S_BADGE_SYSTEM": "Warnungen brauchen Aufmerksamkeit", "S_BADGE_HIST": "Einträge im Verlauf", "S_LOG_OPEN": "Protokoll öffnen", "S_LOG": "Protokoll", "S_LOG_SUB": "alles, was die App über sich schreibt", "S_UPD_AVAIL": "Version %s ist verfügbar.",
		"S_UPD_GO": "Aktualisieren", "S_UPD_ERR": "Update-Prüfung fehlgeschlagen", "S_UPD_DL": "Update wird geladen…",
		"S_SEC_SERVICE": "Dienst", "S_SUB_AUTOSTART": "ausschalten, wenn Sie den Server selbst starten",
		"S_SUB_PORT":    "die Erkennung startet sich selbst neu",
		"S_MODEL_READY": "Modell geladen — wählen Sie es aus, um zu wechseln",
		"S_FIT_OK":      "passt", "S_FIT_WARN": "knapp", "S_FIT_BAD": "zu wenig RAM", "S_RAM": "Arbeitsspeicher:",
		"S_HF_PH":       "Modellname — z. B. qwen2.5 instruct",
		"S_NO_LLM":      "Noch keine Modelle installiert — im Suchfeld unten eines finden und laden.",
		"S_NO_LLM_PROF": "Prompts stehen zur Verfügung, sobald ein Modell installiert ist — der Block „Modell“ oben auf dieser Registerkarte.",
		"S_UPDATED":     "Modell zuletzt aktualisiert", "S_PROF_EDIT": "Bearbeiten", "S_PROF_CLOSE": "Einklappen",
		"S_CONFIRM_DEL": "Modell „%s“ löschen? Es kann erneut geladen werden.", "S_FREE": "frei",
		"S_DEL_ACTIVE":     "Das aktive Modell „%s“ löschen? Die Erkennung hält an, bis Sie ein anderes wählen — herunterladen können Sie es gleich hier.",
		"S_REMOTE_WARN":    "Audio wird an diesen Server gesendet. Der lokale Modus ist aus.",
		"S_REMOTE_ASK":     "Audio wird nicht mehr auf diesem Rechner verarbeitet, sondern an %s gesendet. Fernmodus einschalten?",
		"S_REMOTE_BADGE":   "EXTERN",
		"S_OK":             "Ja", "S_CANCEL": "Abbrechen", "S_DL_START": "Laden", "S_DL_CANCEL": "Download abbrechen",
		"S_DL_ASK":    "Das Modell „%s“ ist nicht geladen (%s). Jetzt laden?",
		"S_NOT_FOUND": "nichts",   "S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autor und Entwickler von {app} — einem lokalen Diktierwerkzeug für Windows: Stimme wird direkt am Cursor zu Text, ohne Cloud, ohne Abo.</p>" +
			"<p>Das Projekt ist offen: Quellcode, Build-Pipeline und aktuelle Releases liegen auf GitHub.</p>" +
			"<ul>" +
			"<li>Repository: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Autorenprofil: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Fehler gefunden oder eine Idee — eröffnen Sie ein Issue im Repository.</p>",
		"S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Was das ist</p>" +
			"<p>Sie sprechen, und die Wörter erscheinen dort, wo der Cursor blinkt. Als hätten Sie sie getippt, nur schneller.</p>" +
			"<p class=\"hdim\">So sieht das aus. Sie klicken in ein Feld, in das Sie sonst tippen: eine E-Mail, einen Chat, ein Suchfeld, ein Dokument. Sie halten zwei Tasten gedrückt und sagen einen Satz mit normaler Stimme. Sie lassen los. Eine Sekunde später steht der Satz im Feld. Nichts umschalten, nichts kopieren: Die App arbeitet über jedem Programm und wartet im Hintergrund auf ihre Tasten.</p>" +
			"<p class=\"hdim\">Die ganze Arbeit erledigt Ihr Rechner. Ihre Stimme wird nicht gespeichert: Sie bleibt ein paar Sekunden im Arbeitsspeicher, wird zu Text und ist dann weg. Ins Internet geht die App nur aus zwei Gründen: wenn Sie ein Erkennungsmodell herunterladen und wenn sie nachsieht, ob eine neue Version erschienen ist. Diese Prüfung lässt sich abschalten, dann braucht die App überhaupt kein Netz.</p>" +
			"<p class=\"hdim\">Das Wort „Modell“ kommt unten oft vor. Es ist einfach eine Datei, die gesprochene Sprache entschlüsseln kann. Sie wird einmal heruntergeladen und liegt in einem Ordner neben der App; je größer sie ist, desto genauer hört sie und desto langsamer arbeitet sie.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Was passiert, wenn Sie die Tasten loslassen</p>" +
			"<div class=\"hflow\">" +
			"<span>Aufnahme</span><em>→</em><span>Erkennung</span><em>→</em><span>Ersetzungen</span><em>→</em><span>Befehle</span><em>→</em><span>Übersetzung</span><em>→</em><span>Bearbeitung</span><em>→</em><span>Einfügen</span>" +
			"</div>" +
			"<p>Meist dauert das ein bis zwei Sekunden, und die einzelnen Schritte muss man nicht kennen. Kommt der Text aber einmal anders heraus als erwartet, sieht man an dieser Kette leicht, welcher Schritt ihn verändert hat.</p>" +
			"<ul>" +
			"<li><b>Aufnahme.</b> Solange die Tasten gedrückt sind, hört die App auf das Mikrofon. Einen zu kurzen Druck überhört sie, damit eine versehentlich berührte Taste kein Diktat startet.</li>" +
			"<li><b>Erkennung.</b> Die Stimme wird zu Buchstaben. Das übernimmt ein Modell. Russisch erledigt meist GigaAM, Polnisch Parakeet, die meisten anderen Sprachen Nemotron, und Whisper übernimmt die übrigen Sprachen und die Übersetzung; das sind nur Namen, von Hand auswählen müssen Sie sie nicht.</li>" +
			"<li><b>Ersetzungen.</b> Ihre persönliche Liste „gehört → gemeint“. Schreibt das Modell hartnäckig „git hub“ statt „GitHub“, korrigieren Sie das einmal, und danach wird es von selbst berichtigt.</li>" +
			"<li><b>Befehle.</b> Ein paar Wendungen wirken wie Schaltflächen. Sagen Sie „neue Zeile“, springt der Text in eine neue Zeile, und die Worte selbst landen nicht im Feld.</li>" +
			"<li><b>Übersetzung.</b> Wenn Sie sie eingeschaltet haben, sprechen Sie in einer Sprache, und der Text erscheint in einer anderen.</li>" +
			"<li><b>Bearbeitung.</b> Ein optionaler Schritt. Ein zweites Modell bringt den Text in Form: Es entfernt „ähm“ und „also quasi“, setzt die Kommas und gibt der Formulierung einen geschäftlichen Ton.</li>" +
			"<li><b>Einfügen.</b> Der Text landet im Feld. Was vor dem Diktat in Ihrer Zwischenablage lag, kommt wieder an seinen Platz.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Während all das läuft, ist unten auf dem Bildschirm eine schmale Leiste mit dem Namen des aktuellen Schritts zu sehen. An ihr erkennen Sie sofort: Die App arbeitet, sie hängt nicht.</p>" +
			"<p class=\"wh\" id=\"s-first\">Der erste Start</p>" +
			"<p>Beim ersten Start führt die App Sie durch drei Schritte. Das dauert eine Minute, danach können Sie diktieren.</p>" +
			"<ul>" +
			"<li><b>In welcher Sprache Sie diktieren werden.</b> Die wichtigste Frage: Passend zur Antwort wählt die App selbst ein geeignetes Modell und lädt es im Hintergrund herunter, der Fortschritt ist zu sehen. Weicht diese Sprache von der Sprache des Fensters ab, erscheint gleich darunter eine Zeile, in der Sie auch die Sprache der Oberfläche wechseln können.</li>" +
			"<li><b>Ausprobieren.</b> Das Tastenkürzel (standardmäßig Strg und die Windows-Taste; zum Ändern drücken Sie, was sich bequem mit einer Hand halten lässt), das Mikrofon mit einem live mitlaufenden Pegelbalken und ein Feld, in das Sie einen Probesatz diktieren. Bewegt sich der Balken beim Sprechen nicht, ist das falsche Mikrofon gewählt. Das Feld wird lebendig, sobald das Modell bereitsteht.</li>" +
			"<li><b>Fertig.</b> Der Schalter „Mit Windows starten“, und der Assistent ist abgeschlossen.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Jeden Schritt können Sie überspringen und später nachholen: All das steht auch in den Einstellungen. Das Modell lädt im Hintergrund weiter, auch wenn Sie schon weitergehen oder den Assistenten beenden, bevor es fertig ist. Wollen Sie den Assistenten noch einmal durchlaufen, starten Sie die App mit dem Zusatz <code>-wizard</code>. Nach einem Update erscheint er nicht mehr.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">Die Leiste unten am Bildschirm</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Sprechen…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">so zeigt die App, dass sie Sie hört</div>" +
			"</div>" +
			"<p>Während des Diktats ist sie das Einzige, was auf dem Bildschirm zu sehen ist. Die Leiste ist klein, fängt die Maus nicht ab und stört nicht bei der Arbeit: Sie zeigt einfach, was gerade mit Ihrem Satz passiert.</p>" +
			"<ul>" +
			"<li><b>„Sprechen…“</b> Die Aufnahme läuft. Die Balken daneben zucken im Takt Ihrer Stimme. Liegen sie still, sind Sie nicht zu hören, und weiterzusprechen hat keinen Sinn.</li>" +
			"<li><b>„Erkenne…“</b> Sie haben die Tasten losgelassen, das Modell arbeitet. Meist ein bis zwei Sekunden. Das allererste Diktat nach dem Start dauert länger: In diesem Moment wird das Modell in den Speicher geladen.</li>" +
			"<li><b>„Übersetze…“</b> Die Übersetzung läuft.</li>" +
			"<li><b>„Bearbeite: Cleanup (1/2)“</b> Das zweite Modell bearbeitet den Text. Sie sehen, welche Anweisung gerade ausgeführt wird und wie viele es insgesamt sind.</li>" +
			"<li><b>„Eingefügt: 128 Zeichen“</b> Fertig, der Text steht schon im Feld. Die Leiste verschwindet nach ein paar Sekunden von selbst.</li>" +
			"<li><b>Das Kreuz.</b> Anders überlegt? Klicken Sie es in jedem Schritt an. Die Aufnahme wird verworfen, ins Feld gelangt nichts.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Die Leiste lässt sich ausschalten, in jede Ecke stellen, an den Cursor binden oder einfach mit der Maus dorthin ziehen, wo sie nicht stört. Das finden Sie auf der Registerkarte „Steuerung & Verhalten“.</p>" +
			"<p class=\"wh\" id=\"s-state\">Der Bildschirm „Status“</p>" +
			"<p>Dieser Bildschirm beantwortet eine Frage: Können Sie gerade diktieren? Wenn nicht, sagt er, warum, und bietet eine Schaltfläche, die dorthin führt, wo sich das beheben lässt.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Bereit zum Diktieren</b><span>halten und sprechen — der Text landet dort, wo der Cursor steht</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">das grüne Quadrat bedeutet, dass die App auf ihre Tasten hört</div>" +
			"</div>" +
			"<p>Die obere Karte gibt es in drei Varianten. <b>Grün</b>: Alles ist in Ordnung, daneben steht eine Erinnerung an die Tasten. <b>Gelb</b>: Sie haben die App über das Symbol im Infobereich neben der Uhr ausgeschaltet, und die Tasten tun nichts, bis Sie direkt auf der Karte „Einschalten“ drücken. <b>Rot</b>: Die Einstellungen widersprechen sich. Die Erkennung auf einem anderen Rechner ist gewählt, aber keine Adresse angegeben.</p>" +
			"<p>Darunter steht, worauf man meist nur kurz schaut:</p>" +
			"<ul>" +
			"<li><b>Untertitel.</b> Die Karte direkt unter der obersten: ob die Untertitel aus sind, ihr Modell laden, laufen oder neu starten, mit dem Modell und wie der Text kommt, ihr Tastenkürzel, wenn eines gesetzt ist, eine Schaltfläche, die sie startet oder beendet, und ein Pfeil zur Registerkarte „Untertitel“.</li>" +
			"<li><b>Diese Woche.</b> Ein Kreisdiagramm: in welche Programme Ihre Diktate gingen und wie viele Zeichen Sie diktiert haben.</li>" +
			"<li><b>Version.</b> Die Build-Nummer und eine Schaltfläche zum Prüfen. Ist eine neuere erschienen, sagt die App es gleich hier.</li>" +
			"<li><b>Speicher.</b> Wie viel davon frei ist, welches Modell gerade darin sitzt und wie viel alle Modelle auf der Festplatte belegen. Die Schaltfläche „Entladen“ gibt den Speicher sofort frei, etwa vor einem anspruchsvollen Spiel.</li>" +
			"<li><b>Schnelleinstellungen.</b> Fünf Karten mit dem Wichtigsten: Mikrofon, Sprache, Leiste, Ton, Start mit Windows. Sie zeigen, wie es gerade eingestellt ist, und der Pfeil führt zur Registerkarte, auf der man es ändert. Schalter gibt es hier mit Absicht nicht, damit nichts im Vorbeigehen verstellt wird.</li>" +
			"<li><b>Modelle.</b> Was installiert ist, was gerade arbeitet, für welche Sprachen es zuständig ist und was nur auf der Festplatte liegt.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">Die Registerkarte „System“</p>" +
			"<p>Das Aussehen der App, ihr Zusammenspiel mit Windows und ein paar Schaltflächen für den Dienst. Im Alltag kommt man hierher etwa einmal im Jahr.</p>" +
			"<p class=\"hsub\">Wie die App aussieht</p>" +
			"<p>Drei Zeilen: Sprache der Oberfläche, Design und Farbe. Auf die Arbeit der App haben sie keinerlei Einfluss.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Sprache der Oberfläche</td><td>Die Sprache des Einstellungsfensters und der Beschriftungen auf der Leiste. Mit der Sprache des Diktats hat sie nichts zu tun: Das Fenster kann auf Deutsch sein, während Sie auf Englisch diktieren.</td></tr>" +
			"<tr><td>Design</td><td>Das Gesamtbild. Es gibt sieben: „Editor“ (Standard: Dunkelgrau mit Blau), „Fluent“ (wie die Einstellungen von Windows 11), „Studio“ (Graphit mit Orange, fast ohne Rundungen), „Sanft“ (helles Lila, abgerundet), „Dokument“ (Papier mit Blau), „Neon“ (Violett mit einem Verlauf von Rosa zu Türkis) und „Terminal“ (grüne Buchstaben auf Schwarz, Klammern um die Schaltflächen). Alle sieben sind in einer Schrift gesetzt, IBM Plex Mono.</td></tr>" +
			"<tr><td>Farbe</td><td>Ein Farbton nur für „Terminal“: Grün, Bernstein, Blau, Rosa. In derselben Farbe erscheinen auch die Leiste und das Symbol im Infobereich. Die anderen Designs haben ihre eigenen Farben, und die Zeile wird ausgeblendet.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Start und Updates</p>" +
			"<p>Hier legen Sie fest, ob die App mit dem System startet und ob sie nach neuen Versionen Ausschau hält.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Mit Windows starten</td><td>Die App startet bei der Anmeldung und wartet auf Ihre Tasten. Das ist ein gewöhnlicher Eintrag im Autostart, genau so einer, wie Sie ihn selbst dort anlegen könnten.</td></tr>" +
			"<tr><td>Beim Start prüfen</td><td>Einmal beim Start fragt die App nach, ob es eine neuere Version gibt. Sonst geht sie nicht ins Internet, abgesehen vom Herunterladen der Modelle. Schalten Sie das aus, dann müssen Sie sich selbst über neue Versionen informieren.</td></tr>" +
			"<tr><td>Aktualisierungen</td><td>Die Schaltfläche prüft sofort. Gibt es eine neuere Version, lädt die App sie herunter und installiert sie über die alte: Einstellungen, Modelle und Verlauf bleiben, wo sie sind.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Dienst</p>" +
			"<p>Drei Schaltflächen für den Fall, dass etwas schiefgegangen ist.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Protokoll</td><td>Öffnet die Datei, in die die App schreibt, was sie getan hat: welches Modell gestartet ist, wie lange die Erkennung gedauert hat, wo ein Fehler aufgetreten ist. Der Grund steht meist in Worten da, nicht als Code. Die Datei wächst nicht immer weiter: Altes wird durch Neues überschrieben.</td></tr>" +
			"<tr><td>config.json neu einlesen</td><td>Alle Einstellungen liegen in einer Textdatei neben der App. Haben Sie sie von Hand bearbeitet oder von einem anderen Rechner mitgebracht, übernimmt diese Schaltfläche die Änderungen ohne Neustart.</td></tr>" +
			"<tr><td>Einstellungen zurücksetzen</td><td>Stellt alles so her wie nach der Installation. Heruntergeladene Modelle und der Verlauf bleiben, wo sie sind: Gelöscht werden nur die Einstellungen.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Wo die Erkennung läuft</p>" +
			"<p>Die Sprache verarbeitet ein eigenes kleines Programm. Normalerweise liegt es neben der App und startet von selbst, und Sie müssen nicht daran denken. Es lässt sich aber auch auf einem anderen Rechner betreiben: zum Beispiel auf dem Rechner zu Hause mit einer starken Grafikkarte, während Sie auf einem leichten Laptop diktieren. Gewählt wird mit dem Schalter auf den Karten.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>Auf diesem Rechner</b><span>CPU-Threads 4 · Autostart ein · Port 8910 · whisper-server.exe neben der App</span></div><button class=\"hbtn\">Bearbeiten</button></div>" +
			"<div class=\"hcap\">die gewählte Karte ist hervorgehoben, die andere ist so lange gedimmt</div>" +
			"</div>" +
			"<p>Die Schaltfläche „Bearbeiten“ öffnet ein Fenster mit vier Einstellungen. Alle haben Standardwerte, und fast immer können Sie sie in Ruhe lassen.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>CPU-Threads</td><td>Wie viele Prozessorkerne die Erkennung bekommt. Mehr Kerne heißt schneller, aber nur bis zu einer Grenze: Darüber bringt es nichts mehr, und andere Programme fangen an zu ruckeln.</td></tr>" +
			"<tr><td>whisper-server automatisch starten</td><td>Die App startet den Erkenner selbst. Ausschalten lohnt sich nur in einem Fall: wenn Sie ihn von Hand starten und nicht wollen, dass die App sich einmischt.</td></tr>" +
			"<tr><td>Port</td><td>Die Nummer, über die die App innerhalb Ihres eigenen Rechners mit dem Erkenner spricht. Von außen ist sie nicht sichtbar. Ändern müssen Sie sie nur, wenn ein anderes Programm diese Nummer belegt hat.</td></tr>" +
			"<tr><td>Pfad zu whisper-server</td><td>Wo der Erkenner liegt. Wird automatisch ausgefüllt. Einen eigenen Pfad einzutragen lohnt sich nur, wenn Sie ihn bewusst woandershin gelegt haben: Sonst startet das Diktat nicht mehr, sobald Sie den Ordner verschieben.</td></tr>" +
			"</table>" +
			"<p>Die zweite Karte hat nur eine Einstellung: die Adresse des Rechners im Netzwerk, zum Beispiel <code>http://192.168.1.50:8080</code>. Bedenken Sie: In dieser Variante geht die aufgenommene Stimme an diesen Rechner, und über eine gewöhnliche Adresse wird sie unverschlüsselt übertragen. Im Heimnetz ist das in Ordnung, über das Internet arbeitet man so besser nicht.</p>" +
			"<div class=\"hbad\">Ist „Auf einem anderen Rechner“ gewählt, aber keine Adresse eingetragen, und wird Ihre Sprache von Whisper bedient, funktioniert das Diktat nicht. Die App meldet das mit einer roten Karte auf „Status“, einer roten Markierung neben „System“ und einer Zeile unten im Fenster.</div>" +
			"<p class=\"wh\" id=\"s-mic\">Die Registerkarte „Mikrofon“</p>" +
			"<p>Mehr als die Hälfte der missglückten Diktate liegt nicht am Modell und nicht an den Einstellungen, sondern am Mikrofon: Es wird das falsche abgehört, oder Sie sind kaum zu hören. Auf dieser Registerkarte finden Sie alles, um das in einer halben Minute zu prüfen.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Mikrofon</td><td>Welches Gerät abgehört wird. Die Variante „Systemstandard“ bedeutet „dasselbe, was in Windows gewählt ist“. Das ist praktisch, wenn Sie mal mit, mal ohne Headset arbeiten: Die App folgt dem System von selbst.</td></tr>" +
			"<tr><td>Liste aktualisieren</td><td>Nötig, wenn Sie das Headset erst nach dem Start der App angeschlossen haben: Die Liste wird neu eingelesen.</td></tr>" +
			"<tr><td>Eingangspegel</td><td>Ein Balken, der live mitläuft. Sagen Sie etwas, und die Balken sollten sich deutlich bewegen. Tun sie es nicht, ist das Mikrofon in Windows stummgeschaltet oder das falsche gewählt.</td></tr>" +
			"<tr><td>Mikrofon prüfen</td><td>Drei Sekunden Aufnahme und eine Auswertung in einfachen Worten: ob es laut genug ist, ob es übersteuert (wenn der Ton bei einem zu nahen Mikrofon „knarzt“) und ob überhaupt Sprache zu hören ist.</td></tr>" +
			"<tr><td>Tonsignale der Aufnahme</td><td>Ein kurzer Ton am Anfang und am Ende der Aufnahme. Besonders nützlich, wenn Sie die Leiste ausgeschaltet haben: Dann hören Sie wenigstens, dass die Aufnahme begonnen hat.</td></tr>" +
			"<tr><td>Signalton</td><td>Welcher Ton genau: Windows-Stimme, Glöckchen, Marimba und noch ein paar weitere. Die Schaltfläche daneben spielt den gewählten ab, damit Sie nicht anhand des Namens raten müssen.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">Die Registerkarte „Verlauf“</p>" +
			"<p>Der Verlauf ist eine Liste dessen, was Sie diktiert haben. Nur Text: Die Stimme wird nirgends und niemals gespeichert. Standardmäßig ist er ausgeschaltet, und das mit Absicht: In Diktate gerät oft Persönliches.</p>" +
			"<p>Einschalten lohnt sich, wenn Sie schon einmal in so einer Lage waren: Sie haben Text ins falsche Fenster eingefügt, und er war weg; Sie müssen einen langen Satz wiederholen; Sie möchten finden, was Sie gestern diktiert haben, und es noch einmal senden.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Wie lange aufbewahren</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">Tage</span></div>" +
			"<div class=\"hrow\"><label>Wie viele Einträge aufbewahren</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">die Frist lässt sich in Minuten, Stunden oder Tagen angeben, mit jeder beliebigen Zahl</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Verlauf der Diktate führen</td><td>Der Hauptschalter. Ausgeschaltet kommen keine neuen Einträge mehr hinzu, und die alten werden gelöscht.</td></tr>" +
			"<tr><td>Wie lange aufbewahren</td><td>Nach welcher Zeit ein Eintrag von selbst verschwindet. Sie können fünf Minuten einstellen, wenn Sie den Verlauf nur für den Fall „an der falschen Stelle eingefügt“ brauchen, oder auch dreißig Tage.</td></tr>" +
			"<tr><td>Wie viele Einträge aufbewahren</td><td>Eine zweite Grenze, nach Anzahl. Sie wirkt zusammen mit der ersten: Was zuerst erreicht ist, greift. Bei „7 Tage und höchstens 50“ verdrängt der einundfünfzigste Eintrag den ältesten, auch wenn die Woche noch nicht um ist.</td></tr>" +
			"<tr><td>Aus diesen Programmen nie aufzeichnen</td><td>Eine Liste von Ausnahmen. Was Sie in diese Fenster diktieren (Bank, Passwortmanager, Arbeits-Chat), landet nicht im Verlauf, das Diktat selbst funktioniert wie gewohnt. Ein Programm können Sie aus der Liste der geöffneten Fenster wählen, ohne nach seinem Dateinamen zu suchen.</td></tr>" +
			"</table>" +
			"<p>Mit jedem Eintrag können Sie vier Dinge tun. <b>Kopieren</b>: Der Text kommt in die Zwischenablage. <b>Einfügen</b>: Die App bittet Sie zu zeigen, wohin genau. Klicken Sie in einem beliebigen Programm auf das gewünschte Feld, und der Text landet dort. <b>Löschen</b>: Entfernt diesen einen Eintrag. <b>Suchen</b>: Die Suche geht sowohl durch den Text als auch durch den Programmnamen. Die Schaltfläche „Leeren“ löscht alles auf einmal und fragt deshalb nach einer Bestätigung.</p>" +
			"<p class=\"hdim\">Rechts an jedem Eintrag steht, wann er verschwindet, etwa „bis 7. Sept.“, und am letzten Tag wird diese Angabe gelb.</p>" +
			"<p class=\"wh\" id=\"s-dict\">Die Registerkarte „Steuerung & Verhalten“</p>" +
			"<p>Wie das Diktat startet und wie der Text ins Feld kommt. Hierher kommt man, wenn etwas nicht so läuft, wie man möchte: Die Tasten sind unbequem, der Text erscheint nicht, Enter wird von selbst gedrückt.</p>" +
			"<p class=\"hsub\">Tasten</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Tastenkürzel</td><td>Drücken Sie „Ändern…“ und halten Sie, was Ihnen bequem ist. Die App fängt diese Kombination im ganzen System ab, nehmen Sie also keine, die schon belegt ist, etwa Strg+C. Am seltensten belegt sind Kombinationen mit der Windows-Taste. Eine Sprache kann auch ein eigenes Kürzel haben; es wird auf der Registerkarte „Sprachen & Modelle“ vergeben.</td></tr>" +
			"<tr><td>Modus</td><td>Bei „halten“ wird aufgenommen, solange die Tasten gedrückt sind: Loslassen, und der Satz geht in die Verarbeitung. Bei „umschalten“ startet der erste Druck die Aufnahme, der zweite beendet sie. Der zweite Modus eignet sich für lange Monologe, wenn die Hand vom Halten der Tasten müde wird.</td></tr>" +
			"<tr><td>Bei Stille beenden</td><td>Nur im Modus „umschalten“. Sobald Sie etwas gesagt haben, beendet eine Pause von der Länge, die unter „Pause vor dem Beenden, ms“ eingestellt ist, die Aufnahme von selbst, als hätten Sie das Tastenkürzel erneut gedrückt, und die Leiste zeigt „Bei Stille beendet“. Vor Ihren ersten Worten endet die Aufnahme nie, Sie können also nach dem Drücken in Ruhe Ihre Gedanken sammeln; sagen Sie gar nichts, beendet sie trotzdem die Grenze „Max. Aufnahme, s“. Ist „Enter nach dem Einfügen drücken (Auto-Senden)“ eingeschaltet, schickt eine Pause die Nachricht ab: Wählen Sie dann eine längere Pause.</td></tr>" +
			"<tr><td>Min. Aufnahme, ms</td><td>Wie kurz ein Druck sein muss, um als versehentlich zu gelten. Berühren Sie die Tasten nur für einen Sekundenbruchteil, bleibt die App still, statt einen Wortfetzen einzufügen.</td></tr>" +
			"<tr><td>Max. Aufnahme, s</td><td>Eine Absicherung gegen eine vergessene Taste: Die Aufnahme bricht von selbst ab, statt eine halbe Stunde lang ins Leere aufzunehmen.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Wie der Text eingefügt wird</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Enter nach dem Einfügen drücken (Auto-Senden)</td><td>Der Text wird nicht nur eingefügt, sondern gleich abgeschickt: Die Nachricht geht raus, die Suche startet. Im Messenger spart das einen Handgriff, in einer E-Mail oder einem Dokument stört es fast sicher.</td></tr>" +
			"<tr><td>Zwischenablage nach Einfügen wiederherstellen</td><td>Das Einfügen läuft über die Zwischenablage, und die App legt zurück, was dort lag: ein Bild, eine Datei, ein kopiertes Stück Text. Lassen Sie es eingeschaltet, wenn Sie oft etwas kopieren und es nicht verlieren wollen.</td></tr>" +
			"<tr><td>Zeichenweise Eingabe (für Felder ohne Einfügen)</td><td>Statt den Text in einem Stück einzufügen, tippt die App ihn ein, als würden Sie selbst schreiben. Das ist langsamer, hilft aber dort, wo Einfügen aus der Zwischenablage gesperrt ist: in manchen Bankformularen, Terminals, alten Programmen.</td></tr>" +
			"<tr><td>Verzögerung vor dem Einfügen</td><td>Eine kleine Pause, bevor der Text an das Fenster geht. Selten nötig: Manchmal nimmt ein Programm den Text nicht schnell genug an, und der Anfang des Satzes geht verloren.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Bildschirmanzeige</p>" +
			"<p>Hier legen Sie fest, ob die Leiste überhaupt erscheint und wo genau.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Leiste anzeigen</td><td>Schalten Sie sie aus, wenn sie stört oder Ihnen die Tonsignale genügen.</td></tr>" +
			"<tr><td>Wo die Leiste erscheint</td><td>Ein Schema des Bildschirms mit neun Plätzen: Klicken Sie auf die Ecke, in der die Leiste nicht stört.</td></tr>" +
			"<tr><td>Am Cursor</td><td>Die Leiste erscheint neben der Stelle, an der Sie diktieren. Meldet das Programm nicht, wo der Cursor steht, erscheint sie neben dem Mauszeiger.</td></tr>" +
			"<tr><td>Eigener Platz</td><td>Die Leiste lässt sich einfach mit der Maus verschieben, und sie merkt sich die Stelle.</td></tr>" +
			"<tr><td>Bildschirm</td><td>Auf welchem Monitor sie erscheint, wenn es mehrere gibt.</td></tr>" +
			"<tr><td>Erkannten Text anzeigen</td><td>Nach dem Einfügen zeigt die Leiste den Satz selbst statt der Zeichenzahl. Praktisch zum Kontrollieren, aber den Satz sieht dann auch jeder, der auf Ihren Bildschirm schaut.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">Die Registerkarte „Sprachen & Modelle“</p>" +
			"<p>Hier sagen Sie, in welcher Sprache Sie diktieren werden, und die App entscheidet, womit sie das erkennt. Es gibt mehrere Modelle: Die einen sind schnell, die anderen genau, wieder andere kennen seltene Sprachen.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Sprache der Rede</td><td>Die Sprache, in der Sie sprechen. Sprechen Sie mehrere, wählen Sie „Automatisch“: Das Modell findet es selbst heraus, wenn auch etwas langsamer und etwas ungenauer, als wenn die Sprache ausdrücklich genannt ist.</td></tr>" +
			"<tr><td>Welches Modell für welche Sprache</td><td>Klicken Sie auf eine Sprache, und es öffnet sich eine Liste der Modelle, die sie beherrschen, mit einer Bewertung von Tempo und Genauigkeit und der Dateigröße. So kann eine Sprache ein schnelles Modell behalten, während eine andere ein genaues bekommt. Die gesonderte Zeile „Automatisch“ gilt für alle: Ihr Modell bedient jede Sprache, der kein eigenes zugewiesen ist, und diese Sprachen stehen grau in der Liste. Dasselbe Modell erkennt auch, welche Sprache gesprochen wird, wenn keine Sprache der Rede gesetzt ist.</td></tr>" +
			"<tr><td>Eigenes Kürzel</td><td>In einer aufgeklappten Sprache gibt es Platz für eine eigene Tastenkombination. Halten Sie sie gedrückt, und das Diktat läuft in dieser Sprache mit ihrem Modell, ihren Ersetzungen und Regeln, egal was als Sprache der Rede eingestellt ist; das Hauptkürzel behält die oben gewählte Sprache. So diktieren Sie Englisch mit einer Kombination und Ihre eigene Sprache mit einer anderen, ohne die Einstellungen zu öffnen. Während der Aufnahme zeigt die Leiste den Sprachcode. Eine Kombination darf weder das Hauptkürzel noch das einer anderen Sprache wiederholen.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Russisch<span class=\"hs\">GigaAM v3 · Genauigkeit 5 von 5 · Tempo 5 von 5 · 232 MB</span></label><span class=\"hchip\"><i></i>aktiv</span></div>" +
			"<div class=\"hrow\"><label>English<span class=\"hs\">wie automatische Erkennung · Nemotron 3.5</span></label><button class=\"hbtn\">Auswählen</button></div>" +
			"<div class=\"hcap\">zu sehen ist, welches Modell eine Sprache bedient und was das an Tempo und Speicherplatz kostet</div>" +
			"</div>" +
			"<p class=\"hsub\">Übersetzung</p>" +
			"<p>Manche Modelle können das Gesagte nicht nur aufschreiben, sondern gleich übersetzen: Sie sprechen in einer Sprache, und im Feld erscheint der Text in einer anderen. Übersetzt wird mit Whisper: ins Englische mit seiner eingebauten Übersetzungsfunktion, in andere Sprachen, indem man ihm die Ausgabesprache als Sprache der Rede nennt; bei weit entfernten Sprachpaaren schwankt die Qualität. Kann das gewählte Modell nicht übersetzen, erscheint neben dem Schalter eine orange Zeile: Ist Whisper installiert, bietet die Leiste an, damit zu übersetzen, sonst bleibt der Text in der Ausgangssprache.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Sprache der Textausgabe ändern</td><td>Schaltet die Übersetzung ein.</td></tr>" +
			"<tr><td>Standardsprache der Textausgabe</td><td>In welche Sprache übersetzt wird, wenn Sie nicht direkt beim Diktat eine andere gewählt haben.</td></tr>" +
			"<tr><td>Nach der Ausgabesprache fragen</td><td>Sie können auf die Frage verzichten und immer in die Standardsprache übersetzen. Sie können jedes Mal gefragt werden: Nach dem Diktat erscheinen auf der Leiste Schaltflächen mit Sprachen. Oder Sie werden mit Timeout gefragt: dieselben Schaltflächen, aber wenn Sie nicht rechtzeitig drücken, gilt die Standardsprache.</td></tr>" +
			"<tr><td>Sprachen im Dialog</td><td>Welche Sprachen als diese Schaltflächen erscheinen. Lassen Sie die zwei oder drei stehen, die Sie wirklich nutzen, sonst dauert das Wählen länger als das Sprechen.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Modelle auf der Festplatte</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Eigenes Modell</td><td>Haben Sie ein eigenes Modell, legen Sie es in den Modellordner und starten Sie die App neu: Es erscheint in den Listen der Sprachen, die es versteht. Die Schaltfläche daneben öffnet den richtigen Ordner, damit Sie ihn nicht im Explorer suchen müssen.</td></tr>" +
			"<tr><td>Aus dem Speicher entladen</td><td>Solange die App läuft, belegt das Modell Arbeitsspeicher. Diese Schaltfläche gibt ihn sofort frei, etwa vor einem anspruchsvollen Spiel. Das nächste Diktat lädt das Modell erneut, und der erste Satz dauert länger als sonst.</td></tr>" +
			"<tr><td>Installierte Modelle prüfen</td><td>Vergleicht die heruntergeladenen Dateien mit ihrer Referenz. Nützlich, wenn die Verbindung beim Herunterladen abgebrochen ist und Sie vermuten, dass eine Datei beschädigt wurde.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">Die Registerkarte „Regeln“</p>" +
			"<p>Drei Listen, die den Text automatisch korrigieren, noch bevor er ins Fenster gelangt. Alle drei sind freiwillig: Die App funktioniert auch ohne sie bestens.</p>" +
			"<p class=\"hsub\">Satzzeichen und Großschreibung</p>" +
			"<p>Woher Punkte, Kommas und Großbuchstaben im Text kommen. Es gibt drei Möglichkeiten: <b>vom Modell</b> (das Modell setzt sie so, wie es sie gehört hat; schnell und meist ausreichend), <b>vom Bearbeitungsmodell</b> (die Zeichen setzt ein zweites Modell; sauberer, aber jedes Diktat dauert ein paar Sekunden länger) und <b>entfernen</b> (der Text kommt als durchgehende Zeile ohne Zeichen; selten nötig, etwa für ein Suchfeld).</p>" +
			"<p class=\"hsub\">Wörterbuch</p>" +
			"<p>Eine Liste von Wörtern, die das Modell meist nicht kennt: Nachnamen, Firmennamen, Fachbegriffe. Das ist kein Befehl „schreib es so“, sondern ein Hinweis: Hört das Modell etwas Ähnliches, wählt es eher Ihre Variante.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">Wörter werden einzeln hinzugefügt und stehen als eigene Karten da</div>" +
			"</div>" +
			"<p class=\"hdim\">Das Wörterbuch lesen nur die Whisper-Modelle. Nemotron, Parakeet, GigaAM und die übrigen Modelle nehmen überhaupt keine Hinweise an: Für sie ist die Liste der Ersetzungen weiter unten das Werkzeug. Der Standardsatz an Wörtern wechselt mit der Sprache, bis Sie eigene eintragen.</p>" +
			"<p class=\"hsub\">Ersetzungen nach der Erkennung</p>" +
			"<p>Auch eine Liste von Wörtern, aber sie arbeitet anders: Das ist kein Hinweis mehr, sondern eine Garantie. Was links steht, wird zu dem, was rechts steht. Jede Ersetzung hat drei Schalter: <b>Nur ganze Wörter</b> (dann rührt „Bahn“ das Wort „Bahnhof“ nicht an, und meist ist genau das gewollt), <b>Groß-/Kleinschreibung beachten</b> (dann gelten „Go“ und „go“ als verschiedene Wörter) und <b>Sprache der Regel</b> (die Ersetzung greift nur, wenn Sie in dieser Sprache diktieren).</p>" +
			"<p class=\"hsub\">Sprachbefehle</p>" +
			"<p>Wendungen, die nicht im Text landen, sondern eine Aktion auslösen. Zeilenumbruch, neuer Absatz und eigener Text werden als ganze Wörter an jeder Stelle des Satzes erkannt. Der Abbruch und die Formbefehle wirken nur am Ende des Diktats, das Löschen nur, wenn es allein gesagt wird. Ein Formbefehl, allein gesagt, gilt für das nächste Diktat, und die Leiste zeigt „Nächste Phrase: …“.</p>" +
			"<ul>" +
			"<li><b>„neue Zeile“</b>: in eine neue Zeile wechseln;</li>" +
			"<li><b>„neuer Absatz“</b>: eine Leerzeile zwischen Absätzen lassen;</li>" +
			"<li><b>„abbrechen“</b>: den Satz verwerfen und nichts einfügen;</li>" +
			"<li><b>„lösch das“</b>: löschen, was das Programm zuletzt eingefügt hat;</li>" +
			"<li><b>„als Überschrift“</b>: ein Zeichen vor den Satz setzen, standardmäßig #, und danach eine Leerzeile;</li>" +
			"<li><b>„als Liste“</b>: den Satz an Kommas, Semikolons und am Wort „und“ in Zeilen mit einem Strich am Anfang teilen;</li>" +
			"<li><b>„großgeschrieben“</b>: den Satz mit einem Großbuchstaben beginnen;</li>" +
			"<li><b>„in Anführungszeichen“</b>: den Satz in Anführungszeichen setzen;</li>" +
			"<li><b>„ohne Punkt“</b>: den Punkt am Ende weglassen;</li>" +
			"<li><b>eigener Befehl</b>: einen vorbereiteten Text einfügen, etwa eine Signatur, eine Adresse oder eine Standardantwort.</li>" +
			"</ul>" +
			"<p>Das Löschen entfernt nur, was das Programm selbst zuletzt eingefügt hat, und nur, wenn seitdem nichts getippt oder geklickt wurde, das Fenster dasselbe ist, der Text keine Zeilenumbrüche hatte, nicht mit Enter gesendet wurde und höchstens fünf Minuten alt ist. Text, der schon im Fenster stand, wird nie angefasst.</p>" +
			"<p class=\"wh\" id=\"s-post\">Die Registerkarte „Nachbearbeitung“</p>" +
			"<p>Die Nachbearbeitung ist ein zweites Modell, das den erkannten Text nach Ihrer Anweisung umschreibt. Zum Beispiel: „entferne Füllwörter und Wiederholungen“, „mach daraus eine geschäftliche E-Mail“, „mach daraus eine Liste“. Die Anweisung wird in gewöhnlichen Worten geschrieben und heißt Prompt.</p>" +
			"<p>Standardmäßig ist sie ausgeschaltet, und zwar aus diesem Grund: Sie ist der einzige Schritt, der das Diktat spürbar verlangsamt. Statt ein, zwei Sekunden dauert das Einfügen mehrere Sekunden, auf einem schwachen Rechner auch zehn.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Nachbearbeitung einschalten</td><td>Der Hauptschalter. Solange er aus ist, funktioniert der Rest der Registerkarte nicht, und der Text wird so eingefügt, wie er erkannt wurde.</td></tr>" +
			"<tr><td>Modell · Lokal</td><td>Ein Editor, der auf Ihrem Rechner läuft. Nichts verlässt den Rechner, aber er braucht Speicher: Solche Modelle sind ab einem Gigabyte groß. Daneben liegt ein Katalog, in dem ein Modell mit einem Klick gewählt und heruntergeladen wird.</td></tr>" +
			"<tr><td>Modell · Externer Server</td><td>Statt Ihres eigenen Modells geht der Text an einen fremden Dienst: OpenAI, Groq oder Ihren eigenen Server. Schneller und klüger, aber der Text verlässt den Rechner. Der Zugriffsschlüssel wird verschlüsselt gespeichert, und die Schaltfläche „Verbindung testen“ prüft, ob Adresse und Schlüssel funktionieren.</td></tr>" +
			"<tr><td>Prompts</td><td>Ihre Anweisungen an den Editor. Die angehakten werden der Reihe nach von oben nach unten angewendet: erst zum Beispiel das Aufräumen, dann der Geschäftsstil. Ist nichts angehakt, tut die Nachbearbeitung nichts, und die App warnt davor mit einer gelben Zeile.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Die Stimme geht nie an einen externen Server, nur der erkannte Text. Aber der geht eben hin, also schalten Sie das bewusst ein.</div>" +
			"<p>Jeder Prompt hat einen <b>Namen</b> (unter ihm erscheint er während der Arbeit auf der Leiste), den <b>Text der Anweisung</b> in gewöhnlichen Worten, ein optionales <b>eigenes Tastenkürzel</b> (ein Diktat damit schickt den Satz nur durch diesen Prompt und lässt die anderen aus) und einen <b>Test</b>: ein Feld, in das Sie einen Satz eintippen und sehen können, was herauskommt, ohne etwas zu diktieren.</p>" +
			"<p class=\"wh\" id=\"s-subs\">Die Registerkarte „Untertitel“</p>" +
			"<p>Untertitel schreiben mit, was der Computer abspielt – ein Gespräch, ein Video, einen Vortrag – und zeigen die Worte in einem kleinen Fenster, das über allen anderen Fenstern bleibt. Es ist dieselbe Erkennung wie beim Diktieren, nur kommt der Ton aus den Lautsprechern statt vom Mikrofon, und das Diktieren funktioniert weiter, während Untertitel laufen. Die eigenen Töne der App werden nicht mitgeschrieben.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Untertitel starten</td><td>Startet und beendet sie. Die Zeile daneben sagt, was gerade passiert: aus, Modell wird geladen, es wird zugehört oder die Erkennung startet neu. Dasselbe tun das Tray-Menü, ein eigenes Tastenkürzel und das Kreuz am Fenster.</td></tr>" +
			"<tr><td>Sprache der Rede</td><td>Die Sprache dessen, was Sie hören, getrennt von der Diktiersprache. Das Modell ist das, das diese Sprache auf der Registerkarte „Sprachen & Modelle“ hat. Ein Streaming-Modell schreibt die Wörter schon beim Sprechen, die anderen jeden Satz nach seinem Ende. Ein Whisper-Modell braucht auf dem Prozessor viele Sekunden pro Satz, sein Text hinkt der Rede deutlich hinterher; aufgestaute Sätze gehen gemeinsam an das Modell und kommen als eine Zeile heraus.</td></tr>" +
			"<tr><td>Übersetzung</td><td>Zeigt im Fenster die Übersetzung statt des erkannten Textes, in der gewählten Sprache; der Fenstertitel lautet dann „Übersetzung“ mit den Sprachen, etwa EN → RU. Die Rede wird schon während des Sprechens in Satzteile geschnitten, und jeder Satzteil wird sofort übersetzt. Ein Satzteil endet am Satzende oder vor einem Bindewort, und der Übersetzer bekommt die vorigen Satzteile und seine eigenen Übersetzungen mit, sodass ein zerteilter Satz zusammenhängend übersetzt wird. Die Übersetzung erscheint gleichmäßig, in ganzen Wörtern, im Lesetempo und schneller, wenn sich viel angesammelt hat, sodass ein großes Stück den Text nicht ruckartig verschiebt. Satzteile eines Redeabschnitts laufen als ein Absatz, nach einer Pause beginnt eine neue Zeile. Ein Satzteil, der nicht übersetzt werden konnte, wird so gezeigt, wie er gehört wurde, in einer blasseren Farbe. „Auch den Originaltext zeigen“ teilt das Fenster mit einer dünnen Linie: oben blasser das Gesagte, unten die Übersetzung; die Linie lässt sich mit der Maus ziehen. Wer übersetzt, legt die Zeile „Wer übersetzt“ fest. „Hochwertig“ ist das Übersetzungsmodell auf der Grafikkarte: einmal 2,3 GB herunterzuladen, für eine Grafikkarte mit mindestens 4 GB eigenem Speicher; es wird beim Start der Untertitel auf die Grafikkarte geladen und bei ihrem Ende wieder entfernt, und ein Link daneben öffnet seine Nutzungsbedingungen. „Ihr eigener Server“ ist der Server für die Nachbearbeitung, der auf der Registerkarte „Nachbearbeitung“ eingerichtet ist, und dann wird jeder Satzteil an diesen Server gesendet. „Das Programm wählt“ nimmt das Übersetzungsmodell, wenn es heruntergeladen ist und die Grafikkarte passt, sonst Ihren Server, wenn einer eingerichtet ist, sonst das auf der Registerkarte „Nachbearbeitung“ installierte Bearbeitungsmodell, das grob übersetzt: Der Sinn kommt an, aber Wörter können danebenliegen. Ein von Hand gewählter Übersetzer wird nie durch einen anderen ersetzt: Kann er nicht arbeiten, sagt der Block, was fehlt, und der Text wird so gezeigt, wie er gehört wurde. Die Übersetzung wird als eigene Datei neben der Mitschrift gespeichert: jeder Satzteil mit der Uhrzeit, zu der er begann; ein Satzteil, der nicht übersetzt werden konnte, wird so geschrieben, wie er gehört wurde. Ohne Übersetzer gibt es keine Übersetzung, und die Untertitel laufen wie bisher.</td></tr>" +
			"<tr><td>Tastenkürzel</td><td>Startet und beendet die Untertitel aus jedem Fenster. Es startet nie ein Diktat und darf kein anderes Kürzel wiederholen.</td></tr>" +
			"<tr><td>Mitschrift speichern</td><td>Jeder Durchlauf schreibt eine Textdatei in den Ordner transcripts neben dem Programm: eine erste Zeile mit Datum, Sprache und Modell, dann jede fertige Zeile mit der Uhrzeit, zu der sie begann. Bis die erste Zeile fertig ist, wird nichts geschrieben. Bei eingeschalteter Übersetzung trägt der Name der Mitschrift original und den Code der gesprochenen Sprache, und daneben entsteht eine zweite Datei mit translation und dem Code der Zielsprache, etwa 2026-11-12_10-18-29_original_EN.txt und 2026-11-12_10-18-29_translation_RU.txt.</td></tr>" +
			"<tr><td>Untertitelfenster</td><td>Wo es steht, wie groß es ist und wie groß die Schrift ist. Der Platz wird auf dem Bildschirmbild gewählt, wie bei der Diktierleiste: Ein Klick auf einen der neun Plätze stellt das Fenster dorthin, und mit „Eigener Platz“ kann es überall stehen, indem Sie die Markierung im Bild oder das Fenster selbst ziehen; bei mehreren Monitoren lässt sich der Bildschirm wählen. Die Größe ändern Sie mit der Maus an Rändern und Ecken; die Höhe rastet auf ganze Zeilen ein, Breite und Zeilenzahl werden für jede Bildschirmauflösung gemerkt. „Textzeilen“ gilt, wo noch keine Größe festgelegt ist, und „Größe zurücksetzen“ vergisst die festgelegten. Das Mausrad über dem Fenster blättert den Text um drei Zeilen je Raste zurück; während Sie Älteres lesen, verschiebt neuer Text es nicht, und die Markierung „zum Neuesten ↓“ im Titel führt zur aktuellen Rede zurück.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Ihre Gesprächspartner sollten wissen, dass ihre Worte mitgeschrieben werden.</div>" +
			"<p class=\"wh\" id=\"s-about\">„Über“ und „Kontakte“</p>" +
			"<p>Auf der Registerkarte „Über“ stehen eine kurze Beschreibung, was das ist und wie es aufgebaut ist, und eine Liste aller fremden Bibliotheken, aus denen die App gebaut ist, mit ihren Lizenzen. Die Namen sind anklickbar und führen zu den Seiten dieser Projekte. Dort steht auch ein Link zum Repository, in dem der Quellcode und alle Versionen liegen.</p>" +
			"<p>„Kontakte“ besteht aus drei Karten: E-Mail, Repository und eine Seite für Fehler und Ideen. Ein Klick auf die E-Mail-Adresse öffnet Ihr E-Mail-Programm mit bereits eingetragenem Empfänger, und das Symbol daneben kopiert die Adresse einfach.</p>" +
			"<p class=\"wh\" id=\"s-files\">Was im Ordner der App liegt</p>" +
			"<p>Die App ist portabel: Alles, was zu ihr gehört, liegt in einem Ordner und breitet sich nirgendwo sonst aus.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Alle Ihre Einstellungen in einer Textdatei. Sie können sie auf einen anderen Rechner kopieren, um nicht alles neu einzurichten.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>Das Protokoll. Der erste Ort, an dem Sie nachsehen sollten, wenn sich die App seltsam verhält.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>Der Diktatverlauf, falls Sie ihn eingeschaltet haben. Gewöhnlicher Text: Wenn Sie möchten, können Sie die Datei einfach löschen.</td></tr>" +
			"<tr><td>models</td><td>Der Ordner mit den heruntergeladenen Modellen, der schwerste Teil, mehrere Gigabyte. Auch ihn können Sie auf einen anderen Rechner mitnehmen, um nicht alles neu herunterzuladen.</td></tr>" +
			"<tr><td>transcripts</td><td>Die Mitschriften der Untertitel, eine Textdatei pro Durchlauf, und ihre Übersetzungen, wenn die Übersetzung eingeschaltet ist. Es sind Ihre Dokumente: Die Deinstallation entfernt sie nur, wenn Sie auch das Löschen von Einstellungen und Daten verlangen.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>Die kleinen Programme, die die Modelle ausführen. Ohne sie ist ein Modell nur eine Datei, die niemand lesen kann.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Was ins Internet geht und was nicht</p>" +
			"<ul>" +
			"<li><b>Die Stimme</b> wird weder gespeichert noch gesendet, es sei denn, Sie haben selbst die Erkennung auf einem anderen Rechner eingeschaltet.</li>" +
			"<li><b>Untertitel</b> werden hier erkannt, genau wie Diktate, und die Mitschrift bleibt im Ordner transcripts; bei Erkennung auf einem anderen Rechner gehen die Sätze an diesen Server; bei Übersetzung über Ihren eigenen Server für die Nachbearbeitung geht auch jede fertige Zeile dorthin.</li>" +
			"<li><b>Der Text</b> verlässt den Rechner nicht, es sei denn, Sie haben selbst einen externen Server für die Nachbearbeitung eingeschaltet.</li>" +
			"<li><b>Der Verlauf</b> liegt als Datei neben der App und wird nach der von Ihnen gesetzten Frist geleert.</li>" +
			"<li><b>Das Netz</b> wird zweimal gebraucht: um ein Modell herunterzuladen und um nach einem Update zu sehen. Die zweite Prüfung lässt sich abschalten.</li>" +
			"<li><b>Der Schlüssel</b> für einen externen Dienst wird mit Mitteln von Windows verschlüsselt gespeichert und landet nie im Protokoll.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Wenn etwas schiefgeht</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Der Text erscheint nicht im Feld</td><td>Das Programm nimmt kein Einfügen aus der Zwischenablage an. Schalten Sie „Zeichenweise Eingabe (für Felder ohne Einfügen)“ auf der Registerkarte „Steuerung & Verhalten“ ein oder fügen Sie eine kleine Verzögerung vor dem Einfügen hinzu.</td></tr>" +
			"<tr><td>Das Diktat startet nicht</td><td>Das Tastenkürzel ist von einem anderen Programm belegt. Legen Sie ein anderes fest, am besten mit der Windows-Taste, solche sind seltener belegt.</td></tr>" +
			"<tr><td>Aufgenommen wurde nur Stille</td><td>Falsches Mikrofon, oder es ist stummgeschaltet. Öffnen Sie die Registerkarte „Mikrofon“, schauen Sie auf den Pegelbalken und drücken Sie „Mikrofon prüfen“.</td></tr>" +
			"<tr><td>Der erste Satz dauert sehr lange</td><td>Das ist so gewollt: Das Modell wird beim ersten Aufruf in den Speicher geladen. Danach dauert die Erkennung Sekunden.</td></tr>" +
			"<tr><td>Die Meldung „Erkennung nicht verfügbar“</td><td>Die Erkennung auf einem anderen Rechner ist gewählt, aber keine Adresse gesetzt. Tragen Sie die Adresse ein oder kehren Sie zu „Auf diesem Rechner“ zurück.</td></tr>" +
			"<tr><td>Die Nachbearbeitung ändert nichts</td><td>Kein Prompt ist angehakt, kein Modell gewählt, oder der externe Server antwortet nicht. Der Grund steht auf der Registerkarte selbst, neben dem Schalter.</td></tr>" +
			"<tr><td>Der Text ist seltsam, Wörter sind ersetzt</td><td>Schauen Sie in „Regeln“: Wahrscheinlich hat eine Ersetzung oder ein Sprachbefehl gegriffen.</td></tr>" +
			"<tr><td>Nichts davon trifft zu</td><td>„System“ → „Protokoll öffnen“. Die App schreibt dort alles über sich auf, und der Grund steht meist in Worten da.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">Keine Antwort gefunden? Schreiben Sie eine E-Mail über die Registerkarte „Kontakte“ oder eröffnen Sie ein Issue im Repository. Beides liest der Autor.</p>",
		"S_WIZ_HELLO":      "Zwei Minuten Einrichtung",
		"S_WIZ_HELLO_TEXT": "{app} macht aus Ihrer Stimme Text direkt an der Eingabemarke: Tastenkürzel halten, Satz sprechen, loslassen — der Text steht da. Alles läuft auf Ihrem Rechner, der Ton verlässt ihn nie.",
		"S_WIZ_LATER":      "Alles, was wir jetzt wählen, lässt sich später in den Einstellungen ändern.",
		"S_WIZ_MODEL_TEXT": "Das Modell für Ihre Sprache wird für Sie ausgewählt und im Hintergrund geladen.",
		"S_WIZ_PICKED":     "Ausgewählt für: %s",
		"S_WIZ_DL_BG":      "Wird im Hintergrund geladen, Sie können weitermachen",
		"S_WIZ_DL_FAIL":    "Das Laden wurde unterbrochen: %s",
		"S_WIZ_RETRY":      "Erneut versuchen",
		"S_WIZ_TRY_DL":     "Das Modell wird noch geladen (%d%%). Das Feld funktioniert, sobald es bereit ist.",
		"S_WIZ_TRY_LOAD":   "Das Modell wird geladen…",
		"S_WIZ_TRY_LIVE":   "Mit %s erscheinen die Wörter in der Leiste, während Sie noch sprechen.",
		"S_WIZ_DONE_DL":    "Das Modell wird noch geladen. Das Diktieren funktioniert von selbst, sobald es da ist.",
		"S_WIZ_INPUT_TEXT": "Diese Tastenkombination halten Sie beim Sprechen. Sagen Sie etwas und prüfen Sie, ob sich der Pegelbalken bewegt.",
		"S_WIZ_T_TRY":      "Ausprobieren",
		"S_WIZ_TRY_PH":     "hier erscheint der Text",
		"S_WIZ_T_DONE":     "Fertig",
		"S_WIZ_DONE_TEXT":  "{app} wohnt im Infobereich: Linksklick auf das Symbol öffnet die Einstellungen, Rechtsklick das Menü. Diktieren können Sie in jedem Fenster mit Eingabemarke.",
		"S_AUTORUN":        "Mit Windows starten",
		"S_AUTORUN_SUB":    "Ein Eintrag im Autostart des aktuellen Benutzers",
		"S_WIZ_SKIP":       "Überspringen",
		"S_WIZ_BACK":       "Zurück",
		"S_WIZ_NEXT":       "Weiter",
		"S_WIZ_FINISH":     "Fertigstellen",
		"S_WIZ_WAIT":       "Warte auf den ersten Satz…",
		"S_WIZ_HEARD":      "Gehört:",
		"S_WIZ_HAVE":       "Alles Nötige ist schon geladen",
		"S_WIZ_TRY_TEXT":   "Setzen Sie die Eingabemarke ins Feld unten, halten Sie %s, sprechen Sie einen Satz und lassen Sie los.",
	}
	settingsStrings["fr"] = map[string]string{
		"S_SUBS_ON": "Les sous-titres tournent",
		"S_SUBS_HERO_HINT": "démarrez-les — un appel ou une vidéo devient du texte à l'écran",
		"S_SUBS_HERO_START": "Démarrer manuellement",
		"S_SUBS_HERO_STOP": "Arrêter",
		"S_UPDATED_TO": "Mis à jour vers la version %s",
		"S_SUBS_HERO_LIVE": "%s · les mots au fil de la parole",
		"S_SUBS_HERO_PHRASES": "%s · le texte après chaque phrase",
		"S_NAV_SUBS": "Sous-titres",
		"S_SUBS_HINT": "L'application note ce que joue l'ordinateur – un appel, une vidéo – et affiche les mots dans une fenêtre au-dessus des autres.",
		"S_SUBS_NOTICE": "Vos interlocuteurs doivent savoir que leurs paroles sont transcrites.",
		"S_SUBS_START": "Démarrer les sous-titres",
		"S_SUBS_STOP": "Arrêter les sous-titres",
		"S_SUBS_OFF": "Les sous-titres sont arrêtés",
		"S_SUBS_LOADING": "Chargement de %s…",
		"S_SUBS_LISTEN_LIVE": "À l'écoute du son de l'ordinateur · %s, les mots au fil de la parole",
		"S_SUBS_LISTEN_PHRASES": "À l'écoute du son de l'ordinateur · %s, le texte après chaque phrase",
		"S_SUBS_RESTARTING": "La reconnaissance redémarre…",
		"S_SUBS_LANG": "Langue parlée",
		"S_SUBS_LANG_SUB": "La langue de ce que vous écoutez, distincte de la langue de dictée",
		"S_SUBS_M_LIVE": "%s écrira : les mots apparaissent au fil de la parole.",
		"S_SUBS_M_PHRASES": "%s écrira : le texte apparaît une fois chaque phrase terminée.",
		"S_SUBS_M_ABSENT": "%s est choisi pour cette langue mais n'est pas installé.",
		"S_SUBS_M_WHISPER": "%s est un modèle Whisper : sur un processeur il lui faut de nombreuses secondes par phrase, le texte aura donc un net retard sur la parole. Un modèle en streaming ou un modèle rapide convient mieux aux sous-titres.",
		"S_SUBS_GO_MODELS": "Ouvrir « Langues et modèles »",
		"S_SUBS_REMOTE": "La reconnaissance tourne sur un autre ordinateur : chaque phrase est envoyée à ce serveur.",
		"S_SUBS_TR": "Traduction",
		"S_SUBS_TR_SUB": "La fenêtre affiche la traduction au lieu du texte reconnu",
		"S_SUBS_TR_OFF": "Désactivée",
		"S_SUBS_TR_LOCAL": "%s traduit : la traduction arrive phrase par phrase, un instant après la parole.",
		"S_SUBS_TR_API": "Votre serveur de post-traitement (%s) traduit : chaque phrase lui est envoyée, et la traduction arrive phrase par phrase, un instant après la parole.",
		"S_SUBS_TR_NONE": "Il n'y a ni modèle de traduction, ni serveur de post-traitement, ni modèle d'édition, il n'y aura donc pas de traduction.",
		"S_SUBS_TR_GO": "Ouvrir « Post-traitement »",
		"S_SUBS_TR_SAME": "La langue de traduction est la même que la langue parlée : il n'y a rien à traduire.",
		"S_SUBS_WHO": "Qui traduit",
		"S_SUBS_WHO_SUB": "Le programme prend le meilleur disponible, ou vous en fixez un",
		"S_SUBS_WHO_AUTO": "Le programme choisit",
		"S_SUBS_WHO_QUALITY": "Haute qualité : un modèle sur la carte graphique",
		"S_SUBS_WHO_SERVER": "Votre propre serveur",
		"S_SUBS_TRM": "Modèle de traduction",
		"S_SUBS_TRM_SUB": "%s, fonctionne sur une carte graphique avec au moins 4 Go de mémoire propre",
		"S_SUBS_TRM_ABSENT": "Non téléchargé",
		"S_SUBS_TRM_GETTING": "Téléchargement… %s %",
		"S_SUBS_TRM_CHECKING": "Vérification du fichier…",
		"S_SUBS_TRM_READY": "Prêt",
		"S_SUBS_TRM_RUNNING": "Prêt, chargé sur la carte graphique",
		"S_SUBS_TRM_BAD": "Le fichier du dossier models est endommagé ou incomplet",
		"S_SUBS_TRM_ERR": "Non téléchargé : %s",
		"S_SUBS_TRM_SMALL": "Cet ordinateur ne peut pas le faire tourner : sa carte graphique a %s Mo de mémoire propre, il en faut 4 Go.",
		"S_SUBS_TRM_INTEGRATED": "Cet ordinateur ne peut pas le faire tourner : il n'a qu'une puce graphique intégrée, qui partage la mémoire de l'ordinateur.",
		"S_SUBS_TRM_NOCARD": "Cet ordinateur ne peut pas le faire tourner : Windows ne signale aucune carte graphique.",
		"S_SUBS_TRM_UNKNOWN": "Cet ordinateur ne peut pas le faire tourner : cette version de Windows ne décrit pas sa carte graphique.",
		"S_SUBS_TRM_DOWN": "N'a pas démarré : %s",
		"S_SUBS_TRM_DOWN_NEXT": "Jusqu'au prochain démarrage des sous-titres, le traducteur suivant prend le relais.",
		"S_SUBS_TRM_DOWN_HEARD": "Jusqu'au prochain démarrage des sous-titres, le texte s'affiche tel qu'il a été entendu.",
		"S_SUBS_TRM_DEL_SURE": "Supprimer ?",
		"S_SUBS_TRM_W_VULKAN": "le pilote graphique n'offre pas Vulkan",
		"S_SUBS_TRM_W_DEVICE": "le moteur ne voit pas cette carte graphique",
		"S_SUBS_TRM_W_MEMORY": "pas assez de mémoire libre sur la carte graphique",
		"S_SUBS_TRM_W_CRASH": "le moteur s'est arrêté pendant le chargement",
		"S_SUBS_TRM_W_TIMEOUT": "le moteur n'a pas répondu à temps",
		"S_SUBS_TRM_W_PORT": "son port est pris par un autre programme",
		"S_SUBS_TRM_W_EXE": "llama-server-gpu.exe manque à côté du programme",
		"S_SUBS_TRM_DL": "Télécharger",
		"S_SUBS_TRM_CANCEL": "Annuler",
		"S_SUBS_TRM_DEL": "Supprimer",
		"S_SUBS_TRM_TERMS": "Conditions d'utilisation du modèle",
		"S_SUBS_TR_QUALITY": "Le modèle de traduction sur la carte graphique (%s) traduit : la traduction arrive phrase par phrase, un instant après la parole.",
		"S_SUBS_TR_MISS_CARD": "« Haute qualité » est choisi, mais la carte graphique de cet ordinateur ne convient pas, le texte s'affiche donc tel qu'il a été entendu.",
		"S_SUBS_TR_MISS_MODEL": "« Haute qualité » est choisi, mais le modèle de traduction n'est pas téléchargé, le texte s'affiche donc tel qu'il a été entendu.",
		"S_SUBS_TR_MISS_START": "« Haute qualité » est choisi, mais le modèle de traduction n'a pas démarré sur la carte graphique, le texte s'affiche donc tel qu'il a été entendu jusqu'au prochain démarrage des sous-titres.",
		"S_SUBS_TR_MISS_SERVER": "Votre propre serveur est choisi, mais aucun serveur de post-traitement n'est configuré, le texte s'affiche donc tel qu'il a été entendu.",
		"S_SUBS_ORIG": "Afficher aussi le texte d'origine",
		"S_SUBS_ORIG_SUB": "la fenêtre affiche au-dessus de la traduction ce qui a été dit ; la ligne entre les deux se tire à la souris",
		"S_SUBS_KEY": "Raccourci clavier",
		"S_SUBS_KEY_SUB": "Démarre et arrête les sous-titres depuis n'importe quelle fenêtre",
		"S_SUBS_SAVE": "Enregistrer la transcription",
		"S_SUBS_SAVE_SUB": "Un fichier texte par session dans le dossier transcripts à côté du programme, et un second avec la traduction quand la traduction est activée",
		"S_SUBS_FOLDER": "Dossier des transcriptions",
		"S_SUBS_WINDOW": "Fenêtre des sous-titres",
		"S_SUBS_ROWS": "Lignes de texte",
		"S_SUBS_SIZE": "Taille du texte",
		"S_SUBS_SIZE_S": "Petite",
		"S_SUBS_SIZE_M": "Moyenne",
		"S_SUBS_SIZE_L": "Grande",
		"S_SUBS_BOX": "Taille de la fenêtre",
		"S_SUBS_BOX_SUB": "tirez les bords et les coins de la fenêtre à la souris ; la taille est retenue pour chaque résolution d'écran, et « Lignes de texte » s'applique tant qu'il n'y en a pas",
		"S_SUBS_BOX_BTN": "Rétablir la taille",
		"S_SUBS_BOX_DONE": "La fenêtre revient à sa taille par défaut",
		"S_SUBS_FREE": "Place à vous",
		"S_SUBS_FREE_SUB": "la fenêtre se déplace où vous voulez, l'emplacement est retenu pour chaque résolution d'écran",
		"S_SUBS_POS": "Où afficher la fenêtre",
		"S_SUBS_POS_SCHEME_SUB": "cliquez sur l’écran — la fenêtre s’y place",
		"S_SUBS_POS_DRAG_SUB": "faites glisser la marque ou la fenêtre elle-même — elle va où vous voulez",
		"S_SUBS_MON": "Écran",
		"S_SUBS_MON_SUB": "sur quel moniteur afficher la fenêtre des sous-titres",
		"S_TITLE": "{app} — Réglages", "S_DICT_HINT": "Termes, noms et abréviations séparés par des virgules — un indice pour l'oreille, pas des commandes. Cela vaut pour Whisper ; le russe via GigaAM l'ignore. Le jeu par défaut suit la langue de reconnaissance jusqu'à ce que vous écriviez le vôtre.",
		"S_TR_DEFAULT": "Changer la langue de sortie du texte", "S_TR_TARGET": "Langue de sortie du texte par défaut", "S_TR_ASK": "Demander la langue de sortie du texte", "S_TR_ASK_NEVER": "Ne pas demander — traduire aussitôt",
		"S_SRCLANG_SUB": "vous la parlez ; elle détermine le modèle de reconnaissance",
		"S_TR_LANGS_SUB": "ces langues deviennent des boutons sur le bandeau à l’insertion",
		"S_TR_UNAVAIL": "%s ne sait pas traduire — le texte restera dans la langue source",
		"S_TR_GET_WHISPER": "Télécharger %s pour la traduction (%d Mo)",
		"S_TR_STANDIN": "%s ne sait pas traduire — le bandeau proposera %s à la place",
		"S_TR_LOCK": "%s ne peut pas être retirée de la liste — c’est la langue de sortie du texte par défaut. Choisissez une autre langue par défaut, et %s pourra être exclue.",
		"S_TR_LOCK_OK": "Compris",
		"S_TR_ONE": "Plusieurs langues sont cochées, mais sans question la traduction ira toujours vers une seule — %s (langue de sortie par défaut). Les autres restent cochées mais sont désactivées.",
		"S_TR_NOMODEL": "%s ne sait pas traduire. Si vous continuez, la traduction sera désactivée et indisponible tant que ce modèle travaille.",
		"S_TR_CONFIRM": "Confirmer",
		"S_TR_ASK_ALWAYS": "Demander à chaque fois", "S_TR_ASK_TIMEOUT": "Demander, avec délai", "S_TR_SECONDS": "Délai, s",
		"S_TR_LANGS": "Langues du dialogue",
		"S_LLM_HINT": "Les profils cochés s'appliquent l'un après l'autre, de haut en bas, lors d'une dictée normale. Rien de coché — le texte est inséré tel quel.",
		"S_PROF_ADD": "Ajouter", "S_PROF_NAME": "Nom", "S_PROF_PROMPT": "Prompt", "S_PROF_TEST": "Test",
		"S_HOTKEY": "Raccourci clavier", "S_CHANGE": "Changer…", "S_UILANG": "Langue de l'interface", "S_AUTO": "Comme le système",
		"S_SEC_SOUND": "Son", "S_BEEP": "Signaux sonores", "S_SOUND": "Son du signal",
		"S_SND_SPEECH": "Voix de Windows", "S_SND_CHIME": "Clochette", "S_SND_SOFT": "Doux", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Appuyer Entrée après insertion (envoi auto)", "S_RESTORE": "Restaurer le presse-papiers après insertion",
		"S_NAV_HISTORY": "Historique", "S_HIST_ON": "Conserver l'historique des dictées", "S_HIST_ON_SUB": "texte seul, sur cet ordinateur ; l'audio n'est jamais conservé",
		"S_HIST_DAYS": "Combien de jours conserver", "S_HIST_MAX": "Combien d'entrées conserver",
		"S_HIST_SKIP": "Ne jamais enregistrer depuis ces programmes", "S_HIST_SKIP_SUB": "séparés par des virgules : keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Ajout d'un programme", "S_SKIP_EDIT_DLG": "Modification du programme", "S_SKIP_NAME": "Nom du programme", "S_SKIP_NAME_SUB": "Le nom de fichier sans chemin : keepass.exe. Une étoile finale attrape toutes les versions : 1password*", "S_SKIP_OPEN": "Programmes ouverts en ce moment", "S_SKIP_REFRESH": "Actualiser la liste", "S_SKIP_PICKED": "%d sur %d choisis", "S_SKIP_NONE": "Rien de choisi", "S_SKIP_EMPTY": "La liste est vide — l'historique est tenu pour tous les programmes", "S_SKIP_ADD_BTN": "Ajouter un programme", "S_SKIP_HINT": "Ce que vous dictez dans ces programmes n'arrive jamais dans l'historique. Le collage lui-même fonctionne normalement.",
		"S_HIST_LIST": "Entrées", "S_HIST_CLEAR": "Vider", "S_HIST_TILL": "jusqu'au %s", "S_HIST_TILL1": "jusqu'à demain", "S_HIST_TILL_FULL": "Sera supprimé le %s — conservation %s", "S_HIST_LIST_HINT": "Ce qui a été dicté : copier, coller dans n'importe quelle fenêtre ou supprimer.", "S_HIST_COPY": "Copier",
		"S_HIST_KEEP": "Durée de conservation",
		"S_UNIT_MIN": "minutes",
		"S_UNIT_HOUR": "heures",
		"S_UNIT_DAY": "jours",
		"S_HIST_FIND": "Chercher dans l'historique…", "S_HIST_EMPTY": "Pas encore d'historique", "S_HIST_ASK": "Supprimer tout l'historique des dictées ?",
		"S_SEC_CMD": "Commandes vocales", "S_CMD_HINT": "Ce que vous dites devient un saut de ligne, un signe, une annulation, l'effacement de la dernière insertion ou une mise en forme de la phrase au lieu d'atterrir dans le texte. Reconnues en mots entiers, appliquées de haut en bas, après les remplacements ; les commandes de mise en forme agissent à la fin de la dictée.",
		"S_CMD_ADD": "Ajouter une commande", "S_CMD_PRESET": "Ajouter les habituelles", "S_CMD_PH": "la phrase que vous direz",
		"S_CMD_NEWLINE": "saut de ligne", "S_CMD_PARAGRAPH": "nouveau paragraphe", "S_CMD_TEXT": "insérer du texte", "S_CMD_CANCEL": "annuler la dictée",
		"S_CMD_TEXT_PH": "quoi insérer", "S_CMD_EMPTY": "Aucune commande pour l'instant", "S_CMD_DEL": "Supprimer la commande",
		"S_CMD_P_NEWLINE": "nouvelle ligne", "S_CMD_P_PARAGRAPH": "nouveau paragraphe", "S_CMD_P_CANCEL": "annuler",
		"S_CMD_UNDO": "effacer la dernière insertion", "S_CMD_HEADING": "mettre en titre", "S_CMD_LIST": "mettre en liste",
		"S_CMD_CAPITAL": "majuscule initiale", "S_CMD_QUOTES": "mettre entre guillemets", "S_CMD_NOSTOP": "retirer le point final",
		"S_CMD_HEADING_PH": "signe avant le titre, p. ex. #", "S_CMD_LIST_PH": "mots qui joignent le dernier élément, p. ex. et", "S_CMD_QUOTES_PH": "signes d'ouverture et de fermeture, p. ex. « »",
		"S_CMD_END_NOTE": "À dire à la fin de la dictée ; dite seule, elle met en forme la suivante.",
		"S_CMD_UNDO_NOTE": "À dire seule. Efface la dernière insertion si vous n'avez ni tapé ni cliqué depuis et que la fenêtre est la même ; le texte déjà présent n'est jamais touché.",
		"S_CMD_P_UNDO": "efface ça", "S_CMD_P_HEADING": "en titre", "S_CMD_P_LIST": "en liste", "S_CMD_P_CAPITAL": "avec majuscule",
		"S_CMD_P_QUOTES": "entre guillemets", "S_CMD_P_NOSTOP": "sans point", "S_CMD_P_AND": "et", "S_CMD_P_MARKS": "«  »",
		"S_SEC_REPLACE": "Remplacements après la reconnaissance", "S_REPLACE_HINT": "Ce qui a été mal entendu devient ce que vous vouliez dire — juste après la reconnaissance, avant les prompts. Appliqués de haut en bas.",
		"S_REPL_WHOLE_FULL": "Mots entiers uniquement", "S_REPL_CASE_FULL": "Respecter la casse", "S_CMD_ACTION": "Action",
		"S_FM_ADD": "Ajouter",
		"S_TIP_REPL_LANG": "La règle ne s’applique que lorsque vous dictez dans la langue choisie. « toutes les langues » — elle s’applique toujours.",
		"S_TIP_REPL_CASE": "Majuscules et minuscules comptent : « git » et « Git » sont des mots différents. Désactivé — la casse est ignorée.",
		"S_TIP_REPL_WHOLE": "Le remplacement ne s’applique que si le texte forme un mot à part. Désactivé — il agit aussi à l’intérieur d’autres mots.",
		"S_TIP_CMD_ACTION": "Ce qui se passe quand vous prononcez la phrase : saut de ligne, nouveau paragraphe, votre texte, annulation de la dictée, effacement de la dernière insertion ou mise en forme de la phrase — titre, liste, majuscule initiale, guillemets ou sans point final.",
		"S_LIST_FILTER_PH": "rechercher…",
		"S_REPL_DEL": "Supprimer le remplacement",
		"S_LIST_NOTHING": "Rien trouvé : « %s »",
		"S_FM_T_REPL_ADD": "Ajout d'un remplacement", "S_FM_T_REPL_EDIT": "Modification du remplacement",
		"S_FM_T_CMD_ADD": "Ajout d'une commande", "S_FM_T_CMD_EDIT": "Modification de la commande",
		"S_MT_DEL": "Suppression du modèle", "S_MT_DEL_PROMPT": "Suppression du prompt", "S_MT_DL": "Téléchargement du modèle",
		"S_MT_TR_OFF": "Désactivation de la traduction", "S_MT_TR_ONE": "Traduction sans demander", "S_MT_TR_LOCK": "Langue de sortie par défaut",
		"S_MT_REMOTE": "Serveur distant", "S_MT_POST": "Serveur externe", "S_MT_HIST": "Effacement de l'historique",
		"S_MT_RESET": "Réinitialisation des réglages", "S_MT_EXE": "Chemin du serveur",
		"S_DICT_ADD": "Ajouter un mot", "S_FM_T_DICT_ADD": "Ajout d'un mot", "S_DICT_EMPTY": "Pas encore de mots",
		"S_DICT_ADD_PH": "un mot, ou plusieurs séparés par des virgules",
		"S_DICT_NOMODEL": "Le modèle actuel %s ne prend pas en charge le dictionnaire — seuls les modèles Whisper le lisent.",
		"S_OV_FREE": "Place à vous", "S_OV_FREE_SUB": "le bandeau se déplace où vous voulez",
		"S_OVPOS_DRAG_SUB": "faites glisser le bandeau à la souris — il va où vous voulez",
		"S_OVMON_N": "Écran %d",
		"S_POST_ENABLE": "Activer le post-traitement",
		"S_API_SUM_URL": "adresse", "S_API_SUM_MODEL": "modèle", "S_API_SUM_KEY": "clé", "S_API_SUM_TIMEOUT": "attente",
		"S_API_SUM_STATE": "état", "S_API_NO_MODEL": "non indiqué",
		"S_API_NONE": "non configuré — le post-traitement reste local",
		"S_POSTAPI_SETUP": "Configurer", "S_API_EDIT": "Modifier", "S_API_KEY_DEL": "Supprimer la clé", "S_API_DLG": "Serveur externe",
		"S_LLM_CATALOG": "Catalogue de modèles", "S_LLM_BLOCK": "Modèles installés", "S_LLM_NONE_HINT": "Aucun modèle installé — téléchargez-en un avec la flèche, il apparaîtra ici", "S_LLM_IN_MEM": "en mémoire", "S_LLM_ON_DISK": "sur le disque", "S_LLM_EJECT": "Décharger de la mémoire", "S_LLM_FOUND": "%d trouvés", "S_LLM_NOSEARCH": "pas encore de recherche", "S_LLM_SEARCH_HINT": "Saisissez un nom de modèle et appuyez sur « Chercher »", "S_LLM_PICK_WAIT": "Disponible une fois le modèle téléchargé", "S_LLM_INSTALLED": "installés",
		"S_LLM_SUM_MODEL": "modèle", "S_LLM_SUM_SIZE": "taille", "S_LLM_SUM_COUNT": "installés", "S_LLM_SUM_RAM": "mémoire",
		"S_DLG_CLOSE": "Fermer", "S_LLM_NOPICK": "non choisi", "S_NO_PROMPTS": "Pas encore de prompts", "S_PROF_DRAG": "faites glisser pour réordonner",
		"S_PROF_NAME_PH": "nom du prompt", "S_PROF_TEST_PH": "écrivez une phrase pour essayer",
		"S_PF_NEW": "Nouveau prompt", "S_PF_EDIT": "Modification du prompt",
		"S_POST_NO_MODEL": "activé, mais aucun modèle choisi", "S_POST_NO_API": "activé, mais le serveur n'est pas configuré", "S_POST_BAD": "le serveur n'a pas répondu : %s", "S_POST_NO_PROMPT": "activé, mais aucun prompt coché", "S_API_TEST": "Tester la connexion", "S_API_TEST_RUN": "Vérification…", "S_API_TEST_OK": "Le serveur a répondu", "S_API_CLEAR": "Effacer", "S_API_CLEAR_ASK": "Supprimer l'adresse, le modèle et la clé du serveur externe ? Le post-traitement revient au modèle local.", "S_RAM_AVAIL": "Mémoire disponible : %s Go sur %s Go", "S_RAM_OF": "%s Go sur %s Go",
		"S_REPL_ADD": "Ajouter un remplacement", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "mots entiers", "S_REPL_CASE": "casse", "S_REPL_EMPTY": "Aucun remplacement pour l'instant",
		"S_PASTE_DELAY": "Délai avant l'insertion", "S_PASTE_DELAY_SUB": "quand le programme n'est pas encore prêt",
		"S_OVPOS": "Où afficher le bandeau", "S_OVPOS_SUB": "au curseur — près de l'endroit où vous tapez ; si l'application ne le montre pas, près du pointeur",
		"S_OVPOS_CARET": "Au curseur",
		"S_OVTEXT": "Afficher le texte reconnu", "S_OVTEXT_SUB": "sur le bandeau après l'insertion, au lieu du nombre de caractères",
		"S_OVERLAY": "Afficher le bandeau", "S_OVERLAY_SUB": "pendant la dictée, l’écran montre que l’enregistrement est en cours", "S_TYPEMODE": "Saisie caractère par caractère",
		"S_RECLANG": "Langue de la parole", "S_RECAUTO": "Auto",
		"S_LANG_KEY": "Raccourci propre", "S_LANG_KEY_SET": "Attribuer…", "S_LANG_KEY_CLEAR": "Retirer le raccourci", "S_LANG_KEY_HINT": "Maintenez-le pour dicter dans cette langue avec son modèle. Le raccourci principal garde la langue de la parole choisie plus haut.",
		"S_DL": "Télécharger", "S_DEL": "Supprimer",
		"S_M_BASE": "rapide, PC modestes", "S_M_SMALL": "équilibré", "S_M_MED": "Whisper pour toutes les langues ; celui qui traduit", "S_M_TURBO": "précision max sur CPU", "S_M_GIGAAM": "plus précis en russe, ponctue de lui-même", "S_M_PARAKEET": "25 langues européennes, ponctue d'elle-même",
		"S_THREADS":  "Threads CPU", "S_MINMS": "Enreg. min, ms", "S_MAXSEC": "Enreg. max, s",
		"S_AUTOSTART": "Démarrer whisper-server automatiquement", "S_PORT": "Port", "S_SERVEREXE": "Chemin de whisper-server", "S_SERVEREXE_SUB": "rempli automatiquement ; à changer seulement si le serveur est ailleurs", "S_EXE_RESET": "Réinitialiser", "S_EXE_WARN": "L'application trouve whisper-server à côté d'elle. Avec un chemin écrit à la main, déplacer le dossier empêchera la reconnaissance de démarrer. Modifier ?", "S_RESET_ALL": "Réinitialiser les réglages", "S_RESET_ALL_SUB": "tout revient d'usine, sauf les modèles et l'historique", "S_RESET_ALL_BTN": "Réinitialiser", "S_RESET_ALL_ASK": "Remettre tous les réglages d'usine ? Les modèles, l'historique et les prompts restent en place.", "S_RELOAD_CFG": "Relire config.json", "S_RELOAD_CFG_SUB": "si vous avez modifié le fichier à la main", "S_RELOAD_CFG_BTN": "Relire", "S_UPD_FOUND": "La version %s est sortie", "S_THEME": "Couleur", "S_THEME_SUB": "la couleur de la fenêtre, du bandeau et de l'icône", "S_THEME_GREEN": "Vert", "S_THEME_AMBER": "Ambre", "S_THEME_BLUE": "Bleu", "S_THEME_PINK": "Rose", "S_THEME_EDITOR": "Éditeur", "S_THEME_NEON": "Néon", "S_WND_MAX": "Occuper l'écran", "S_WND_RESTORE": "Taille précédente", "S_WND_MIN": "Réduire dans la zone de notification", "S_WND_CLOSE": "Fermer la fenêtre", "S_SKIN": "Design", "S_SKIN_SUB": "police, formes, effets et animations", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Doux", "S_SKIN_PAPER": "Document", "S_SKIN_FLUENT": "Fluent", "S_SKIN_STUDIO": "Studio",
		"S_SERVERURL": "Serveur externe (URL)", "S_URLHINT": "Si défini, le serveur local ne démarre pas",
		"S_STT_SRV": "Serveur de reconnaissance",
		"S_STT_SRV_HINT": "Les modèles Whisper tournent dans un programme séparé. Il peut fonctionner sur cet ordinateur ou sur un autre — choisissez lequel utiliser.",
		"S_SRV_LOCAL": "Sur cet ordinateur",
		"S_SRV_REMOTE": "Sur un autre ordinateur",
		"S_SRV_REMOTE_HINT": "Le même whisper-server, lancé ailleurs : serveur domestique, machine avec carte graphique, ordinateur voisin.",
		"S_SRV_K_AUTO": "démarrage auto",
		"S_SRV_K_FILE": "fichier",
		"S_SRV_K_ADDR": "adresse",
		"S_SRV_K_CHECK": "dernier test",
		"S_SRV_NEAR": "whisper-server.exe à côté de l’application",
		"S_SRV_NOADDR": "non défini",
		"S_SRV_NOCHECK": "jamais testé",
		"S_SRV_LOCAL_DLG": "Serveur de reconnaissance local",
		"S_SRV_ADDR": "Adresse du serveur",
		"S_SRV_ADDR_SUB": "adresse du whisper-server sur l’autre machine, port compris",
		"S_SRV_ON": "activé",
		"S_SRV_OFF": "désactivé",
		"S_SRV_K_THREADS": "threads CPU",
		"S_SRV_K_PORT": "port",
		"S_SRV_DOWN": "Reconnaissance indisponible",
		"S_SRV_DOWN_WHY": "le serveur de reconnaissance distant n’est pas configuré — renseignez son adresse dans les réglages",
		"S_SRV_DOWN_GO": "Ouvrir les réglages du serveur",
		"S_SRV_WARN_NOW": "La dictée ne fonctionne pas pour l’instant : le serveur distant est choisi, mais son adresse n’est pas renseignée.",
		"S_SRV_WARN_LATER": "Dès qu’un modèle Whisper sera choisi, la reconnaissance ne fonctionnera pas : l’adresse du serveur distant n’est pas renseignée.",
		"S_SAVED":      "Enregistré",
		"S_ABOUT_HTML": "<p><b>Voix → texte à la position du curseur.</b></p><p>Placez le curseur, maintenez le raccourci, parlez, relâchez — le texte s'insère.</p><p>Entièrement local et hors ligne. Technologies : <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b> ; modèles depuis Hugging Face.</p><p>Les logs ne dépassent jamais ~2 Mo.</p>",
		  "S_SEARCH": "Trouver un réglage…",
		"S_GRP_GENERAL": "Général", "S_GRP_SPEECH": "Traitement de la parole", "S_GRP_INFO": "Infos", "S_NAV_POST": "Post-traitement", "S_NAV_HELP": "Aide", "S_NAV_CONTACTS": "Contacts", "S_HIST_ADD": "Ajouter", "S_CONTACT_MAIL": "E-mail", "S_DICT_MODEL": "Modèle de reconnaissance", "S_LIB_ACC": "précision", "S_LIB_SPD": "vitesse",
		"S_HELP_TOC": "Sur cette page",
		"S_HELP_TOC_SHOW": "Afficher le sommaire — la fenêtre s’élargit",
		"S_HELP_TOC_HIDE": "Masquer le sommaire et rétablir la largeur",
		"S_CONTACT_TITLE": "Nous écrire",
		"S_ABOUT_DEPS": "Modules externes",
		"S_ABOUT_DEPS_HINT": "Du code tiers intégré à l’application, avec ses licences. Un clic sur le nom ouvre la page du projet.",
		"S_DEP_WHISPER": "exécute les modèles Whisper",
		"S_DEP_LLAMA": "post-traitement, modèles GGUF, traduction des sous-titres sur la carte graphique (llama-server-gpu.exe)",
		"S_DEP_VULKAN": "en-têtes et bibliothèque d'import du moteur pour la carte graphique ; l'environnement Vulkan lui-même vient avec le pilote graphique",
		"S_DEP_SHERPA": "moteur de GigaAM, Parakeet, Canary, Qwen, Moonshine",
		"S_DEP_GGML": "bibliothèque de tenseurs dans whisper.cpp et llama.cpp",
		"S_DEP_ONNX": "exécute les modèles dans sherpa-onnx",
		"S_DEP_WEBVIEW": "fenêtre des réglages sur WebView2",
		"S_DEP_WV2RT": "composant Windows qui dessine cette fenêtre",
		"S_DEP_MALGO": "capture du micro",
		"S_DEP_MINIAUDIO": "couche audio dans malgo",
		"S_DEP_WS": "liaison avec sherpa-server",
		"S_DEP_XSYS": "appels WinAPI depuis Go",
		"S_DEP_WINLOADER": "chargement des DLL dans go-webview2",
		"S_DEP_PLEX": "police de l’interface",
		"S_DEP_HF": "catalogue d’où viennent les modèles",
		"S_CONTACT_HINT": "Un bug, une idée, une question sur un réglage — écrivez un e-mail si c’est personnel, ou ouvrez une issue si c’est un bug.",
		"S_CONTACT_REPO": "Dépôt",
		"S_CONTACT_ISSUES": "Bugs et idées",
		"S_CONTACT_WRITE": "Écrire un e-mail",
		"S_CONTACT_OPEN": "Ouvrir",
		"S_STATE_ACTIVE": "Reconnaît",
		"S_STATE_USED": "Modèles en service",
		"S_STATE_INST": "Installés localement",
		"S_STATE_INST_SUB": "modèles sur le disque, prêts à être affectés",
		"S_PRESETS": "Quel modèle pour quelle langue",
		"S_PRESETS_HINT": "Le modèle de la ligne « Auto » sert toutes les langues qui n’ont pas le leur — elles sont en gris dans la liste. Cliquez sur une langue pour lui en donner un. Ce même modèle reconnaît aussi la langue parlée quand aucune langue source n’est choisie.",
		"S_MFOLDER": "Votre propre modèle",
		"S_DICT_SAVE": "Enregistrer",
		"S_OWNM_SUB": "Ajoutez un modèle local de reconnaissance vocale",
		"S_OWNM_ONEFILE": "Un fichier",
		"S_OWNM_FOLDERF": "Dossier avec les fichiers du modèle",
		"S_OWNM_S1": "Ouvrez le dossier des modèles",
		"S_OWNM_S1S": "Dossier de destination :",
		"S_OWNM_S2": "Copiez le modèle",
		"S_OWNM_S2S": "Choisissez l’une des structures prises en charge",
		"S_OWNM_S3": "Redémarrez l’application",
		"S_OWNM_S3S": "Le modèle apparaît pour les langues qu’il prend en charge",
		"S_AS_AUTO": "comme la détection automatique",
		"S_REC_CHIP": "recommandé",
		"S_BACK_AUTO": "Revenir à la détection automatique",
		"S_LANGS_COUNT": "langues : %d",
		"S_LANGS_UNKNOWN": "langues : inconnues",
		"S_TR_EN": "traduit vers %d langues",
		"S_TR_LIST": "traduit : %s",
		"S_DL_GOING": "téléchargement :",
		"S_OPEN_FOLDER": "Ouvrir le dossier",
		"S_UNLOAD": "Décharger de la mémoire",
		"S_UNLOAD_SUB": "libère la mémoire ; la prochaine dictée recharge le modèle",
		"S_UNLOAD_GO": "Décharger",
		"S_UNLOADED": "Déchargé",
		"S_NOT_FOR_LANG": "%s ne reconnaît pas cette langue",
		"S_MANUAL_NOTE": "Téléchargement impossible depuis l’application — la licence interdit la redistribution. Téléchargez l’archive vous-même et décompressez-la dans models/moonshine-uk.",
		"S_MANUAL_LINK": "Télécharger soi-même",
		"S_HF_FIT": "seulement ceux qui conviennent à cet ordinateur",
		"S_HF_HIDDEN": "masqués : %s",
		"S_WIZ_SKIP_DL": "Télécharger plus tard",
		"S_WIZ_SKIP_NOTE": "Sans modèle, la dictée ne fonctionnera pas. Téléchargez-en un dans « Langues et modèles ».",
		"S_M_GIGAAM2": "la génération précédente du modèle russe : même vitesse, mais sans ponctuation",
		"S_M_MOONUK": "modèle ukrainien Moonshine : rapide et léger, sans ponctuation",
		"S_M_LOCAL": "trouvé dans le dossier models ; propriétés inconnues, donc pas de barres",
		"S_ALL_LANGS": "toutes les langues",
		"S_OVPOS_SCHEME_SUB": "cliquez sur l’écran — le bandeau s’y place",
		"S_OVDRAG": "Glissez-le où vous voulez",
		"S_OVMON": "Écran",
		"S_OVMON_SUB": "sur quel moniteur afficher le bandeau",
		"S_OVMON_CURSOR": "L’écran avec le curseur",
		"S_M_NEMOTRON": "écrit pendant que vous parlez : le texte apparaît sur le bandeau en direct ; 40 langues, ponctue d’elle-même",
		"S_M_TINY": "le plus petit et le plus rapide, pour les machines très modestes ; nettement moins précis",
		"S_STATE_LOADED": "En mémoire en ce moment",
		"S_STATE_LOADED_SUB": "les modèles se déchargent seuls après inactivité",
		"S_STATE_WEEK": "Cette semaine",
		"S_ST_SUMMARY": "Résumé", "S_ST_OVERLAY": "Bandeau à l'écran", "S_ST_BEEP": "Signal sonore", "S_ST_AUTORUN": "Démarrage avec Windows", "S_ST_POST": "Post-traitement", "S_ST_LOCAL": "en local", "S_ST_CHECKED": "vérifié %s", "S_ST_GB": "%s Go", "S_ST_ON_M": "activé", "S_ST_OFF_M": "désactivé", "S_ST_MIC_OK": "le signal est bon", "S_ST_MIC_BAD": "le micro ne répond pas", "S_ST_CHECK": "Vérifier", "S_ST_RECOG": "reconnu par %s", "S_ST_VER": "Version %s", "S_ST_LATEST": "la plus récente", "S_ST_OUTDATED": "obsolète", "S_ST_UPD_OK": "vous avez la dernière version", "S_ST_UPD_DL": "Téléchargement de la mise à jour…",
		"S_ST_QUICK": "Réglages rapides",
		"S_ST_MODELS": "Modèles",
		"S_ST_USAGE": "Cette semaine",
		"S_ST_READY": "Prêt à dicter",
		"S_ST_OFF": "Désactivé dans la zone de notification",
		"S_ST_OFF_SUB": "le raccourci ne fait rien tant que vous ne le réactivez pas",
		"S_ST_ENABLE": "Activer",
		"S_ST_GOTO": "Ouvrir ce réglage sur son onglet",
		"S_ST_HOTKEY_GO": "Changer le raccourci",
		"S_ST_UPD_LAST": "Version %s — la plus récente",
		"S_ST_UPD_HAVE": "La version %s est disponible",
		"S_ST_MEM": "%s Go libres sur %s",
		"S_ST_MEM_SUB": "en mémoire : %s · sur le disque : %d modèles, %s Go",
		"S_ST_MEM_NONE": "rien",
		"S_ST_LANG": "Langue parlée",
		"S_ST_ASR": "Reconnaissance",
		"S_ST_ON": "activé",
		"S_ST_OFF_W": "désactivé",
		"S_ST_ON_F": "activée",
		"S_ST_OFF_F": "désactivée",
		"S_ST_ACTIVE": "active",
		"S_ST_IDLE": "non lancée",
		"S_ST_DISK": "sur le disque, %s",
		"S_ST_USAGE_SUB": "%d caractères · %d aujourd'hui · %d caractères en moyenne",
		"S_WEEK_OTHER": "autres",
		"S_ST_NO_WEEK": "aucune dictée cette semaine",
		"S_ST_AUTORUN_SUB": "l'application ne démarrera pas toute seule",
		"S_ST_OVERLAY_SUB": "visible pendant l'enregistrement",
		"S_REPL_LANG": "Langue de la règle",
		"S_REPL_LANG_ALL": "toutes les langues",
		"S_M_CANARY": "anglais, allemand, espagnol, français — et traduit entre eux de lui-même",
		"S_M_QWEN3": "environ 30 langues, ponctue de lui-même ; le plus lourd et le plus précis du catalogue",
		"S_POSTAPI": "Serveur externe",
		"S_POST_HINT": "Corrige le texte reconnu selon vos prompts : retire les mots parasites, répare la ponctuation, change le style. Désactivé — le texte est inséré tel quel.",
		"S_POST_MODEL": "Modèle",
		"S_SRC_LOCAL": "Locale",
		"S_SRC_USED": "utilisée",
		"S_HF_GO": "Rechercher",
		"S_POSTAPI_HINT": "Vide par défaut — tout le post-traitement est local. Saisissez une adresse et les prompts s’exécutent sur un serveur externe : OpenAI, Groq, votre vLLM — tout ce qui a une API compatible.",
		"S_POSTAPI_URL": "Adresse",
		"S_POSTAPI_URL_SUB": "vide = le modèle local ; exemple : https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Modèle",
		"S_POSTAPI_KEY": "Clé API",
		"S_POSTAPI_KEY_SET": "clé enregistrée (chiffrée par Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "pas de clé",
		"S_POSTAPI_SAVE": "Enregistrer la clé",
		"S_POSTAPI_TIMEOUT": "Délai de réponse", "S_SEC_SHORT": "s",
		"S_POSTAPI_WARN": "⚠ Le texte reconnu de vos dictées partira vers cette adresse. Le son ne part jamais. La clé est stockée chiffrée.",
		"S_POSTAPI_ASK": "Envoyer le texte reconnu à %s ? Le son reste sur cet ordinateur, mais le texte le quittera.",
		"S_POSTAPI_BADGE": "serveur externe",
		"S_NOT_INSTALLED": "non installé",
		"S_NAV_STATE": "État", "S_NAV_DICT": "Commandes et comportement", "S_NAV_MIC": "Microphone", "S_NAV_MODELS": "Langues et modèles",
		"S_NAV_TEXT": "Règles", "S_NAV_TR": "Traduction", "S_NAV_SYSTEM": "Système", "S_NAV_ABOUT": "À propos",
		"S_STATE_HINT": "maintenez et parlez — le texte arrive là où se trouve le curseur",
		"S_STATE_PROC": "Post-traitement",
		"S_CHANGE_MODEL": "Changer", "S_PICK_MODEL": "Choisir", "S_STATE_GET": "Télécharger",
		"S_RETRY": "Réessayer", "S_BERR_OPEN": "Ouvrir les réglages du serveur",
		"S_STATE_LAST": "Dernière dictée", "S_STATE_COPY": "Copier", "S_STATE_MEM": "Mémoire",
		"S_STATE_MEM_SUB": "les modèles restent chargés, la première phrase part sans délai",
		"S_HOTMODE":       "Mode", "S_HOTMODE_HOLD": "maintenir", "S_HOTMODE_TOGGLE": "bascule",
		"S_SUB_HOTMODE": "maintenez les touches, ou appuyez une fois pour lancer et une fois pour arrêter",
		"S_SILENCE":     "Arrêter sur un silence", "S_SILENCE_MS": "Pause avant l'arrêt, ms",
		"S_SUB_SILENCE": "mode « bascule » : une fois que vous avez parlé, une pause termine l'enregistrement comme si vous appuyiez de nouveau sur les touches",
		"S_SUB_SILENCE_MS": "durée de silence tolérée ; l'enregistrement ne s'arrête jamais avant vos premiers mots",
		"S_SUB_MINMS":   "ignore les appuis accidentels",
		"S_SUB_ENTER":   "envoie le message aussitôt",
		"S_SUB_CLIP":    "les images et les fichiers reviennent tels quels",
		"S_SUB_TYPE":    "utile là où un champ refuse le collage",
		"S_SEC_OVERLAY": "Bandeau à l'écran",
		"S_MIC_CHECK":   "Vérifier le microphone", "S_MIC_CHECK_SUB": "trois secondes d'enregistrement, puis un verdict : niveau, saturation, présence de parole", "S_MIC_CHECKING": "Vérification…",
		"S_MCHECK": "Vérifier les modèles installés", "S_MCHECK_SUB": "compare les fichiers des modèles aux empreintes de référence", "S_MCHECK_GO": "Vérifier", "S_MCHECK_RUN": "Vérification…",
		"S_HIST_INSERT": "Coller",
		"S_MIC": "Microphone", "S_MIC_DEFAULT": "Par défaut du système", "S_MIC_REFRESH": "Actualiser la liste",
		"S_MIC_LEVEL": "Niveau d'entrée", "S_MIC_QUIET": "silence",
		"S_SUB_THREADS": "plus de threads n'est pas toujours plus rapide — mesurez sur votre machine",
		"S_SEC_LLM":     "Modèle d'édition",
		"S_PUNCT":       "Ponctuation et majuscules", "S_SUB_PUNCT": "d'où viennent la ponctuation et les majuscules",
		"S_PUNCT_MODEL": "du modèle", "S_PUNCT_LLM": "par le modèle d'édition", "S_PUNCT_OFF": "retirer",
		"S_SUB_DICT": "Dictionnaire", "S_SUB_PROMPTS": "Prompts",
		"S_SUB_TRTARGET": "le texte y est traduit ; le dialogue sur le bandeau la propose en premier",
		"S_REMOTE_ABOUT": "Un serveur distant est configuré : l'audio y est envoyé, et la promesse ci-dessus ne tient pas tant qu'il est actif.",
		"S_UPD":          "Mises à jour", "S_UPD_CHECK": "Rechercher des mises à jour", "S_UPD_AUTO": "Vérifier au démarrage",
		"S_SUB_UPD":  "la seule requête réseau en dehors du téléchargement des modèles",
		"S_UPD_NONE": "Vous avez la dernière version", "S_BADGE_MODELS": "Modèles installés", "S_BADGE_MISS": "Un modèle n'est pas téléchargé", "S_BADGE_SYSTEM": "Avertissements à examiner", "S_BADGE_HIST": "Entrées dans l'historique", "S_LOG_OPEN": "Ouvrir le journal", "S_LOG": "Journal", "S_LOG_SUB": "tout ce que l'application écrit sur elle-même", "S_UPD_AVAIL": "La version %s est disponible.",
		"S_UPD_GO": "Mettre à jour", "S_UPD_ERR": "Échec de la vérification", "S_UPD_DL": "Téléchargement de la mise à jour…",
		"S_SEC_SERVICE": "Service", "S_SUB_AUTOSTART": "désactivez si vous lancez le serveur vous-même",
		"S_SUB_PORT":    "la reconnaissance redémarre toute seule",
		"S_MODEL_READY": "Modèle téléchargé — choisissez-le pour basculer",
		"S_FIT_OK":      "tient", "S_FIT_WARN": "juste", "S_FIT_BAD": "mémoire insuffisante", "S_RAM": "Mémoire de l'ordinateur :",
		"S_HF_PH":       "Nom du modèle — par ex. qwen2.5 instruct",
		"S_NO_LLM":      "Aucun modèle installé pour l'instant — trouvez-en un dans le champ de recherche ci-dessous.",
		"S_NO_LLM_PROF": "Les prompts deviennent disponibles dès qu’un modèle est installé — le bloc « Modèle » en haut de cet onglet.",
		"S_UPDATED":     "Dernière mise à jour du modèle", "S_PROF_EDIT": "Modifier", "S_PROF_CLOSE": "Replier",
		"S_CONFIRM_DEL": "Supprimer le modèle « %s » ? Il pourra être téléchargé à nouveau.", "S_FREE": "libre",
		"S_DEL_ACTIVE":     "Supprimer le modèle actif « %s » ? La reconnaissance s'arrête jusqu'à ce que vous en choisissiez un autre — vous pouvez le télécharger ici même.",
		"S_REMOTE_WARN":    "L'audio sera envoyé à ce serveur. Le mode local est désactivé.",
		"S_REMOTE_ASK":     "L'audio ne sera plus traité sur cet ordinateur : il sera envoyé à %s. Activer le mode distant ?",
		"S_REMOTE_BADGE":   "DISTANT",
		"S_OK":             "Oui", "S_CANCEL": "Annuler", "S_DL_START": "Télécharger", "S_DL_CANCEL": "Annuler le téléchargement",
		"S_DL_ASK":    "Le modèle « %s » n'est pas téléchargé (%s). Commencer le téléchargement ?",
		"S_NOT_FOUND": "rien",   "S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Auteur et développeur de {app} — un outil de dictée local pour Windows : la voix devient du texte à l'endroit du curseur, sans nuage ni abonnement.</p>" +
			"<p>Le projet est ouvert : code source, chaîne de compilation et versions récentes sont sur GitHub.</p>" +
			"<ul>" +
			"<li>Dépôt : <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Profil de l'auteur : <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Un bug ou une idée — ouvrez une issue dans le dépôt.</p>",
		"S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Ce que c'est</p>" +
			"<p>Vous parlez, et les mots apparaissent là où clignote le curseur. Comme si vous les aviez tapés, en plus rapide.</p>" +
			"<p class=\"hdim\">Concrètement, cela se passe ainsi. Vous cliquez dans un champ où vous tapez d'habitude : un e-mail, un chat, une barre de recherche, un document. Vous maintenez deux touches et dites une phrase de votre voix normale. Vous relâchez. Une seconde plus tard, la phrase est dans le champ. Pas besoin de changer de fenêtre ni de copier quoi que ce soit : l'application fonctionne par-dessus n'importe quel programme et attend ses touches en arrière-plan.</p>" +
			"<p class=\"hdim\">Tout le calcul se fait sur votre ordinateur. Votre voix n'est pas enregistrée : elle reste quelques secondes en mémoire, devient du texte et disparaît. L'application ne va sur internet qu'en deux occasions : quand vous téléchargez un modèle de reconnaissance et quand elle vérifie si une nouvelle version est sortie. Cette vérification peut être désactivée, et l'application n'a alors plus du tout besoin du réseau.</p>" +
			"<p class=\"hdim\">Le mot « modèle » reviendra souvent plus bas. C'est simplement un fichier qui sait déchiffrer la parole. Il se télécharge une fois et reste dans un dossier à côté de l'application ; plus il est gros, mieux il entend et plus il est lent.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Ce qui se passe quand vous relâchez les touches</p>" +
			"<div class=\"hflow\">" +
			"<span>enregistrement</span><em>→</em><span>reconnaissance</span><em>→</em><span>remplacements</span><em>→</em><span>commandes</span><em>→</em><span>traduction</span><em>→</em><span>édition</span><em>→</em><span>insertion</span>" +
			"</div>" +
			"<p>En général, cela prend une ou deux secondes, et il est inutile de connaître les étapes. Mais si un jour le texte ne sort pas comme vous l'attendiez, cette chaîne permet de voir facilement quelle étape l'a modifié.</p>" +
			"<ul>" +
			"<li><b>Enregistrement.</b> Tant que les touches sont maintenues, l'application écoute le microphone. Elle ignore un appui trop court, pour qu'une touche effleurée par erreur ne lance pas de dictée.</li>" +
			"<li><b>Reconnaissance.</b> La voix devient des lettres. C'est le travail d'un modèle. Le russe est en général confié à GigaAM, le polonais à Parakeet, la plupart des autres langues à Nemotron, et Whisper prend en charge les langues restantes ainsi que la traduction ; ce ne sont que des noms, et vous n'avez pas à les choisir vous-même.</li>" +
			"<li><b>Remplacements.</b> Votre liste personnelle « entendu → voulu ». Si le modèle s'obstine à écrire « git hub » au lieu de « GitHub », vous le corrigez une fois, et ensuite la correction se fait toute seule.</li>" +
			"<li><b>Commandes.</b> Quelques phrases fonctionnent comme des boutons. Dites « nouvelle ligne », et le texte passe à la ligne, sans que ces mots eux-mêmes arrivent dans le champ.</li>" +
			"<li><b>Traduction.</b> Si vous l'avez activée, vous parlez une langue et le texte arrive dans une autre.</li>" +
			"<li><b>Édition.</b> Une étape facultative. Un deuxième modèle met le texte au propre : il retire les « euh » et les « en fait », place les virgules, donne à la formulation un ton professionnel.</li>" +
			"<li><b>Insertion.</b> Le texte arrive dans le champ. Ce qui se trouvait dans votre presse-papiers avant la dictée y est remis.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Pendant tout ce temps, un bandeau étroit en bas de l'écran affiche le nom de l'étape en cours. Il vous montre d'un coup d'œil que l'application travaille et qu'elle n'est pas bloquée.</p>" +
			"<p class=\"wh\" id=\"s-first\">Premier lancement</p>" +
			"<p>Au premier démarrage, l'application vous guide en trois étapes. Cela prend une minute, et ensuite vous pouvez dicter.</p>" +
			"<ul>" +
			"<li><b>Dans quelle langue vous allez dicter.</b> La question la plus importante : d'après votre réponse, l'application choisit elle-même un modèle adapté et le télécharge en arrière-plan, avec la progression affichée. Si cette langue diffère de celle de la fenêtre, une ligne apparaît juste là pour changer aussi la langue de l'interface.</li>" +
			"<li><b>Essai.</b> Le raccourci clavier (par défaut Ctrl et la touche Windows ; pour le changer, appuyez sur ce qui est confortable à maintenir d'une main), le microphone avec une barre de niveau en direct, et un champ où dicter une phrase d'essai. Si la barre ne bouge pas quand vous parlez, ce n'est pas le bon microphone qui est choisi. Le champ s'anime dès que le modèle est en place.</li>" +
			"<li><b>C'est prêt.</b> L'interrupteur « Démarrer avec Windows », et l'assistant est terminé.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Vous pouvez passer n'importe quelle étape et y revenir plus tard : tout cela se trouve aussi dans les réglages. Le modèle continue de se télécharger en arrière-plan, même si vous passez à la suite ou terminez l'assistant avant la fin. Si vous voulez refaire l'assistant, lancez l'application en ajoutant <code>-wizard</code>. Après une mise à jour, il ne réapparaît pas.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">Le bandeau en bas de l'écran</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Parlez…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">c'est ainsi que l'application montre qu'elle vous entend</div>" +
			"</div>" +
			"<p>Pendant la dictée, c'est la seule chose visible à l'écran. Le bandeau est petit, n'intercepte pas la souris et ne gêne pas votre travail : il vous dit simplement ce qui arrive à votre phrase en ce moment.</p>" +
			"<ul>" +
			"<li><b>« Parlez… »</b> L'enregistrement est en cours. Les barres à côté bougent au rythme de votre voix. Si elles restent immobiles, on ne vous entend pas, et il est inutile de continuer à parler.</li>" +
			"<li><b>« Reconnaissance… »</b> Vous avez relâché les touches, le modèle travaille. En général une ou deux secondes. La toute première dictée après le lancement est plus longue : c'est à ce moment que le modèle se charge en mémoire.</li>" +
			"<li><b>« Traduction… »</b> La traduction est en cours.</li>" +
			"<li><b>« Édition : Cleanup (1/2) »</b> Le deuxième modèle retouche le texte. Vous voyez quelle consigne est en train d'être appliquée et combien il y en a en tout.</li>" +
			"<li><b>« Inséré : 128 caractères »</b> C'est fait, le texte est déjà dans le champ. Le bandeau s'éteint tout seul au bout de quelques secondes.</li>" +
			"<li><b>La croix.</b> Vous avez changé d'avis ? Cliquez dessus à n'importe quelle étape. L'enregistrement est jeté, et rien n'arrive dans le champ.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Le bandeau peut être désactivé, placé dans n'importe quel coin, attaché au curseur ou simplement déplacé à la souris là où il ne gêne pas. Cela se règle dans l'onglet « Commandes et comportement ».</p>" +
			"<p class=\"wh\" id=\"s-state\">L'écran « État »</p>" +
			"<p>Cet écran répond à une seule question : pouvez-vous dicter tout de suite ? Sinon, il dit pourquoi et propose un bouton qui mène là où cela se répare.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Prêt à dicter</b><span>maintenez et parlez — le texte arrive là où se trouve le curseur</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">le carré vert signifie que l'application écoute ses touches</div>" +
			"</div>" +
			"<p>La carte du haut existe en trois versions. <b>Verte</b> : tout va bien, avec un rappel des touches à côté. <b>Jaune</b> : vous avez désactivé l'application depuis son icône dans la zone de notification, près de l'horloge, et les touches ne font rien tant que vous n'appuyez pas sur « Activer », directement sur la carte. <b>Rouge</b> : les réglages se contredisent. La reconnaissance sur un autre ordinateur est choisie, mais aucune adresse n'est indiquée.</p>" +
			"<p>En dessous se trouve ce qu'on regarde d'habitude en passant :</p>" +
			"<ul>" +
			"<li><b>Sous-titres.</b> La carte juste sous celle du haut : si les sous-titres sont arrêtés, chargent leur modèle, tournent ou redémarrent, avec le modèle et la façon dont arrive le texte, leur raccourci s'il y en a un, un bouton qui les démarre ou les arrête, et une flèche vers l'onglet « Sous-titres ».</li>" +
			"<li><b>Cette semaine.</b> Un graphique circulaire : vers quels programmes sont allées vos dictées et combien de caractères vous avez dictés.</li>" +
			"<li><b>Version.</b> Le numéro de build et un bouton « Vérifier ». Si une version plus récente est sortie, l'application le dit ici même.</li>" +
			"<li><b>Mémoire.</b> Combien il en reste de libre, quel modèle s'y trouve en ce moment et combien de place tous les modèles prennent sur le disque. Le bouton « Décharger » libère la mémoire immédiatement, par exemple avant un jeu gourmand.</li>" +
			"<li><b>Réglages rapides.</b> Cinq cartes avec l'essentiel : microphone, langue, bandeau, son, démarrage avec Windows. Elles montrent comment c'est réglé en ce moment, et la flèche mène à l'onglet où cela se change. Il n'y a volontairement pas d'interrupteurs ici, pour que rien ne bascule par mégarde.</li>" +
			"<li><b>Modèles.</b> Ce qui est installé, ce qui fonctionne en ce moment, de quelles langues chacun s'occupe et ce qui dort simplement sur le disque.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">L'onglet « Système »</p>" +
			"<p>L'apparence de l'application, ses rapports avec Windows et quelques boutons de service. Dans la vie de tous les jours, on vient ici à peu près une fois par an.</p>" +
			"<p class=\"hsub\">À quoi ressemble l'application</p>" +
			"<p>Trois lignes : langue de l'interface, design et couleur. Elles n'ont aucun effet sur le fonctionnement.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Langue de l'interface</td><td>La langue de la fenêtre des réglages et des textes du bandeau. Elle n'a rien à voir avec la langue de dictée : la fenêtre peut être en français pendant que vous dictez en anglais.</td></tr>" +
			"<tr><td>Design</td><td>L'allure générale. Il y en a sept : « Éditeur » (par défaut : gris foncé et bleu), « Fluent » (comme les Paramètres de Windows 11), « Studio » (graphite et orange, presque sans arrondis), « Doux » (lilas clair, arrondi), « Document » (papier et bleu), « Néon » (violet avec un dégradé du rose au turquoise) et « Terminal » (lettres vertes sur fond noir, crochets autour des boutons). Les sept utilisent la même police, IBM Plex Mono.</td></tr>" +
			"<tr><td>Couleur</td><td>Une teinte réservée à « Terminal » : vert, ambre, bleu, rose. Le bandeau et l'icône de la zone de notification prennent la même couleur. Les autres designs ont leurs propres couleurs, et la ligne est masquée.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Démarrage et mises à jour</p>" +
			"<p>C'est ici que vous décidez si l'application démarre avec le système et si elle guette les nouvelles versions.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Démarrer avec Windows</td><td>L'application se lance à l'ouverture de votre session et attend vos touches. C'est une entrée ordinaire dans la liste de démarrage, la même que celle que vous pourriez y ajouter vous-même.</td></tr>" +
			"<tr><td>Vérifier au démarrage</td><td>Une fois, au démarrage, l'application demande si une version plus récente est sortie. En dehors du téléchargement des modèles, elle ne va jamais sur internet. Désactivez cette option, et vous suivrez vous-même la sortie des nouvelles versions.</td></tr>" +
			"<tr><td>Mises à jour</td><td>Le bouton vérifie tout de suite. S'il existe une version plus récente, l'application la télécharge et l'installe par-dessus l'actuelle : réglages, modèles et historique restent en place.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Service</p>" +
			"<p>Trois boutons pour le cas où quelque chose s'est mal passé.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Journal</td><td>Ouvre le fichier où l'application note ce qu'elle a fait : quel modèle a démarré, combien de temps a pris la reconnaissance, où une erreur s'est produite. La cause est en général écrite en mots, pas en codes. Le fichier ne grossit pas indéfiniment : les anciennes lignes sont remplacées par les nouvelles.</td></tr>" +
			"<tr><td>Relire config.json</td><td>Tous les réglages tiennent dans un seul fichier texte à côté de l'application. Si vous l'avez modifié à la main ou rapporté d'un autre ordinateur, ce bouton prend en compte les changements sans redémarrage.</td></tr>" +
			"<tr><td>Réinitialiser les réglages</td><td>Remet tout comme après l'installation. Les modèles téléchargés et l'historique restent en place : seuls les réglages sont effacés.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Où tourne la reconnaissance</p>" +
			"<p>La parole est traitée par un petit programme à part. D'habitude, il se trouve à côté de l'application, démarre tout seul, et vous n'avez pas à vous en soucier. Mais il peut aussi être installé sur un autre ordinateur : par exemple sur la machine de la maison, qui a une carte graphique puissante, pendant que vous dictez depuis un portable léger. Le choix se fait avec le sélecteur sur les cartes.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>Sur cet ordinateur</b><span>threads CPU 4 · démarrage auto activé · port 8910 · whisper-server.exe à côté de l'application</span></div><button class=\"hbtn\">Modifier</button></div>" +
			"<div class=\"hcap\">la carte choisie est mise en évidence, l'autre est estompée pendant ce temps</div>" +
			"</div>" +
			"<p>Le bouton « Modifier » ouvre une fenêtre avec quatre réglages. Ils ont tous des valeurs par défaut, et vous pouvez presque toujours les laisser tels quels.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Threads CPU</td><td>Combien de cœurs du processeur donner à la reconnaissance. Plus de cœurs, c'est plus rapide, mais seulement jusqu'à un certain point : au-delà, on ne gagne plus rien, et les autres programmes commencent à ralentir.</td></tr>" +
			"<tr><td>Démarrer whisper-server automatiquement</td><td>L'application lance elle-même le moteur de reconnaissance. Il n'y a qu'une raison de désactiver cette option : vous le lancez à la main et ne voulez pas que l'application s'en mêle.</td></tr>" +
			"<tr><td>Port</td><td>Le numéro par lequel l'application parle au moteur de reconnaissance à l'intérieur de votre propre ordinateur. Il n'est pas visible de l'extérieur. Ne le changez que si un autre programme a pris ce numéro.</td></tr>" +
			"<tr><td>Chemin de whisper-server</td><td>L'endroit où se trouve le moteur de reconnaissance. Il se remplit tout seul. N'écrivez le vôtre que si vous l'avez délibérément placé ailleurs : sinon, dès que vous déplacerez le dossier, la dictée ne démarrera plus.</td></tr>" +
			"</table>" +
			"<p>La deuxième carte n'a qu'un réglage : l'adresse de l'ordinateur sur le réseau, par exemple <code>http://192.168.1.50:8080</code>. Sachez que dans ce cas votre voix enregistrée part vers cette machine, et qu'avec une adresse ordinaire elle circule sans chiffrement. Sur un réseau domestique, cela ne pose pas de problème ; par internet, mieux vaut ne pas travailler ainsi.</p>" +
			"<div class=\"hbad\">Si « Sur un autre ordinateur » est choisi sans qu'une adresse soit renseignée, et que votre langue est servie par Whisper, la dictée ne fonctionnera pas. L'application le signale par une carte rouge sur « État », une marque rouge à côté de « Système » et une ligne en bas de la fenêtre.</div>" +
			"<p class=\"wh\" id=\"s-mic\">L'onglet « Microphone »</p>" +
			"<p>Plus de la moitié des dictées ratées ne viennent ni du modèle ni des réglages, mais du microphone : ce n'est pas le bon qui est écouté, ou l'on vous entend à peine. Cet onglet réunit tout ce qu'il faut pour le vérifier en trente secondes.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Microphone</td><td>Quel appareil écouter. L'option « Par défaut du système » signifie « celui qui est choisi dans Windows », ce qui est pratique si vous travaillez tantôt avec un casque, tantôt sans : l'application suit le système toute seule.</td></tr>" +
			"<tr><td>Actualiser la liste</td><td>Utile si vous avez branché un casque après le lancement de l'application : la liste est relue.</td></tr>" +
			"<tr><td>Niveau d'entrée</td><td>Une barre en direct. Dites quelque chose : les barres doivent bouger nettement. Si elles ne bougent pas, le microphone est coupé dans Windows, ou ce n'est pas le bon qui est choisi.</td></tr>" +
			"<tr><td>Vérifier le microphone</td><td>Trois secondes d'enregistrement et un verdict en mots simples : le niveau est-il suffisant, y a-t-il de la saturation (quand le son « grésille » parce que le micro est trop près) et entend-on bien de la parole.</td></tr>" +
			"<tr><td>Signaux sonores</td><td>Un son bref au début et à la fin de l'enregistrement. Particulièrement utile si vous avez désactivé le bandeau : on entend au moins que l'enregistrement a commencé.</td></tr>" +
			"<tr><td>Son du signal</td><td>Quel son exactement : voix de Windows, clochette, marimba et quelques autres. Le bouton à côté joue le son choisi, pour que vous n'ayez pas à deviner d'après le nom.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">L'onglet « Historique »</p>" +
			"<p>L'historique est la liste de ce que vous avez dicté. Le texte seulement : la voix n'est jamais enregistrée, nulle part. Par défaut, il est désactivé, et c'est voulu : les dictées contiennent souvent des choses personnelles.</p>" +
			"<p>Il vaut la peine de l'activer si vous vous êtes déjà retrouvé dans l'une de ces situations : vous avez collé un texte dans la mauvaise fenêtre et il a été perdu ; vous devez répéter une longue phrase ; vous voulez retrouver ce que vous avez dicté hier et l'envoyer à nouveau.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Durée de conservation</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">jours</span></div>" +
			"<div class=\"hrow\"><label>Combien d'entrées conserver</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">la durée se règle en minutes, en heures ou en jours, avec n'importe quel nombre</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Conserver l'historique des dictées</td><td>L'interrupteur principal. Désactivé, plus aucune nouvelle entrée n'apparaît, et les anciennes sont effacées.</td></tr>" +
			"<tr><td>Durée de conservation</td><td>Au bout de combien de temps une entrée disparaît toute seule. Vous pouvez mettre cinq minutes si l'historique ne vous sert qu'au cas où un texte se colle au mauvais endroit, ou bien trente jours.</td></tr>" +
			"<tr><td>Combien d'entrées conserver</td><td>Une deuxième limite, en nombre. Elle agit avec la première : celle qui est atteinte d'abord l'emporte. Avec « 7 jours et pas plus de 50 », la cinquante et unième entrée chasse la plus ancienne, même si la semaine n'est pas encore écoulée.</td></tr>" +
			"<tr><td>Ne jamais enregistrer depuis ces programmes</td><td>Une liste d'exceptions. Tout ce que vous dictez dans ces fenêtres (une banque, un gestionnaire de mots de passe, un chat professionnel) reste hors de l'historique, tandis que la dictée elle-même fonctionne normalement. Vous pouvez choisir un programme dans la liste des fenêtres ouvertes, sans chercher le nom de son fichier.</td></tr>" +
			"</table>" +
			"<p>Avec chaque entrée, vous pouvez faire quatre choses. <b>Copier</b> : le texte va dans le presse-papiers. <b>Coller</b> : l'application vous demande de montrer où. Cliquez dans le champ voulu, dans n'importe quel programme, et le texte s'y place. <b>Supprimer</b> : retire cette seule entrée. <b>Chercher</b> : la recherche porte à la fois sur le texte et sur le nom du programme. Le bouton « Vider » efface tout d'un coup, c'est pourquoi il demande une confirmation.</p>" +
			"<p class=\"hdim\">À droite de chaque entrée est indiqué quand elle disparaîtra, par exemple « jusqu'au 7 sept. », et le dernier jour cette mention devient jaune.</p>" +
			"<p class=\"wh\" id=\"s-dict\">L'onglet « Commandes et comportement »</p>" +
			"<p>Comment la dictée démarre et comment le texte arrive dans le champ. On vient ici quand quelque chose ne marche pas comme on le voudrait : les touches sont peu pratiques, le texte n'apparaît pas, Entrée est pressée toute seule.</p>" +
			"<p class=\"hsub\">Touches</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Raccourci clavier</td><td>Appuyez sur « Changer… » et maintenez ce qui vous convient. L'application intercepte cette combinaison dans tout le système, donc n'en prenez pas une déjà utilisée, comme Ctrl+C. Les combinaisons avec la touche Windows sont les moins souvent prises. Une langue peut aussi avoir son propre raccourci, défini dans l’onglet « Langues et modèles ».</td></tr>" +
			"<tr><td>Mode</td><td>Avec « maintenir », l'enregistrement dure tant que les touches sont enfoncées : vous relâchez, et la phrase part au traitement. Avec « bascule », le premier appui lance l'enregistrement et le second l'arrête. Le second mode convient aux longs monologues, quand la main se fatigue de tenir les touches.</td></tr>" +
			"<tr><td>Arrêter sur un silence</td><td>Uniquement en mode « bascule ». Une fois que vous avez dit quelque chose, une pause de la durée réglée dans « Pause avant l'arrêt, ms » termine l'enregistrement d'elle-même, comme si vous aviez appuyé de nouveau sur le raccourci, et le bandeau affiche « Arrêté sur un silence ». L'enregistrement ne s'arrête jamais avant vos premiers mots : vous pouvez donc rassembler vos idées après avoir appuyé ; si vous ne dites rien du tout, la limite « Enreg. max, s » y met fin quand même. Avec « Appuyer Entrée après insertion (envoi auto) » activé, une pause envoie le message : choisissez alors une pause plus longue.</td></tr>" +
			"<tr><td>Enreg. min, ms</td><td>La durée en dessous de laquelle un appui compte comme accidentel. Effleurez les touches une fraction de seconde, et l'application reste muette au lieu d'insérer un bout de mot.</td></tr>" +
			"<tr><td>Enreg. max, s</td><td>Une sécurité contre une touche oubliée : l'enregistrement s'arrête tout seul au lieu d'enregistrer dans le vide pendant une demi-heure.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Comment le texte est inséré</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Appuyer Entrée après insertion (envoi auto)</td><td>Le texte n'est pas seulement inséré, il est envoyé aussitôt : le message part, la recherche se lance. Dans une messagerie, cela vous épargne un geste ; dans un e-mail ou un document, cela gênera presque à coup sûr.</td></tr>" +
			"<tr><td>Restaurer le presse-papiers après insertion</td><td>L'insertion passe par le presse-papiers, et l'application y remet ce qui s'y trouvait : une image, un fichier, un morceau de texte copié. Laissez cette option activée si vous copiez souvent des choses et ne voulez pas les perdre.</td></tr>" +
			"<tr><td>Saisie caractère par caractère</td><td>Au lieu de coller le texte d'un bloc, l'application le tape comme si vous le saisissiez vous-même. C'est plus lent, mais cela dépanne là où le collage depuis le presse-papiers est interdit : certains formulaires bancaires, des terminaux, de vieux programmes.</td></tr>" +
			"<tr><td>Délai avant l'insertion</td><td>Une courte pause avant de remettre le texte à la fenêtre. Rarement utile : il arrive qu'un programme n'ait pas le temps de recevoir le texte, et le début de la phrase se perd.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Bandeau à l'écran</p>" +
			"<p>C'est ici que vous décidez s'il faut afficher le bandeau, et où exactement.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Afficher le bandeau</td><td>Désactivez-le s'il vous gêne ou si les signaux sonores vous suffisent.</td></tr>" +
			"<tr><td>Où afficher le bandeau</td><td>Un schéma de l'écran avec neuf emplacements : cliquez sur le coin où le bandeau ne gênera pas.</td></tr>" +
			"<tr><td>Au curseur</td><td>Le bandeau apparaît près de l'endroit où vous dictez. Si le programme n'indique pas où se trouve le curseur, il apparaît près du pointeur de la souris.</td></tr>" +
			"<tr><td>Place à vous</td><td>Vous pouvez simplement faire glisser le bandeau à la souris, et il retient l'emplacement.</td></tr>" +
			"<tr><td>Écran</td><td>Sur quel moniteur l'afficher, si vous en avez plusieurs.</td></tr>" +
			"<tr><td>Afficher le texte reconnu</td><td>Après l'insertion, le bandeau montre la phrase elle-même plutôt que le nombre de caractères. Pratique pour vérifier, mais toute personne qui regarde votre écran verra aussi la phrase.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">L'onglet « Langues et modèles »</p>" +
			"<p>Ici, vous dites dans quelle langue vous allez dicter, et l'application décide avec quoi la reconnaître. Il y a plusieurs modèles : certains sont rapides, d'autres précis, d'autres encore connaissent des langues rares.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Langue de la parole</td><td>La langue dans laquelle vous parlez. Si vous en parlez plusieurs, choisissez « Auto » : le modèle s'y retrouve tout seul, un peu plus lentement et un peu moins précisément toutefois que lorsque la langue est nommée explicitement.</td></tr>" +
			"<tr><td>Quel modèle pour quelle langue</td><td>Cliquez sur une langue : la liste des modèles qui la connaissent se déplie, avec une note de vitesse et de précision et la taille du fichier. Ainsi, une langue peut garder un modèle rapide pendant qu'une autre en reçoit un précis. La ligne « Auto », à part, est commune : son modèle sert toutes les langues qui n'ont pas le leur, et celles-ci apparaissent en gris dans la liste. C'est aussi ce modèle qui reconnaît la langue parlée quand aucune langue de la parole n'est définie.</td></tr>" +
			"<tr><td>Raccourci propre</td><td>Dans une langue dépliée, il y a une place pour sa propre combinaison de touches. Maintenez-la, et la dictée se fait dans cette langue avec son modèle, ses remplacements et ses règles, quelle que soit la langue de la parole réglée ; le raccourci principal garde la langue choisie plus haut. C’est ainsi que l’on dicte l’anglais avec une combinaison et sa propre langue avec une autre, sans ouvrir les réglages. Pendant l’enregistrement, le bandeau affiche le code de la langue. Une combinaison ne peut pas répéter la principale ni celle d’une autre langue.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Russe<span class=\"hs\">GigaAM v3 · précision 5 sur 5 · vitesse 5 sur 5 · 232 Mo</span></label><span class=\"hchip\"><i></i>active</span></div>" +
			"<div class=\"hrow\"><label>Anglais<span class=\"hs\">comme la détection automatique · Nemotron 3.5</span></label><button class=\"hbtn\">Choisir</button></div>" +
			"<div class=\"hcap\">on voit quel modèle s'occupe d'une langue et ce qu'il coûte en vitesse et en place sur le disque</div>" +
			"</div>" +
			"<p class=\"hsub\">Traduction</p>" +
			"<p>Certains modèles savent faire plus que transcrire la parole : ils la traduisent aussitôt. Vous parlez une langue, et le texte arrive dans le champ dans une autre. C'est Whisper qui traduit : vers l'anglais avec sa fonction de traduction intégrée, vers les autres langues en lui indiquant la langue de sortie comme langue de la parole ; pour les paires de langues éloignées, la qualité est inégale. Si le modèle choisi ne sait pas traduire, une ligne orange apparaît à côté de l'interrupteur : si Whisper est installé, le bandeau propose de traduire avec lui, sinon le texte reste dans la langue source.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Changer la langue de sortie du texte</td><td>Active la traduction.</td></tr>" +
			"<tr><td>Langue de sortie du texte par défaut</td><td>La langue vers laquelle traduire si vous n'en avez pas choisi une autre pendant la dictée elle-même.</td></tr>" +
			"<tr><td>Demander la langue de sortie du texte</td><td>Vous pouvez ne pas demander et toujours traduire vers la langue par défaut. Vous pouvez demander à chaque fois : après la dictée, des boutons de langues apparaissent sur le bandeau. Ou demander avec délai : les mêmes boutons, mais si vous n'avez pas appuyé à temps, c'est la langue par défaut qui est prise.</td></tr>" +
			"<tr><td>Langues du dialogue</td><td>Quelles langues proposer sous forme de ces boutons. Gardez-en deux ou trois que vous utilisez vraiment, sinon choisir prendra plus de temps que parler.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Modèles sur le disque</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Votre propre modèle</td><td>Si vous avez votre propre modèle, placez-le dans le dossier des modèles et redémarrez l'application : il apparaît dans les listes des langues qu'il comprend. Le bouton à côté ouvre le bon dossier, pour que vous n'ayez pas à le chercher dans l'Explorateur de fichiers.</td></tr>" +
			"<tr><td>Décharger de la mémoire</td><td>Tant que l'application tourne, le modèle occupe de la mémoire vive. Ce bouton la libère immédiatement, par exemple avant un jeu gourmand. La dictée suivante recharge le modèle, et la première phrase prend plus de temps que d'habitude.</td></tr>" +
			"<tr><td>Vérifier les modèles installés</td><td>Compare les fichiers téléchargés à leur référence. Utile si la connexion a été coupée pendant un téléchargement et que vous soupçonnez un fichier abîmé.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">L'onglet « Règles »</p>" +
			"<p>Trois listes qui corrigent le texte automatiquement, avant même qu'il n'atteigne la fenêtre. Toutes trois sont facultatives : l'application fonctionne très bien sans elles.</p>" +
			"<p class=\"hsub\">Ponctuation et majuscules</p>" +
			"<p>D'où viennent les points, les virgules et les majuscules du texte. Il y a trois options : <b>du modèle</b> (le modèle les place comme il les a entendus ; rapide et en général suffisant), <b>par le modèle d'édition</b> (un deuxième modèle ajoute les signes ; plus soigné, mais chaque dictée prend quelques secondes de plus) et <b>retirer</b> (le texte sort d'un seul tenant, sans signes ; rarement utile, par exemple pour une barre de recherche).</p>" +
			"<p class=\"hsub\">Dictionnaire</p>" +
			"<p>Une liste de mots que le modèle ne connaît généralement pas : noms de famille, noms d'entreprises, termes professionnels. Ce n'est pas un ordre « écris comme ça », mais un indice : quand il entend quelque chose de proche, le modèle choisit plus volontiers votre variante.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">les mots s'ajoutent un par un et vivent sur des cartes séparées</div>" +
			"</div>" +
			"<p class=\"hdim\">Seuls les modèles Whisper lisent le dictionnaire. Nemotron, Parakeet, GigaAM et les autres modèles n'acceptent aucun indice : pour eux, l'outil, c'est la liste des remplacements plus bas. Le jeu de mots par défaut suit la langue tant que vous n'écrivez pas le vôtre.</p>" +
			"<p class=\"hsub\">Remplacements après la reconnaissance</p>" +
			"<p>C'est aussi une liste de mots, mais elle fonctionne autrement : ce n'est plus un indice, c'est une garantie. Ce qui apparaît à gauche devient ce qui est à droite. Chaque remplacement a trois réglages : <b>Mots entiers uniquement</b> (alors « chat » ne touchera pas à « château », et c'est en général ce qu'on veut), <b>Respecter la casse</b> (alors « Pierre » et « pierre » comptent comme des mots différents) et <b>Langue de la règle</b> (le remplacement ne s'applique que lorsque vous dictez dans cette langue).</p>" +
			"<p class=\"hsub\">Commandes vocales</p>" +
			"<p>Des phrases qui n'arrivent pas dans le texte mais déclenchent une action. Le saut de ligne, le nouveau paragraphe et votre propre texte sont reconnus en mots entiers à n'importe quel endroit de la phrase. L'annulation et les commandes de mise en forme n'agissent qu'à la fin de la dictée, l'effacement seulement s'il est dit seul. Une commande de mise en forme dite seule s'applique à la dictée suivante, et le bandeau affiche « Phrase suivante : … ».</p>" +
			"<ul>" +
			"<li><b>« nouvelle ligne »</b> : passer à la ligne ;</li>" +
			"<li><b>« nouveau paragraphe »</b> : laisser une ligne vide entre les paragraphes ;</li>" +
			"<li><b>« annuler »</b> : jeter la phrase et ne rien insérer ;</li>" +
			"<li><b>« efface ça »</b> : effacer ce que le programme a inséré en dernier ;</li>" +
			"<li><b>« en titre »</b> : placer un signe devant la phrase, # par défaut, et une ligne vide après ;</li>" +
			"<li><b>« en liste »</b> : couper la phrase aux virgules, aux points-virgules et au mot « et » en lignes qui commencent par un tiret ;</li>" +
			"<li><b>« avec majuscule »</b> : commencer la phrase par une majuscule ;</li>" +
			"<li><b>« entre guillemets »</b> : mettre la phrase entre guillemets ;</li>" +
			"<li><b>« sans point »</b> : retirer le point final ;</li>" +
			"<li><b>votre propre commande</b> : insérer un texte tout prêt, comme une signature, une adresse ou une réponse type.</li>" +
			"</ul>" +
			"<p>L'effacement ne retire que ce que le programme lui-même a inséré en dernier, et seulement si rien n'a été tapé ni cliqué depuis, que la fenêtre est la même, que le texte n'avait pas de sauts de ligne, qu'il n'a pas été envoyé avec Entrée et qu'il date de moins de cinq minutes. Le texte déjà présent dans la fenêtre n'est jamais touché.</p>" +
			"<p class=\"wh\" id=\"s-post\">L'onglet « Post-traitement »</p>" +
			"<p>Le post-traitement, c'est un deuxième modèle qui réécrit le texte reconnu selon votre consigne. Par exemple : « retire les mots parasites et les répétitions », « fais-en un e-mail professionnel », « présente ça sous forme de liste ». La consigne s'écrit avec des mots ordinaires et s'appelle un prompt.</p>" +
			"<p>Par défaut, il est désactivé, et voici pourquoi : c'est la seule étape qui ralentit nettement la dictée. Au lieu d'une ou deux secondes, l'insertion en prend plusieurs, et sur une machine modeste jusqu'à une dizaine.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Activer le post-traitement</td><td>L'interrupteur principal. Tant qu'il est désactivé, le reste de l'onglet ne fonctionne pas, et le texte est inséré tel qu'il a été reconnu.</td></tr>" +
			"<tr><td>Modèle · Locale</td><td>Un éditeur qui tourne sur votre ordinateur. Rien n'en sort, mais il faut de la mémoire : ces modèles pèsent un gigaoctet ou plus. À côté se trouve un catalogue où l'on choisit et télécharge un modèle en un clic.</td></tr>" +
			"<tr><td>Modèle · Serveur externe</td><td>Au lieu de votre propre modèle, le texte est envoyé à un service extérieur : OpenAI, Groq ou votre propre serveur. Plus rapide et plus intelligent, mais le texte quitte votre ordinateur. La clé d'accès est stockée chiffrée, et le bouton « Tester la connexion » vérifie que l'adresse et la clé fonctionnent.</td></tr>" +
			"<tr><td>Prompts</td><td>Vos consignes pour l'éditeur. Les prompts cochés s'appliquent l'un après l'autre, de haut en bas : d'abord, disons, le nettoyage, puis le style professionnel. Si rien n'est coché, le post-traitement ne fait rien, et l'application vous en avertit par une ligne jaune.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">La voix ne part jamais vers un serveur externe, seulement le texte reconnu. Mais ce texte, lui, part bel et bien : activez donc cette option en connaissance de cause.</div>" +
			"<p>Chaque prompt a un <b>nom</b> (c'est sous ce nom qu'il apparaît sur le bandeau pendant qu'il travaille), un <b>texte de consigne</b> en mots ordinaires, un <b>raccourci propre</b> facultatif (dicter avec ce raccourci fait passer la phrase par ce seul prompt, sans les autres) et un <b>test</b> : un champ où vous pouvez taper une phrase et voir ce qui en sort, sans rien dicter.</p>" +
			"<p class=\"wh\" id=\"s-subs\">L'onglet « Sous-titres »</p>" +
			"<p>Les sous-titres notent ce que joue l'ordinateur – un appel, une vidéo, un cours – et affichent les mots dans une petite fenêtre qui reste au-dessus des autres. C'est la même reconnaissance que pour la dictée, sauf que le son vient des haut-parleurs et non du micro, et la dictée continue de fonctionner pendant les sous-titres. Les sons de l'application elle-même ne sont pas notés.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Démarrer les sous-titres</td><td>Les démarre et les arrête. La ligne à côté dit ce qui se passe : arrêtés, chargement du modèle, à l'écoute ou redémarrage de la reconnaissance. Le menu de la zone de notification, un raccourci à eux et la croix de la fenêtre font de même.</td></tr>" +
			"<tr><td>Langue parlée</td><td>La langue de ce que vous écoutez, distincte de la langue de dictée. Le modèle est celui que l'onglet « Langues et modèles » donne à cette langue. Un modèle en streaming écrit les mots au fil de la parole, les autres chaque phrase une fois celle-ci terminée. Un modèle Whisper a besoin de nombreuses secondes par phrase sur un processeur, son texte a donc un net retard sur la parole ; les phrases accumulées lui sont envoyées ensemble et ressortent en une seule ligne.</td></tr>" +
			"<tr><td>Traduction</td><td>Affiche dans la fenêtre la traduction au lieu du texte reconnu, dans la langue choisie ; le titre de la fenêtre indique alors « Traduction » et les langues, par exemple EN → RU. La parole est découpée en phrases au fil de l'eau, et chaque phrase est traduite aussitôt. Un morceau se termine à la fin d'une phrase ou avant un mot de liaison, et le traducteur reçoit avec lui les morceaux précédents et ses propres traductions, si bien qu'une phrase coupée en morceaux est traduite d'un seul tenant. La traduction apparaît régulièrement, par mots entiers, au rythme de la lecture, et plus vite quand beaucoup s'est accumulé, si bien qu'un gros morceau ne fait pas sauter le texte. Les phrases d'un même passage s'enchaînent en un paragraphe, et une pause commence une nouvelle ligne. Une phrase qui n'a pas pu être traduite est affichée telle qu'elle a été entendue, d'une couleur plus pâle. « Afficher aussi le texte d'origine » partage la fenêtre d'une fine ligne : en haut, plus pâle, ce qui a été dit, en bas la traduction ; la ligne se tire à la souris. La ligne « Qui traduit » décide qui le fait. « Haute qualité » est le modèle de traduction sur la carte graphique : un téléchargement de 2,3 Go, pour une carte graphique avec au moins 4 Go de mémoire propre ; il est chargé sur la carte au démarrage des sous-titres et la quitte à leur arrêt, et un lien à côté ouvre ses conditions d'utilisation. « Votre propre serveur » est le serveur de post-traitement configuré dans l'onglet « Post-traitement », et chaque phrase est alors envoyée à ce serveur. « Le programme choisit » prend le modèle de traduction s'il est téléchargé et si la carte graphique convient, sinon votre serveur s'il est configuré, sinon le modèle d'édition installé dans l'onglet « Post-traitement », qui traduit grossièrement : le sens passe, mais des mots peuvent être inexacts. Un traducteur choisi à la main n'est jamais remplacé par un autre : s'il ne peut pas travailler, le bloc dit ce qui manque et le texte s'affiche tel qu'il a été entendu. La traduction est enregistrée dans un fichier à part, à côté de la transcription : chaque morceau avec l'heure à laquelle il a commencé ; un morceau qui n'a pas pu être traduit est écrit tel qu'il a été entendu. Sans traducteur, pas de traduction, et les sous-titres fonctionnent comme avant.</td></tr>" +
			"<tr><td>Raccourci clavier</td><td>Démarre et arrête les sous-titres depuis n'importe quelle fenêtre. Il ne lance jamais de dictée et ne peut pas reprendre un autre raccourci.</td></tr>" +
			"<tr><td>Enregistrer la transcription</td><td>Chaque session écrit un fichier texte dans le dossier transcripts à côté du programme : une première ligne avec la date, la langue et le modèle, puis chaque ligne terminée avec l'heure à laquelle elle a commencé. Rien n'est écrit tant que la première ligne n'est pas terminée. Quand la traduction est activée, le nom de la transcription porte original et le code de la langue parlée, et un second fichier apparaît à côté, avec translation et le code de la langue cible, par exemple 2026-11-12_10-18-29_original_EN.txt et 2026-11-12_10-18-29_translation_RU.txt.</td></tr>" +
			"<tr><td>Fenêtre des sous-titres</td><td>Où elle se trouve, quelle taille elle a et quelle taille ont les lettres. L'emplacement se choisit sur l'image d'un écran, comme pour le bandeau de dictée : un clic sur l'une des neuf places y met la fenêtre, et avec « Place à vous » elle peut aller n'importe où, en tirant la marque sur l'image ou la fenêtre elle-même ; avec plusieurs moniteurs, on choisit l'écran. La taille se change à la souris par les bords et les coins ; la hauteur se cale sur un nombre entier de lignes, et la largeur et les lignes sont retenues pour chaque résolution d'écran. « Lignes de texte » s'applique là où aucune taille n'est encore fixée, et « Rétablir la taille » oublie celles qui le sont. La molette au-dessus de la fenêtre fait défiler le texte vers l'arrière de trois lignes par cran ; pendant que vous lisez l'ancien, le nouveau texte ne le déplace pas, et la marque « au plus récent ↓ » dans le titre ramène à la parole en cours.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Vos interlocuteurs doivent savoir que leurs paroles sont transcrites.</div>" +
			"<p class=\"wh\" id=\"s-about\">« À propos » et « Contacts »</p>" +
			"<p>L'onglet « À propos » contient un court texte sur ce qu'est l'application et comment elle est construite, ainsi que la liste de toutes les bibliothèques tierces dont elle est faite, avec leurs licences. Les noms sont cliquables et mènent aux pages de ces projets. On y trouve aussi un lien vers le dépôt, où se trouvent le code source et toutes les versions.</p>" +
			"<p>« Contacts » propose trois cartes : l'e-mail, le dépôt et une page pour les bugs et les idées. Un clic sur l'adresse e-mail ouvre votre logiciel de messagerie avec le destinataire déjà rempli, et l'icône à côté copie simplement l'adresse.</p>" +
			"<p class=\"wh\" id=\"s-files\">Ce qu'il y a dans le dossier de l'application</p>" +
			"<p>L'application est portable : tout ce qui lui appartient se trouve dans un seul dossier et ne se répand nulle part ailleurs.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Tous vos réglages dans un seul fichier texte. Vous pouvez le copier sur un autre ordinateur pour ne pas tout reconfigurer.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>Le journal. Le premier endroit où regarder si l'application se comporte bizarrement.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>L'historique des dictées, si vous l'avez activé. Du texte ordinaire : si vous le souhaitez, vous pouvez simplement supprimer le fichier.</td></tr>" +
			"<tr><td>models</td><td>Le dossier des modèles téléchargés, la partie la plus lourde : plusieurs gigaoctets. Lui aussi peut être déplacé sur un autre ordinateur, pour ne pas tout retélécharger.</td></tr>" +
			"<tr><td>transcripts</td><td>Les transcriptions des sous-titres, un fichier texte par session, et leurs traductions quand la traduction est activée. Ce sont vos documents : le désinstallateur ne les supprime que si vous lui demandez aussi d'effacer les réglages et les données.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>Les petits programmes qui font tourner les modèles. Sans eux, un modèle n'est qu'un fichier que rien ne sait lire.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Ce qui part sur internet et ce qui n'y part pas</p>" +
			"<ul>" +
			"<li><b>La voix</b> n'est ni enregistrée ni envoyée, sauf si vous avez vous-même activé la reconnaissance sur un autre ordinateur.</li>" +
			"<li><b>Les sous-titres</b> sont reconnus ici, comme la dictée, et la transcription reste dans le dossier transcripts ; avec la reconnaissance sur un autre ordinateur, les phrases partent vers ce serveur ; avec la traduction par votre propre serveur de post-traitement, chaque ligne terminée y part aussi.</li>" +
			"<li><b>Le texte</b> ne quitte pas la machine, sauf si vous avez vous-même activé un serveur externe de post-traitement.</li>" +
			"<li><b>L'historique</b> est un fichier à côté de l'application et se vide selon la durée que vous avez choisie.</li>" +
			"<li><b>Le réseau</b> sert deux fois : pour télécharger un modèle et pour vérifier les mises à jour. Cette vérification peut être désactivée.</li>" +
			"<li><b>La clé</b> d'un service externe est stockée chiffrée par Windows et n'apparaît jamais dans le journal.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Si quelque chose ne va pas</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Le texte n'apparaît pas dans le champ</td><td>Le programme n'accepte pas le collage depuis le presse-papiers. Activez « Saisie caractère par caractère » dans l'onglet « Commandes et comportement », ou ajoutez un petit délai avant l'insertion.</td></tr>" +
			"<tr><td>La dictée ne démarre pas</td><td>Le raccourci clavier est pris par un autre programme. Choisissez-en un autre, de préférence avec la touche Windows : ces combinaisons sont plus rarement prises.</td></tr>" +
			"<tr><td>Seul du silence a été enregistré</td><td>Ce n'est pas le bon microphone, ou il est coupé. Ouvrez l'onglet « Microphone », regardez la barre de niveau et appuyez sur « Vérifier le microphone ».</td></tr>" +
			"<tr><td>La première phrase est très longue</td><td>C'est prévu ainsi : le modèle se charge en mémoire à la première utilisation. Ensuite, la reconnaissance prend quelques secondes.</td></tr>" +
			"<tr><td>Le message « Reconnaissance indisponible »</td><td>La reconnaissance sur un autre ordinateur est choisie, mais aucune adresse n'est définie. Saisissez l'adresse ou revenez à « Sur cet ordinateur ».</td></tr>" +
			"<tr><td>Le post-traitement ne change rien</td><td>Aucun prompt n'est coché, aucun modèle n'est choisi, ou le serveur externe ne répond pas. La cause est écrite sur l'onglet lui-même, à côté de l'interrupteur.</td></tr>" +
			"<tr><td>Le texte est étrange, des mots sont remplacés</td><td>Jetez un œil à l'onglet « Règles » : très probablement, un remplacement ou une commande vocale s'est déclenché.</td></tr>" +
			"<tr><td>Rien de tout cela</td><td>« Système » → « Ouvrir le journal ». L'application y note tout ce qui la concerne, et la cause est en général écrite en toutes lettres.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">Vous n'avez pas trouvé de réponse ? Écrivez un e-mail depuis l'onglet « Contacts » ou ouvrez une issue dans le dépôt. L'auteur lit les deux.</p>",
		"S_WIZ_HELLO":      "Deux minutes de configuration",
		"S_WIZ_HELLO_TEXT": "{app} transforme votre voix en texte à l'endroit du curseur : maintenez le raccourci, dites une phrase, relâchez — le texte est là. Tout tourne sur votre machine, le son n'en sort jamais.",
		"S_WIZ_LATER":      "Tout ce que nous choisissons ici se modifie ensuite dans les paramètres.",
		"S_WIZ_MODEL_TEXT": "Le modèle adapté à votre langue est choisi pour vous et se télécharge en arrière-plan.",
		"S_WIZ_PICKED":     "Choisi pour : %s",
		"S_WIZ_DL_BG":      "Téléchargement en arrière-plan, vous pouvez continuer",
		"S_WIZ_DL_FAIL":    "Le téléchargement s’est arrêté : %s",
		"S_WIZ_RETRY":      "Réessayer",
		"S_WIZ_TRY_DL":     "Le modèle est encore en cours de téléchargement (%d%%). Le champ fonctionnera dès qu’il sera prêt.",
		"S_WIZ_TRY_LOAD":   "Chargement du modèle…",
		"S_WIZ_TRY_LIVE":   "Avec %s, les mots s’affichent sur le bandeau pendant que vous parlez encore.",
		"S_WIZ_DONE_DL":    "Le modèle est encore en cours de téléchargement. La dictée fonctionnera d’elle-même dès qu’il sera en place.",
		"S_WIZ_INPUT_TEXT": "C'est la combinaison que vous maintiendrez en parlant. Dites quelque chose et vérifiez que la barre de niveau bouge.",
		"S_WIZ_T_TRY":      "Essai",
		"S_WIZ_TRY_PH":     "le texte apparaîtra ici",
		"S_WIZ_T_DONE":     "C'est prêt",
		"S_WIZ_DONE_TEXT":  "{app} vit dans la zone de notification : clic gauche sur l'icône pour les paramètres, clic droit pour le menu. Vous pouvez dicter dans toute fenêtre dotée d'un curseur de saisie.",
		"S_AUTORUN":        "Démarrer avec Windows",
		"S_AUTORUN_SUB":    "Une entrée au démarrage de l'utilisateur courant",
		"S_WIZ_SKIP":       "Passer",
		"S_WIZ_BACK":       "Retour",
		"S_WIZ_NEXT":       "Suivant",
		"S_WIZ_FINISH":     "Terminer",
		"S_WIZ_WAIT":       "J'attends la première phrase…",
		"S_WIZ_HEARD":      "Entendu :",
		"S_WIZ_HAVE":       "Tout le nécessaire est déjà téléchargé",
		"S_WIZ_TRY_TEXT":   "Placez le curseur dans le champ ci-dessous, maintenez %s, dites une phrase et relâchez.",
	}
	settingsStrings["es"] = map[string]string{
		"S_SUBS_ON": "Los subtítulos están en marcha",
		"S_SUBS_HERO_HINT": "inícialos — una llamada o un vídeo se vuelven texto en pantalla",
		"S_SUBS_HERO_START": "Iniciar manualmente",
		"S_SUBS_HERO_STOP": "Detener",
		"S_UPDATED_TO": "Actualizado a la versión %s",
		"S_SUBS_HERO_LIVE": "%s · palabras a medida que se dicen",
		"S_SUBS_HERO_PHRASES": "%s · texto tras cada frase",
		"S_NAV_SUBS": "Subtítulos",
		"S_SUBS_HINT": "La aplicación transcribe lo que suena en el ordenador —una llamada, un vídeo— y muestra las palabras en una ventana encima de las demás.",
		"S_SUBS_NOTICE": "Las personas con las que hablas deben saber que se transcribe lo que dicen.",
		"S_SUBS_START": "Iniciar subtítulos",
		"S_SUBS_STOP": "Detener subtítulos",
		"S_SUBS_OFF": "Los subtítulos están apagados",
		"S_SUBS_LOADING": "Cargando %s…",
		"S_SUBS_LISTEN_LIVE": "Escuchando el sonido del ordenador · %s, palabras a medida que se dicen",
		"S_SUBS_LISTEN_PHRASES": "Escuchando el sonido del ordenador · %s, texto tras cada frase",
		"S_SUBS_RESTARTING": "El reconocimiento se está reiniciando…",
		"S_SUBS_LANG": "Idioma del habla",
		"S_SUBS_LANG_SUB": "El idioma de lo que escuchas, aparte del idioma de dictado",
		"S_SUBS_M_LIVE": "Lo escribirá %s: las palabras aparecen a medida que se dicen.",
		"S_SUBS_M_PHRASES": "Lo escribirá %s: el texto aparece cuando termina cada frase.",
		"S_SUBS_M_ABSENT": "Para este idioma está elegido %s, pero no está instalado.",
		"S_SUBS_M_WHISPER": "%s es un modelo Whisper: en un procesador necesita muchos segundos por frase, así que el texto irá muy por detrás del habla. Para subtítulos va mejor un modelo de streaming o uno rápido.",
		"S_SUBS_GO_MODELS": "Abrir «Idiomas y modelos»",
		"S_SUBS_REMOTE": "El reconocimiento se hace en otro ordenador: cada frase se envía a ese servidor.",
		"S_SUBS_TR": "Traducción",
		"S_SUBS_TR_SUB": "La ventana muestra la traducción en lugar del texto reconocido",
		"S_SUBS_TR_OFF": "Desactivada",
		"S_SUBS_TR_LOCAL": "%s traduce: la traducción llega frase a frase, un momento después del habla.",
		"S_SUBS_TR_API": "Traduce tu servidor de posprocesado (%s): cada frase se envía a él, y la traducción llega frase a frase, un momento después del habla.",
		"S_SUBS_TR_NONE": "No hay modelo de traducción, ni servidor de posprocesado, ni modelo editor, así que no habrá traducción.",
		"S_SUBS_TR_GO": "Abrir «Posprocesado»",
		"S_SUBS_TR_SAME": "El idioma de la traducción es el mismo que el idioma del habla: no hay nada que traducir.",
		"S_SUBS_WHO": "Quién traduce",
		"S_SUBS_WHO_SUB": "El programa toma el mejor disponible, o fijas uno",
		"S_SUBS_WHO_AUTO": "Elige el programa",
		"S_SUBS_WHO_QUALITY": "Alta calidad: un modelo en la tarjeta gráfica",
		"S_SUBS_WHO_SERVER": "Tu propio servidor",
		"S_SUBS_TRM": "Modelo de traducción",
		"S_SUBS_TRM_SUB": "%s, funciona en una tarjeta gráfica con al menos 4 GB de memoria propia",
		"S_SUBS_TRM_ABSENT": "No descargado",
		"S_SUBS_TRM_GETTING": "Descargando… %s %",
		"S_SUBS_TRM_CHECKING": "Comprobando el archivo…",
		"S_SUBS_TRM_READY": "Listo",
		"S_SUBS_TRM_RUNNING": "Listo, cargado en la tarjeta gráfica",
		"S_SUBS_TRM_BAD": "El archivo de la carpeta models está dañado o incompleto",
		"S_SUBS_TRM_ERR": "No descargado: %s",
		"S_SUBS_TRM_SMALL": "Este ordenador no puede usarlo: su tarjeta gráfica tiene %s MB de memoria propia y hacen falta 4 GB.",
		"S_SUBS_TRM_INTEGRATED": "Este ordenador no puede usarlo: solo tiene gráficos integrados, que comparten la memoria del ordenador.",
		"S_SUBS_TRM_NOCARD": "Este ordenador no puede usarlo: Windows no informa de ninguna tarjeta gráfica.",
		"S_SUBS_TRM_UNKNOWN": "Este ordenador no puede usarlo: esta versión de Windows no describe su tarjeta gráfica.",
		"S_SUBS_TRM_DOWN": "No se inició: %s",
		"S_SUBS_TRM_DOWN_NEXT": "Hasta el próximo inicio de los subtítulos traduce el siguiente de la lista.",
		"S_SUBS_TRM_DOWN_HEARD": "Hasta el próximo inicio de los subtítulos el texto se muestra tal como se oyó.",
		"S_SUBS_TRM_DEL_SURE": "¿Eliminar?",
		"S_SUBS_TRM_W_VULKAN": "el controlador gráfico no ofrece Vulkan",
		"S_SUBS_TRM_W_DEVICE": "el motor no ve esta tarjeta gráfica",
		"S_SUBS_TRM_W_MEMORY": "no hay suficiente memoria libre en la tarjeta gráfica",
		"S_SUBS_TRM_W_CRASH": "el motor se detuvo durante la carga",
		"S_SUBS_TRM_W_TIMEOUT": "el motor no respondió a tiempo",
		"S_SUBS_TRM_W_PORT": "su puerto lo ocupa otro programa",
		"S_SUBS_TRM_W_EXE": "falta llama-server-gpu.exe junto al programa",
		"S_SUBS_TRM_DL": "Descargar",
		"S_SUBS_TRM_CANCEL": "Cancelar",
		"S_SUBS_TRM_DEL": "Eliminar",
		"S_SUBS_TRM_TERMS": "Condiciones de uso del modelo",
		"S_SUBS_TR_QUALITY": "Traduce el modelo de traducción en la tarjeta gráfica (%s): la traducción llega frase a frase, un momento después del habla.",
		"S_SUBS_TR_MISS_CARD": "Está elegida la opción «Alta calidad», pero la tarjeta gráfica de este ordenador no sirve, así que el texto se muestra tal como se oyó.",
		"S_SUBS_TR_MISS_MODEL": "Está elegida la opción «Alta calidad», pero el modelo de traducción no está descargado, así que el texto se muestra tal como se oyó.",
		"S_SUBS_TR_MISS_START": "Está elegida la opción «Alta calidad», pero el modelo de traducción no se inició en la tarjeta gráfica, así que hasta el próximo inicio de los subtítulos el texto se muestra tal como se oyó.",
		"S_SUBS_TR_MISS_SERVER": "Está elegido tu propio servidor, pero no hay ningún servidor de posprocesado configurado, así que el texto se muestra tal como se oyó.",
		"S_SUBS_ORIG": "Mostrar también el texto original",
		"S_SUBS_ORIG_SUB": "la ventana muestra sobre la traducción lo que se dijo; la línea entre ambos se arrastra con el ratón",
		"S_SUBS_KEY": "Atajo de teclado",
		"S_SUBS_KEY_SUB": "Inicia y detiene los subtítulos desde cualquier ventana",
		"S_SUBS_SAVE": "Guardar la transcripción",
		"S_SUBS_SAVE_SUB": "Un archivo de texto por sesión en la carpeta transcripts junto al programa, y un segundo con la traducción cuando la traducción está activada",
		"S_SUBS_FOLDER": "Carpeta de transcripciones",
		"S_SUBS_WINDOW": "Ventana de subtítulos",
		"S_SUBS_ROWS": "Filas de texto",
		"S_SUBS_SIZE": "Tamaño del texto",
		"S_SUBS_SIZE_S": "Pequeño",
		"S_SUBS_SIZE_M": "Mediano",
		"S_SUBS_SIZE_L": "Grande",
		"S_SUBS_BOX": "Tamaño de la ventana",
		"S_SUBS_BOX_SUB": "arrastra los bordes y las esquinas de la ventana con el ratón; el tamaño se recuerda para cada resolución de pantalla y «Filas de texto» vale mientras no lo haya",
		"S_SUBS_BOX_BTN": "Restablecer el tamaño",
		"S_SUBS_BOX_DONE": "La ventana vuelve a su tamaño predeterminado",
		"S_SUBS_FREE": "Sitio propio",
		"S_SUBS_FREE_SUB": "la ventana se puede arrastrar a cualquier parte, el sitio se recuerda para cada resolución de pantalla",
		"S_SUBS_POS": "Dónde mostrar la ventana",
		"S_SUBS_POS_SCHEME_SUB": "haz clic en la pantalla: la ventana se coloca ahí",
		"S_SUBS_POS_DRAG_SUB": "arrastra la marca o la propia ventana: se coloca donde quieras",
		"S_SUBS_MON": "Pantalla",
		"S_SUBS_MON_SUB": "en qué monitor mostrar la ventana de subtítulos",
		"S_TITLE": "{app} — Ajustes", "S_DICT_HINT": "Términos, nombres y abreviaturas separados por comas — una pista para el oído, no comandos. Vale para Whisper; el ruso a través de GigaAM lo ignora. El conjunto por defecto sigue al idioma de reconocimiento hasta que escribas el tuyo.",
		"S_TR_DEFAULT": "Cambiar el idioma de salida del texto", "S_TR_TARGET": "Idioma de salida del texto por defecto", "S_TR_ASK": "Preguntar el idioma de salida del texto", "S_TR_ASK_NEVER": "No preguntar — traducir enseguida",
		"S_SRCLANG_SUB": "usted lo habla; determina el modelo de reconocimiento",
		"S_TR_LANGS_SUB": "estos idiomas serán botones en la barra al insertar",
		"S_TR_UNAVAIL": "%s no sabe traducir — el texto se quedará en el idioma de origen",
		"S_TR_GET_WHISPER": "Descargar %s para traducir (%d MB)",
		"S_TR_STANDIN": "%s no sabe traducir — la barra ofrecerá %s en su lugar",
		"S_TR_LOCK": "%s no se puede quitar de la lista — es el idioma de salida del texto por defecto. Elija otro idioma por defecto y entonces %s podrá excluirse.",
		"S_TR_LOCK_OK": "Entendido",
		"S_TR_ONE": "Hay varios idiomas marcados, pero sin preguntar la traducción irá siempre a uno — %s (idioma de salida por defecto). Los demás quedan marcados pero deshabilitados.",
		"S_TR_NOMODEL": "%s no sabe traducir. Si continúa, la traducción se apagará y no estará disponible mientras trabaje este modelo.",
		"S_TR_CONFIRM": "Confirmar",
		"S_TR_ASK_ALWAYS": "Preguntar cada vez", "S_TR_ASK_TIMEOUT": "Preguntar, con tiempo límite", "S_TR_SECONDS": "Tiempo, s",
		"S_TR_LANGS": "Idiomas del diálogo",
		"S_LLM_HINT": "Los perfiles marcados se aplican uno tras otro, de arriba abajo, en el dictado normal. Nada marcado: el texto se inserta tal cual.",
		"S_PROF_ADD": "Añadir", "S_PROF_NAME": "Nombre", "S_PROF_PROMPT": "Prompt", "S_PROF_TEST": "Prueba",
		"S_HOTKEY": "Atajo de teclado", "S_CHANGE": "Cambiar…", "S_UILANG": "Idioma de la interfaz", "S_AUTO": "Como el sistema",
		"S_SEC_SOUND": "Sonido", "S_BEEP": "Señales sonoras", "S_SOUND": "Sonido de señal",
		"S_SND_SPEECH": "Voz de Windows", "S_SND_CHIME": "Campanilla", "S_SND_SOFT": "Suave", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Pulsar Enter tras insertar (envío auto)", "S_RESTORE": "Restaurar portapapeles tras insertar",
		"S_NAV_HISTORY": "Historial", "S_HIST_ON": "Guardar el historial de dictados", "S_HIST_ON_SUB": "solo texto, en este equipo; el audio nunca se guarda",
		"S_HIST_DAYS": "Cuántos días guardar", "S_HIST_MAX": "Cuántas entradas guardar",
		"S_HIST_SKIP": "No registrar nunca desde estos programas", "S_HIST_SKIP_SUB": "separados por comas: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Añadir un programa", "S_SKIP_EDIT_DLG": "Editar el programa", "S_SKIP_NAME": "Nombre del programa", "S_SKIP_NAME_SUB": "El nombre del archivo sin ruta: keepass.exe. Un asterisco final atrapa todas las versiones: 1password*", "S_SKIP_OPEN": "Programas abiertos ahora", "S_SKIP_REFRESH": "Actualizar la lista", "S_SKIP_PICKED": "%d de %d elegidos", "S_SKIP_NONE": "Nada elegido", "S_SKIP_EMPTY": "La lista está vacía — el historial se guarda de todos los programas", "S_SKIP_ADD_BTN": "Añadir un programa", "S_SKIP_HINT": "Lo que dicte en estos programas nunca llega al historial. El pegado funciona como siempre.",
		"S_HIST_LIST": "Entradas", "S_HIST_CLEAR": "Vaciar", "S_HIST_TILL": "hasta %s", "S_HIST_TILL1": "hasta mañana", "S_HIST_TILL_FULL": "Se borrará el %s — se guarda %s", "S_HIST_LIST_HINT": "Lo que se dictó: copiar, pegar en cualquier ventana o borrar.", "S_HIST_COPY": "Copiar",
		"S_HIST_KEEP": "Cuánto tiempo guardar",
		"S_UNIT_MIN": "minutos",
		"S_UNIT_HOUR": "horas",
		"S_UNIT_DAY": "días",
		"S_HIST_FIND": "Buscar en el historial…", "S_HIST_EMPTY": "Todavía no hay historial", "S_HIST_ASK": "¿Eliminar todo el historial de dictados?",
		"S_SEC_CMD": "Comandos de voz", "S_CMD_HINT": "Lo que dices se convierte en un salto de línea, un signo, una cancelación, el borrado de la última inserción o una forma para la frase en vez de acabar en el texto. Se buscan como palabras completas y se aplican de arriba abajo, después de los reemplazos; los comandos de forma actúan al final del dictado.",
		"S_CMD_ADD": "Añadir un comando", "S_CMD_PRESET": "Añadir los habituales", "S_CMD_PH": "la frase que dirás",
		"S_CMD_NEWLINE": "salto de línea", "S_CMD_PARAGRAPH": "nuevo párrafo", "S_CMD_TEXT": "insertar texto", "S_CMD_CANCEL": "cancelar el dictado",
		"S_CMD_TEXT_PH": "qué insertar", "S_CMD_EMPTY": "Todavía no hay comandos", "S_CMD_DEL": "Eliminar el comando",
		"S_CMD_P_NEWLINE": "nueva línea", "S_CMD_P_PARAGRAPH": "nuevo párrafo", "S_CMD_P_CANCEL": "cancelar",
		"S_CMD_UNDO": "borrar la última inserción", "S_CMD_HEADING": "convertir en título", "S_CMD_LIST": "convertir en lista",
		"S_CMD_CAPITAL": "mayúscula inicial", "S_CMD_QUOTES": "poner entre comillas", "S_CMD_NOSTOP": "quitar el punto final",
		"S_CMD_HEADING_PH": "signo antes del título, p. ej. #", "S_CMD_LIST_PH": "palabras que unen el último elemento, p. ej. y", "S_CMD_QUOTES_PH": "signos de apertura y cierre, p. ej. «»",
		"S_CMD_END_NOTE": "Dilo al final del dictado; dicho solo, da forma al siguiente.",
		"S_CMD_UNDO_NOTE": "Dilo solo. Borra lo último que se insertó si desde entonces no has escrito ni hecho clic y la ventana es la misma; el texto que ya estaba nunca se toca.",
		"S_CMD_P_UNDO": "borra eso", "S_CMD_P_HEADING": "como título", "S_CMD_P_LIST": "como lista", "S_CMD_P_CAPITAL": "con mayúscula",
		"S_CMD_P_QUOTES": "entre comillas", "S_CMD_P_NOSTOP": "sin punto", "S_CMD_P_AND": "y, e", "S_CMD_P_MARKS": "«»",
		"S_SEC_REPLACE": "Reemplazos tras el reconocimiento", "S_REPLACE_HINT": "Lo que se oyó mal se convierte en lo que querías decir, justo tras el reconocimiento y antes de los prompts. Se aplican de arriba abajo.",
		"S_REPL_WHOLE_FULL": "Solo palabras completas", "S_REPL_CASE_FULL": "Distinguir mayúsculas y minúsculas", "S_CMD_ACTION": "Acción",
		"S_FM_ADD": "Añadir",
		"S_TIP_REPL_LANG": "La regla solo actúa cuando dictas en el idioma elegido. «todos los idiomas» — actúa siempre.",
		"S_TIP_REPL_CASE": "Las mayúsculas y minúsculas importan: «git» y «Git» son palabras distintas. Desactivado — se ignora.",
		"S_TIP_REPL_WHOLE": "El reemplazo solo actúa si el texto forma una palabra aparte. Desactivado — también dentro de otras palabras.",
		"S_TIP_CMD_ACTION": "Qué ocurre al decir la frase: salto de línea, nuevo párrafo, tu propio texto, cancelar el dictado, borrar la última inserción o dar forma a la frase: título, lista, mayúscula inicial, comillas o sin punto final.",
		"S_LIST_FILTER_PH": "buscar…",
		"S_REPL_DEL": "Eliminar el reemplazo",
		"S_LIST_NOTHING": "No se encontró nada: «%s»",
		"S_FM_T_REPL_ADD": "Añadir un reemplazo", "S_FM_T_REPL_EDIT": "Editar el reemplazo",
		"S_FM_T_CMD_ADD": "Añadir una orden", "S_FM_T_CMD_EDIT": "Editar la orden",
		"S_MT_DEL": "Eliminar el modelo", "S_MT_DEL_PROMPT": "Eliminar el prompt", "S_MT_DL": "Descargar el modelo",
		"S_MT_TR_OFF": "Apagar la traducción", "S_MT_TR_ONE": "Traducir sin preguntar", "S_MT_TR_LOCK": "Idioma de salida por defecto",
		"S_MT_REMOTE": "Servidor remoto", "S_MT_POST": "Servidor externo", "S_MT_HIST": "Vaciar el historial",
		"S_MT_RESET": "Restablecer los ajustes", "S_MT_EXE": "Ruta del servidor",
		"S_DICT_ADD": "Añadir una palabra", "S_FM_T_DICT_ADD": "Añadir una palabra", "S_DICT_EMPTY": "Aún no hay palabras",
		"S_DICT_ADD_PH": "una palabra o varias separadas por comas",
		"S_DICT_NOMODEL": "El modelo actual %s no admite el diccionario: solo los modelos Whisper lo leen.",
		"S_OV_FREE": "Sitio propio", "S_OV_FREE_SUB": "la barra se puede arrastrar a cualquier parte",
		"S_OVPOS_DRAG_SUB": "arrastra la barra con el ratón: se coloca donde quieras",
		"S_OVMON_N": "Pantalla %d",
		"S_POST_ENABLE": "Activar el posprocesado",
		"S_API_SUM_URL": "dirección", "S_API_SUM_MODEL": "modelo", "S_API_SUM_KEY": "clave", "S_API_SUM_TIMEOUT": "espera",
		"S_API_SUM_STATE": "estado", "S_API_NO_MODEL": "sin indicar",
		"S_API_NONE": "sin configurar: el posprocesado va en local",
		"S_POSTAPI_SETUP": "Configurar", "S_API_EDIT": "Cambiar", "S_API_KEY_DEL": "Eliminar la clave", "S_API_DLG": "Servidor externo",
		"S_LLM_CATALOG": "Catálogo de modelos", "S_LLM_BLOCK": "Modelos instalados", "S_LLM_NONE_HINT": "Ningún modelo instalado — descargue uno encontrado con la flecha y aparecerá aquí", "S_LLM_IN_MEM": "en memoria", "S_LLM_ON_DISK": "en disco", "S_LLM_EJECT": "Descargar de la memoria", "S_LLM_FOUND": "%d encontrados", "S_LLM_NOSEARCH": "sin búsqueda todavía", "S_LLM_SEARCH_HINT": "Escriba el nombre del modelo y pulse «Buscar»", "S_LLM_PICK_WAIT": "Disponible cuando el modelo se descargue", "S_LLM_INSTALLED": "instalados",
		"S_LLM_SUM_MODEL": "modelo", "S_LLM_SUM_SIZE": "tamaño", "S_LLM_SUM_COUNT": "instalados", "S_LLM_SUM_RAM": "memoria",
		"S_DLG_CLOSE": "Cerrar", "S_LLM_NOPICK": "sin elegir", "S_NO_PROMPTS": "Aún no hay prompts", "S_PROF_DRAG": "arrastra para reordenar",
		"S_PROF_NAME_PH": "cómo llamar al prompt", "S_PROF_TEST_PH": "escribe una frase para probar",
		"S_PF_NEW": "Nuevo prompt", "S_PF_EDIT": "Editar el prompt",
		"S_POST_NO_MODEL": "activado, pero no hay modelo elegido", "S_POST_NO_API": "activado, pero el servidor no está configurado", "S_POST_BAD": "el servidor no respondió: %s", "S_POST_NO_PROMPT": "activado, pero ningún prompt marcado", "S_API_TEST": "Probar la conexión", "S_API_TEST_RUN": "Comprobando…", "S_API_TEST_OK": "El servidor respondió", "S_API_CLEAR": "Borrar", "S_API_CLEAR_ASK": "¿Borrar la dirección, el modelo y la clave del servidor externo? El posprocesado vuelve al modelo local.", "S_RAM_AVAIL": "Memoria disponible: %s GB de %s GB", "S_RAM_OF": "%s GB de %s GB",
		"S_REPL_ADD": "Añadir un reemplazo", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "palabras completas", "S_REPL_CASE": "mayúsculas", "S_REPL_EMPTY": "Todavía no hay reemplazos",
		"S_PASTE_DELAY": "Retraso antes de insertar", "S_PASTE_DELAY_SUB": "cuando el programa aún no acepta el texto",
		"S_OVPOS": "Dónde mostrar la barra", "S_OVPOS_SUB": "en el cursor: junto al punto de escritura; si la aplicación no lo expone, junto al puntero",
		"S_OVPOS_CARET": "En el cursor",
		"S_OVTEXT": "Mostrar el texto reconocido", "S_OVTEXT_SUB": "en la barra tras insertar, en vez del número de caracteres",
		"S_OVERLAY": "Mostrar la barra", "S_OVERLAY_SUB": "durante el dictado la pantalla muestra que se está grabando", "S_TYPEMODE": "Escritura carácter a carácter",
		"S_RECLANG": "Idioma del habla", "S_RECAUTO": "Auto",
		"S_LANG_KEY": "Atajo propio", "S_LANG_KEY_SET": "Asignar…", "S_LANG_KEY_CLEAR": "Quitar el atajo", "S_LANG_KEY_HINT": "Manténgalo pulsado para dictar en este idioma con su modelo. El atajo principal sigue usando el idioma del habla elegido arriba.",
		"S_DL": "Descargar", "S_DEL": "Borrar",
		"S_M_BASE": "rápido, PCs modestos", "S_M_SMALL": "equilibrado", "S_M_MED": "Whisper para cualquier idioma; el que traduce", "S_M_TURBO": "máx. precisión en CPU", "S_M_GIGAAM": "más preciso en ruso, puntúa por sí solo", "S_M_PARAKEET": "25 idiomas europeos, puntúa por sí solo",
		"S_THREADS":  "Hilos de CPU", "S_MINMS": "Grab. mín, ms", "S_MAXSEC": "Grab. máx, s",
		"S_AUTOSTART": "Iniciar whisper-server automáticamente", "S_PORT": "Puerto", "S_SERVEREXE": "Ruta de whisper-server", "S_SERVEREXE_SUB": "se rellena solo; cámbiela únicamente si el servidor está en otro sitio", "S_EXE_RESET": "Restablecer", "S_EXE_WARN": "La aplicación encuentra whisper-server junto a ella. Con una ruta escrita a mano, al mover la carpeta el reconocimiento dejará de arrancar. ¿Cambiarla?", "S_RESET_ALL": "Restablecer los ajustes", "S_RESET_ALL_SUB": "todo vuelve de fábrica, salvo modelos e historial", "S_RESET_ALL_BTN": "Restablecer", "S_RESET_ALL_ASK": "¿Devolver todos los ajustes de fábrica? Los modelos, el historial y los prompts se quedan.", "S_RELOAD_CFG": "Releer config.json", "S_RELOAD_CFG_SUB": "si editó el archivo a mano", "S_RELOAD_CFG_BTN": "Releer", "S_UPD_FOUND": "Está la versión %s", "S_THEME": "Color", "S_THEME_SUB": "el color de la ventana, la barra y el icono", "S_THEME_GREEN": "Verde", "S_THEME_AMBER": "Ámbar", "S_THEME_BLUE": "Azul", "S_THEME_PINK": "Rosa", "S_THEME_EDITOR": "Editor", "S_THEME_NEON": "Neón", "S_WND_MAX": "Ocupar la pantalla", "S_WND_RESTORE": "Volver al tamaño anterior", "S_WND_MIN": "Ocultar en la bandeja", "S_WND_CLOSE": "Cerrar la ventana", "S_SKIN": "Diseño", "S_SKIN_SUB": "tipografía, formas, efectos y animación", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Suave", "S_SKIN_PAPER": "Documento", "S_SKIN_FLUENT": "Fluent", "S_SKIN_STUDIO": "Estudio",
		"S_SERVERURL": "Servidor externo (URL)", "S_URLHINT": "Si se define, no se inicia el servidor local",
		"S_STT_SRV": "Servidor de reconocimiento",
		"S_STT_SRV_HINT": "Los modelos Whisper los ejecuta un programa aparte. Puede funcionar en este ordenador o en otro: elija cuál usar.",
		"S_SRV_LOCAL": "En este ordenador",
		"S_SRV_REMOTE": "En otro ordenador",
		"S_SRV_REMOTE_HINT": "El mismo whisper-server, arrancado en otro sitio: un servidor doméstico, una máquina con tarjeta gráfica, el ordenador de al lado.",
		"S_SRV_K_AUTO": "inicio automático",
		"S_SRV_K_FILE": "archivo",
		"S_SRV_K_ADDR": "dirección",
		"S_SRV_K_CHECK": "última comprobación",
		"S_SRV_NEAR": "whisper-server.exe junto a la aplicación",
		"S_SRV_NOADDR": "sin definir",
		"S_SRV_NOCHECK": "sin comprobar",
		"S_SRV_LOCAL_DLG": "Servidor de reconocimiento local",
		"S_SRV_ADDR": "Dirección del servidor",
		"S_SRV_ADDR_SUB": "dirección de whisper-server en la otra máquina, con el puerto",
		"S_SRV_ON": "activado",
		"S_SRV_OFF": "desactivado",
		"S_SRV_K_THREADS": "hilos de CPU",
		"S_SRV_K_PORT": "puerto",
		"S_SRV_DOWN": "Reconocimiento no disponible",
		"S_SRV_DOWN_WHY": "el servidor de reconocimiento remoto no está configurado: indique su dirección en los ajustes",
		"S_SRV_DOWN_GO": "Abrir los ajustes del servidor",
		"S_SRV_WARN_NOW": "El dictado no funciona ahora mismo: está elegido el servidor remoto, pero no tiene dirección.",
		"S_SRV_WARN_LATER": "En cuanto se elija un modelo Whisper, el reconocimiento no funcionará: falta la dirección del servidor remoto.",
		"S_SAVED":      "Guardado",
		"S_ABOUT_HTML": "<p><b>Voz → texto en la posición del cursor.</b></p><p>Coloque el cursor, mantenga el atajo, hable, suelte — el texto se inserta.</p><p>Totalmente local y sin conexión. Tecnologías: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; modelos de Hugging Face.</p><p>Los logs nunca superan ~2 MB.</p>",
		  "S_SEARCH": "Buscar un ajuste…",
		"S_GRP_GENERAL": "General", "S_GRP_SPEECH": "Procesamiento de voz", "S_GRP_INFO": "Información", "S_NAV_POST": "Posprocesado", "S_NAV_HELP": "Ayuda", "S_NAV_CONTACTS": "Contactos", "S_HIST_ADD": "Añadir", "S_CONTACT_MAIL": "Correo", "S_DICT_MODEL": "Modelo de reconocimiento", "S_LIB_ACC": "precisión", "S_LIB_SPD": "velocidad",
		"S_HELP_TOC": "En esta página",
		"S_HELP_TOC_SHOW": "Mostrar el índice: la ventana se ensancha",
		"S_HELP_TOC_HIDE": "Ocultar el índice y devolver el ancho de la ventana",
		"S_CONTACT_TITLE": "Contactar",
		"S_ABOUT_DEPS": "Módulos externos",
		"S_ABOUT_DEPS_HINT": "Código ajeno integrado en la aplicación y sus licencias. Al pulsar el nombre se abre la página del proyecto.",
		"S_DEP_WHISPER": "ejecuta los modelos Whisper",
		"S_DEP_LLAMA": "posprocesado, modelos GGUF, traducción de subtítulos en la tarjeta gráfica (llama-server-gpu.exe)",
		"S_DEP_VULKAN": "cabeceras y biblioteca de importación del motor para la tarjeta gráfica; el entorno de Vulkan llega con el controlador gráfico",
		"S_DEP_SHERPA": "motor de GigaAM, Parakeet, Canary, Qwen, Moonshine",
		"S_DEP_GGML": "biblioteca de tensores dentro de whisper.cpp y llama.cpp",
		"S_DEP_ONNX": "ejecuta los modelos dentro de sherpa-onnx",
		"S_DEP_WEBVIEW": "ventana de ajustes sobre WebView2",
		"S_DEP_WV2RT": "componente de Windows que dibuja esa ventana",
		"S_DEP_MALGO": "captura del micrófono",
		"S_DEP_MINIAUDIO": "capa de audio dentro de malgo",
		"S_DEP_WS": "conexión con sherpa-server",
		"S_DEP_XSYS": "llamadas a WinAPI desde Go",
		"S_DEP_WINLOADER": "carga de DLL dentro de go-webview2",
		"S_DEP_PLEX": "tipografía de la interfaz",
		"S_DEP_HF": "catálogo del que se descargan los modelos",
		"S_CONTACT_HINT": "Un fallo, una idea, una duda sobre un ajuste: escriba un correo si es personal, o abra un issue si es un fallo.",
		"S_CONTACT_REPO": "Repositorio",
		"S_CONTACT_ISSUES": "Fallos e ideas",
		"S_CONTACT_WRITE": "Escribir un correo",
		"S_CONTACT_OPEN": "Abrir",
		"S_STATE_ACTIVE": "Reconoce",
		"S_STATE_USED": "Modelos en uso",
		"S_STATE_INST": "Instalados localmente",
		"S_STATE_INST_SUB": "modelos en disco, listos para asignar",
		"S_PRESETS": "Qué modelo para qué idioma",
		"S_PRESETS_HINT": "El modelo de la fila «Auto» sirve a todos los idiomas que no tienen el suyo — en la lista aparecen en gris. Haga clic en un idioma para darle uno propio. Ese mismo modelo reconoce qué idioma se habla cuando no hay idioma de origen fijado.",
		"S_MFOLDER": "Modelo propio",
		"S_DICT_SAVE": "Guardar",
		"S_OWNM_SUB": "Añada un modelo local de reconocimiento de voz",
		"S_OWNM_ONEFILE": "Un archivo",
		"S_OWNM_FOLDERF": "Carpeta con los archivos del modelo",
		"S_OWNM_S1": "Abra la carpeta de modelos",
		"S_OWNM_S1S": "Carpeta de destino:",
		"S_OWNM_S2": "Copie el modelo",
		"S_OWNM_S2S": "Elija una de las estructuras compatibles",
		"S_OWNM_S3": "Reinicie la aplicación",
		"S_OWNM_S3S": "El modelo aparecerá para los idiomas que admite",
		"S_AS_AUTO": "como la detección automática",
		"S_REC_CHIP": "recomendado",
		"S_BACK_AUTO": "Volver a la detección automática",
		"S_LANGS_COUNT": "idiomas: %d",
		"S_LANGS_UNKNOWN": "idiomas: desconocidos",
		"S_TR_EN": "traduce a %d idiomas",
		"S_TR_LIST": "traduce: %s",
		"S_DL_GOING": "descargando:",
		"S_OPEN_FOLDER": "Abrir la carpeta",
		"S_UNLOAD": "Descargar de la memoria",
		"S_UNLOAD_SUB": "libera la memoria; el siguiente dictado vuelve a cargar el modelo",
		"S_UNLOAD_GO": "Liberar",
		"S_UNLOADED": "Liberado",
		"S_NOT_FOR_LANG": "%s no reconoce este idioma",
		"S_MANUAL_NOTE": "No se puede descargar desde la aplicación — la licencia prohíbe la redistribución. Descargue el archivo usted mismo y descomprímalo en models/moonshine-uk.",
		"S_MANUAL_LINK": "Descargar uno mismo",
		"S_HF_FIT": "solo los que caben en este equipo",
		"S_HF_HIDDEN": "ocultos: %s",
		"S_WIZ_SKIP_DL": "Descargar más tarde",
		"S_WIZ_SKIP_NOTE": "Sin modelo, el dictado no funcionará. Puede descargarlo en «Idiomas y modelos».",
		"S_M_GIGAAM2": "la generación anterior del modelo ruso: misma velocidad, pero sin puntuación",
		"S_M_MOONUK": "modelo ucraniano Moonshine: rápido y ligero, sin puntuación",
		"S_M_LOCAL": "encontrado en la carpeta models; propiedades desconocidas, por eso sin barras",
		"S_ALL_LANGS": "todos los idiomas",
		"S_OVPOS_SCHEME_SUB": "haz clic en la pantalla: la barra se coloca ahí",
		"S_OVDRAG": "Arrástrala a donde quieras",
		"S_OVMON": "Pantalla",
		"S_OVMON_SUB": "en qué monitor mostrar la barra",
		"S_OVMON_CURSOR": "La pantalla con el cursor",
		"S_M_NEMOTRON": "escribe mientras hablas: el texto aparece en la barra en directo; 40 idiomas, puntúa por sí solo",
		"S_M_TINY": "el más pequeño y rápido, para equipos muy modestos; bastante menos preciso",
		"S_STATE_LOADED": "En memoria ahora mismo",
		"S_STATE_LOADED_SUB": "los modelos se descargan solos tras la inactividad",
		"S_STATE_WEEK": "Esta semana",
		"S_ST_SUMMARY": "Resumen", "S_ST_OVERLAY": "Barra en pantalla", "S_ST_BEEP": "Señal sonora", "S_ST_AUTORUN": "Inicio con Windows", "S_ST_POST": "Posprocesado", "S_ST_LOCAL": "local", "S_ST_CHECKED": "comprobado %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "activado", "S_ST_OFF_M": "desactivado", "S_ST_MIC_OK": "la señal es buena", "S_ST_MIC_BAD": "el micrófono no responde", "S_ST_CHECK": "Comprobar", "S_ST_RECOG": "reconocido por %s", "S_ST_VER": "Versión %s", "S_ST_LATEST": "la última", "S_ST_OUTDATED": "desactualizada", "S_ST_UPD_OK": "tiene la última versión", "S_ST_UPD_DL": "Descargando la actualización…",
		"S_ST_QUICK": "Ajustes rápidos",
		"S_ST_MODELS": "Modelos",
		"S_ST_USAGE": "Esta semana",
		"S_ST_READY": "Listo para dictar",
		"S_ST_OFF": "Desactivado en la bandeja",
		"S_ST_OFF_SUB": "el atajo no hace nada hasta que lo vuelva a activar",
		"S_ST_ENABLE": "Activar",
		"S_ST_GOTO": "Abrir este ajuste en su pestaña",
		"S_ST_HOTKEY_GO": "Cambiar el atajo",
		"S_ST_UPD_LAST": "Versión %s — la última",
		"S_ST_UPD_HAVE": "Versión %s disponible",
		"S_ST_MEM": "%s GB libres de %s",
		"S_ST_MEM_SUB": "en memoria: %s · en disco: %d modelos, %s GB",
		"S_ST_MEM_NONE": "nada",
		"S_ST_LANG": "Idioma del habla",
		"S_ST_ASR": "Reconocimiento",
		"S_ST_ON": "activado",
		"S_ST_OFF_W": "desactivado",
		"S_ST_ON_F": "activada",
		"S_ST_OFF_F": "desactivada",
		"S_ST_ACTIVE": "activa",
		"S_ST_IDLE": "no se inicia",
		"S_ST_DISK": "en disco, %s",
		"S_ST_USAGE_SUB": "%d caracteres · hoy %d · %d caracteres de media",
		"S_WEEK_OTHER": "otros",
		"S_ST_NO_WEEK": "esta semana sin dictados",
		"S_ST_AUTORUN_SUB": "la aplicación no se iniciará sola",
		"S_ST_OVERLAY_SUB": "visible al grabar",
		"S_REPL_LANG": "Idioma de la regla",
		"S_REPL_LANG_ALL": "todos los idiomas",
		"S_M_CANARY": "inglés, alemán, español, francés — y traduce entre ellos por sí solo",
		"S_M_QWEN3": "unos 30 idiomas, puntúa por sí solo; el más pesado y preciso del catálogo",
		"S_POSTAPI": "Servidor externo",
		"S_POST_HINT": "Corrige el texto reconocido según sus prompts: quita muletillas, arregla la puntuación, cambia el estilo. Apagado — el texto se inserta tal cual.",
		"S_POST_MODEL": "Modelo",
		"S_SRC_LOCAL": "Local",
		"S_SRC_USED": "en uso",
		"S_HF_GO": "Buscar",
		"S_POSTAPI_HINT": "Vacío por defecto — todo el posprocesado es local. Escriba una dirección y los prompts se ejecutan en un servidor externo: OpenAI, Groq, su propio vLLM — cualquiera con API compatible.",
		"S_POSTAPI_URL": "Dirección",
		"S_POSTAPI_URL_SUB": "vacío = el modelo local; ejemplo: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Modelo",
		"S_POSTAPI_KEY": "Clave API",
		"S_POSTAPI_KEY_SET": "clave guardada (cifrada con Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "sin clave",
		"S_POSTAPI_SAVE": "Guardar clave",
		"S_POSTAPI_TIMEOUT": "Tiempo de espera", "S_SEC_SHORT": "s",
		"S_POSTAPI_WARN": "⚠ El texto reconocido de sus dictados irá a esta dirección. El audio nunca sale. La clave se guarda cifrada.",
		"S_POSTAPI_ASK": "¿Enviar el texto reconocido a %s? El audio se queda en este equipo, pero el texto saldrá.",
		"S_POSTAPI_BADGE": "servidor externo",
		"S_NOT_INSTALLED": "no instalado",
		"S_NAV_STATE": "Estado", "S_NAV_DICT": "Control y comportamiento", "S_NAV_MIC": "Micrófono", "S_NAV_MODELS": "Idiomas y modelos",
		"S_NAV_TEXT": "Reglas", "S_NAV_TR": "Traducción", "S_NAV_SYSTEM": "Sistema", "S_NAV_ABOUT": "Acerca de",
		"S_STATE_HINT": "mantén y habla — el texto aparece donde está el cursor",
		"S_STATE_PROC": "Posprocesado",
		"S_CHANGE_MODEL": "Cambiar", "S_PICK_MODEL": "Elegir", "S_STATE_GET": "Descargar",
		"S_RETRY": "Reintentar", "S_BERR_OPEN": "Abrir los ajustes del servidor",
		"S_STATE_LAST": "Último dictado", "S_STATE_COPY": "Copiar", "S_STATE_MEM": "Memoria",
		"S_STATE_MEM_SUB": "los modelos siguen cargados, la primera frase no espera",
		"S_HOTMODE":       "Modo", "S_HOTMODE_HOLD": "mantener", "S_HOTMODE_TOGGLE": "alternar",
		"S_SUB_HOTMODE": "mantén las teclas, o pulsa una vez para empezar y otra para parar",
		"S_SILENCE":     "Detener por silencio", "S_SILENCE_MS": "Pausa antes de detener, ms",
		"S_SUB_SILENCE": "modo «alternar»: cuando ya has hablado, una pausa termina la grabación como si pulsaras las teclas otra vez",
		"S_SUB_SILENCE_MS": "cuánto puedes estar en silencio; nunca se detiene antes de tus primeras palabras",
		"S_SUB_MINMS":   "ignora pulsaciones accidentales",
		"S_SUB_ENTER":   "envía el mensaje enseguida",
		"S_SUB_CLIP":    "las imágenes y los archivos vuelven tal como estaban",
		"S_SUB_TYPE":    "ayuda donde un campo no admite pegar",
		"S_SEC_OVERLAY": "Aviso en pantalla",
		"S_MIC_CHECK":   "Comprobar el micrófono", "S_MIC_CHECK_SUB": "tres segundos de grabación y un veredicto: nivel, saturación, si hay voz", "S_MIC_CHECKING": "Comprobando…",
		"S_MCHECK": "Comprobar los modelos instalados", "S_MCHECK_SUB": "compara los archivos de los modelos con los hashes de referencia", "S_MCHECK_GO": "Comprobar", "S_MCHECK_RUN": "Comprobando…",
		"S_HIST_INSERT": "Pegar",
		"S_MIC": "Micrófono", "S_MIC_DEFAULT": "Predeterminado del sistema", "S_MIC_REFRESH": "Actualizar la lista",
		"S_MIC_LEVEL": "Nivel de entrada", "S_MIC_QUIET": "silencio",
		"S_SUB_THREADS": "más hilos no siempre es más rápido — mídelo en tu equipo",
		"S_SEC_LLM":     "Modelo editor",
		"S_PUNCT":       "Puntuación y mayúsculas", "S_SUB_PUNCT": "de dónde salen la puntuación y las mayúsculas",
		"S_PUNCT_MODEL": "del modelo", "S_PUNCT_LLM": "del modelo editor", "S_PUNCT_OFF": "quitar",
		"S_SUB_DICT": "Diccionario", "S_SUB_PROMPTS": "Prompts",
		"S_SUB_TRTARGET": "el texto se traduce a él; el diálogo en la barra lo ofrece primero",
		"S_REMOTE_ABOUT": "Hay un servidor remoto configurado: el audio se envía allí y la promesa anterior no se cumple mientras esté activo.",
		"S_UPD":          "Actualizaciones", "S_UPD_CHECK": "Buscar actualizaciones", "S_UPD_AUTO": "Comprobar al iniciar",
		"S_SUB_UPD":  "la única petición de red aparte de la descarga de modelos",
		"S_UPD_NONE": "Tienes la última versión", "S_BADGE_MODELS": "Modelos instalados", "S_BADGE_MISS": "Falta descargar un modelo", "S_BADGE_SYSTEM": "Avisos que requieren atención", "S_BADGE_HIST": "Entradas en el historial", "S_LOG_OPEN": "Abrir el registro", "S_LOG": "Registro", "S_LOG_SUB": "todo lo que la aplicación escribe sobre sí misma", "S_UPD_AVAIL": "La versión %s está disponible.",
		"S_UPD_GO": "Actualizar", "S_UPD_ERR": "Fallo al comprobar", "S_UPD_DL": "Descargando la actualización…",
		"S_SEC_SERVICE": "Servicio", "S_SUB_AUTOSTART": "desactívalo si arrancas el servidor tú mismo",
		"S_SUB_PORT":    "el reconocedor se reinicia solo",
		"S_MODEL_READY": "Modelo descargado — elígelo para cambiar",
		"S_FIT_OK":      "cabe", "S_FIT_WARN": "justo", "S_FIT_BAD": "falta memoria", "S_RAM": "Memoria del equipo:",
		"S_HF_PH":       "Nombre del modelo — p. ej. qwen2.5 instruct",
		"S_NO_LLM":      "Todavía no hay modelos instalados — busca uno en el campo de búsqueda de abajo.",
		"S_NO_LLM_PROF": "Los prompts estarán disponibles en cuanto haya un modelo instalado — el bloque «Modelo» arriba en esta pestaña.",
		"S_UPDATED":     "Última actualización del modelo", "S_PROF_EDIT": "Editar", "S_PROF_CLOSE": "Plegar",
		"S_CONFIRM_DEL": "¿Eliminar el modelo «%s»? Se puede descargar de nuevo.", "S_FREE": "libre",
		"S_DEL_ACTIVE":     "¿Eliminar el modelo activo «%s»? El reconocimiento se detiene hasta que elijas otro, y puedes descargarlo aquí mismo.",
		"S_REMOTE_WARN":    "El audio se enviará a este servidor. El modo local está desactivado.",
		"S_REMOTE_ASK":     "El audio dejará de procesarse en este equipo y se enviará a %s. ¿Activar el modo remoto?",
		"S_REMOTE_BADGE":   "REMOTO",
		"S_OK":             "Sí", "S_CANCEL": "Cancelar", "S_DL_START": "Descargar", "S_DL_CANCEL": "Cancelar la descarga",
		"S_DL_ASK":    "El modelo «%s» no está descargado (%s). ¿Empezar la descarga?",
		"S_NOT_FOUND": "nada",   "S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autor y desarrollador de {app} — una herramienta de dictado local para Windows: la voz se convierte en texto justo en el cursor, sin nubes ni suscripciones.</p>" +
			"<p>El proyecto es abierto: el código, la compilación y las versiones están en GitHub.</p>" +
			"<ul>" +
			"<li>Repositorio: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Perfil del autor: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>¿Has encontrado un fallo o tienes una idea? Abre un issue en el repositorio.</p>",
		"S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Qué es</p>" +
			"<p>Hablas y las palabras aparecen donde parpadea el cursor. Como si las hubieras escrito tú, solo que más rápido.</p>" +
			"<p class=\"hdim\">Funciona así. Haces clic en un campo donde normalmente escribes: un correo, un chat, un campo de búsqueda, un documento. Mantienes pulsadas dos teclas y dices una frase con tu voz de siempre. Las sueltas. Un segundo después, la frase está en el campo. No hace falta cambiar de ventana ni copiar nada: la aplicación funciona por encima de cualquier programa y espera sus teclas en segundo plano.</p>" +
			"<p class=\"hdim\">Todo el cálculo lo hace tu ordenador. La voz no se guarda: pasa unos segundos en la memoria, se convierte en texto y desaparece. La aplicación solo se conecta a internet en dos ocasiones: cuando descargas un modelo de reconocimiento y cuando comprueba si ha salido una versión nueva. Esa comprobación se puede desactivar, y entonces no necesita la red para nada.</p>" +
			"<p class=\"hdim\">La palabra «modelo» aparecerá a menudo más abajo. Es simplemente un archivo que sabe entender el habla. Se descarga una vez y se queda en una carpeta junto a la aplicación; cuanto más grande es, mejor oye y más despacio trabaja.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Qué pasa cuando sueltas las teclas</p>" +
			"<div class=\"hflow\">" +
			"<span>grabación</span><em>→</em><span>reconocimiento</span><em>→</em><span>reemplazos</span><em>→</em><span>comandos</span><em>→</em><span>traducción</span><em>→</em><span>edición</span><em>→</em><span>inserción</span>" +
			"</div>" +
			"<p>Normalmente tarda uno o dos segundos, y no hace falta conocer los pasos. Pero si algún día el texto sale distinto de lo que esperabas, esta cadena te ayuda a ver qué paso lo cambió.</p>" +
			"<ul>" +
			"<li><b>Grabación.</b> Mientras mantienes las teclas, la aplicación escucha el micrófono. Ignora las pulsaciones demasiado cortas, para que una tecla rozada sin querer no empiece un dictado.</li>" +
			"<li><b>Reconocimiento.</b> La voz se convierte en letras. De eso se encarga un modelo. Para el ruso suele trabajar GigaAM; para el polaco, Parakeet; para la mayoría de los demás idiomas, Nemotron, y Whisper se ocupa de los idiomas restantes y de la traducción. Son solo nombres: no hace falta elegirlos a mano.</li>" +
			"<li><b>Reemplazos.</b> Tu lista personal de «lo oído → lo correcto». Si el modelo se empeña en escribir «git hub» en vez de «GitHub», lo corriges una vez y a partir de ahí se arregla solo.</li>" +
			"<li><b>Comandos.</b> Algunas frases funcionan como botones. Di «nueva línea» y el texto salta a una línea nueva, sin que esas palabras lleguen al campo.</li>" +
			"<li><b>Traducción.</b> Si la activaste, hablas en un idioma y el texto llega en otro.</li>" +
			"<li><b>Edición.</b> Un paso opcional. Un segundo modelo pule el texto: quita los «eh» y los «o sea», pone las comas y da a la redacción un tono formal.</li>" +
			"<li><b>Inserción.</b> El texto llega al campo. Lo que tenías en el portapapeles antes del dictado vuelve a su sitio.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Mientras ocurre todo esto, abajo en la pantalla se ve una barra estrecha con el nombre del paso actual. Así sabes de un vistazo que la aplicación está trabajando y no colgada.</p>" +
			"<p class=\"wh\" id=\"s-first\">El primer inicio</p>" +
			"<p>La primera vez que se abre, la aplicación te guía en tres pasos. Lleva un minuto, y después ya puedes dictar.</p>" +
			"<ul>" +
			"<li><b>En qué idioma vas a dictar.</b> La pregunta más importante: según tu respuesta, la aplicación elige por sí sola un modelo adecuado y lo descarga en segundo plano, y se ve el progreso. Si ese idioma no coincide con el de la ventana, ahí mismo aparece una fila para cambiar también el idioma de la interfaz.</li>" +
			"<li><b>Prueba.</b> El atajo de teclado (por defecto Ctrl y la tecla Windows; para cambiarlo, pulsa lo que te resulte cómodo mantener con una mano), el micrófono con un medidor de nivel en directo y un campo donde dictar una frase de prueba. Si el medidor no se mueve cuando hablas, está elegido el micrófono equivocado. El campo cobra vida en cuanto el modelo está listo.</li>" +
			"<li><b>Listo.</b> El interruptor «Iniciar con Windows», y el asistente termina.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Puedes saltarte cualquier paso y volver a él más tarde: todo está también en los ajustes. El modelo sigue descargándose en segundo plano aunque pases al siguiente paso o termines el asistente antes de que acabe. Si quieres repetir el asistente, abre la aplicación añadiendo <code>-wizard</code>. Después de una actualización no vuelve a aparecer.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">La barra en la parte inferior de la pantalla</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Hable…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">así muestra la aplicación que te está oyendo</div>" +
			"</div>" +
			"<p>Mientras dictas, es lo único que se ve en la pantalla. La barra es pequeña, no captura el ratón y no estorba tu trabajo: simplemente te cuenta qué está pasando ahora mismo con tu frase.</p>" +
			"<ul>" +
			"<li><b>«Hable…»</b> Se está grabando. Las columnas de al lado suben y bajan al ritmo de tu voz. Si se quedan quietas, no se te oye, y no tiene sentido seguir hablando.</li>" +
			"<li><b>«Reconociendo…»</b> Soltaste las teclas y el modelo está trabajando. Suele tardar uno o dos segundos. El primer dictado después de abrir la aplicación tarda más: en ese momento el modelo se está cargando en la memoria.</li>" +
			"<li><b>«Traduciendo…»</b> Se está traduciendo.</li>" +
			"<li><b>«Editando: Cleanup (1/2)»</b> El segundo modelo está retocando el texto. Se ve qué instrucción se está aplicando y cuántas hay en total.</li>" +
			"<li><b>«Insertado: 128 caracteres»</b> Listo, el texto ya está en el campo. La barra se apaga sola al cabo de un par de segundos.</li>" +
			"<li><b>La cruz.</b> ¿Cambiaste de idea? Haz clic en ella en cualquier paso. La grabación se descarta y no llega nada al campo.</li>" +
			"</ul>" +
			"<p class=\"hdim\">La barra se puede desactivar, colocar en cualquier esquina, situar junto al cursor o simplemente arrastrar con el ratón adonde no moleste. Todo eso está en la pestaña «Control y comportamiento».</p>" +
			"<p class=\"wh\" id=\"s-state\">La pantalla «Estado»</p>" +
			"<p>Esta pantalla responde a una sola pregunta: ¿puedes dictar ahora mismo? Si no, te dice por qué y te da un botón que lleva adonde se arregla.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Listo para dictar</b><span>mantén y habla — el texto aparece donde está el cursor</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">el cuadrado verde significa que la aplicación está atenta a sus teclas</div>" +
			"</div>" +
			"<p>La tarjeta de arriba puede ser de tres tipos. <b>Verde</b>: todo va bien, con un recordatorio de las teclas al lado. <b>Amarilla</b>: desactivaste la aplicación desde su icono en la bandeja, junto al reloj, y las teclas no hacen nada hasta que pulses «Activar» ahí mismo, en la tarjeta. <b>Roja</b>: los ajustes se contradicen. Está elegido el reconocimiento en otro ordenador, pero no se ha indicado ninguna dirección.</p>" +
			"<p>Debajo está lo que normalmente se mira de pasada:</p>" +
			"<ul>" +
			"<li><b>Subtítulos.</b> La tarjeta justo debajo de la primera: si los subtítulos están apagados, cargando su modelo, en marcha o reiniciándose, con el modelo y cómo llega el texto, su atajo si hay uno, un botón que los inicia o los detiene y una flecha a la pestaña «Subtítulos».</li>" +
			"<li><b>Esta semana.</b> Un gráfico circular: a qué programas fueron tus dictados y cuántos caracteres dictaste.</li>" +
			"<li><b>Versión.</b> El número de compilación y el botón «Comprobar». Si ha salido una más nueva, la aplicación te lo dice aquí mismo.</li>" +
			"<li><b>Memoria.</b> Cuánta queda libre, qué modelo está cargado en ella ahora y cuánto ocupan todos los modelos en el disco. El botón «Liberar» vacía la memoria al momento, por ejemplo antes de un juego exigente.</li>" +
			"<li><b>Ajustes rápidos.</b> Cinco tarjetas con lo más usado: micrófono, idioma, barra, sonido e inicio con Windows. Muestran cómo está configurado ahora, y la flecha lleva a la pestaña donde se cambia. Aquí no hay interruptores a propósito, para que nada se toque sin querer.</li>" +
			"<li><b>Modelos.</b> Qué está instalado, qué funciona ahora mismo, de qué idiomas se encarga cada uno y qué está simplemente guardado en el disco.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">La pestaña «Sistema»</p>" +
			"<p>El aspecto de la aplicación, su relación con Windows y algunos botones de servicio. En el día a día se entra aquí más o menos una vez al año.</p>" +
			"<p class=\"hsub\">Cómo se ve la aplicación</p>" +
			"<p>Tres filas: idioma de la interfaz, diseño y color. No influyen en absoluto en el funcionamiento.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Idioma de la interfaz</td><td>El idioma de la ventana de ajustes y de los textos de la barra. No tiene nada que ver con el idioma de dictado: la ventana puede estar en español mientras dictas en inglés.</td></tr>" +
			"<tr><td>Diseño</td><td>El aspecto general. Hay siete: «Editor» (por defecto: gris oscuro con azul), «Fluent» (como la Configuración de Windows 11), «Estudio» (grafito con naranja, casi sin esquinas redondeadas), «Suave» (lila claro, redondeado), «Documento» (papel con azul), «Neón» (violeta con un degradado de rosa a turquesa) y «Terminal» (letras verdes sobre negro, corchetes alrededor de los botones). Los siete usan la misma tipografía, IBM Plex Mono.</td></tr>" +
			"<tr><td>Color</td><td>Un tono solo para «Terminal»: verde, ámbar, azul, rosa. La barra y el icono de la bandeja toman el mismo color. Los demás diseños tienen sus propios colores, y la fila se oculta.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Inicio y actualizaciones</p>" +
			"<p>Aquí decides si la aplicación arranca con el sistema y si está pendiente de las versiones nuevas.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Iniciar con Windows</td><td>La aplicación se abre al iniciar sesión y espera tus teclas. Es una entrada normal en la lista de inicio, igual que la que podrías añadir tú mismo.</td></tr>" +
			"<tr><td>Comprobar al iniciar</td><td>Una vez, al arrancar, la aplicación pregunta si ha salido una versión más nueva. Aparte de descargar modelos, nunca se conecta a internet. Si lo desactivas, tendrás que enterarte tú de las versiones nuevas.</td></tr>" +
			"<tr><td>Actualizaciones</td><td>El botón comprueba ahora mismo. Si hay una versión más nueva, la aplicación la descarga y la instala encima de la actual: los ajustes, los modelos y el historial se quedan donde están.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Servicio</p>" +
			"<p>Tres botones para cuando algo ha salido mal.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Registro</td><td>Abre el archivo donde la aplicación anota lo que ha hecho: qué modelo arrancó, cuánto tardó el reconocimiento, dónde se produjo un error. La causa suele estar escrita con palabras, no con códigos. El archivo no crece sin fin: lo nuevo va sustituyendo a lo antiguo.</td></tr>" +
			"<tr><td>Releer config.json</td><td>Todos los ajustes están en un solo archivo de texto junto a la aplicación. Si lo editaste a mano o lo trajiste de otro ordenador, este botón recoge los cambios sin reiniciar.</td></tr>" +
			"<tr><td>Restablecer los ajustes</td><td>Lo deja todo como recién instalado. Los modelos descargados y el historial se quedan donde están: solo se borran los ajustes.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Dónde funciona el reconocimiento</p>" +
			"<p>El habla la procesa un pequeño programa aparte. Normalmente está junto a la aplicación, arranca solo y no tienes que preocuparte por él. Pero también puede estar en otro ordenador: por ejemplo, en el de casa, que tiene una tarjeta gráfica potente, mientras tú dictas desde un portátil ligero. Se elige con el selector de las tarjetas.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>En este ordenador</b><span>hilos de CPU 4 · inicio automático activado · puerto 8910 · whisper-server.exe junto a la aplicación</span></div><button class=\"hbtn\">Editar</button></div>" +
			"<div class=\"hcap\">la tarjeta elegida está resaltada, y la otra, mientras tanto, atenuada</div>" +
			"</div>" +
			"<p>El botón «Editar» abre una ventana con cuatro ajustes. Todos vienen con valores por defecto, y casi siempre puedes dejarlos como están.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Hilos de CPU</td><td>Cuántos núcleos del procesador dedicar al reconocimiento. Más núcleos es más rápido, pero solo hasta cierto punto: a partir de ahí no se gana nada y los demás programas empiezan a ir lentos.</td></tr>" +
			"<tr><td>Iniciar whisper-server automáticamente</td><td>La aplicación arranca el reconocedor por su cuenta. Solo hay un motivo para desactivarlo: lo arrancas tú a mano y no quieres que la aplicación intervenga.</td></tr>" +
			"<tr><td>Puerto</td><td>El número por el que la aplicación habla con el reconocedor dentro de tu propio ordenador. Desde fuera no se ve. Cámbialo solo si otro programa ha ocupado ese número.</td></tr>" +
			"<tr><td>Ruta de whisper-server</td><td>Dónde está el reconocedor. Se rellena solo. Escribe la tuya solo si lo pusiste a propósito en otro sitio: si no, en cuanto muevas la carpeta, el dictado dejará de arrancar.</td></tr>" +
			"</table>" +
			"<p>La segunda tarjeta tiene un solo ajuste: la dirección del ordenador en la red, por ejemplo <code>http://192.168.1.50:8080</code>. Ten en cuenta que, en este caso, tu voz grabada va a esa máquina y, con una dirección normal, viaja sin cifrar. En una red doméstica no pasa nada, pero por internet es mejor no trabajar así.</p>" +
			"<div class=\"hbad\">Si está elegido «En otro ordenador» pero no se ha escrito la dirección, y tu idioma lo reconoce Whisper, el dictado no funcionará. La aplicación te lo dice con una tarjeta roja en «Estado», una marca roja junto a «Sistema» y una línea en la parte inferior de la ventana.</div>" +
			"<p class=\"wh\" id=\"s-mic\">La pestaña «Micrófono»</p>" +
			"<p>Más de la mitad de los dictados fallidos no se deben al modelo ni a los ajustes, sino al micrófono: se está escuchando el que no es o apenas se te oye. En esta pestaña está todo lo necesario para comprobarlo en medio minuto.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Micrófono</td><td>Qué dispositivo escuchar. La opción «Predeterminado del sistema» significa «el que esté elegido en Windows», algo cómodo si unas veces usas auriculares y otras no: la aplicación sigue al sistema por su cuenta.</td></tr>" +
			"<tr><td>Actualizar la lista</td><td>Hace falta si conectaste unos auriculares después de abrir la aplicación: la lista se vuelve a leer.</td></tr>" +
			"<tr><td>Nivel de entrada</td><td>Un medidor en directo. Di algo: las columnas tienen que moverse claramente. Si no se mueven, el micrófono está silenciado en Windows o está elegido el que no es.</td></tr>" +
			"<tr><td>Comprobar el micrófono</td><td>Tres segundos de grabación y un veredicto con palabras sencillas: si el volumen es suficiente, si hay saturación (cuando el sonido «cruje» porque el micrófono está demasiado cerca) y si de verdad se oye voz.</td></tr>" +
			"<tr><td>Señales sonoras</td><td>Un sonido corto al empezar y al terminar la grabación. Es especialmente útil si desactivaste la barra: así al menos oyes que la grabación ha empezado.</td></tr>" +
			"<tr><td>Sonido de señal</td><td>Qué sonido exactamente: voz de Windows, campanilla, marimba y algunos más. El botón de al lado reproduce el elegido, para que no tengas que adivinarlo por el nombre.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">La pestaña «Historial»</p>" +
			"<p>El historial es la lista de lo que has dictado. Solo el texto: la voz no se guarda nunca, en ningún sitio. Por defecto está desactivado, y es a propósito: en los dictados suele haber cosas personales.</p>" +
			"<p>Merece la pena activarlo si alguna vez te ha pasado algo así: pegaste un texto en la ventana equivocada y se perdió; necesitas repetir una frase larga; quieres encontrar lo que dictaste ayer y enviarlo otra vez.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Cuánto tiempo guardar</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">días</span></div>" +
			"<div class=\"hrow\"><label>Cuántas entradas guardar</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">el plazo se puede fijar en minutos, horas o días, con cualquier número</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Guardar el historial de dictados</td><td>El interruptor principal. Si lo desactivas, dejan de aparecer entradas nuevas y las antiguas se borran.</td></tr>" +
			"<tr><td>Cuánto tiempo guardar</td><td>Cuánto tarda una entrada en desaparecer sola. Puedes poner cinco minutos si el historial solo te sirve por si un texto se pega donde no debía, o treinta días.</td></tr>" +
			"<tr><td>Cuántas entradas guardar</td><td>Un segundo límite, por cantidad. Funciona junto con el primero: manda el que se alcance antes. Con «7 días y no más de 50», la entrada número cincuenta y uno desplaza a la más antigua, aunque todavía no haya pasado la semana.</td></tr>" +
			"<tr><td>No registrar nunca desde estos programas</td><td>Una lista de excepciones. Lo que dictes en estas ventanas (el banco, un gestor de contraseñas, el chat del trabajo) no entra en el historial, y el dictado funciona como siempre. Puedes elegir el programa en la lista de ventanas abiertas, sin buscar el nombre de su archivo.</td></tr>" +
			"</table>" +
			"<p>Con cada entrada puedes hacer cuatro cosas. <b>Copiar</b>: el texto va al portapapeles. <b>Pegar</b>: la aplicación te pide que le indiques dónde. Haz clic en el campo que quieras, en cualquier programa, y el texto aparece ahí. <b>Borrar</b>: quita esa entrada. <b>Buscar</b>: busca tanto en el texto como en el nombre del programa. El botón «Vaciar» lo borra todo de una vez, y por eso pide confirmación.</p>" +
			"<p class=\"hdim\">A la derecha de cada entrada pone cuándo desaparecerá, por ejemplo «hasta 7 sept», y el último día ese aviso se vuelve amarillo.</p>" +
			"<p class=\"wh\" id=\"s-dict\">La pestaña «Control y comportamiento»</p>" +
			"<p>Cómo empieza el dictado y cómo llega el texto al campo. Se entra aquí cuando algo no funciona como te gustaría: las teclas son incómodas, el texto no aparece, Enter se pulsa solo.</p>" +
			"<p class=\"hsub\">Teclas</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Atajo de teclado</td><td>Pulsa «Cambiar…» y mantén la combinación que te resulte cómoda. La aplicación captura esa combinación en todo el sistema, así que no elijas una que ya esté ocupada, como Ctrl+C. Las combinaciones con la tecla Windows son las que menos suelen estar ocupadas. Un idioma también puede tener su propio atajo, que se asigna en la pestaña «Idiomas y modelos».</td></tr>" +
			"<tr><td>Modo</td><td>Con «mantener», graba mientras tienes pulsadas las teclas: las sueltas y la frase pasa a procesarse. Con «alternar», la primera pulsación empieza la grabación y la segunda la detiene. El segundo modo va bien para monólogos largos, cuando la mano se cansa de mantener las teclas.</td></tr>" +
			"<tr><td>Detener por silencio</td><td>Solo en el modo «alternar». Cuando ya has dicho algo, una pausa tan larga como la fijada en «Pausa antes de detener, ms» termina la grabación por sí sola, como si pulsaras el atajo otra vez, y la barra muestra «Detenido por silencio». Nunca se detiene antes de tus primeras palabras, así que tras pulsar puedes ordenar tus ideas con calma; si no dices nada, el límite «Grab. máx, s» la termina igualmente. Con «Pulsar Enter tras insertar (envío auto)» activado, una pausa envía el mensaje, así que elige una pausa más larga.</td></tr>" +
			"<tr><td>Grab. mín, ms</td><td>Por debajo de qué duración una pulsación cuenta como accidental. Si rozas las teclas una fracción de segundo, la aplicación no hace nada en vez de insertar un trozo de palabra.</td></tr>" +
			"<tr><td>Grab. máx, s</td><td>Una protección contra una tecla olvidada: la grabación se corta sola en vez de seguir grabando a la nada durante media hora.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Cómo se inserta el texto</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Pulsar Enter tras insertar (envío auto)</td><td>El texto no solo se inserta, sino que se envía al momento: el mensaje sale, la búsqueda se lanza. En una aplicación de mensajería te ahorra un gesto; en un correo o un documento, casi seguro que estorba.</td></tr>" +
			"<tr><td>Restaurar portapapeles tras insertar</td><td>La inserción pasa por el portapapeles, y la aplicación vuelve a dejar ahí lo que había: una imagen, un archivo, un trozo de texto copiado. Déjalo activado si copias cosas a menudo y no quieres perderlas.</td></tr>" +
			"<tr><td>Escritura carácter a carácter</td><td>En lugar de pegar el texto de una vez, la aplicación lo teclea como si escribieras tú. Es más lento, pero te saca del apuro donde no se permite pegar desde el portapapeles: algunos formularios de bancos, terminales, programas antiguos.</td></tr>" +
			"<tr><td>Retraso antes de insertar</td><td>Una pequeña pausa antes de entregar el texto a la ventana. Rara vez hace falta: a veces el programa no llega a tiempo a recibir el texto y se pierde el principio de la frase.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Aviso en pantalla</p>" +
			"<p>Aquí decides si mostrar la barra y dónde exactamente.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Mostrar la barra</td><td>Desactívala si te molesta o si te bastan las señales sonoras.</td></tr>" +
			"<tr><td>Dónde mostrar la barra</td><td>Un esquema de la pantalla con nueve posiciones: haz clic en la esquina donde la barra no vaya a molestar.</td></tr>" +
			"<tr><td>En el cursor</td><td>La barra aparece junto al lugar donde estás dictando. Si el programa no indica dónde está el cursor, aparece junto al puntero del ratón.</td></tr>" +
			"<tr><td>Sitio propio</td><td>Puedes simplemente arrastrar la barra con el ratón, y recordará el sitio.</td></tr>" +
			"<tr><td>Pantalla</td><td>En qué monitor mostrarla, si tienes varios.</td></tr>" +
			"<tr><td>Mostrar el texto reconocido</td><td>Después de insertar, la barra muestra la propia frase en lugar del número de caracteres. Es cómodo para comprobarla, pero cualquiera que mire tu pantalla verá también la frase.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">La pestaña «Idiomas y modelos»</p>" +
			"<p>Aquí dices en qué idioma vas a dictar, y la aplicación decide con qué reconocerlo. Hay varios modelos: unos son rápidos, otros precisos, otros conocen idiomas poco comunes.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Idioma del habla</td><td>El idioma en el que hablas. Si hablas varios, elige «Auto»: el modelo lo averigua solo, aunque algo más despacio y con algo menos de precisión que cuando el idioma se indica directamente.</td></tr>" +
			"<tr><td>Qué modelo para qué idioma</td><td>Haz clic en un idioma y se despliega la lista de modelos que lo conocen, con una nota de velocidad y de precisión y el tamaño del archivo. Así, un idioma puede quedarse con un modelo rápido mientras otro recibe uno preciso. La fila aparte «Auto» es común: su modelo sirve a todos los idiomas que no tienen uno propio, y esos idiomas aparecen en gris en la lista. También es el que averigua qué idioma se habla cuando no está fijado el idioma del habla.</td></tr>" +
			"<tr><td>Atajo propio</td><td>Dentro de un idioma desplegado hay sitio para su propia combinación de teclas. Manténgala pulsada y el dictado se hace en ese idioma con su modelo, sus reemplazos y sus reglas, diga lo que diga el idioma del habla; el atajo principal conserva el idioma elegido arriba. Así se dicta en inglés con una combinación y en el idioma propio con otra, sin abrir los ajustes. Mientras graba, la barra muestra el código del idioma. Una combinación no puede repetir la principal ni la de otro idioma.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Ruso<span class=\"hs\">GigaAM v3 · precisión 5 de 5 · velocidad 5 de 5 · 232 MB</span></label><span class=\"hchip\"><i></i>activa</span></div>" +
			"<div class=\"hrow\"><label>Inglés<span class=\"hs\">como la detección automática · Nemotron 3.5</span></label><button class=\"hbtn\">Elegir</button></div>" +
			"<div class=\"hcap\">se ve qué modelo se encarga de cada idioma y lo que cuesta en velocidad y en espacio en disco</div>" +
			"</div>" +
			"<p class=\"hsub\">Traducción</p>" +
			"<p>Algunos modelos no solo saben transcribir lo que dices, sino también traducirlo al momento: hablas en un idioma y el texto llega al campo en otro. Quien traduce es Whisper: al inglés, con su propia función de traducción; a los demás idiomas, indicándole el idioma de salida como si fuera el idioma del habla; entre idiomas muy alejados la calidad es irregular. Si el modelo elegido no sabe traducir, junto al interruptor aparece una línea naranja: si Whisper está instalado, la barra ofrecerá traducir con él; si no, el texto se quedará en el idioma del habla.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Cambiar el idioma de salida del texto</td><td>Activa la traducción.</td></tr>" +
			"<tr><td>Idioma de salida del texto por defecto</td><td>A qué idioma traducir si no elegiste otro durante el propio dictado.</td></tr>" +
			"<tr><td>Preguntar el idioma de salida del texto</td><td>Puedes no preguntar y traducir siempre al idioma por defecto. Puedes preguntar cada vez: después del dictado, en la barra aparecen botones con idiomas. O puedes preguntar con tiempo límite: los mismos botones, pero si no pulsas ninguno a tiempo, se usa el idioma por defecto.</td></tr>" +
			"<tr><td>Idiomas del diálogo</td><td>Qué idiomas ofrecer como esos botones. Deja los dos o tres que de verdad usas; si no, elegir llevará más tiempo que hablar.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Modelos en disco</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Modelo propio</td><td>Si tienes un modelo propio, ponlo en la carpeta de modelos y reinicia la aplicación: aparecerá en las listas de los idiomas que entiende. El botón de al lado abre esa carpeta, para que no tengas que buscarla en el Explorador de archivos.</td></tr>" +
			"<tr><td>Descargar de la memoria</td><td>Mientras la aplicación está abierta, el modelo ocupa memoria RAM. Este botón la libera al momento, por ejemplo antes de un juego exigente. El siguiente dictado vuelve a cargar el modelo, y la primera frase tardará más de lo normal.</td></tr>" +
			"<tr><td>Comprobar los modelos instalados</td><td>Compara los archivos descargados con su referencia. Sirve si la conexión se cortó durante una descarga y sospechas que un archivo se dañó.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">La pestaña «Reglas»</p>" +
			"<p>Tres listas que corrigen el texto automáticamente, antes incluso de que llegue a la ventana. Las tres son opcionales: la aplicación funciona perfectamente sin ellas.</p>" +
			"<p class=\"hsub\">Puntuación y mayúsculas</p>" +
			"<p>De dónde salen los puntos, las comas y las mayúsculas del texto. Hay tres opciones: <b>del modelo</b> (el modelo los pone tal como los oyó; rápido y normalmente suficiente), <b>del modelo editor</b> (un segundo modelo añade los signos; queda más cuidado, pero cada dictado tarda unos segundos más) y <b>quitar</b> (el texto sale seguido, sin signos; rara vez hace falta, por ejemplo para un campo de búsqueda).</p>" +
			"<p class=\"hsub\">Diccionario</p>" +
			"<p>Una lista de palabras que el modelo no suele conocer: apellidos, nombres de empresas, términos profesionales. No es una orden de «escríbelo así», sino una pista: cuando oye algo parecido, el modelo tiende a elegir tu versión.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">las palabras se añaden de una en una y viven como tarjetas separadas</div>" +
			"</div>" +
			"<p class=\"hdim\">El diccionario solo lo leen los modelos Whisper. Nemotron, Parakeet, GigaAM y los demás modelos no aceptan pistas de ningún tipo: para ellos, la herramienta es la lista de reemplazos de más abajo. El conjunto de palabras por defecto cambia con el idioma hasta que escribas el tuyo.</p>" +
			"<p class=\"hsub\">Reemplazos tras el reconocimiento</p>" +
			"<p>También es una lista de palabras, pero funciona de otra manera: ya no es una pista, sino una garantía. Lo que aparezca a la izquierda se convertirá en lo de la derecha. Cada reemplazo tiene tres ajustes: <b>Solo palabras completas</b> (así «sol» no tocará «soldado», que normalmente es lo que quieres), <b>Distinguir mayúsculas y minúsculas</b> (así «León» y «león» cuentan como palabras distintas) e <b>Idioma de la regla</b> (el reemplazo solo actúa cuando dictas en ese idioma).</p>" +
			"<p class=\"hsub\">Comandos de voz</p>" +
			"<p>Frases que no entran en el texto, sino que ejecutan una acción. El salto de línea, el nuevo párrafo y tu propio texto se buscan como palabras completas en cualquier lugar de la frase. La cancelación y los comandos de forma solo actúan al final del dictado, y el borrado solo si se dice solo. Un comando de forma dicho solo se aplica al dictado siguiente, y la barra muestra «Siguiente frase: …».</p>" +
			"<ul>" +
			"<li><b>«nueva línea»</b>: saltar a la línea siguiente;</li>" +
			"<li><b>«nuevo párrafo»</b>: dejar una línea en blanco entre párrafos;</li>" +
			"<li><b>«cancelar»</b>: descartar la frase y no insertar nada;</li>" +
			"<li><b>«borra eso»</b>: borrar lo último que insertó el programa;</li>" +
			"<li><b>«como título»</b>: poner un signo delante de la frase, # por defecto, y una línea en blanco después;</li>" +
			"<li><b>«como lista»</b>: dividir la frase por las comas, los puntos y coma y la palabra «y» en líneas que empiezan con un guion;</li>" +
			"<li><b>«con mayúscula»</b>: empezar la frase con mayúscula;</li>" +
			"<li><b>«entre comillas»</b>: poner la frase entre comillas;</li>" +
			"<li><b>«sin punto»</b>: quitar el punto final;</li>" +
			"<li><b>tu propio comando</b>: insertar un texto preparado, como una firma, una dirección o una respuesta habitual.</li>" +
			"</ul>" +
			"<p>El borrado solo quita lo último que insertó el propio programa, y solo si desde entonces no has escrito ni hecho clic, la ventana es la misma, el texto no tenía saltos de línea, no se envió con Enter y no han pasado más de cinco minutos. El texto que ya estaba en la ventana nunca se toca.</p>" +
			"<p class=\"wh\" id=\"s-post\">La pestaña «Posprocesado»</p>" +
			"<p>El posprocesado es un segundo modelo que reescribe el texto reconocido según tu instrucción. Por ejemplo: «quita las muletillas y las repeticiones», «convierte esto en un correo formal», «ponlo en forma de lista». La instrucción se escribe con palabras normales y se llama prompt.</p>" +
			"<p>Por defecto está desactivado, y este es el motivo: es el único paso que ralentiza claramente el dictado. En lugar de uno o dos segundos, la inserción tardará varios, y en un equipo poco potente, hasta diez.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Activar el posprocesado</td><td>El interruptor principal. Mientras está desactivado, el resto de la pestaña no funciona y el texto se inserta tal como se reconoció.</td></tr>" +
			"<tr><td>Modelo · Local</td><td>Un editor que funciona en tu ordenador. No sale nada de él, pero necesita memoria: estos modelos pesan 1 GB o más. Al lado hay un catálogo donde se elige y se descarga un modelo con un clic.</td></tr>" +
			"<tr><td>Modelo · Servidor externo</td><td>En lugar de tu propio modelo, el texto se envía a un servicio externo: OpenAI, Groq o tu propio servidor. Es más rápido y más listo, pero el texto sale de tu ordenador. La clave de acceso se guarda cifrada, y el botón «Probar la conexión» comprueba que la dirección y la clave funcionan.</td></tr>" +
			"<tr><td>Prompts</td><td>Tus instrucciones para el editor. Los marcados se aplican uno tras otro, de arriba abajo: primero, por ejemplo, la limpieza y luego el estilo formal. Si no hay nada marcado, el posprocesado no hace nada, y la aplicación te avisa con una línea amarilla.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">La voz nunca va a un servidor externo, solo el texto reconocido. Pero ese texto sí sale, así que actívalo con conocimiento de causa.</div>" +
			"<p>Cada prompt tiene un <b>nombre</b> (con él aparece en la barra mientras trabaja), el <b>texto de la instrucción</b> con palabras normales, un <b>atajo propio</b> opcional (al dictar con él, la frase pasa solo por este prompt, sin los demás) y una <b>prueba</b>: un campo donde puedes escribir una frase y ver qué sale, sin dictar nada.</p>" +
			"<p class=\"wh\" id=\"s-subs\">La pestaña «Subtítulos»</p>" +
			"<p>Los subtítulos transcriben lo que suena en el ordenador —una llamada, un vídeo, una clase— y muestran las palabras en una ventana pequeña que se queda encima de las demás. Es el mismo reconocimiento que el del dictado, solo que el sonido sale de los altavoces y no del micrófono, y el dictado sigue funcionando mientras hay subtítulos. Los sonidos de la propia aplicación no se transcriben.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Iniciar subtítulos</td><td>Los inicia y los detiene. La línea de al lado dice qué pasa: apagados, cargando el modelo, escuchando o reiniciando el reconocimiento. El menú de la bandeja, un atajo propio y la cruz de la ventana hacen lo mismo.</td></tr>" +
			"<tr><td>Idioma del habla</td><td>El idioma de lo que escuchas, aparte del idioma de dictado. El modelo es el que la pestaña «Idiomas y modelos» da a ese idioma. Un modelo de streaming escribe las palabras a medida que se dicen; los demás, cada frase cuando termina. Un modelo Whisper necesita muchos segundos por frase en un procesador, así que su texto va muy por detrás del habla; las frases acumuladas le llegan juntas y salen como una sola línea.</td></tr>" +
			"<tr><td>Traducción</td><td>Muestra en la ventana la traducción en lugar del texto reconocido, en el idioma que elijas; el título de la ventana dice entonces «Traducción» y los idiomas, por ejemplo EN → RU. El habla se corta en frases sobre la marcha, y cada frase se traduce al momento. Un fragmento termina al final de una oración o antes de una palabra de enlace, y el traductor recibe con él los fragmentos anteriores y sus propias traducciones, así que una oración cortada en partes se traduce de forma coherente. La traducción aparece de forma pareja, por palabras enteras, al ritmo de lectura, y más rápido cuando se ha acumulado mucho, así que un trozo grande no hace saltar el texto. Las frases de un mismo tramo de habla forman un párrafo, y una pausa empieza una fila nueva. Una frase que no se pudo traducir se muestra tal como se oyó, en un color más tenue. «Mostrar también el texto original» divide la ventana con una línea fina: arriba, más tenue, lo que se dijo, y abajo la traducción; la línea se arrastra con el ratón. La fila «Quién traduce» decide quién lo hace. «Alta calidad» es el modelo de traducción en la tarjeta gráfica: una descarga de 2,3 GB, para una tarjeta gráfica con al menos 4 GB de memoria propia; se carga en la tarjeta cuando empiezan los subtítulos y sale de ella cuando terminan, y un enlace al lado abre sus condiciones de uso. «Tu propio servidor» es el servidor de posprocesado configurado en la pestaña «Posprocesado», y entonces cada frase se envía a ese servidor. «Elige el programa» toma el modelo de traducción si está descargado y la tarjeta gráfica sirve; si no, tu servidor si hay uno configurado; si no, el modelo editor instalado en la pestaña «Posprocesado», que traduce de forma aproximada: el sentido se entiende, pero algunas palabras pueden no ser exactas. Un traductor elegido a mano nunca se sustituye por otro: si no puede trabajar, el bloque dice lo que falta y el texto se muestra tal como se oyó. La traducción se guarda en un archivo aparte junto a la transcripción: cada fragmento con la hora en que empezó; un fragmento que no se pudo traducir se escribe tal como se oyó. Sin traductor no hay traducción, y los subtítulos funcionan como antes.</td></tr>" +
			"<tr><td>Atajo de teclado</td><td>Inicia y detiene los subtítulos desde cualquier ventana. Nunca empieza un dictado y no puede repetir otro atajo.</td></tr>" +
			"<tr><td>Guardar la transcripción</td><td>Cada sesión escribe un archivo de texto en la carpeta transcripts junto al programa: una primera línea con la fecha, el idioma y el modelo, y luego cada línea terminada con la hora en que empezó. No se escribe nada hasta que la primera línea está terminada. Con la traducción activada, el nombre de la transcripción lleva original y el código del idioma hablado, y a su lado aparece un segundo archivo con translation y el código del idioma de destino, por ejemplo 2026-11-12_10-18-29_original_EN.txt y 2026-11-12_10-18-29_translation_RU.txt.</td></tr>" +
			"<tr><td>Ventana de subtítulos</td><td>Dónde está, qué tamaño tiene y qué tamaño tienen las letras. El lugar se elige en el dibujo de una pantalla, igual que para la barra de dictado: un clic en uno de los nueve lugares pone la ventana ahí, y con «Sitio propio» puede ir a cualquier parte, arrastrando la marca del dibujo o la propia ventana; con varios monitores se puede elegir la pantalla. El tamaño se cambia con el ratón por los bordes y las esquinas; la altura se ajusta a filas enteras, y el ancho y las filas se recuerdan para cada resolución de pantalla. «Filas de texto» vale donde aún no hay tamaño, y «Restablecer el tamaño» olvida los que hay. La rueda del ratón sobre la ventana desplaza el texto hacia atrás tres filas por paso; mientras lees lo anterior, el texto nuevo no lo mueve, y la marca «a lo último ↓» del título vuelve al habla actual.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Las personas con las que hablas deben saber que se transcribe lo que dicen.</div>" +
			"<p class=\"wh\" id=\"s-about\">«Acerca de» y «Contactos»</p>" +
			"<p>La pestaña «Acerca de» tiene un breve texto sobre qué es la aplicación y cómo está hecha, y la lista de todas las bibliotecas de terceros con las que está construida, con sus licencias. Los nombres son enlaces y llevan a las páginas de esos proyectos. También hay un enlace al repositorio, donde están el código fuente y todas las versiones.</p>" +
			"<p>«Contactos» tiene tres tarjetas: el correo, el repositorio y una página para fallos e ideas. Al hacer clic en la dirección de correo se abre tu programa de correo con el destinatario ya puesto, y el icono de al lado simplemente copia la dirección.</p>" +
			"<p class=\"wh\" id=\"s-files\">Qué hay en la carpeta de la aplicación</p>" +
			"<p>La aplicación es portátil: todo lo suyo está en una sola carpeta y no se esparce por ningún otro sitio.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Todos tus ajustes en un solo archivo de texto. Puedes copiarlo a otro ordenador para no tener que configurarlo todo de nuevo.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>El registro. Lo primero que conviene mirar si la aplicación se comporta de forma rara.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>El historial de dictados, si lo activaste. Texto normal: si quieres, puedes simplemente borrar el archivo.</td></tr>" +
			"<tr><td>models</td><td>La carpeta con los modelos descargados, la parte más pesada: varios GB. También se puede llevar a otro ordenador para no tener que descargarlo todo otra vez.</td></tr>" +
			"<tr><td>transcripts</td><td>Las transcripciones de los subtítulos, un archivo de texto por sesión, y sus traducciones cuando la traducción está activada. Son tus documentos: el desinstalador solo los borra si también le pides borrar los ajustes y los datos.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>Los pequeños programas que ejecutan los modelos. Sin ellos, un modelo es solo un archivo sin nada que lo lea.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Qué sale a internet y qué no</p>" +
			"<ul>" +
			"<li><b>La voz</b> no se guarda ni se envía, salvo que tú mismo hayas activado el reconocimiento en otro ordenador.</li>" +
			"<li><b>Los subtítulos</b> se reconocen aquí, igual que el dictado, y la transcripción se queda en la carpeta transcripts; con el reconocimiento en otro ordenador, las frases van a ese servidor; con la traducción a través de tu propio servidor de posprocesado, cada línea terminada va también allí.</li>" +
			"<li><b>El texto</b> no sale del equipo, salvo que tú mismo hayas activado un servidor externo de posprocesado.</li>" +
			"<li><b>El historial</b> es un archivo junto a la aplicación y se limpia según el plazo que tú fijes.</li>" +
			"<li><b>La red</b> se necesita dos veces: para descargar un modelo y para comprobar si hay actualizaciones. La segunda se puede desactivar.</li>" +
			"<li><b>La clave</b> de un servicio externo se guarda cifrada por Windows y nunca llega al registro.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Si algo sale mal</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>El texto no aparece en el campo</td><td>El programa no acepta pegar desde el portapapeles. Activa «Escritura carácter a carácter» en la pestaña «Control y comportamiento» o añade un pequeño retraso antes de insertar.</td></tr>" +
			"<tr><td>El dictado no empieza</td><td>El atajo de teclado lo usa otro programa. Elige otro, mejor con la tecla Windows: esas combinaciones se ocupan menos.</td></tr>" +
			"<tr><td>Solo se grabó silencio</td><td>El micrófono no es el correcto o está silenciado. Abre la pestaña «Micrófono», mira el medidor de nivel y pulsa «Comprobar el micrófono».</td></tr>" +
			"<tr><td>La primera frase tarda mucho</td><td>Es lo previsto: el modelo se carga en la memoria la primera vez que se usa. Después, el reconocimiento lleva segundos.</td></tr>" +
			"<tr><td>Pone «Reconocimiento no disponible»</td><td>Está elegido el reconocimiento en otro ordenador, pero no hay dirección. Escribe la dirección o vuelve a «En este ordenador».</td></tr>" +
			"<tr><td>El posprocesado no cambia nada</td><td>No hay ningún prompt marcado, no hay modelo elegido o el servidor externo no responde. El motivo está escrito en la propia pestaña, junto al interruptor.</td></tr>" +
			"<tr><td>El texto sale raro, con palabras cambiadas</td><td>Mira en «Reglas»: lo más probable es que haya actuado un reemplazo o un comando de voz.</td></tr>" +
			"<tr><td>Nada de lo anterior</td><td>«Sistema» → «Abrir el registro». La aplicación anota ahí todo lo que hace, y el motivo suele estar explicado con palabras.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">¿No encontraste la respuesta? Escribe un correo desde la pestaña «Contactos» o abre un issue en el repositorio. El autor lee las dos cosas.</p>",
		"S_WIZ_HELLO":      "Dos minutos de configuración",
		"S_WIZ_HELLO_TEXT": "{app} convierte tu voz en texto justo donde está el cursor: mantén la combinación, di una frase, suelta y el texto aparece. Todo se procesa en tu equipo; el audio nunca sale de él.",
		"S_WIZ_LATER":      "Todo lo que elijamos ahora se puede cambiar después en los ajustes.",
		"S_WIZ_MODEL_TEXT": "El modelo para tu idioma se elige por ti y se descarga en segundo plano.",
		"S_WIZ_PICKED":     "Elegido para: %s",
		"S_WIZ_DL_BG":      "Se descarga en segundo plano, puedes seguir",
		"S_WIZ_DL_FAIL":    "La descarga se detuvo: %s",
		"S_WIZ_RETRY":      "Volver a intentar",
		"S_WIZ_TRY_DL":     "El modelo aún se está descargando (%d%%). El campo funcionará en cuanto esté listo.",
		"S_WIZ_TRY_LOAD":   "Cargando el modelo…",
		"S_WIZ_TRY_LIVE":   "Con %s las palabras aparecen en la barra mientras todavía hablas.",
		"S_WIZ_DONE_DL":    "El modelo aún se está descargando. El dictado empezará a funcionar solo en cuanto esté listo.",
		"S_WIZ_INPUT_TEXT": "Esta es la combinación que mantendrás pulsada al hablar. Di algo y comprueba que la barra de nivel se mueve.",
		"S_WIZ_T_TRY":      "Prueba",
		"S_WIZ_TRY_PH":     "el texto aparecerá aquí",
		"S_WIZ_T_DONE":     "Listo",
		"S_WIZ_DONE_TEXT":  "{app} vive en la bandeja: clic izquierdo en el icono para los ajustes, clic derecho para el menú. Puedes dictar en cualquier ventana que tenga cursor de texto.",
		"S_AUTORUN":        "Iniciar con Windows",
		"S_AUTORUN_SUB":    "Una entrada en el inicio del usuario actual",
		"S_WIZ_SKIP":       "Omitir",
		"S_WIZ_BACK":       "Atrás",
		"S_WIZ_NEXT":       "Siguiente",
		"S_WIZ_FINISH":     "Finalizar",
		"S_WIZ_WAIT":       "Esperando la primera frase…",
		"S_WIZ_HEARD":      "Escuché:",
		"S_WIZ_HAVE":       "Todo lo necesario ya está descargado",
		"S_WIZ_TRY_TEXT":   "Coloca el cursor en el campo de abajo, mantén %s, di una frase y suelta.",
	}
	settingsStrings["it"] = map[string]string{
		"S_SUBS_ON": "I sottotitoli sono attivi",
		"S_SUBS_HERO_HINT": "avviali — una chiamata o un video diventano testo sullo schermo",
		"S_SUBS_HERO_START": "Avvia manualmente",
		"S_SUBS_HERO_STOP": "Ferma",
		"S_UPDATED_TO": "Aggiornato alla versione %s",
		"S_SUBS_HERO_LIVE": "%s · parole mentre vengono dette",
		"S_SUBS_HERO_PHRASES": "%s · testo dopo ogni frase",
		"S_NAV_SUBS": "Sottotitoli",
		"S_SUBS_HINT": "L'app trascrive ciò che il computer riproduce — una chiamata, un video — e mostra le parole in una finestra sopra le altre.",
		"S_SUBS_NOTICE": "Le persone con cui parli devono sapere che le loro parole vengono trascritte.",
		"S_SUBS_START": "Avvia i sottotitoli",
		"S_SUBS_STOP": "Ferma i sottotitoli",
		"S_SUBS_OFF": "I sottotitoli sono spenti",
		"S_SUBS_LOADING": "Caricamento di %s…",
		"S_SUBS_LISTEN_LIVE": "In ascolto del suono del computer · %s, parole mentre vengono dette",
		"S_SUBS_LISTEN_PHRASES": "In ascolto del suono del computer · %s, testo dopo ogni frase",
		"S_SUBS_RESTARTING": "Il riconoscimento si sta riavviando…",
		"S_SUBS_LANG": "Lingua del parlato",
		"S_SUBS_LANG_SUB": "La lingua di ciò che ascolti, distinta dalla lingua di dettatura",
		"S_SUBS_M_LIVE": "Scriverà %s: le parole appaiono mentre vengono dette.",
		"S_SUBS_M_PHRASES": "Scriverà %s: il testo appare dopo la fine di ogni frase.",
		"S_SUBS_M_ABSENT": "Per questa lingua è scelto %s, ma non è installato.",
		"S_SUBS_M_WHISPER": "%s è un modello Whisper: sul processore gli servono molti secondi per frase, quindi il testo resterà molto indietro rispetto al parlato. Per i sottotitoli va meglio un modello in streaming o uno veloce.",
		"S_SUBS_GO_MODELS": "Apri «Lingue e modelli»",
		"S_SUBS_REMOTE": "Il riconoscimento avviene su un altro computer: ogni frase viene inviata a quel server.",
		"S_SUBS_TR": "Traduzione",
		"S_SUBS_TR_SUB": "La finestra mostra la traduzione al posto del testo riconosciuto",
		"S_SUBS_TR_OFF": "Disattivata",
		"S_SUBS_TR_LOCAL": "%s traduce: la traduzione arriva frase per frase, un attimo dopo il parlato.",
		"S_SUBS_TR_API": "Il tuo server di post-elaborazione (%s) traduce: ogni frase gli viene inviata, e la traduzione arriva frase per frase, un attimo dopo il parlato.",
		"S_SUBS_TR_NONE": "Non c'è né un modello di traduzione, né un server di post-elaborazione, né un modello editor, quindi non ci sarà traduzione.",
		"S_SUBS_TR_GO": "Apri «Post-elaborazione»",
		"S_SUBS_TR_SAME": "La lingua della traduzione è la stessa del parlato: non c'è nulla da tradurre.",
		"S_SUBS_WHO": "Chi traduce",
		"S_SUBS_WHO_SUB": "Il programma prende il migliore disponibile, oppure ne fissi uno",
		"S_SUBS_WHO_AUTO": "Sceglie il programma",
		"S_SUBS_WHO_QUALITY": "Alta qualità: un modello sulla scheda grafica",
		"S_SUBS_WHO_SERVER": "Il tuo server",
		"S_SUBS_TRM": "Modello di traduzione",
		"S_SUBS_TRM_SUB": "%s, funziona su una scheda grafica con almeno 4 GB di memoria propria",
		"S_SUBS_TRM_ABSENT": "Non scaricato",
		"S_SUBS_TRM_GETTING": "Download in corso… %s%",
		"S_SUBS_TRM_CHECKING": "Controllo del file…",
		"S_SUBS_TRM_READY": "Pronto",
		"S_SUBS_TRM_RUNNING": "Pronto, caricato sulla scheda grafica",
		"S_SUBS_TRM_BAD": "Il file nella cartella models è danneggiato o incompleto",
		"S_SUBS_TRM_ERR": "Non scaricato: %s",
		"S_SUBS_TRM_SMALL": "Questo computer non può usarlo: la sua scheda grafica ha %s MB di memoria propria, ne servono 4 GB.",
		"S_SUBS_TRM_INTEGRATED": "Questo computer non può usarlo: ha solo la grafica integrata, che condivide la memoria del computer.",
		"S_SUBS_TRM_NOCARD": "Questo computer non può usarlo: Windows non segnala alcuna scheda grafica.",
		"S_SUBS_TRM_UNKNOWN": "Questo computer non può usarlo: questa versione di Windows non descrive la sua scheda grafica.",
		"S_SUBS_TRM_DOWN": "Non è partito: %s",
		"S_SUBS_TRM_DOWN_NEXT": "Fino al prossimo avvio dei sottotitoli traduce il successivo in ordine.",
		"S_SUBS_TRM_DOWN_HEARD": "Fino al prossimo avvio dei sottotitoli il testo è mostrato come è stato sentito.",
		"S_SUBS_TRM_DEL_SURE": "Eliminare?",
		"S_SUBS_TRM_W_VULKAN": "il driver grafico non offre Vulkan",
		"S_SUBS_TRM_W_DEVICE": "il motore non vede questa scheda grafica",
		"S_SUBS_TRM_W_MEMORY": "non c'è abbastanza memoria libera sulla scheda grafica",
		"S_SUBS_TRM_W_CRASH": "il motore si è fermato durante il caricamento",
		"S_SUBS_TRM_W_TIMEOUT": "il motore non ha risposto in tempo",
		"S_SUBS_TRM_W_PORT": "la sua porta è occupata da un altro programma",
		"S_SUBS_TRM_W_EXE": "manca llama-server-gpu.exe accanto al programma",
		"S_SUBS_TRM_DL": "Scarica",
		"S_SUBS_TRM_CANCEL": "Annulla",
		"S_SUBS_TRM_DEL": "Elimina",
		"S_SUBS_TRM_TERMS": "Condizioni d'uso del modello",
		"S_SUBS_TR_QUALITY": "Traduce il modello di traduzione sulla scheda grafica (%s): la traduzione arriva frase per frase, un attimo dopo il parlato.",
		"S_SUBS_TR_MISS_CARD": "È scelta l'opzione «Alta qualità», ma la scheda grafica di questo computer non va bene, quindi il testo è mostrato come è stato sentito.",
		"S_SUBS_TR_MISS_MODEL": "È scelta l'opzione «Alta qualità», ma il modello di traduzione non è scaricato, quindi il testo è mostrato come è stato sentito.",
		"S_SUBS_TR_MISS_START": "È scelta l'opzione «Alta qualità», ma il modello di traduzione non è partito sulla scheda grafica, quindi fino al prossimo avvio dei sottotitoli il testo è mostrato come è stato sentito.",
		"S_SUBS_TR_MISS_SERVER": "È scelto il tuo server, ma non è impostato alcun server di post-elaborazione, quindi il testo è mostrato come è stato sentito.",
		"S_SUBS_ORIG": "Mostra anche il testo originale",
		"S_SUBS_ORIG_SUB": "la finestra mostra sopra la traduzione ciò che è stato detto; la linea tra i due si trascina con il mouse",
		"S_SUBS_KEY": "Scorciatoia da tastiera",
		"S_SUBS_KEY_SUB": "Avvia e ferma i sottotitoli da qualsiasi finestra",
		"S_SUBS_SAVE": "Salva la trascrizione",
		"S_SUBS_SAVE_SUB": "Un file di testo per sessione nella cartella transcripts accanto al programma, e un secondo con la traduzione quando la traduzione è attiva",
		"S_SUBS_FOLDER": "Cartella delle trascrizioni",
		"S_SUBS_WINDOW": "Finestra dei sottotitoli",
		"S_SUBS_ROWS": "Righe di testo",
		"S_SUBS_SIZE": "Dimensione del testo",
		"S_SUBS_SIZE_S": "Piccolo",
		"S_SUBS_SIZE_M": "Medio",
		"S_SUBS_SIZE_L": "Grande",
		"S_SUBS_BOX": "Dimensione della finestra",
		"S_SUBS_BOX_SUB": "trascina i bordi e gli angoli della finestra con il mouse; la dimensione viene ricordata per ogni risoluzione dello schermo, e «Righe di testo» vale finché non c'è",
		"S_SUBS_BOX_BTN": "Ripristina la dimensione",
		"S_SUBS_BOX_DONE": "La finestra torna alla dimensione predefinita",
		"S_SUBS_FREE": "Posto tuo",
		"S_SUBS_FREE_SUB": "la finestra si può trascinare ovunque, la posizione viene ricordata per ogni risoluzione dello schermo",
		"S_SUBS_POS": "Dove mostrare la finestra",
		"S_SUBS_POS_SCHEME_SUB": "fai clic sullo schermo: la finestra va lì",
		"S_SUBS_POS_DRAG_SUB": "trascina il segno o la finestra stessa: va dove vuoi",
		"S_SUBS_MON": "Schermo",
		"S_SUBS_MON_SUB": "su quale monitor mostrare la finestra dei sottotitoli",
		"S_TITLE": "{app} — Impostazioni", "S_DICT_HINT": "Termini, nomi e abbreviazioni separati da virgole — un suggerimento per l'ascolto, non comandi. Vale per Whisper; il russo tramite GigaAM lo ignora. Il set predefinito segue la lingua di riconoscimento finché non scrivi il tuo.",
		"S_TR_DEFAULT": "Cambiare la lingua di uscita del testo", "S_TR_TARGET": "Lingua di uscita del testo predefinita", "S_TR_ASK": "Chiedere la lingua di uscita del testo", "S_TR_ASK_NEVER": "Non chiedere — tradurre subito",
		"S_SRCLANG_SUB": "la parlate voi; determina il modello di riconoscimento",
		"S_TR_LANGS_SUB": "queste lingue diventano pulsanti sulla barra all’inserimento",
		"S_TR_UNAVAIL": "%s non sa tradurre — il testo resterà nella lingua di partenza",
		"S_TR_GET_WHISPER": "Scarica %s per la traduzione (%d MB)",
		"S_TR_STANDIN": "%s non sa tradurre — la barra proporrà %s",
		"S_TR_LOCK": "%s non può essere tolta dall’elenco — è la lingua di uscita del testo predefinita. Scegliete un’altra lingua predefinita e allora %s potrà essere esclusa.",
		"S_TR_LOCK_OK": "Capito",
		"S_TR_ONE": "Sono spuntate più lingue, ma senza domanda la traduzione andrà sempre in una sola — %s (lingua di uscita predefinita). Le altre restano spuntate ma disabilitate.",
		"S_TR_NOMODEL": "%s non sa tradurre. Se continuate, la traduzione sarà spenta e non disponibile finché lavora questo modello.",
		"S_TR_CONFIRM": "Confermare",
		"S_TR_ASK_ALWAYS": "Chiedere ogni volta", "S_TR_ASK_TIMEOUT": "Chiedere, con timeout", "S_TR_SECONDS": "Timeout, s",
		"S_TR_LANGS": "Lingue nel dialogo",
		"S_LLM_HINT": "I profili spuntati si applicano uno dopo l'altro, dall'alto in basso, nella dettatura normale. Niente spuntato: il testo viene inserito così com'è.",
		"S_PROF_ADD": "Aggiungi", "S_PROF_NAME": "Nome", "S_PROF_PROMPT": "Prompt", "S_PROF_TEST": "Prova",
		"S_HOTKEY": "Scorciatoia da tastiera", "S_CHANGE": "Cambia…", "S_UILANG": "Lingua dell'interfaccia", "S_AUTO": "Come il sistema",
		"S_SEC_SOUND": "Suono", "S_BEEP": "Segnali acustici", "S_SOUND": "Suono del segnale",
		"S_SND_SPEECH": "Voce di Windows", "S_SND_CHIME": "Campanello", "S_SND_SOFT": "Morbido", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Premi Invio dopo l'inserimento (invio auto)", "S_RESTORE": "Ripristina appunti dopo l'inserimento",
		"S_NAV_HISTORY": "Cronologia", "S_HIST_ON": "Conservare la cronologia delle dettature", "S_HIST_ON_SUB": "solo testo, su questo computer; l'audio non viene mai salvato",
		"S_HIST_DAYS": "Quanti giorni conservare", "S_HIST_MAX": "Quante voci conservare",
		"S_HIST_SKIP": "Non registrare mai da questi programmi", "S_HIST_SKIP_SUB": "separati da virgole: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Aggiunta di un programma", "S_SKIP_EDIT_DLG": "Modifica del programma", "S_SKIP_NAME": "Nome del programma", "S_SKIP_NAME_SUB": "Il nome del file senza percorso: keepass.exe. Un asterisco finale prende tutte le versioni: 1password*", "S_SKIP_OPEN": "Programmi aperti adesso", "S_SKIP_REFRESH": "Aggiorna l'elenco", "S_SKIP_PICKED": "%d di %d scelti", "S_SKIP_NONE": "Niente scelto", "S_SKIP_EMPTY": "L'elenco è vuoto — la cronologia si tiene da tutti i programmi", "S_SKIP_ADD_BTN": "Aggiungi un programma", "S_SKIP_HINT": "Quello che detti in questi programmi non finisce nella cronologia. L'incollaggio funziona come sempre.",
		"S_HIST_LIST": "Voci", "S_HIST_CLEAR": "Svuota", "S_HIST_TILL": "fino al %s", "S_HIST_TILL1": "fino a domani", "S_HIST_TILL_FULL": "Sarà eliminato il %s — conservazione %s", "S_HIST_LIST_HINT": "Quello che è stato dettato: copiare, incollare in qualsiasi finestra o cancellare.", "S_HIST_COPY": "Copia",
		"S_HIST_KEEP": "Per quanto conservare",
		"S_UNIT_MIN": "minuti",
		"S_UNIT_HOUR": "ore",
		"S_UNIT_DAY": "giorni",
		"S_HIST_FIND": "Cerca nella cronologia…", "S_HIST_EMPTY": "Ancora nessuna cronologia", "S_HIST_ASK": "Eliminare tutta la cronologia delle dettature?",
		"S_SEC_CMD": "Comandi vocali", "S_CMD_HINT": "Ciò che dici diventa un a capo, un segno, un annullamento, la cancellazione dell'ultimo inserimento o una forma per la frase invece di finire nel testo. Riconosciuti come parole intere, applicati dall'alto in basso, dopo le sostituzioni; i comandi di forma agiscono alla fine della dettatura.",
		"S_CMD_ADD": "Aggiungi un comando", "S_CMD_PRESET": "Aggiungi quelli soliti", "S_CMD_PH": "la frase che dirai",
		"S_CMD_NEWLINE": "a capo", "S_CMD_PARAGRAPH": "nuovo paragrafo", "S_CMD_TEXT": "inserire testo", "S_CMD_CANCEL": "annullare la dettatura",
		"S_CMD_TEXT_PH": "cosa inserire", "S_CMD_EMPTY": "Ancora nessun comando", "S_CMD_DEL": "Elimina il comando",
		"S_CMD_P_NEWLINE": "nuova riga", "S_CMD_P_PARAGRAPH": "nuovo paragrafo", "S_CMD_P_CANCEL": "annulla",
		"S_CMD_UNDO": "cancellare l'ultimo inserimento", "S_CMD_HEADING": "trasformare in titolo", "S_CMD_LIST": "trasformare in elenco",
		"S_CMD_CAPITAL": "iniziale maiuscola", "S_CMD_QUOTES": "mettere tra virgolette", "S_CMD_NOSTOP": "togliere il punto finale",
		"S_CMD_HEADING_PH": "segno prima del titolo, ad es. #", "S_CMD_LIST_PH": "parole che uniscono l'ultima voce, ad es. e", "S_CMD_QUOTES_PH": "segni di apertura e chiusura, ad es. «»",
		"S_CMD_END_NOTE": "Dillo alla fine della dettatura; detto da solo, dà forma alla successiva.",
		"S_CMD_UNDO_NOTE": "Dillo da solo. Cancella l'ultimo inserimento se da allora non hai scritto né cliccato e la finestra è la stessa; il testo che c'era già non viene mai toccato.",
		"S_CMD_P_UNDO": "cancella questo", "S_CMD_P_HEADING": "come titolo", "S_CMD_P_LIST": "come elenco", "S_CMD_P_CAPITAL": "con maiuscola",
		"S_CMD_P_QUOTES": "tra virgolette", "S_CMD_P_NOSTOP": "senza punto", "S_CMD_P_AND": "e, ed", "S_CMD_P_MARKS": "«»",
		"S_SEC_REPLACE": "Sostituzioni dopo il riconoscimento", "S_REPLACE_HINT": "Ciò che è stato sentito male diventa ciò che intendevi — subito dopo il riconoscimento, prima dei prompt. Applicate dall'alto in basso.",
		"S_REPL_WHOLE_FULL": "Solo parole intere", "S_REPL_CASE_FULL": "Distingui le maiuscole", "S_CMD_ACTION": "Azione",
		"S_FM_ADD": "Aggiungi",
		"S_TIP_REPL_LANG": "La regola agisce solo quando detti nella lingua scelta. «tutte le lingue» — agisce sempre.",
		"S_TIP_REPL_CASE": "Maiuscole e minuscole contano: «git» e «Git» sono parole diverse. Disattivato — il caso è ignorato.",
		"S_TIP_REPL_WHOLE": "La sostituzione agisce solo se il testo è una parola a sé. Disattivato — anche dentro altre parole.",
		"S_TIP_CMD_ACTION": "Cosa succede quando pronunci la frase: a capo, nuovo paragrafo, un tuo testo, annullamento della dettatura, cancellazione dell'ultimo inserimento o una forma per la frase — titolo, elenco, iniziale maiuscola, virgolette o senza punto finale.",
		"S_LIST_FILTER_PH": "cerca…",
		"S_REPL_DEL": "Elimina la sostituzione",
		"S_LIST_NOTHING": "Nessun risultato: «%s»",
		"S_FM_T_REPL_ADD": "Aggiunta di una sostituzione", "S_FM_T_REPL_EDIT": "Modifica della sostituzione",
		"S_FM_T_CMD_ADD": "Aggiunta di un comando", "S_FM_T_CMD_EDIT": "Modifica del comando",
		"S_MT_DEL": "Eliminazione del modello", "S_MT_DEL_PROMPT": "Eliminazione del prompt", "S_MT_DL": "Download del modello",
		"S_MT_TR_OFF": "Disattivazione della traduzione", "S_MT_TR_ONE": "Traduzione senza chiedere", "S_MT_TR_LOCK": "Lingua di uscita predefinita",
		"S_MT_REMOTE": "Server remoto", "S_MT_POST": "Server esterno", "S_MT_HIST": "Svuotamento della cronologia",
		"S_MT_RESET": "Ripristino delle impostazioni", "S_MT_EXE": "Percorso del server",
		"S_DICT_ADD": "Aggiungi una parola", "S_FM_T_DICT_ADD": "Aggiunta di una parola", "S_DICT_EMPTY": "Ancora nessuna parola",
		"S_DICT_ADD_PH": "una parola o più separate da virgole",
		"S_DICT_NOMODEL": "Il modello attuale %s non supporta il dizionario: solo i modelli Whisper lo leggono.",
		"S_OV_FREE": "Posto tuo", "S_OV_FREE_SUB": "la barra si può trascinare ovunque",
		"S_OVPOS_DRAG_SUB": "trascina la barra con il mouse: va dove vuoi",
		"S_OVMON_N": "Schermo %d",
		"S_POST_ENABLE": "Attiva la post-elaborazione",
		"S_API_SUM_URL": "indirizzo", "S_API_SUM_MODEL": "modello", "S_API_SUM_KEY": "chiave", "S_API_SUM_TIMEOUT": "attesa",
		"S_API_SUM_STATE": "stato", "S_API_NO_MODEL": "non indicato",
		"S_API_NONE": "non configurato: la post-elaborazione resta locale",
		"S_POSTAPI_SETUP": "Configura", "S_API_EDIT": "Modifica", "S_API_KEY_DEL": "Elimina la chiave", "S_API_DLG": "Server esterno",
		"S_LLM_CATALOG": "Catalogo dei modelli", "S_LLM_BLOCK": "Modelli installati", "S_LLM_NONE_HINT": "Nessun modello installato — scaricane uno trovato con la freccia e comparirà qui", "S_LLM_IN_MEM": "in memoria", "S_LLM_ON_DISK": "su disco", "S_LLM_EJECT": "Scarica dalla memoria", "S_LLM_FOUND": "trovati %d", "S_LLM_NOSEARCH": "nessuna ricerca", "S_LLM_SEARCH_HINT": "Digita il nome del modello e premi «Cerca»", "S_LLM_PICK_WAIT": "Disponibile quando il modello sarà scaricato", "S_LLM_INSTALLED": "installati",
		"S_LLM_SUM_MODEL": "modello", "S_LLM_SUM_SIZE": "dimensione", "S_LLM_SUM_COUNT": "installati", "S_LLM_SUM_RAM": "memoria",
		"S_DLG_CLOSE": "Chiudi", "S_LLM_NOPICK": "non scelto", "S_NO_PROMPTS": "Ancora nessun prompt", "S_PROF_DRAG": "trascina per riordinare",
		"S_PROF_NAME_PH": "come chiamare il prompt", "S_PROF_TEST_PH": "scrivi una frase per provare",
		"S_PF_NEW": "Nuovo prompt", "S_PF_EDIT": "Modifica del prompt",
		"S_POST_NO_MODEL": "attiva, ma nessun modello scelto", "S_POST_NO_API": "attiva, ma il server non è configurato", "S_POST_BAD": "il server non ha risposto: %s", "S_POST_NO_PROMPT": "attiva, ma nessun prompt selezionato", "S_API_TEST": "Prova connessione", "S_API_TEST_RUN": "Controllo…", "S_API_TEST_OK": "Il server ha risposto", "S_API_CLEAR": "Cancella", "S_API_CLEAR_ASK": "Cancellare indirizzo, modello e chiave del server esterno? La post-elaborazione torna al modello locale.", "S_RAM_AVAIL": "Memoria disponibile: %s GB di %s GB", "S_RAM_OF": "%s GB di %s GB",
		"S_REPL_ADD": "Aggiungi una sostituzione", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "parole intere", "S_REPL_CASE": "maiuscole", "S_REPL_EMPTY": "Ancora nessuna sostituzione",
		"S_PASTE_DELAY": "Ritardo prima di inserire", "S_PASTE_DELAY_SUB": "quando il programma non è ancora pronto",
		"S_OVPOS": "Dove mostrare la barra", "S_OVPOS_SUB": "al cursore — accanto al punto in cui scrivi; se l'app non lo espone, accanto al puntatore",
		"S_OVPOS_CARET": "Al cursore",
		"S_OVTEXT": "Mostrare il testo riconosciuto", "S_OVTEXT_SUB": "sulla barra dopo l'inserimento, invece del numero di caratteri",
		"S_OVERLAY": "Mostra la barra", "S_OVERLAY_SUB": "durante la dettatura lo schermo mostra che è in corso la registrazione", "S_TYPEMODE": "Digitazione carattere per carattere",
		"S_RECLANG": "Lingua del parlato", "S_RECAUTO": "Auto",
		"S_LANG_KEY": "Scorciatoia propria", "S_LANG_KEY_SET": "Assegna…", "S_LANG_KEY_CLEAR": "Rimuovi la scorciatoia", "S_LANG_KEY_HINT": "Tienila premuta per dettare in questa lingua con il suo modello. La scorciatoia principale continua a usare la lingua del parlato scelta sopra.",
		"S_DL": "Scarica", "S_DEL": "Elimina",
		"S_M_BASE": "veloce, PC modesti", "S_M_SMALL": "bilanciato", "S_M_MED": "Whisper per qualsiasi lingua; quello che traduce", "S_M_TURBO": "massima precisione su CPU", "S_M_GIGAAM": "più preciso in russo, punteggia da sé", "S_M_PARAKEET": "25 lingue europee, punteggia da sé",
		"S_THREADS":  "Thread CPU", "S_MINMS": "Registrazione min, ms", "S_MAXSEC": "Registrazione max, s",
		"S_AUTOSTART": "Avvia whisper-server automaticamente", "S_PORT": "Porta", "S_SERVEREXE": "Percorso whisper-server", "S_SERVEREXE_SUB": "si compila da sé; cambialo solo se il server sta altrove", "S_EXE_RESET": "Ripristina", "S_EXE_WARN": "L'app trova whisper-server accanto a sé. Con un percorso scritto a mano, spostando la cartella il riconoscimento non parte più. Cambiarlo?", "S_RESET_ALL": "Ripristina le impostazioni", "S_RESET_ALL_SUB": "tutto torna di fabbrica, tranne modelli e cronologia", "S_RESET_ALL_BTN": "Ripristina", "S_RESET_ALL_ASK": "Riportare tutte le impostazioni di fabbrica? Modelli, cronologia e prompt restano.", "S_RELOAD_CFG": "Rileggi config.json", "S_RELOAD_CFG_SUB": "se hai modificato il file a mano", "S_RELOAD_CFG_BTN": "Rileggi", "S_UPD_FOUND": "È uscita la versione %s", "S_THEME": "Colore", "S_THEME_SUB": "il colore della finestra, della barra e dell'icona", "S_THEME_GREEN": "Verde", "S_THEME_AMBER": "Ambra", "S_THEME_BLUE": "Blu", "S_THEME_PINK": "Rosa", "S_THEME_EDITOR": "Editor", "S_THEME_NEON": "Neon", "S_WND_MAX": "Riempi lo schermo", "S_WND_RESTORE": "Torna alla dimensione precedente", "S_WND_MIN": "Riduci nell'area di notifica", "S_WND_CLOSE": "Chiudi la finestra", "S_SKIN": "Design", "S_SKIN_SUB": "carattere, forme, effetti e movimento", "S_SKIN_TERMINAL": "Terminale", "S_SKIN_SOFT": "Morbido", "S_SKIN_PAPER": "Documento", "S_SKIN_FLUENT": "Fluent", "S_SKIN_STUDIO": "Studio",
		"S_SERVERURL": "Server esterno (URL)", "S_URLHINT": "Se impostato, il server locale non parte",
		"S_STT_SRV": "Server di riconoscimento",
		"S_STT_SRV_HINT": "I modelli Whisper girano in un programma separato. Può funzionare su questo computer o su un altro — scegli quale usare.",
		"S_SRV_LOCAL": "Su questo computer",
		"S_SRV_REMOTE": "Su un altro computer",
		"S_SRV_REMOTE_HINT": "Lo stesso whisper-server, avviato altrove: un server di casa, una macchina con scheda grafica, il computer accanto.",
		"S_SRV_K_AUTO": "avvio automatico",
		"S_SRV_K_FILE": "file",
		"S_SRV_K_ADDR": "indirizzo",
		"S_SRV_K_CHECK": "ultimo controllo",
		"S_SRV_NEAR": "whisper-server.exe accanto all’app",
		"S_SRV_NOADDR": "non impostato",
		"S_SRV_NOCHECK": "mai controllato",
		"S_SRV_LOCAL_DLG": "Server di riconoscimento locale",
		"S_SRV_ADDR": "Indirizzo del server",
		"S_SRV_ADDR_SUB": "indirizzo di whisper-server sull’altra macchina, porta inclusa",
		"S_SRV_ON": "attivo",
		"S_SRV_OFF": "disattivato",
		"S_SRV_K_THREADS": "thread CPU",
		"S_SRV_K_PORT": "porta",
		"S_SRV_DOWN": "Riconoscimento non disponibile",
		"S_SRV_DOWN_WHY": "il server di riconoscimento remoto non è configurato — imposta il suo indirizzo nelle impostazioni",
		"S_SRV_DOWN_GO": "Apri le impostazioni del server",
		"S_SRV_WARN_NOW": "La dettatura non funziona adesso: è scelto il server remoto, ma manca il suo indirizzo.",
		"S_SRV_WARN_LATER": "Appena verrà scelto un modello Whisper, il riconoscimento non funzionerà: manca l’indirizzo del server remoto.",
		"S_SAVED":      "Salvato",
		"S_ABOUT_HTML": "<p><b>Voce → testo alla posizione del cursore.</b></p><p>Posiziona il cursore, tieni la scorciatoia, parla, rilascia — il testo viene inserito.</p><p>Completamente locale e offline. Tecnologie: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; modelli da Hugging Face.</p><p>I log non superano mai ~2 MB.</p>",
		  "S_SEARCH": "Trova un'impostazione…",
		"S_GRP_GENERAL": "Generale", "S_GRP_SPEECH": "Elaborazione vocale", "S_GRP_INFO": "Info", "S_NAV_POST": "Post-elaborazione", "S_NAV_HELP": "Guida", "S_NAV_CONTACTS": "Contatti", "S_HIST_ADD": "Aggiungi", "S_CONTACT_MAIL": "E-mail", "S_DICT_MODEL": "Modello di riconoscimento", "S_LIB_ACC": "precisione", "S_LIB_SPD": "velocità",
		"S_HELP_TOC": "In questa pagina",
		"S_HELP_TOC_SHOW": "Mostra l’indice — la finestra si allarga",
		"S_HELP_TOC_HIDE": "Nascondi l’indice e ripristina la larghezza",
		"S_CONTACT_TITLE": "Contattare",
		"S_ABOUT_DEPS": "Moduli esterni",
		"S_ABOUT_DEPS_HINT": "Codice di terzi incluso nell’app e le sue licenze. Un clic sul nome apre la pagina del progetto.",
		"S_DEP_WHISPER": "esegue i modelli Whisper",
		"S_DEP_LLAMA": "post-elaborazione, modelli GGUF, traduzione dei sottotitoli sulla scheda grafica (llama-server-gpu.exe)",
		"S_DEP_VULKAN": "intestazioni e libreria di importazione del motore per la scheda grafica; il runtime di Vulkan arriva con il driver grafico",
		"S_DEP_SHERPA": "motore di GigaAM, Parakeet, Canary, Qwen, Moonshine",
		"S_DEP_GGML": "libreria di tensori dentro whisper.cpp e llama.cpp",
		"S_DEP_ONNX": "esegue i modelli dentro sherpa-onnx",
		"S_DEP_WEBVIEW": "finestra delle impostazioni su WebView2",
		"S_DEP_WV2RT": "componente di Windows che disegna quella finestra",
		"S_DEP_MALGO": "cattura dal microfono",
		"S_DEP_MINIAUDIO": "livello audio dentro malgo",
		"S_DEP_WS": "collegamento con sherpa-server",
		"S_DEP_XSYS": "chiamate WinAPI da Go",
		"S_DEP_WINLOADER": "caricamento delle DLL dentro go-webview2",
		"S_DEP_PLEX": "carattere dell’interfaccia",
		"S_DEP_HF": "catalogo da cui si scaricano i modelli",
		"S_CONTACT_HINT": "Un bug, un’idea, una domanda su un’impostazione: scrivi un’email se è personale, o apri una issue se è un bug.",
		"S_CONTACT_REPO": "Repository",
		"S_CONTACT_ISSUES": "Bug e idee",
		"S_CONTACT_WRITE": "Scrivere un’email",
		"S_CONTACT_OPEN": "Apri",
		"S_STATE_ACTIVE": "Riconosce",
		"S_STATE_USED": "Modelli in uso",
		"S_STATE_INST": "Installati localmente",
		"S_STATE_INST_SUB": "modelli su disco, pronti da assegnare",
		"S_PRESETS": "Quale modello per quale lingua",
		"S_PRESETS_HINT": "Il modello della riga «Auto» serve tutte le lingue che non ne hanno uno proprio — nell’elenco sono in grigio. Fate clic su una lingua per assegnarle il suo. Lo stesso modello riconosce anche quale lingua si parla quando la lingua di partenza non è impostata.",
		"S_MFOLDER": "Modello proprio",
		"S_DICT_SAVE": "Salva",
		"S_OWNM_SUB": "Aggiungete un modello locale di riconoscimento vocale",
		"S_OWNM_ONEFILE": "Un file",
		"S_OWNM_FOLDERF": "Cartella con i file del modello",
		"S_OWNM_S1": "Aprite la cartella dei modelli",
		"S_OWNM_S1S": "Cartella di destinazione:",
		"S_OWNM_S2": "Copiate il modello",
		"S_OWNM_S2S": "Scegliete una delle strutture supportate",
		"S_OWNM_S3": "Riavviate l’applicazione",
		"S_OWNM_S3S": "Il modello apparirà per le lingue che supporta",
		"S_AS_AUTO": "come il rilevamento automatico",
		"S_REC_CHIP": "consigliato",
		"S_BACK_AUTO": "Tornare al rilevamento automatico",
		"S_LANGS_COUNT": "lingue: %d",
		"S_LANGS_UNKNOWN": "lingue: sconosciute",
		"S_TR_EN": "traduce in %d lingue",
		"S_TR_LIST": "traduce: %s",
		"S_DL_GOING": "scaricamento:",
		"S_OPEN_FOLDER": "Apri la cartella",
		"S_UNLOAD": "Scarica dalla memoria",
		"S_UNLOAD_SUB": "libera la memoria; la prossima dettatura ricarica il modello",
		"S_UNLOAD_GO": "Libera",
		"S_UNLOADED": "Liberato",
		"S_NOT_FOR_LANG": "%s non riconosce questa lingua",
		"S_MANUAL_NOTE": "Non scaricabile dall’app — la licenza vieta la ridistribuzione. Scaricate l’archivio da soli e scompattatelo in models/moonshine-uk.",
		"S_MANUAL_LINK": "Scarica tu stesso",
		"S_HF_FIT": "solo quelli adatti a questo computer",
		"S_HF_HIDDEN": "nascosti: %s",
		"S_WIZ_SKIP_DL": "Scarica più tardi",
		"S_WIZ_SKIP_NOTE": "Senza modello la dettatura non funzionerà. Potete scaricarlo in «Lingue e modelli».",
		"S_M_GIGAAM2": "la generazione precedente del modello russo: stessa velocità, ma senza punteggiatura",
		"S_M_MOONUK": "modello ucraino Moonshine: veloce e leggero, senza punteggiatura",
		"S_M_LOCAL": "trovato nella cartella models; proprietà sconosciute, perciò niente barre",
		"S_ALL_LANGS": "tutte le lingue",
		"S_OVPOS_SCHEME_SUB": "fai clic sullo schermo: la barra va lì",
		"S_OVDRAG": "Trascinala dove vuoi",
		"S_OVMON": "Schermo",
		"S_OVMON_SUB": "su quale monitor mostrare la barra",
		"S_OVMON_CURSOR": "Lo schermo con il cursore",
		"S_M_NEMOTRON": "scrive mentre parli: il testo appare sulla barra in diretta; 40 lingue, punteggia da sé",
		"S_M_TINY": "il più piccolo e veloce, per macchine molto modeste; sensibilmente meno preciso",
		"S_STATE_LOADED": "In memoria in questo momento",
		"S_STATE_LOADED_SUB": "i modelli si scaricano da soli dopo l’inattività",
		"S_STATE_WEEK": "Questa settimana",
		"S_ST_SUMMARY": "Riepilogo", "S_ST_OVERLAY": "Barra a schermo", "S_ST_BEEP": "Segnale acustico", "S_ST_AUTORUN": "Avvio con Windows", "S_ST_POST": "Post-elaborazione", "S_ST_LOCAL": "in locale", "S_ST_CHECKED": "verificato %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "attivo", "S_ST_OFF_M": "spento", "S_ST_MIC_OK": "il segnale è buono", "S_ST_MIC_BAD": "il microfono non risponde", "S_ST_CHECK": "Controlla", "S_ST_RECOG": "riconosciuto da %s", "S_ST_VER": "Versione %s", "S_ST_LATEST": "l'ultima", "S_ST_OUTDATED": "non aggiornata", "S_ST_UPD_OK": "hai l'ultima versione", "S_ST_UPD_DL": "Scarico l'aggiornamento…",
		"S_ST_QUICK": "Impostazioni rapide",
		"S_ST_MODELS": "Modelli",
		"S_ST_USAGE": "Questa settimana",
		"S_ST_READY": "Pronto a dettare",
		"S_ST_OFF": "Spento nell'area di notifica",
		"S_ST_OFF_SUB": "la scorciatoia non fa nulla finché non la riattivi",
		"S_ST_ENABLE": "Accendi",
		"S_ST_GOTO": "Apri questa impostazione nella sua scheda",
		"S_ST_HOTKEY_GO": "Cambia la scorciatoia",
		"S_ST_UPD_LAST": "Versione %s — l'ultima",
		"S_ST_UPD_HAVE": "Versione %s disponibile",
		"S_ST_MEM": "%s GB liberi su %s",
		"S_ST_MEM_SUB": "in memoria: %s · su disco: %d modelli, %s GB",
		"S_ST_MEM_NONE": "niente",
		"S_ST_LANG": "Lingua parlata",
		"S_ST_ASR": "Riconoscimento",
		"S_ST_ON": "attivo",
		"S_ST_OFF_W": "spento",
		"S_ST_ON_F": "attiva",
		"S_ST_OFF_F": "spenta",
		"S_ST_ACTIVE": "attiva",
		"S_ST_IDLE": "non avviata",
		"S_ST_DISK": "su disco, %s",
		"S_ST_USAGE_SUB": "%d caratteri · oggi %d · %d caratteri in media",
		"S_WEEK_OTHER": "altri",
		"S_ST_NO_WEEK": "nessuna dettatura questa settimana",
		"S_ST_AUTORUN_SUB": "l'app non si avvia da sola",
		"S_ST_OVERLAY_SUB": "visibile durante la registrazione",
		"S_REPL_LANG": "Lingua della regola",
		"S_REPL_LANG_ALL": "tutte le lingue",
		"S_M_CANARY": "inglese, tedesco, spagnolo, francese — e traduce tra loro da sé",
		"S_M_QWEN3": "circa 30 lingue, punteggia da sé; il più pesante e preciso del catalogo",
		"S_POSTAPI": "Server esterno",
		"S_POST_HINT": "Corregge il testo riconosciuto secondo i vostri prompt: toglie gli intercalari, sistema la punteggiatura, cambia lo stile. Spento — il testo è inserito così com’è.",
		"S_POST_MODEL": "Modello",
		"S_SRC_LOCAL": "Locale",
		"S_SRC_USED": "in uso",
		"S_HF_GO": "Cerca",
		"S_POSTAPI_HINT": "Vuoto di default — tutta la post-elaborazione è locale. Inserite un indirizzo e i prompt girano su un server esterno: OpenAI, Groq, il vostro vLLM — qualsiasi cosa con API compatibile.",
		"S_POSTAPI_URL": "Indirizzo",
		"S_POSTAPI_URL_SUB": "vuoto = il modello locale; esempio: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Modello",
		"S_POSTAPI_KEY": "Chiave API",
		"S_POSTAPI_KEY_SET": "chiave salvata (cifrata con Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "nessuna chiave",
		"S_POSTAPI_SAVE": "Salva chiave",
		"S_POSTAPI_TIMEOUT": "Attesa della risposta", "S_SEC_SHORT": "s",
		"S_POSTAPI_WARN": "⚠ Il testo riconosciuto delle dettature andrà a questo indirizzo. L’audio non esce mai. La chiave è salvata cifrata.",
		"S_POSTAPI_ASK": "Inviare il testo riconosciuto a %s? L’audio resta su questo computer, ma il testo lo lascerà.",
		"S_POSTAPI_BADGE": "server esterno",
		"S_NOT_INSTALLED": "non installato",
		"S_NAV_STATE": "Stato", "S_NAV_DICT": "Controlli e comportamento", "S_NAV_MIC": "Microfono", "S_NAV_MODELS": "Lingue e modelli",
		"S_NAV_TEXT": "Regole", "S_NAV_TR": "Traduzione", "S_NAV_SYSTEM": "Sistema", "S_NAV_ABOUT": "Informazioni",
		"S_STATE_HINT": "tieni premuto e parla — il testo arriva dov'è il cursore",
		"S_STATE_PROC": "Post-elaborazione",
		"S_CHANGE_MODEL": "Cambia", "S_PICK_MODEL": "Scegli", "S_STATE_GET": "Scarica",
		"S_RETRY": "Riprova", "S_BERR_OPEN": "Apri le impostazioni del server",
		"S_STATE_LAST": "Ultima dettatura", "S_STATE_COPY": "Copia", "S_STATE_MEM": "Memoria",
		"S_STATE_MEM_SUB": "i modelli restano caricati, la prima frase non aspetta",
		"S_HOTMODE":       "Modo", "S_HOTMODE_HOLD": "tenere premuto", "S_HOTMODE_TOGGLE": "interruttore",
		"S_SUB_HOTMODE": "tieni premuti i tasti, oppure premi una volta per iniziare e una per fermare",
		"S_SILENCE":     "Ferma al silenzio", "S_SILENCE_MS": "Pausa prima di fermare, ms",
		"S_SUB_SILENCE": "modo «interruttore»: quando hai già parlato, una pausa chiude la registrazione come se premessi di nuovo i tasti",
		"S_SUB_SILENCE_MS": "quanto puoi restare in silenzio; non si ferma mai prima delle tue prime parole",
		"S_SUB_MINMS":   "ignora le pressioni accidentali",
		"S_SUB_ENTER":   "invia subito il messaggio",
		"S_SUB_CLIP":    "immagini e file tornano come erano",
		"S_SUB_TYPE":    "serve dove un campo rifiuta di incollare",
		"S_SEC_OVERLAY": "Avviso a schermo",
		"S_MIC_CHECK":   "Controlla il microfono", "S_MIC_CHECK_SUB": "tre secondi di registrazione e un verdetto: livello, distorsione, se c'è voce", "S_MIC_CHECKING": "Controllo…",
		"S_MCHECK": "Controlla i modelli installati", "S_MCHECK_SUB": "confronta i file dei modelli con gli hash di riferimento", "S_MCHECK_GO": "Controlla", "S_MCHECK_RUN": "Controllo…",
		"S_HIST_INSERT": "Incolla",
		"S_MIC": "Microfono", "S_MIC_DEFAULT": "Predefinito di sistema", "S_MIC_REFRESH": "Aggiorna l'elenco",
		"S_MIC_LEVEL": "Livello d'ingresso", "S_MIC_QUIET": "silenzio",
		"S_SUB_THREADS": "più thread non è sempre più veloce — misura sulla tua macchina",
		"S_SEC_LLM":     "Modello editor",
		"S_PUNCT":       "Punteggiatura e maiuscole", "S_SUB_PUNCT": "da dove arrivano punteggiatura e maiuscole",
		"S_PUNCT_MODEL": "dal modello", "S_PUNCT_LLM": "dal modello editor", "S_PUNCT_OFF": "togliere",
		"S_SUB_DICT": "Dizionario", "S_SUB_PROMPTS": "Prompt",
		"S_SUB_TRTARGET": "il testo viene tradotto in essa; il dialogo sulla barra la propone per prima",
		"S_REMOTE_ABOUT": "È impostato un server remoto: l'audio viene inviato lì e la promessa qui sopra non vale finché è attivo.",
		"S_UPD":          "Aggiornamenti", "S_UPD_CHECK": "Cerca aggiornamenti", "S_UPD_AUTO": "Controlla all'avvio",
		"S_SUB_UPD":  "l'unica richiesta di rete oltre allo scaricamento dei modelli",
		"S_UPD_NONE": "Hai l'ultima versione", "S_BADGE_MODELS": "Modelli installati", "S_BADGE_MISS": "Un modello non è scaricato", "S_BADGE_SYSTEM": "Avvisi da controllare", "S_BADGE_HIST": "Voci nella cronologia", "S_LOG_OPEN": "Apri il registro", "S_LOG": "Registro", "S_LOG_SUB": "tutto ciò che l'app scrive su di sé", "S_UPD_AVAIL": "È disponibile la versione %s.",
		"S_UPD_GO": "Aggiorna", "S_UPD_ERR": "Controllo non riuscito", "S_UPD_DL": "Scaricamento dell'aggiornamento…",
		"S_SEC_SERVICE": "Servizio", "S_SUB_AUTOSTART": "disattiva se avvii il server da solo",
		"S_SUB_PORT":    "il riconoscitore si riavvia da sé",
		"S_MODEL_READY": "Modello scaricato — scegli per passare a lui",
		"S_FIT_OK":      "ci sta", "S_FIT_WARN": "al limite", "S_FIT_BAD": "memoria insufficiente", "S_RAM": "Memoria del computer:",
		"S_HF_PH":       "Nome del modello — per es. qwen2.5 instruct",
		"S_NO_LLM":      "Ancora nessun modello installato — cercane uno nel campo di ricerca qui sotto.",
		"S_NO_LLM_PROF": "I prompt diventano disponibili appena è installato un modello — il blocco «Modello» in alto in questa scheda.",
		"S_UPDATED":     "Ultimo aggiornamento del modello", "S_PROF_EDIT": "Modifica", "S_PROF_CLOSE": "Riduci",
		"S_CONFIRM_DEL": "Eliminare il modello «%s»? Si potrà scaricare di nuovo.", "S_FREE": "liberi",
		"S_DEL_ACTIVE":     "Eliminare il modello attivo «%s»? Il riconoscimento si ferma finché non ne scegli un altro, che puoi scaricare qui stesso.",
		"S_REMOTE_WARN":    "L'audio verrà inviato a questo server. La modalità locale è spenta.",
		"S_REMOTE_ASK":     "L'audio non sarà più elaborato su questo computer e verrà inviato a %s. Attivare la modalità remota?",
		"S_REMOTE_BADGE":   "REMOTO",
		"S_OK":             "Sì", "S_CANCEL": "Annulla", "S_DL_START": "Scarica", "S_DL_CANCEL": "Annulla lo scaricamento",
		"S_DL_ASK":    "Il modello «%s» non è scaricato (%s). Iniziare lo scaricamento?",
		"S_NOT_FOUND": "niente",   "S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autore e sviluppatore di {app} — uno strumento di dettatura locale per Windows: la voce diventa testo proprio dov'è il cursore, senza cloud e senza abbonamenti.</p>" +
			"<p>Il progetto è aperto: codice sorgente, build e ultime release stanno su GitHub.</p>" +
			"<ul>" +
			"<li>Repository: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Profilo dell'autore: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Hai trovato un bug o hai un'idea — apri una issue nel repository.</p>",
		"S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Che cos'è</p>" +
			"<p>Parli e le parole compaiono dove lampeggia il cursore. Come se le avessi scritte tu, solo più in fretta.</p>" +
			"<p class=\"hdim\">Funziona così. Fai clic in un campo dove di solito scrivi: un'e-mail, una chat, un campo di ricerca, un documento. Tieni premuti due tasti e dici una frase con la tua voce di sempre. Li rilasci. Un secondo dopo la frase è nel campo. Non serve passare a un'altra finestra né copiare nulla: l'app funziona sopra qualsiasi programma e aspetta i suoi tasti in background.</p>" +
			"<p class=\"hdim\">Tutto il lavoro lo fa il tuo computer. La voce non viene salvata: resta qualche secondo in memoria, diventa testo e sparisce. L'app si collega a internet solo in due occasioni: quando scarichi un modello di riconoscimento e quando controlla se è uscita una nuova versione. Il controllo si può disattivare, e allora la rete non le serve affatto.</p>" +
			"<p class=\"hdim\">Più avanti incontrerai spesso la parola «modello». È semplicemente un file che sa capire il parlato. Si scarica una volta e resta in una cartella accanto all'app; più è grande, meglio sente e più lentamente lavora.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Che cosa succede quando rilasci i tasti</p>" +
			"<div class=\"hflow\">" +
			"<span>registrazione</span><em>→</em><span>riconoscimento</span><em>→</em><span>sostituzioni</span><em>→</em><span>comandi</span><em>→</em><span>traduzione</span><em>→</em><span>modifica</span><em>→</em><span>inserimento</span>" +
			"</div>" +
			"<p>Di solito ci vogliono uno o due secondi, e non c'è bisogno di conoscere i passaggi. Ma se un giorno il testo esce diverso da come te lo aspettavi, questa catena ti aiuta a capire chi l'ha cambiato.</p>" +
			"<ul>" +
			"<li><b>Registrazione.</b> Finché tieni premuti i tasti, l'app ascolta il microfono. Ignora le pressioni troppo brevi, così un tasto sfiorato per sbaglio non avvia una dettatura.</li>" +
			"<li><b>Riconoscimento.</b> La voce diventa lettere. Se ne occupa un modello. Per il russo di solito lavora GigaAM, per il polacco Parakeet, per la maggior parte delle altre lingue Nemotron, e Whisper si occupa delle lingue restanti e della traduzione; sono solo nomi, e non è necessario sceglierli a mano.</li>" +
			"<li><b>Sostituzioni.</b> Il tuo elenco personale «sentito → voluto». Se il modello si ostina a scrivere «git hub» invece di «GitHub», lo correggi una volta e da lì in poi si sistema da solo.</li>" +
			"<li><b>Comandi.</b> Alcune frasi funzionano come pulsanti. Di' «nuova riga» e il testo va a capo, senza che quelle parole finiscano nel campo.</li>" +
			"<li><b>Traduzione.</b> Se l'hai attivata, parli in una lingua e il testo arriva in un'altra.</li>" +
			"<li><b>Modifica.</b> Un passaggio facoltativo. Un secondo modello dà una ripulita al testo: toglie gli «ehm» e i «cioè», mette le virgole, dà alla frase un tono formale.</li>" +
			"<li><b>Inserimento.</b> Il testo arriva nel campo. Quello che avevi negli appunti prima della dettatura torna al suo posto.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Mentre succede tutto questo, in fondo allo schermo si vede una barra sottile con il nome del passaggio in corso. Così capisci subito che l'app è al lavoro e non bloccata.</p>" +
			"<p class=\"wh\" id=\"s-first\">Il primo avvio</p>" +
			"<p>Al primo avvio l'app ti accompagna in tre passaggi. Ci vuole un minuto, e poi puoi dettare.</p>" +
			"<ul>" +
			"<li><b>In che lingua detterai.</b> La domanda più importante: in base alla risposta l'app sceglie da sola un modello adatto e lo scarica in background, con l'avanzamento visibile. Se questa lingua è diversa da quella della finestra, lì stesso compare una riga per cambiare anche la lingua dell'interfaccia.</li>" +
			"<li><b>Prova.</b> La scorciatoia da tastiera (per impostazione predefinita Ctrl e il tasto Windows; per cambiarla, premi quello che ti viene comodo tenere con una mano), il microfono con un indicatore di livello dal vivo e un campo dove dettare una frase di prova. Se l'indicatore non si muove mentre parli, è scelto il microfono sbagliato. Il campo si anima non appena il modello è pronto.</li>" +
			"<li><b>Fatto.</b> L'interruttore «Avvia con Windows», e la procedura guidata è conclusa.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Puoi saltare qualsiasi passaggio e tornarci più tardi: tutto si trova anche nelle impostazioni. Il modello continua a scaricarsi in background anche se passi oltre o chiudi la procedura guidata prima che abbia finito. Se vuoi rifare la procedura guidata, avvia l'app aggiungendo <code>-wizard</code>. Dopo un aggiornamento non ricompare.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">La barra in fondo allo schermo</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Parla…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">così l'app ti fa vedere che ti sente</div>" +
			"</div>" +
			"<p>Mentre detti, è l'unica cosa che si vede sullo schermo. La barra è piccola, non cattura il mouse e non intralcia il lavoro: ti racconta semplicemente che cosa sta succedendo in questo momento alla tua frase.</p>" +
			"<ul>" +
			"<li><b>«Parla…»</b> È in corso la registrazione. Le colonnine accanto saltano al ritmo della voce. Se restano ferme, non ti si sente, e non ha senso continuare a parlare.</li>" +
			"<li><b>«Riconoscimento…»</b> Hai rilasciato i tasti e il modello è al lavoro. Di solito uno o due secondi. La primissima dettatura dopo l'avvio dura di più: in quel momento il modello si sta caricando in memoria.</li>" +
			"<li><b>«Traduzione…»</b> È in corso la traduzione.</li>" +
			"<li><b>«Modifica: Cleanup (1/2)»</b> Il testo lo sta correggendo il secondo modello. Si vede quale istruzione è in corso e quante sono in tutto.</li>" +
			"<li><b>«Inserito: 128 caratteri»</b> Fatto, il testo è già nel campo. La barra si spegne da sola dopo un paio di secondi.</li>" +
			"<li><b>La crocetta.</b> Hai cambiato idea? Cliccala in qualsiasi fase. La registrazione viene scartata e nel campo non arriva nulla.</li>" +
			"</ul>" +
			"<p class=\"hdim\">La barra si può spegnere, mettere in qualsiasi angolo, legare al cursore o semplicemente trascinare con il mouse dove non dà fastidio. Lo trovi nella scheda «Controlli e comportamento».</p>" +
			"<p class=\"wh\" id=\"s-state\">La schermata «Stato»</p>" +
			"<p>Questa schermata risponde a una sola domanda: puoi dettare proprio adesso? Se no, ti dice perché e ti dà un pulsante che porta dove lo si sistema.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Pronto a dettare</b><span>tieni premuto e parla — il testo arriva dov'è il cursore</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">il quadratino verde vuol dire che l'app è in ascolto dei suoi tasti</div>" +
			"</div>" +
			"<p>Il riquadro in alto può essere di tre tipi. <b>Verde</b>: va tutto bene, e accanto c'è un promemoria dei tasti. <b>Giallo</b>: hai spento l'app dalla sua icona nell'area di notifica, vicino all'orologio, e i tasti non fanno nulla finché non premi «Accendi» lì stesso, nel riquadro. <b>Rosso</b>: le impostazioni si contraddicono. È scelto il riconoscimento su un altro computer, ma l'indirizzo non è indicato.</p>" +
			"<p>Più sotto c'è quello che di solito si guarda di sfuggita:</p>" +
			"<ul>" +
			"<li><b>Sottotitoli.</b> Il riquadro subito sotto quello in alto: se i sottotitoli sono spenti, stanno caricando il modello, sono attivi o si stanno riavviando, con il modello e come arriva il testo, la loro scorciatoia se è impostata, un pulsante che li avvia o li ferma e una freccia verso la scheda «Sottotitoli».</li>" +
			"<li><b>Questa settimana.</b> Un grafico a torta: in quali programmi sono finite le tue dettature e quanti caratteri hai dettato.</li>" +
			"<li><b>Versione.</b> Il numero della build e il pulsante «Controlla». Se ne è uscita una più nuova, l'app te lo dice proprio qui.</li>" +
			"<li><b>Memoria.</b> Quanta ne resta libera, quale modello c'è dentro adesso e quanto pesano tutti i modelli su disco. Il pulsante «Libera» svuota la memoria subito, per esempio prima di un gioco pesante.</li>" +
			"<li><b>Impostazioni rapide.</b> Cinque riquadri con le cose più usate: microfono, lingua, barra, suono e avvio con Windows. Mostrano com'è impostato adesso, e la freccia porta alla scheda dove si cambia. Qui di proposito non ci sono interruttori, così non tocchi niente per sbaglio.</li>" +
			"<li><b>Modelli.</b> Che cosa è installato, che cosa sta lavorando in questo momento, di quali lingue si occupa ciascuno e che cosa sta semplicemente su disco.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">La scheda «Sistema»</p>" +
			"<p>L'aspetto dell'app, i suoi rapporti con Windows e qualche pulsante di servizio. Di solito ci si passa una volta all'anno.</p>" +
			"<p class=\"hsub\">Come appare l'app</p>" +
			"<p>Tre righe: lingua dell'interfaccia, design e colore. Sul funzionamento non influiscono per niente.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Lingua dell'interfaccia</td><td>La lingua della finestra delle impostazioni e delle scritte sulla barra. Con la lingua di dettatura non ha niente a che fare: la finestra può essere in italiano mentre detti in inglese.</td></tr>" +
			"<tr><td>Design</td><td>L'aspetto generale. Sono sette: «Editor» (predefinito: grigio scuro con blu), «Fluent» (come le Impostazioni di Windows 11), «Studio» (grafite con arancione, quasi senza angoli arrotondati), «Morbido» (lilla chiaro, arrotondato), «Documento» (carta con blu), «Neon» (viola con una sfumatura dal rosa al turchese) e «Terminale» (lettere verdi su nero, parentesi quadre attorno ai pulsanti). Tutti e sette usano lo stesso carattere, IBM Plex Mono.</td></tr>" +
			"<tr><td>Colore</td><td>Una tinta solo per «Terminale»: verde, ambra, blu, rosa. Dello stesso colore diventano anche la barra e l'icona nell'area di notifica. Gli altri design hanno colori propri, e la riga si nasconde.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Avvio e aggiornamenti</p>" +
			"<p>Qui decidi se l'app parte insieme al sistema e se tiene d'occhio le nuove versioni.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Avvia con Windows</td><td>L'app si avvia quando accedi al sistema e aspetta i tuoi tasti. È una normale voce nell'avvio automatico, uguale a quella che potresti aggiungere tu stesso.</td></tr>" +
			"<tr><td>Controlla all'avvio</td><td>Una volta, all'avvio, l'app chiede se è uscita una versione più nuova. A parte lo scaricamento dei modelli, non va mai su internet. Se lo disattivi, le nuove versioni dovrai scoprirle da solo.</td></tr>" +
			"<tr><td>Aggiornamenti</td><td>Il pulsante controlla subito. Se c'è una versione più nuova, l'app la scarica e la installa sopra quella attuale: impostazioni, modelli e cronologia restano al loro posto.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Servizio</p>" +
			"<p>Tre pulsanti per quando qualcosa è andato storto.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Registro</td><td>Apre il file in cui l'app annota che cosa ha fatto: quale modello è partito, quanto è durato il riconoscimento, dove si è verificato un errore. Il motivo di solito è scritto a parole, non con codici. Il file non cresce all'infinito: le righe nuove prendono il posto delle vecchie.</td></tr>" +
			"<tr><td>Rileggi config.json</td><td>Tutte le impostazioni stanno in un unico file di testo accanto all'app. Se l'hai modificato a mano o l'hai portato da un altro computer, questo pulsante carica le modifiche senza riavvio.</td></tr>" +
			"<tr><td>Ripristina le impostazioni</td><td>Riporta tutto com'era dopo l'installazione. I modelli scaricati e la cronologia restano al loro posto: si cancellano solo le impostazioni.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Dove gira il riconoscimento</p>" +
			"<p>Il parlato lo elabora un piccolo programma a parte. Di solito sta accanto all'app, si avvia da solo e non devi pensarci. Ma può stare anche su un altro computer: per esempio su quello di casa, che ha una scheda grafica potente, mentre detti da un portatile leggero. Si sceglie con l'interruttore sui riquadri.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>Su questo computer</b><span>thread CPU 4 · avvio automatico attivo · porta 8910 · whisper-server.exe accanto all’app</span></div><button class=\"hbtn\">Modifica</button></div>" +
			"<div class=\"hcap\">il riquadro scelto è evidenziato, e l'altro intanto è attenuato</div>" +
			"</div>" +
			"<p>Il pulsante «Modifica» apre una finestra con quattro impostazioni. Hanno tutte valori predefiniti, e quasi sempre puoi lasciarle come sono.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Thread CPU</td><td>Quanti core del processore dedicare al riconoscimento. Più core vuol dire più velocità, ma solo fino a un certo punto: oltre non si guadagna nulla, e gli altri programmi cominciano a rallentare.</td></tr>" +
			"<tr><td>Avvia whisper-server automaticamente</td><td>L'app avvia il riconoscitore da sola. C'è un solo motivo per disattivarlo: lo avvii tu a mano e non vuoi che l'app intervenga.</td></tr>" +
			"<tr><td>Porta</td><td>Il numero con cui l'app parla con il riconoscitore all'interno del tuo stesso computer. Da fuori non si vede. Va cambiato solo se quel numero l'ha occupato un altro programma.</td></tr>" +
			"<tr><td>Percorso whisper-server</td><td>Dove si trova il riconoscitore. Si compila da sé. Scrivi il tuo solo se l'hai messo apposta da un'altra parte: altrimenti, appena sposti la cartella, la dettatura smette di partire.</td></tr>" +
			"</table>" +
			"<p>Il secondo riquadro ha una sola impostazione: l'indirizzo del computer in rete, per esempio <code>http://192.168.1.50:8080</code>. Tieni presente che in questo caso la voce registrata va su quella macchina e, con un indirizzo normale, viaggia senza cifratura. In una rete di casa va bene, ma via internet è meglio non lavorare così.</p>" +
			"<div class=\"hbad\">Se è scelto «Su un altro computer» ma l'indirizzo non è stato scritto, e la tua lingua è affidata a Whisper, la dettatura non funzionerà. L'app te lo segnala con un riquadro rosso su «Stato», un segno rosso accanto a «Sistema» e una riga in fondo alla finestra.</div>" +
			"<p class=\"wh\" id=\"s-mic\">La scheda «Microfono»</p>" +
			"<p>Più della metà delle dettature fallite non dipende dal modello né dalle impostazioni, ma dal microfono: se ne ascolta uno sbagliato, oppure ti si sente appena. In questa scheda c'è tutto per verificarlo in mezzo minuto.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Microfono</td><td>Quale dispositivo ascoltare. L'opzione «Predefinito di sistema» significa «quello scelto in Windows»: comodo se a volte usi le cuffie e a volte no, perché l'app segue il sistema da sola.</td></tr>" +
			"<tr><td>Aggiorna l'elenco</td><td>Serve se hai collegato le cuffie dopo aver avviato l'app: l'elenco viene riletto.</td></tr>" +
			"<tr><td>Livello d'ingresso</td><td>Un indicatore dal vivo. Di' qualcosa: le colonnine devono muoversi in modo evidente. Se non si muovono, il microfono è disattivato in Windows oppure è scelto quello sbagliato.</td></tr>" +
			"<tr><td>Controlla il microfono</td><td>Tre secondi di registrazione e un verdetto in parole semplici: se il volume basta, se c'è distorsione (quando il suono «gracchia» perché il microfono è troppo vicino) e se si sente davvero una voce.</td></tr>" +
			"<tr><td>Segnali acustici</td><td>Un breve suono all'inizio e alla fine della registrazione. Utile soprattutto se hai spento la barra: così almeno senti che la registrazione è partita.</td></tr>" +
			"<tr><td>Suono del segnale</td><td>Quale suono, esattamente: voce di Windows, campanello, marimba e qualche altro. Il pulsante accanto riproduce quello scelto, così non devi indovinarlo dal nome.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">La scheda «Cronologia»</p>" +
			"<p>La cronologia è l'elenco di quello che hai dettato. Solo il testo: la voce non viene salvata mai, da nessuna parte. Per impostazione predefinita è disattivata, ed è una scelta voluta: nelle dettature finiscono spesso cose personali.</p>" +
			"<p>Conviene attivarla se ti è capitato almeno una volta qualcosa del genere: hai incollato il testo nella finestra sbagliata ed è andato perso; devi ripetere una frase lunga; vuoi ritrovare quello che hai dettato ieri e inviarlo di nuovo.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Per quanto conservare</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">giorni</span></div>" +
			"<div class=\"hrow\"><label>Quante voci conservare</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">la durata si può impostare in minuti, ore o giorni, con qualsiasi numero</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Conservare la cronologia delle dettature</td><td>L'interruttore principale. Se lo spegni, le nuove voci smettono di comparire e quelle vecchie vengono cancellate.</td></tr>" +
			"<tr><td>Per quanto conservare</td><td>Dopo quanto tempo una voce sparisce da sola. Puoi mettere cinque minuti, se la cronologia ti serve solo per quando il testo finisce nel posto sbagliato, oppure trenta giorni.</td></tr>" +
			"<tr><td>Quante voci conservare</td><td>Un secondo limite, per numero. Funziona insieme al primo: vale quello che si raggiunge prima. Con «7 giorni e non più di 50», la cinquantunesima voce spinge fuori la più vecchia, anche se la settimana non è ancora passata.</td></tr>" +
			"<tr><td>Non registrare mai da questi programmi</td><td>Un elenco di eccezioni. Quello che detti in queste finestre (la banca, un gestore di password, la chat di lavoro) non finisce nella cronologia, mentre la dettatura funziona come sempre. Puoi scegliere il programma dall'elenco delle finestre aperte, senza andare a cercare il nome del file.</td></tr>" +
			"</table>" +
			"<p>Con ogni voce puoi fare quattro cose. <b>Copia</b>: il testo va negli appunti. <b>Incolla</b>: l'app ti chiede di indicarle dove. Fai clic sul campo che vuoi, in qualsiasi programma, e il testo arriva lì. <b>Elimina</b>: toglie quella sola voce. <b>Cerca</b>: la ricerca guarda sia nel testo sia nel nome del programma. Il pulsante «Svuota» cancella tutto in una volta, e per questo chiede conferma.</p>" +
			"<p class=\"hdim\">A destra di ogni voce c'è scritto quando sparirà, per esempio «fino al 7 set», e l'ultimo giorno questa scritta diventa gialla.</p>" +
			"<p class=\"wh\" id=\"s-dict\">La scheda «Controlli e comportamento»</p>" +
			"<p>Come parte la dettatura e come il testo arriva nel campo. Ci si viene quando qualcosa non va come vorresti: i tasti sono scomodi, il testo non compare, Invio si preme da solo.</p>" +
			"<p class=\"hsub\">Tasti</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Scorciatoia da tastiera</td><td>Premi «Cambia…» e tieni premuta la combinazione che ti è comoda. L'app intercetta questa combinazione in tutto il sistema, quindi non sceglierne una già occupata, come Ctrl+C. Le combinazioni con il tasto Windows sono quelle occupate più di rado. Una lingua può avere anche una scorciatoia propria, assegnata nella scheda «Lingue e modelli».</td></tr>" +
			"<tr><td>Modo</td><td>Con «tenere premuto» registra finché tieni premuti i tasti: li rilasci e la frase passa all'elaborazione. Con «interruttore» la prima pressione avvia la registrazione e la seconda la ferma. Il secondo modo è comodo per i monologhi lunghi, quando la mano si stanca di tenere i tasti.</td></tr>" +
			"<tr><td>Ferma al silenzio</td><td>Solo in modo «interruttore». Quando hai già detto qualcosa, una pausa lunga quanto quella impostata in «Pausa prima di fermare, ms» chiude la registrazione da sola, come se premessi di nuovo la scorciatoia, e la barra mostra «Fermato dal silenzio». Non si ferma mai prima delle tue prime parole, quindi dopo aver premuto puoi raccogliere le idee con calma; se non dici niente, il limite «Registrazione max, s» la chiude comunque. Con «Premi Invio dopo l'inserimento (invio auto)» attivo, una pausa invia il messaggio, quindi scegli una pausa più lunga.</td></tr>" +
			"<tr><td>Registrazione min, ms</td><td>Sotto quale durata una pressione conta come accidentale. Se sfiori i tasti per una frazione di secondo, l'app non fa nulla invece di inserire un pezzo di parola.</td></tr>" +
			"<tr><td>Registrazione max, s</td><td>Una protezione contro un tasto dimenticato: la registrazione si interrompe da sola invece di registrare il vuoto per mezz'ora.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Come viene inserito il testo</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Premi Invio dopo l'inserimento (invio auto)</td><td>Il testo non viene solo inserito, ma inviato subito: il messaggio parte, la ricerca si avvia. In un'app di messaggistica ti risparmia un gesto; in un'e-mail o in un documento quasi sicuramente darà fastidio.</td></tr>" +
			"<tr><td>Ripristina appunti dopo l'inserimento</td><td>L'inserimento passa dagli appunti, e l'app rimette lì quello che c'era: un'immagine, un file, un pezzo di testo copiato. Lascialo attivo se copi spesso qualcosa e non vuoi perderlo.</td></tr>" +
			"<tr><td>Digitazione carattere per carattere</td><td>Invece di incollare il testo in un colpo solo, l'app lo digita come se scrivessi tu. È più lento, ma ti toglie dai guai dove non si può incollare dagli appunti: alcuni moduli bancari, i terminali, i programmi vecchi.</td></tr>" +
			"<tr><td>Ritardo prima di inserire</td><td>Una breve pausa prima di consegnare il testo alla finestra. Serve di rado: capita che un programma non faccia in tempo a ricevere il testo e l'inizio della frase vada perso.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Avviso a schermo</p>" +
			"<p>Qui decidi se mostrare la barra e dove esattamente.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Mostra la barra</td><td>Spegnila se ti dà fastidio o se ti bastano i segnali acustici.</td></tr>" +
			"<tr><td>Dove mostrare la barra</td><td>Uno schema dello schermo con nove posizioni: fai clic sull'angolo dove la barra non darà fastidio.</td></tr>" +
			"<tr><td>Al cursore</td><td>La barra compare accanto al punto in cui stai dettando. Se il programma non dice dove si trova il cursore, compare accanto al puntatore del mouse.</td></tr>" +
			"<tr><td>Posto tuo</td><td>Puoi semplicemente trascinare la barra con il mouse, e ricorderà quel punto.</td></tr>" +
			"<tr><td>Schermo</td><td>Su quale monitor mostrarla, se ne hai più di uno.</td></tr>" +
			"<tr><td>Mostrare il testo riconosciuto</td><td>Dopo l'inserimento la barra mostra la frase stessa invece del numero di caratteri. Comodo per controllare, ma la frase la vede chiunque guardi il tuo schermo.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">La scheda «Lingue e modelli»</p>" +
			"<p>Qui dici in che lingua detterai, e l'app decide con che cosa riconoscerla. I modelli sono diversi: alcuni veloci, altri precisi, altri ancora conoscono lingue rare.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Lingua del parlato</td><td>La lingua in cui parli. Se ne parli più di una, scegli «Auto»: il modello la riconosce da solo, anche se un po' più lentamente e con un po' meno precisione rispetto a quando la lingua è indicata direttamente.</td></tr>" +
			"<tr><td>Quale modello per quale lingua</td><td>Fai clic su una lingua e si apre l'elenco dei modelli che la conoscono, con una valutazione di velocità e precisione e la dimensione del file. Così una lingua può tenere un modello veloce mentre un'altra ne riceve uno preciso. La riga a parte «Auto» è comune: il suo modello serve tutte le lingue che non ne hanno uno proprio, e nell'elenco queste compaiono in grigio. È anche quella che capisce quale lingua si sta parlando quando la lingua del parlato non è impostata.</td></tr>" +
			"<tr><td>Scorciatoia propria</td><td>Dentro una lingua aperta c’è posto per una combinazione di tasti tutta sua. Tienila premuta e la dettatura avviene in quella lingua con il suo modello, le sue sostituzioni e le sue regole, qualunque sia la lingua del parlato impostata; la scorciatoia principale mantiene la lingua scelta sopra. Così si detta in inglese con una combinazione e nella propria lingua con un’altra, senza aprire le impostazioni. Durante la registrazione la barra mostra il codice della lingua. Una combinazione non può ripetere quella principale né quella di un’altra lingua.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Russo<span class=\"hs\">GigaAM v3 · precisione 5 su 5 · velocità 5 su 5 · 232 MB</span></label><span class=\"hchip\"><i></i>attiva</span></div>" +
			"<div class=\"hrow\"><label>Inglese<span class=\"hs\">come il rilevamento automatico · Nemotron 3.5</span></label><button class=\"hbtn\">Scegli</button></div>" +
			"<div class=\"hcap\">si vede quale modello si occupa di ogni lingua e quanto costa in velocità e spazio su disco</div>" +
			"</div>" +
			"<p class=\"hsub\">Traduzione</p>" +
			"<p>Alcuni modelli non si limitano a trascrivere il parlato, ma lo traducono subito: parli in una lingua e nel campo arriva il testo in un'altra. A tradurre è Whisper: in inglese con la sua funzione di traduzione integrata, nelle altre lingue indicandogli la lingua di uscita come se fosse quella del parlato; tra lingue molto lontane la qualità è irregolare. Se il modello scelto non sa tradurre, accanto all'interruttore compare una riga arancione: se Whisper è installato, la barra proporrà di usarlo per tradurre, altrimenti il testo resterà nella lingua del parlato.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Cambiare la lingua di uscita del testo</td><td>Attiva la traduzione.</td></tr>" +
			"<tr><td>Lingua di uscita del testo predefinita</td><td>In quale lingua tradurre se non ne hai scelta un'altra durante la dettatura stessa.</td></tr>" +
			"<tr><td>Chiedere la lingua di uscita del testo</td><td>Puoi non chiedere e tradurre sempre nella lingua predefinita. Puoi chiedere ogni volta: dopo la dettatura, sulla barra compaiono pulsanti con le lingue. Oppure puoi chiedere con un timeout: gli stessi pulsanti, ma se non ne premi uno in tempo, si usa la lingua predefinita.</td></tr>" +
			"<tr><td>Lingue nel dialogo</td><td>Quali lingue offrire con questi pulsanti. Lasciane due o tre che usi davvero, altrimenti scegliere richiederà più tempo che parlare.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Modelli su disco</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Modello proprio</td><td>Se hai un modello tuo, mettilo nella cartella dei modelli e riavvia l'app: comparirà negli elenchi delle lingue che capisce. Il pulsante accanto apre la cartella giusta, così non devi cercarla in Esplora file.</td></tr>" +
			"<tr><td>Scarica dalla memoria</td><td>Finché l'app è in esecuzione, il modello occupa la RAM. Questo pulsante la libera subito, per esempio prima di un gioco pesante. La dettatura successiva ricarica il modello, e la prima frase richiederà più tempo del solito.</td></tr>" +
			"<tr><td>Controlla i modelli installati</td><td>Confronta i file scaricati con quelli di riferimento. Torna utile se la connessione è caduta durante lo scaricamento e sospetti che un file si sia rovinato.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">La scheda «Regole»</p>" +
			"<p>Tre elenchi che correggono il testo in automatico, ancora prima che arrivi nella finestra. Sono tutti e tre facoltativi: l'app funziona benissimo anche senza.</p>" +
			"<p class=\"hsub\">Punteggiatura e maiuscole</p>" +
			"<p>Da dove arrivano nel testo i punti, le virgole e le maiuscole. Le opzioni sono tre: <b>dal modello</b> (li mette come li ha sentiti; veloce e di solito sufficiente), <b>dal modello editor</b> (i segni li aggiunge un secondo modello; più curato, ma ogni dettatura dura qualche secondo in più) e <b>togliere</b> (il testo esce tutto di fila, senza segni; serve di rado, per esempio per un campo di ricerca).</p>" +
			"<p class=\"hsub\">Dizionario</p>" +
			"<p>Un elenco di parole che di solito il modello non conosce: cognomi, nomi di aziende, termini tecnici. Non è un ordine «scrivi così», ma un suggerimento: quando sente qualcosa di simile, il modello sceglie più volentieri la tua versione.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">le parole si aggiungono una alla volta e restano come etichette separate</div>" +
			"</div>" +
			"<p class=\"hdim\">Il dizionario lo leggono solo i modelli Whisper. Nemotron, Parakeet, GigaAM e gli altri modelli non accettano suggerimenti di alcun tipo: per loro lo strumento è l'elenco delle sostituzioni più sotto. L'insieme di parole predefinito cambia insieme alla lingua, finché non scrivi il tuo.</p>" +
			"<p class=\"hsub\">Sostituzioni dopo il riconoscimento</p>" +
			"<p>Anche questo è un elenco di parole, ma funziona in un altro modo: non è più un suggerimento, è una garanzia. Quello che compare a sinistra diventa quello che c'è a destra. Ogni sostituzione ha tre impostazioni: <b>Solo parole intere</b> (così «pane» non toccherà «panettone», che di solito è proprio ciò che vuoi), <b>Distingui le maiuscole</b> (così «Rosa» e «rosa» contano come parole diverse) e <b>Lingua della regola</b> (la sostituzione scatta solo quando detti in quella lingua).</p>" +
			"<p class=\"hsub\">Comandi vocali</p>" +
			"<p>Frasi che non finiscono nel testo, ma eseguono un'azione. L'a capo, il nuovo paragrafo e un tuo testo vengono riconosciuti come parole intere in qualsiasi punto della frase. L'annullamento e i comandi di forma agiscono solo alla fine della dettatura, la cancellazione solo se detta da sola. Un comando di forma detto da solo vale per la dettatura successiva, e la barra mostra «Prossima frase: …».</p>" +
			"<ul>" +
			"<li><b>«nuova riga»</b>: andare a capo;</li>" +
			"<li><b>«nuovo paragrafo»</b>: lasciare una riga vuota tra i paragrafi;</li>" +
			"<li><b>«annulla»</b>: scartare la frase e non inserire nulla;</li>" +
			"<li><b>«cancella questo»</b>: cancellare ciò che il programma ha inserito per ultimo;</li>" +
			"<li><b>«come titolo»</b>: mettere un segno davanti alla frase, # di default, e una riga vuota dopo;</li>" +
			"<li><b>«come elenco»</b>: dividere la frase alle virgole, ai punti e virgola e alla parola «e» in righe che iniziano con un trattino;</li>" +
			"<li><b>«con maiuscola»</b>: iniziare la frase con la maiuscola;</li>" +
			"<li><b>«tra virgolette»</b>: mettere la frase tra virgolette;</li>" +
			"<li><b>«senza punto»</b>: togliere il punto finale;</li>" +
			"<li><b>un tuo comando</b>: inserire un testo pronto, come una firma, un indirizzo o una risposta di rito.</li>" +
			"</ul>" +
			"<p>La cancellazione toglie solo ciò che il programma stesso ha inserito per ultimo, e solo se da allora non hai scritto né cliccato, la finestra è la stessa, il testo non andava a capo, non è stato inviato con Invio e non sono passati più di cinque minuti. Il testo che c'era già nella finestra non viene mai toccato.</p>" +
			"<p class=\"wh\" id=\"s-post\">La scheda «Post-elaborazione»</p>" +
			"<p>La post-elaborazione è un secondo modello che riscrive il testo riconosciuto secondo una tua istruzione. Per esempio: «togli gli intercalari e le ripetizioni», «trasformalo in un'e-mail formale», «mettilo sotto forma di elenco». L'istruzione si scrive con parole normali e si chiama prompt.</p>" +
			"<p>Per impostazione predefinita è spenta, ed ecco perché: è l'unico passaggio che rallenta in modo evidente la dettatura. Invece di uno o due secondi, l'inserimento richiederà qualche secondo, e su un computer poco potente anche una decina.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Attiva la post-elaborazione</td><td>L'interruttore principale. Finché è spento, il resto della scheda non funziona, e il testo viene inserito così com'è stato riconosciuto.</td></tr>" +
			"<tr><td>Modello · Locale</td><td>Un editor che gira sul tuo computer. Dal computer non esce nulla, ma serve memoria: questi modelli pesano 1 GB o più. Accanto c'è un catalogo dove il modello si sceglie e si scarica con un clic.</td></tr>" +
			"<tr><td>Modello · Server esterno</td><td>Invece del tuo modello, il testo viene inviato a un servizio esterno: OpenAI, Groq o un tuo server. Più veloce e più intelligente, ma il testo lascia il computer. La chiave di accesso è salvata cifrata, e il pulsante «Prova connessione» verifica che indirizzo e chiave funzionino.</td></tr>" +
			"<tr><td>Prompt</td><td>Le tue istruzioni per l'editor. Quelli spuntati si applicano uno dopo l'altro, dall'alto in basso: prima, per esempio, la pulizia, poi lo stile formale. Se non è spuntato nulla, la post-elaborazione non fa nulla, e l'app te lo segnala con una riga gialla.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">La voce non va mai al server esterno, solo il testo riconosciuto. Ma quello sì che esce, quindi attivalo con cognizione di causa.</div>" +
			"<p>Ogni prompt ha un <b>nome</b> (con quello compare sulla barra mentre lavora), il <b>testo dell'istruzione</b> in parole normali, una <b>scorciatoia propria</b> facoltativa (dettando con questa, la frase passa solo per questo prompt, saltando gli altri) e una <b>prova</b>: un campo dove puoi scrivere una frase e vedere che cosa ne esce, senza dettare nulla.</p>" +
			"<p class=\"wh\" id=\"s-subs\">La scheda «Sottotitoli»</p>" +
			"<p>I sottotitoli trascrivono ciò che il computer riproduce — una chiamata, un video, una lezione — e mostrano le parole in una piccola finestra che resta sopra le altre. È lo stesso riconoscimento della dettatura, solo che il suono arriva dagli altoparlanti invece che dal microfono, e la dettatura continua a funzionare mentre i sottotitoli sono attivi. I suoni dell'app stessa non vengono trascritti.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Avvia i sottotitoli</td><td>Li avvia e li ferma. La riga accanto dice cosa succede: spenti, caricamento del modello, in ascolto o riconoscimento in riavvio. Lo stesso fanno il menu nell'area di notifica, una scorciatoia tutta loro e la croce sulla finestra.</td></tr>" +
			"<tr><td>Lingua del parlato</td><td>La lingua di ciò che ascolti, distinta dalla lingua di dettatura. Il modello è quello che la scheda «Lingue e modelli» assegna a quella lingua. Un modello in streaming scrive le parole mentre vengono dette, gli altri ogni frase dopo la sua fine. Un modello Whisper sul processore ha bisogno di molti secondi per frase, quindi il suo testo resta molto indietro rispetto al parlato; le frasi accumulate gli arrivano insieme ed escono come una sola riga.</td></tr>" +
			"<tr><td>Traduzione</td><td>Mostra nella finestra la traduzione al posto del testo riconosciuto, nella lingua scelta; il titolo della finestra dice allora «Traduzione» e le lingue, per esempio EN → RU. Il parlato viene tagliato in frasi mentre procede, e ogni frase viene tradotta subito. Un pezzo finisce alla fine di una frase o prima di una parola di collegamento, e il traduttore riceve insieme a lui i pezzi precedenti e le proprie traduzioni, così una frase spezzata in parti viene tradotta in modo coerente. La traduzione compare in modo regolare, a parole intere, al ritmo della lettura, e più in fretta quando se n'è accumulata molta, così un pezzo grande non fa saltare il testo. Le frasi di uno stesso tratto di parlato scorrono come un unico paragrafo, e una pausa inizia una riga nuova. Una frase che non è stato possibile tradurre viene mostrata così come è stata sentita, con un colore più tenue. «Mostra anche il testo originale» divide la finestra con una linea sottile: sopra, più tenue, ciò che è stato detto, sotto la traduzione; la linea si trascina con il mouse. La riga «Chi traduce» decide chi lo fa. «Alta qualità» è il modello di traduzione sulla scheda grafica: un download di 2,3 GB, per una scheda grafica con almeno 4 GB di memoria propria; viene caricato sulla scheda all'avvio dei sottotitoli e la lascia quando si fermano, e un link accanto apre le sue condizioni d'uso. «Il tuo server» è il server di post-elaborazione impostato nella scheda «Post-elaborazione», e allora ogni frase viene inviata a quel server. «Sceglie il programma» prende il modello di traduzione se è scaricato e la scheda grafica va bene, altrimenti il tuo server se è impostato, altrimenti il modello editor installato nella scheda «Post-elaborazione», che traduce in modo approssimativo: il senso arriva, ma le parole possono essere imprecise. Un traduttore scelto a mano non viene mai sostituito da un altro: se non può lavorare, il blocco dice cosa manca e il testo è mostrato come è stato sentito. La traduzione viene salvata in un file a parte accanto alla trascrizione: ogni pezzo con l'ora in cui è cominciato; un pezzo che non è stato possibile tradurre viene scritto così come è stato sentito. Senza traduttore non c'è traduzione, e i sottotitoli funzionano come prima.</td></tr>" +
			"<tr><td>Scorciatoia da tastiera</td><td>Avvia e ferma i sottotitoli da qualsiasi finestra. Non avvia mai una dettatura e non può ripetere un'altra scorciatoia.</td></tr>" +
			"<tr><td>Salva la trascrizione</td><td>Ogni sessione scrive un file di testo nella cartella transcripts accanto al programma: una prima riga con data, lingua e modello, poi ogni riga finita con l'ora in cui è cominciata. Finché la prima riga non è finita, non viene scritto nulla. Con la traduzione attiva, il nome della trascrizione porta original e il codice della lingua parlata, e accanto compare un secondo file con translation e il codice della lingua di destinazione, per esempio 2026-11-12_10-18-29_original_EN.txt e 2026-11-12_10-18-29_translation_RU.txt.</td></tr>" +
			"<tr><td>Finestra dei sottotitoli</td><td>Dove sta, quanto è grande e quanto sono grandi le lettere. La posizione si sceglie sull'immagine di uno schermo, come per la barra di dettatura: un clic su uno dei nove posti mette lì la finestra, e con «Posto tuo» può stare ovunque, trascinando il segno sull'immagine o la finestra stessa; con più monitor si può scegliere lo schermo. La dimensione si cambia con il mouse dai bordi e dagli angoli; l'altezza si ferma su righe intere, e larghezza e righe vengono ricordate per ogni risoluzione dello schermo. «Righe di testo» vale dove non c'è ancora una dimensione, e «Ripristina la dimensione» dimentica quelle impostate. La rotella del mouse sopra la finestra scorre il testo indietro di tre righe a scatto; mentre leggi il testo precedente, quello nuovo non lo sposta, e il segno «all'ultimo ↓» nel titolo riporta al parlato in corso.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Le persone con cui parli devono sapere che le loro parole vengono trascritte.</div>" +
			"<p class=\"wh\" id=\"s-about\">«Informazioni» e «Contatti»</p>" +
			"<p>Nella scheda «Informazioni» c'è un breve testo su che cos'è l'app e su come è fatta, e l'elenco di tutte le librerie di terzi con cui è costruita, con le loro licenze. I nomi sono cliccabili e portano alle pagine di quei progetti. C'è anche un link al repository, dove si trovano il codice sorgente e tutte le versioni.</p>" +
			"<p>«Contatti» ha tre riquadri: l'e-mail, il repository e una pagina per bug e idee. Un clic sull'indirizzo e-mail apre il tuo programma di posta con il destinatario già compilato, e l'icona accanto copia semplicemente l'indirizzo.</p>" +
			"<p class=\"wh\" id=\"s-files\">Che cosa c'è nella cartella dell'app</p>" +
			"<p>L'app è portabile: tutte le sue cose stanno in una sola cartella e non si spargono da nessun'altra parte.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Tutte le tue impostazioni in un unico file di testo. Puoi copiarlo su un altro computer per non dover configurare tutto da capo.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>Il registro. Il primo posto da guardare se l'app si comporta in modo strano.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>La cronologia delle dettature, se l'hai attivata. Testo normale: se vuoi, puoi anche solo eliminare il file.</td></tr>" +
			"<tr><td>models</td><td>La cartella con i modelli scaricati, la parte più pesante: qualche GB. Anche questa si può spostare su un altro computer, per non dover riscaricare tutto.</td></tr>" +
			"<tr><td>transcripts</td><td>Le trascrizioni dei sottotitoli, un file di testo per sessione, e le loro traduzioni quando la traduzione è attiva. Sono documenti tuoi: il programma di disinstallazione li elimina solo se gli chiedi di cancellare anche impostazioni e dati.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>I piccoli programmi che fanno girare i modelli. Senza di loro un modello è solo un file che nessuno sa leggere.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Che cosa va su internet e che cosa no</p>" +
			"<ul>" +
			"<li><b>La voce</b> non viene salvata né inviata, a meno che tu stesso non abbia attivato il riconoscimento su un altro computer.</li>" +
			"<li><b>I sottotitoli</b> vengono riconosciuti qui, come la dettatura, e la trascrizione resta nella cartella transcripts; con il riconoscimento su un altro computer, le frasi vanno a quel server; con la traduzione tramite il tuo server di post-elaborazione, ogni riga finita va anche lì.</li>" +
			"<li><b>Il testo</b> non lascia il computer, a meno che tu stesso non abbia attivato un server esterno di post-elaborazione.</li>" +
			"<li><b>La cronologia</b> è un file accanto all'app e si cancella secondo la durata che imposti tu.</li>" +
			"<li><b>La rete</b> serve due volte: per scaricare un modello e per controllare gli aggiornamenti. Il secondo controllo si può disattivare.</li>" +
			"<li><b>La chiave</b> di un servizio esterno è salvata cifrata con gli strumenti di Windows e non finisce mai nel registro.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Se qualcosa va storto</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Il testo non compare nel campo</td><td>Il programma non accetta di incollare dagli appunti. Attiva «Digitazione carattere per carattere» nella scheda «Controlli e comportamento» oppure aggiungi un breve ritardo prima di inserire.</td></tr>" +
			"<tr><td>La dettatura non parte</td><td>La scorciatoia da tastiera è occupata da un altro programma. Impostane un'altra, meglio con il tasto Windows: quelle vengono occupate più di rado.</td></tr>" +
			"<tr><td>È stato registrato solo silenzio</td><td>Il microfono è sbagliato o disattivato. Apri la scheda «Microfono», guarda l'indicatore di livello e premi «Controlla il microfono».</td></tr>" +
			"<tr><td>La prima frase ci mette moltissimo</td><td>È previsto: il modello si carica in memoria la prima volta che serve. Dopo, il riconoscimento è questione di secondi.</td></tr>" +
			"<tr><td>Compare «Riconoscimento non disponibile»</td><td>È scelto il riconoscimento su un altro computer, ma l'indirizzo non è impostato. Scrivi l'indirizzo o torna a «Su questo computer».</td></tr>" +
			"<tr><td>La post-elaborazione non cambia nulla</td><td>Non è spuntato nessun prompt, non è scelto un modello oppure il server esterno non risponde. Il motivo è scritto nella scheda stessa, accanto all'interruttore.</td></tr>" +
			"<tr><td>Il testo è strano, alcune parole sono cambiate</td><td>Dai un'occhiata a «Regole»: molto probabilmente è scattata una sostituzione o un comando vocale.</td></tr>" +
			"<tr><td>Nessuno dei casi precedenti</td><td>«Sistema» → «Apri il registro». L'app ci scrive tutto di sé, e il motivo di solito è spiegato a parole.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">Non hai trovato la risposta? Scrivi un'e-mail dalla scheda «Contatti» o apri una issue nel repository. L'autore legge entrambe.</p>",
		"S_WIZ_HELLO":      "Due minuti di configurazione",
		"S_WIZ_HELLO_TEXT": "{app} trasforma la voce in testo proprio dove si trova il cursore: tieni premuta la combinazione, pronuncia una frase, rilascia — il testo è lì. Tutto gira sul tuo computer, l'audio non esce mai.",
		"S_WIZ_LATER":      "Tutto quello che scegliamo ora si può cambiare dopo nelle impostazioni.",
		"S_WIZ_MODEL_TEXT": "Il modello per la tua lingua viene scelto per te e si scarica in background.",
		"S_WIZ_PICKED":     "Scelto per: %s",
		"S_WIZ_DL_BG":      "Si scarica in background, puoi andare avanti",
		"S_WIZ_DL_FAIL":    "Il download si è interrotto: %s",
		"S_WIZ_RETRY":      "Riprova",
		"S_WIZ_TRY_DL":     "Il modello si sta ancora scaricando (%d%%). Il campo funzionerà appena sarà pronto.",
		"S_WIZ_TRY_LOAD":   "Caricamento del modello…",
		"S_WIZ_TRY_LIVE":   "Con %s le parole compaiono sulla barra mentre stai ancora parlando.",
		"S_WIZ_DONE_DL":    "Il modello si sta ancora scaricando. La dettatura inizierà a funzionare da sola appena sarà pronto.",
		"S_WIZ_INPUT_TEXT": "Questa è la combinazione da tenere premuta mentre parli. Di' qualcosa e controlla che la barra del livello si muova.",
		"S_WIZ_T_TRY":      "Prova",
		"S_WIZ_TRY_PH":     "il testo apparirà qui",
		"S_WIZ_T_DONE":     "Fatto",
		"S_WIZ_DONE_TEXT":  "{app} vive nell'area di notifica: clic sinistro sull'icona per le impostazioni, destro per il menu. Puoi dettare in qualsiasi finestra con un cursore di testo.",
		"S_AUTORUN":        "Avvia con Windows",
		"S_AUTORUN_SUB":    "Una voce nell'avvio automatico dell'utente corrente",
		"S_WIZ_SKIP":       "Salta",
		"S_WIZ_BACK":       "Indietro",
		"S_WIZ_NEXT":       "Avanti",
		"S_WIZ_FINISH":     "Fine",
		"S_WIZ_WAIT":       "Aspetto la prima frase…",
		"S_WIZ_HEARD":      "Sentito:",
		"S_WIZ_HAVE":       "Tutto il necessario è già scaricato",
		"S_WIZ_TRY_TEXT":   "Metti il cursore nel campo qui sotto, tieni premuto %s, di' una frase e rilascia.",
	}
	settingsStrings["pl"] = map[string]string{
		"S_SUBS_ON": "Napisy działają",
		"S_SUBS_HERO_HINT": "włącz je — rozmowa albo film zamienią się w tekst na ekranie",
		"S_SUBS_HERO_START": "Włącz ręcznie",
		"S_SUBS_HERO_STOP": "Wyłącz",
		"S_UPDATED_TO": "Zaktualizowano do wersji %s",
		"S_SUBS_HERO_LIVE": "%s · słowa w trakcie mówienia",
		"S_SUBS_HERO_PHRASES": "%s · tekst po każdym zdaniu",
		"S_NAV_SUBS": "Napisy",
		"S_SUBS_HINT": "Aplikacja zapisuje to, co odtwarza komputer — rozmowę, film — i pokazuje słowa w oknie nad innymi.",
		"S_SUBS_NOTICE": "Twoi rozmówcy powinni wiedzieć, że ich słowa są zapisywane.",
		"S_SUBS_START": "Włącz napisy",
		"S_SUBS_STOP": "Wyłącz napisy",
		"S_SUBS_OFF": "Napisy są wyłączone",
		"S_SUBS_LOADING": "Wczytuję %s…",
		"S_SUBS_LISTEN_LIVE": "Słucham dźwięku komputera · %s, słowa w trakcie mówienia",
		"S_SUBS_LISTEN_PHRASES": "Słucham dźwięku komputera · %s, tekst po każdym zdaniu",
		"S_SUBS_RESTARTING": "Rozpoznawanie uruchamia się ponownie…",
		"S_SUBS_LANG": "Język mowy",
		"S_SUBS_LANG_SUB": "Język tego, czego słuchasz, osobno od języka dyktowania",
		"S_SUBS_M_LIVE": "Pisać będzie %s: słowa pojawiają się w trakcie mówienia.",
		"S_SUBS_M_PHRASES": "Pisać będzie %s: tekst pojawia się po zakończeniu każdego zdania.",
		"S_SUBS_M_ABSENT": "Dla tego języka wybrano %s, ale nie jest zainstalowany.",
		"S_SUBS_M_WHISPER": "%s to model Whisper: na procesorze potrzebuje wielu sekund na każde zdanie, więc tekst będzie mocno spóźniony względem mowy. Do napisów lepiej nadaje się model strumieniowy albo szybki.",
		"S_SUBS_GO_MODELS": "Otwórz „Języki i modele”",
		"S_SUBS_REMOTE": "Rozpoznawanie działa na innym komputerze: każde zdanie jest wysyłane na ten serwer.",
		"S_SUBS_TR": "Tłumaczenie",
		"S_SUBS_TR_SUB": "Okno pokazuje tłumaczenie zamiast rozpoznanego tekstu",
		"S_SUBS_TR_OFF": "Wyłączone",
		"S_SUBS_TR_LOCAL": "Tłumaczy %s: tłumaczenie pojawia się fraza po frazie, chwilę po mowie.",
		"S_SUBS_TR_API": "Tłumaczy twój serwer obróbki końcowej (%s): każda fraza jest na niego wysyłana, a tłumaczenie pojawia się fraza po frazie, chwilę po mowie.",
		"S_SUBS_TR_NONE": "Nie ma modelu tłumaczenia, serwera obróbki końcowej ani modelu redaktora, więc tłumaczenia nie będzie.",
		"S_SUBS_TR_GO": "Otwórz „Obróbka końcowa”",
		"S_SUBS_TR_SAME": "Język tłumaczenia jest taki sam jak język mowy: nie ma czego tłumaczyć.",
		"S_SUBS_WHO": "Kto tłumaczy",
		"S_SUBS_WHO_SUB": "Program bierze najlepszego dostępnego albo ty ustalasz jednego",
		"S_SUBS_WHO_AUTO": "Wybiera program",
		"S_SUBS_WHO_QUALITY": "Dokładny: model na karcie graficznej",
		"S_SUBS_WHO_SERVER": "Twój własny serwer",
		"S_SUBS_TRM": "Model tłumaczenia",
		"S_SUBS_TRM_SUB": "%s, działa na karcie graficznej z co najmniej 4 GB własnej pamięci",
		"S_SUBS_TRM_ABSENT": "Nie pobrany",
		"S_SUBS_TRM_GETTING": "Pobieranie… %s%",
		"S_SUBS_TRM_CHECKING": "Sprawdzanie pliku…",
		"S_SUBS_TRM_READY": "Gotowy",
		"S_SUBS_TRM_RUNNING": "Gotowy, załadowany na kartę graficzną",
		"S_SUBS_TRM_BAD": "Plik w folderze models jest uszkodzony lub niepełny",
		"S_SUBS_TRM_ERR": "Nie pobrany: %s",
		"S_SUBS_TRM_SMALL": "Ten komputer go nie uruchomi: karta graficzna ma %s MB własnej pamięci, a potrzeba 4 GB.",
		"S_SUBS_TRM_INTEGRATED": "Ten komputer go nie uruchomi: ma tylko zintegrowaną grafikę, która dzieli pamięć z komputerem.",
		"S_SUBS_TRM_NOCARD": "Ten komputer go nie uruchomi: Windows nie zgłasza żadnej karty graficznej.",
		"S_SUBS_TRM_UNKNOWN": "Ten komputer go nie uruchomi: ta wersja Windows nie opisuje swojej karty graficznej.",
		"S_SUBS_TRM_DOWN": "Nie uruchomił się: %s",
		"S_SUBS_TRM_DOWN_NEXT": "Do następnego uruchomienia napisów tłumaczy następny w kolejce.",
		"S_SUBS_TRM_DOWN_HEARD": "Do następnego uruchomienia napisów tekst jest pokazywany tak, jak go usłyszano.",
		"S_SUBS_TRM_DEL_SURE": "Usunąć?",
		"S_SUBS_TRM_W_VULKAN": "sterownik grafiki nie obsługuje Vulkan",
		"S_SUBS_TRM_W_DEVICE": "silnik nie widzi tej karty graficznej",
		"S_SUBS_TRM_W_MEMORY": "na karcie graficznej brakuje wolnej pamięci",
		"S_SUBS_TRM_W_CRASH": "silnik zatrzymał się podczas ładowania",
		"S_SUBS_TRM_W_TIMEOUT": "silnik nie odpowiedział na czas",
		"S_SUBS_TRM_W_PORT": "jego port zajmuje inny program",
		"S_SUBS_TRM_W_EXE": "obok programu brakuje llama-server-gpu.exe",
		"S_SUBS_TRM_DL": "Pobierz",
		"S_SUBS_TRM_CANCEL": "Anuluj",
		"S_SUBS_TRM_DEL": "Usuń",
		"S_SUBS_TRM_TERMS": "Warunki korzystania z modelu",
		"S_SUBS_TR_QUALITY": "Tłumaczy model tłumaczenia na karcie graficznej (%s): tłumaczenie pojawia się fraza po frazie, chwilę po mowie.",
		"S_SUBS_TR_MISS_CARD": "Wybrano „Dokładny”, ale karta graficzna tego komputera się nie nadaje, więc tekst jest pokazywany tak, jak go usłyszano.",
		"S_SUBS_TR_MISS_MODEL": "Wybrano „Dokładny”, ale model tłumaczenia nie jest pobrany, więc tekst jest pokazywany tak, jak go usłyszano.",
		"S_SUBS_TR_MISS_START": "Wybrano „Dokładny”, ale model tłumaczenia nie uruchomił się na karcie graficznej, więc do następnego uruchomienia napisów tekst jest pokazywany tak, jak go usłyszano.",
		"S_SUBS_TR_MISS_SERVER": "Wybrano twój własny serwer, ale serwer obróbki końcowej nie jest ustawiony, więc tekst jest pokazywany tak, jak go usłyszano.",
		"S_SUBS_ORIG": "Pokazuj też tekst oryginalny",
		"S_SUBS_ORIG_SUB": "okno pokazuje nad tłumaczeniem to, co zostało powiedziane; linię między nimi można przeciągać myszą",
		"S_SUBS_KEY": "Skrót klawiszowy",
		"S_SUBS_KEY_SUB": "Włącza i wyłącza napisy z każdego okna",
		"S_SUBS_SAVE": "Zapisuj transkrypcję",
		"S_SUBS_SAVE_SUB": "Jeden plik tekstowy na każde uruchomienie w folderze transcripts obok programu, a przy włączonym tłumaczeniu drugi, z tłumaczeniem",
		"S_SUBS_FOLDER": "Folder transkrypcji",
		"S_SUBS_WINDOW": "Okno napisów",
		"S_SUBS_ROWS": "Wiersze tekstu",
		"S_SUBS_SIZE": "Rozmiar tekstu",
		"S_SUBS_SIZE_S": "Mały",
		"S_SUBS_SIZE_M": "Średni",
		"S_SUBS_SIZE_L": "Duży",
		"S_SUBS_BOX": "Rozmiar okna",
		"S_SUBS_BOX_SUB": "przeciągnij krawędzie i rogi okna myszą; rozmiar jest zapamiętywany dla każdej rozdzielczości ekranu, a „Wiersze tekstu” obowiązują, dopóki go nie ma",
		"S_SUBS_BOX_BTN": "Przywróć rozmiar",
		"S_SUBS_BOX_DONE": "Okno wróci do domyślnego rozmiaru",
		"S_SUBS_FREE": "Własne miejsce",
		"S_SUBS_FREE_SUB": "okno można przeciągnąć gdziekolwiek, miejsce jest zapamiętywane dla każdej rozdzielczości ekranu",
		"S_SUBS_POS": "Gdzie pokazywać okno",
		"S_SUBS_POS_SCHEME_SUB": "kliknij ekran — okno stanie w tym miejscu",
		"S_SUBS_POS_DRAG_SUB": "przeciągnij znacznik albo samo okno — stanie gdziekolwiek",
		"S_SUBS_MON": "Ekran",
		"S_SUBS_MON_SUB": "na którym monitorze pokazywać okno napisów",
		"S_TITLE": "{app} — Ustawienia", "S_DICT_HINT": "Terminy, nazwy i skróty oddzielone przecinkami — podpowiedź dla słuchu, nie polecenia. Działa dla Whispera; rosyjska mowa przez GigaAM tego nie używa. Zestaw domyślny podąża za językiem rozpoznawania, dopóki nie wpiszesz własnego.",
		"S_TR_DEFAULT": "Zmień język wyjściowego tekstu", "S_TR_TARGET": "Domyślny język wyjściowego tekstu", "S_TR_ASK": "Pytaj o język wyjściowego tekstu", "S_TR_ASK_NEVER": "Nie pytaj — tłumacz od razu",
		"S_SRCLANG_SUB": "mówisz w nim; określa model rozpoznawania",
		"S_TR_LANGS_SUB": "te języki będą przyciskami na pasku przy wstawianiu",
		"S_TR_UNAVAIL": "%s nie umie tłumaczyć — tekst zostanie w języku źródłowym",
		"S_TR_GET_WHISPER": "Pobierz %s do tłumaczenia (%d MB)",
		"S_TR_STANDIN": "%s nie umie tłumaczyć — pasek zaproponuje %s",
		"S_TR_LOCK": "%s nie można usunąć z listy — to domyślny język wyjściowego tekstu. Wybierz inny język domyślny, a wtedy %s będzie można wykluczyć.",
		"S_TR_LOCK_OK": "Rozumiem",
		"S_TR_ONE": "Zaznaczono kilka języków, ale bez pytania tłumaczenie zawsze pójdzie w jeden — %s (domyślny język wyjściowy). Pozostałe zostaną zaznaczone, ale wyłączone.",
		"S_TR_NOMODEL": "%s nie umie tłumaczyć. Jeśli kontynuujesz, tłumaczenie zostanie wyłączone i niedostępne, dopóki pracuje ten model.",
		"S_TR_CONFIRM": "Potwierdź",
		"S_TR_ASK_ALWAYS": "Pytaj za każdym razem", "S_TR_ASK_TIMEOUT": "Pytaj, z limitem czasu", "S_TR_SECONDS": "Limit, s",
		"S_TR_LANGS": "Języki w oknie dialogowym",
		"S_LLM_HINT": "Zaznaczone profile działają po kolei, z góry na dół, przy zwykłym dyktowaniu. Nic nie zaznaczono — tekst wstawia się bez zmian.",
		"S_PROF_ADD": "Dodaj", "S_PROF_NAME": "Nazwa", "S_PROF_PROMPT": "Prompt", "S_PROF_TEST": "Test",
		"S_HOTKEY": "Skrót klawiszowy", "S_CHANGE": "Zmień…", "S_UILANG": "Język interfejsu", "S_AUTO": "Jak w systemie",
		"S_SEC_SOUND": "Dźwięk", "S_BEEP": "Sygnały dźwiękowe", "S_SOUND": "Dźwięk sygnału",
		"S_SND_SPEECH": "Głos Windows", "S_SND_CHIME": "Dzwonek", "S_SND_SOFT": "Miękki", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Naciśnij Enter po wstawieniu (auto-wysyłka)", "S_RESTORE": "Przywróć schowek po wstawieniu",
		"S_NAV_HISTORY": "Historia", "S_HIST_ON": "Przechowuj historię dyktowań", "S_HIST_ON_SUB": "tylko tekst, na tym komputerze; dźwięk nigdy nie jest zapisywany",
		"S_HIST_DAYS": "Ile dni przechowywać", "S_HIST_MAX": "Ile wpisów przechowywać",
		"S_HIST_SKIP": "Nigdy nie zapisuj z tych programów", "S_HIST_SKIP_SUB": "po przecinku: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Dodawanie programu", "S_SKIP_EDIT_DLG": "Edycja programu", "S_SKIP_NAME": "Nazwa programu", "S_SKIP_NAME_SUB": "Nazwa pliku bez ścieżki: keepass.exe. Gwiazdka na końcu łapie wszystkie wersje: 1password*", "S_SKIP_OPEN": "Programy otwarte teraz", "S_SKIP_REFRESH": "Odśwież listę", "S_SKIP_PICKED": "Wybrano %d z %d", "S_SKIP_NONE": "Nic nie wybrano", "S_SKIP_EMPTY": "Lista jest pusta — historia zapisuje się ze wszystkich programów", "S_SKIP_ADD_BTN": "Dodaj program", "S_SKIP_HINT": "To, co dyktujesz do tych programów, nie trafia do historii. Samo wklejanie działa jak zwykle.",
		"S_HIST_LIST": "Wpisy", "S_HIST_CLEAR": "Wyczyść", "S_HIST_TILL": "do %s", "S_HIST_TILL1": "do jutra", "S_HIST_TILL_FULL": "Zostanie usunięte %s — przechowywanie %s", "S_HIST_LIST_HINT": "To, co podyktowano: skopiować, wkleić do dowolnego okna lub usunąć.", "S_HIST_COPY": "Kopiuj",
		"S_HIST_KEEP": "Jak długo przechowywać",
		"S_UNIT_MIN": "minut",
		"S_UNIT_HOUR": "godzin",
		"S_UNIT_DAY": "dni",
		"S_HIST_FIND": "Szukaj w historii…", "S_HIST_EMPTY": "Na razie brak historii", "S_HIST_ASK": "Usunąć całą historię dyktowań?",
		"S_SEC_CMD": "Komendy głosowe", "S_CMD_HINT": "To, co powiesz, zamienia się w złamanie wiersza, znak, anulowanie, usunięcie ostatniego wstawienia albo formatowanie frazy, zamiast trafić do tekstu. Rozpoznawane jako całe słowa, stosowane od góry do dołu, po zamianach; komendy formatowania działają na końcu dyktowania.",
		"S_CMD_ADD": "Dodaj komendę", "S_CMD_PRESET": "Dodaj typowe", "S_CMD_PH": "fraza, którą wypowiesz",
		"S_CMD_NEWLINE": "złamanie wiersza", "S_CMD_PARAGRAPH": "nowy akapit", "S_CMD_TEXT": "wstawić tekst", "S_CMD_CANCEL": "anulować dyktowanie",
		"S_CMD_TEXT_PH": "co wstawić", "S_CMD_EMPTY": "Na razie brak komend", "S_CMD_DEL": "Usuń komendę",
		"S_CMD_P_NEWLINE": "nowy wiersz", "S_CMD_P_PARAGRAPH": "nowy akapit", "S_CMD_P_CANCEL": "anuluj",
		"S_CMD_UNDO": "usunąć ostatnie wstawienie", "S_CMD_HEADING": "zrobić nagłówek", "S_CMD_LIST": "zrobić listę",
		"S_CMD_CAPITAL": "wielka pierwsza litera", "S_CMD_QUOTES": "ująć w cudzysłów", "S_CMD_NOSTOP": "usunąć kropkę na końcu",
		"S_CMD_HEADING_PH": "znak przed nagłówkiem, np. #", "S_CMD_LIST_PH": "słowa łączące ostatni punkt, np. i", "S_CMD_QUOTES_PH": "znak otwierający i zamykający, np. „”",
		"S_CMD_END_NOTE": "Powiedz na końcu dyktowania; powiedziana osobno, komenda sformatuje następne.",
		"S_CMD_UNDO_NOTE": "Powiedz osobno. Usuwa ostatnie wstawienie, jeśli od tamtej pory nic nie wpisano ani nie kliknięto, a okno jest to samo; tekst, który już tam był, nigdy nie jest ruszany.",
		"S_CMD_P_UNDO": "usuń to", "S_CMD_P_HEADING": "jako nagłówek", "S_CMD_P_LIST": "jako lista", "S_CMD_P_CAPITAL": "wielką literą",
		"S_CMD_P_QUOTES": "w cudzysłowie", "S_CMD_P_NOSTOP": "bez kropki", "S_CMD_P_AND": "i, oraz", "S_CMD_P_MARKS": "„”",
		"S_SEC_REPLACE": "Zamiany po rozpoznaniu", "S_REPLACE_HINT": "To, co zostało źle usłyszane, staje się tym, co miałeś na myśli — zaraz po rozpoznaniu, przed promptami. Stosowane od góry do dołu.",
		"S_REPL_WHOLE_FULL": "Tylko całe słowa", "S_REPL_CASE_FULL": "Uwzględniaj wielkość liter", "S_CMD_ACTION": "Działanie",
		"S_FM_ADD": "Dodaj",
		"S_TIP_REPL_LANG": "Reguła działa tylko, gdy dyktujesz w wybranym języku. „wszystkie języki” — działa zawsze.",
		"S_TIP_REPL_CASE": "Wielkość liter ma znaczenie: „git” i „Git” to różne słowa. Wyłączone — wielkość liter jest pomijana.",
		"S_TIP_REPL_WHOLE": "Zamiana działa tylko, gdy tekst stoi jako osobne słowo. Wyłączone — trafia też wewnątrz innych słów.",
		"S_TIP_CMD_ACTION": "Co się stanie, gdy wypowiesz frazę: nowa linia, nowy akapit, własny tekst, anulowanie dyktowania, usunięcie ostatniego wstawienia albo formatowanie frazy — nagłówek, lista, wielka litera, cudzysłów lub bez kropki.",
		"S_LIST_FILTER_PH": "szukaj…",
		"S_REPL_DEL": "Usuń zamianę",
		"S_LIST_NOTHING": "Nic nie znaleziono: „%s”",
		"S_FM_T_REPL_ADD": "Dodawanie zamiany", "S_FM_T_REPL_EDIT": "Edycja zamiany",
		"S_FM_T_CMD_ADD": "Dodawanie polecenia", "S_FM_T_CMD_EDIT": "Edycja polecenia",
		"S_MT_DEL": "Usuwanie modelu", "S_MT_DEL_PROMPT": "Usuwanie promptu", "S_MT_DL": "Pobieranie modelu",
		"S_MT_TR_OFF": "Wyłączenie tłumaczenia", "S_MT_TR_ONE": "Tłumaczenie bez pytania", "S_MT_TR_LOCK": "Domyślny język tekstu",
		"S_MT_REMOTE": "Zdalny serwer", "S_MT_POST": "Zewnętrzny serwer", "S_MT_HIST": "Czyszczenie historii",
		"S_MT_RESET": "Reset ustawień", "S_MT_EXE": "Ścieżka serwera",
		"S_DICT_ADD": "Dodaj słowo", "S_FM_T_DICT_ADD": "Dodawanie słowa", "S_DICT_EMPTY": "Brak słów",
		"S_DICT_ADD_PH": "słowo lub kilka po przecinku",
		"S_DICT_NOMODEL": "Bieżący model %s nie obsługuje słownika — czytają go tylko modele Whisper.",
		"S_OV_FREE": "Własne miejsce", "S_OV_FREE_SUB": "pasek można przeciągnąć gdziekolwiek",
		"S_OVPOS_DRAG_SUB": "przeciągnij pasek myszą — stanie gdziekolwiek",
		"S_OVMON_N": "Ekran %d",
		"S_POST_ENABLE": "Włącz przetwarzanie końcowe",
		"S_API_SUM_URL": "adres", "S_API_SUM_MODEL": "model", "S_API_SUM_KEY": "klucz", "S_API_SUM_TIMEOUT": "oczekiwanie",
		"S_API_SUM_STATE": "stan", "S_API_NO_MODEL": "nie podano",
		"S_API_NONE": "nieskonfigurowany — przetwarzanie idzie lokalnie",
		"S_POSTAPI_SETUP": "Skonfiguruj", "S_API_EDIT": "Zmień", "S_API_KEY_DEL": "Usuń klucz", "S_API_DLG": "Serwer zewnętrzny",
		"S_LLM_CATALOG": "Katalog modeli", "S_LLM_BLOCK": "Zainstalowane modele", "S_LLM_NONE_HINT": "Nie zainstalowano żadnego modelu — pobierz znaleziony strzałką, pojawi się tutaj", "S_LLM_IN_MEM": "w pamięci", "S_LLM_ON_DISK": "na dysku", "S_LLM_EJECT": "Zwolnij z pamięci", "S_LLM_FOUND": "znaleziono %d", "S_LLM_NOSEARCH": "nie szukano", "S_LLM_SEARCH_HINT": "Wpisz nazwę modelu i naciśnij „Szukaj”", "S_LLM_PICK_WAIT": "Dostępny po pobraniu modelu", "S_LLM_INSTALLED": "zainstalowane",
		"S_LLM_SUM_MODEL": "model", "S_LLM_SUM_SIZE": "rozmiar", "S_LLM_SUM_COUNT": "zainstalowano", "S_LLM_SUM_RAM": "pamięć",
		"S_DLG_CLOSE": "Zamknij", "S_LLM_NOPICK": "nie wybrano", "S_NO_PROMPTS": "Brak promptów", "S_PROF_DRAG": "przeciągnij, aby zmienić kolejność",
		"S_PROF_NAME_PH": "jak nazwać prompt", "S_PROF_TEST_PH": "wpisz zdanie do sprawdzenia",
		"S_PF_NEW": "Nowy prompt", "S_PF_EDIT": "Edycja promptu",
		"S_POST_NO_MODEL": "włączone, ale nie wybrano modelu", "S_POST_NO_API": "włączone, ale serwer nie jest skonfigurowany", "S_POST_BAD": "serwer nie odpowiedział: %s", "S_POST_NO_PROMPT": "włączone, ale nie zaznaczono żadnego promptu", "S_API_TEST": "Test połączenia", "S_API_TEST_RUN": "Sprawdzam…", "S_API_TEST_OK": "Serwer odpowiedział", "S_API_CLEAR": "Wyczyść", "S_API_CLEAR_ASK": "Usunąć adres, model i klucz serwera zewnętrznego? Przetwarzanie wróci do modelu lokalnego.", "S_RAM_AVAIL": "Dostępna pamięć: %s GB z %s GB", "S_RAM_OF": "%s GB z %s GB",
		"S_REPL_ADD": "Dodaj zamianę", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "całe słowa", "S_REPL_CASE": "wielkość liter", "S_REPL_EMPTY": "Na razie brak zamian",
		"S_PASTE_DELAY": "Opóźnienie przed wstawieniem", "S_PASTE_DELAY_SUB": "gdy program nie zdąża przyjąć tekstu",
		"S_OVPOS": "Gdzie pokazywać pasek", "S_OVPOS_SUB": "przy kursorze — obok miejsca pisania; jeśli aplikacja go nie pokazuje, obok wskaźnika myszy",
		"S_OVPOS_CARET": "Przy kursorze",
		"S_OVTEXT": "Pokazywać rozpoznany tekst", "S_OVTEXT_SUB": "na pasku po wstawieniu, zamiast liczby znaków",
		"S_OVERLAY": "Pokazuj pasek", "S_OVERLAY_SUB": "podczas dyktowania ekran pokazuje, że trwa nagrywanie", "S_TYPEMODE": "Wpisywanie znak po znaku",
		"S_RECLANG": "Język mowy źródłowej", "S_RECAUTO": "Auto",
		"S_LANG_KEY": "Własny skrót", "S_LANG_KEY_SET": "Przypisz…", "S_LANG_KEY_CLEAR": "Usuń skrót", "S_LANG_KEY_HINT": "Przytrzymaj go, aby dyktować w tym języku jego modelem. Główny skrót nadal używa języka mowy źródłowej wybranego wyżej.",
		"S_DL": "Pobierz", "S_DEL": "Usuń",
		"S_M_BASE": "szybki, słabe komputery", "S_M_SMALL": "zrównoważony", "S_M_MED": "Whisper dla każdego języka; ten, który tłumaczy", "S_M_TURBO": "najlepsza dokładność na CPU", "S_M_GIGAAM": "dokładniejszy po rosyjsku, sam stawia interpunkcję", "S_M_PARAKEET": "25 języków europejskich, sama stawia interpunkcję",
		"S_THREADS":  "Wątki CPU", "S_MINMS": "Min. nagranie, ms", "S_MAXSEC": "Maks. nagranie, s",
		"S_AUTOSTART": "Uruchamiaj whisper-server automatycznie", "S_PORT": "Port", "S_SERVEREXE": "Ścieżka whisper-server", "S_SERVEREXE_SUB": "wypełnia się sama; zmieniaj tylko, jeśli serwer leży gdzie indziej", "S_EXE_RESET": "Przywróć", "S_EXE_WARN": "Aplikacja znajduje whisper-server obok siebie. Przy ręcznie wpisanej ścieżce po przeniesieniu folderu rozpoznawanie przestanie się uruchamiać. Zmienić?", "S_RESET_ALL": "Przywróć ustawienia", "S_RESET_ALL_SUB": "wszystko wraca do fabrycznych, poza modelami i historią", "S_RESET_ALL_BTN": "Przywróć", "S_RESET_ALL_ASK": "Przywrócić wszystkie ustawienia fabryczne? Modele, historia i prompty zostają.", "S_RELOAD_CFG": "Wczytaj config.json ponownie", "S_RELOAD_CFG_SUB": "jeśli plik był zmieniany ręcznie", "S_RELOAD_CFG_BTN": "Wczytaj", "S_UPD_FOUND": "Jest wersja %s", "S_THEME": "Kolor", "S_THEME_SUB": "kolor okna, paska i ikony w zasobniku", "S_THEME_GREEN": "Zielony", "S_THEME_AMBER": "Bursztynowy", "S_THEME_BLUE": "Niebieski", "S_THEME_PINK": "Różowy", "S_THEME_EDITOR": "Edytor", "S_THEME_NEON": "Neon", "S_WND_MAX": "Na cały ekran", "S_WND_RESTORE": "Poprzedni rozmiar", "S_WND_MIN": "Ukryj w zasobniku", "S_WND_CLOSE": "Zamknij okno", "S_SKIN": "Wygląd", "S_SKIN_SUB": "krój pisma, kształty, efekty i animacja", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Miękki", "S_SKIN_PAPER": "Dokument", "S_SKIN_FLUENT": "Fluent", "S_SKIN_STUDIO": "Studio",
		"S_SERVERURL": "Serwer zewnętrzny (URL)", "S_URLHINT": "Jeśli ustawiony, lokalny serwer nie startuje",
		"S_STT_SRV": "Serwer rozpoznawania",
		"S_STT_SRV_HINT": "Modele Whisper uruchamia osobny program. Może działać na tym komputerze albo na innym — wybierz, którego użyć.",
		"S_SRV_LOCAL": "Na tym komputerze",
		"S_SRV_REMOTE": "Na innym komputerze",
		"S_SRV_REMOTE_HINT": "Ten sam whisper-server, uruchomiony gdzie indziej: serwer domowy, komputer z kartą graficzną, maszyna obok.",
		"S_SRV_K_AUTO": "autostart",
		"S_SRV_K_FILE": "plik",
		"S_SRV_K_ADDR": "adres",
		"S_SRV_K_CHECK": "ostatni test",
		"S_SRV_NEAR": "whisper-server.exe obok aplikacji",
		"S_SRV_NOADDR": "nie podano",
		"S_SRV_NOCHECK": "nie sprawdzano",
		"S_SRV_LOCAL_DLG": "Lokalny serwer rozpoznawania",
		"S_SRV_ADDR": "Adres serwera",
		"S_SRV_ADDR_SUB": "adres whisper-server na drugiej maszynie, razem z portem",
		"S_SRV_ON": "włączony",
		"S_SRV_OFF": "wyłączony",
		"S_SRV_K_THREADS": "wątki CPU",
		"S_SRV_K_PORT": "port",
		"S_SRV_DOWN": "Rozpoznawanie niedostępne",
		"S_SRV_DOWN_WHY": "zdalny serwer rozpoznawania nie jest skonfigurowany — podaj adres w ustawieniach",
		"S_SRV_DOWN_GO": "Otwórz ustawienia serwera",
		"S_SRV_WARN_NOW": "Dyktowanie teraz nie działa: wybrany jest serwer zdalny, ale nie podano jego adresu.",
		"S_SRV_WARN_LATER": "Gdy zostanie wybrany model Whisper, rozpoznawanie nie zadziała: brak adresu serwera zdalnego.",
		"S_SAVED":      "Zapisano",
		"S_ABOUT_HTML": "<p><b>Głos → tekst w pozycji kursora.</b></p><p>Ustaw kursor, przytrzymaj skrót, mów, puść — tekst zostanie wstawiony.</p><p>W pełni lokalnie i offline. Technologie: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; modele z Hugging Face.</p><p>Logi nigdy nie przekraczają ~2 MB.</p>",
		  "S_SEARCH": "Znajdź ustawienie…",
		"S_GRP_GENERAL": "Ogólne", "S_GRP_SPEECH": "Przetwarzanie mowy", "S_GRP_INFO": "Informacje", "S_NAV_POST": "Obróbka końcowa", "S_NAV_HELP": "Pomoc", "S_NAV_CONTACTS": "Kontakty", "S_HIST_ADD": "Dodaj", "S_CONTACT_MAIL": "E-mail", "S_DICT_MODEL": "Model rozpoznawania", "S_LIB_ACC": "dokładność", "S_LIB_SPD": "szybkość",
		"S_HELP_TOC": "Na tej stronie",
		"S_HELP_TOC_SHOW": "Pokaż spis treści — okno stanie się szersze",
		"S_HELP_TOC_HIDE": "Ukryj spis treści i przywróć szerokość okna",
		"S_CONTACT_TITLE": "Kontakt",
		"S_ABOUT_DEPS": "Moduły zewnętrzne",
		"S_ABOUT_DEPS_HINT": "Cudzy kod wbudowany w aplikację i jego licencje. Kliknięcie nazwy otwiera stronę projektu.",
		"S_DEP_WHISPER": "uruchamia modele Whisper",
		"S_DEP_LLAMA": "obróbka tekstu, modele GGUF, tłumaczenie napisów na karcie graficznej (llama-server-gpu.exe)",
		"S_DEP_VULKAN": "nagłówki i biblioteka importu silnika dla karty graficznej; samo środowisko Vulkan przychodzi ze sterownikiem grafiki",
		"S_DEP_SHERPA": "silnik GigaAM, Parakeet, Canary, Qwen, Moonshine",
		"S_DEP_GGML": "biblioteka tensorów w whisper.cpp i llama.cpp",
		"S_DEP_ONNX": "uruchamia modele w sherpa-onnx",
		"S_DEP_WEBVIEW": "okno ustawień na WebView2",
		"S_DEP_WV2RT": "składnik Windows rysujący to okno",
		"S_DEP_MALGO": "przechwytywanie dźwięku z mikrofonu",
		"S_DEP_MINIAUDIO": "warstwa audio w malgo",
		"S_DEP_WS": "połączenie z sherpa-server",
		"S_DEP_XSYS": "wywołania WinAPI z Go",
		"S_DEP_WINLOADER": "ładowanie DLL w go-webview2",
		"S_DEP_PLEX": "czcionka interfejsu",
		"S_DEP_HF": "katalog, z którego pobierane są modele",
		"S_CONTACT_HINT": "Błąd, pomysł, pytanie o ustawienie — napisz e-mail, jeśli sprawa jest osobista, albo załóż issue, jeśli to błąd.",
		"S_CONTACT_REPO": "Repozytorium",
		"S_CONTACT_ISSUES": "Błędy i pomysły",
		"S_CONTACT_WRITE": "Napisz e-mail",
		"S_CONTACT_OPEN": "Otwórz",
		"S_STATE_ACTIVE": "Rozpoznaje",
		"S_STATE_USED": "Modele w użyciu",
		"S_STATE_INST": "Zainstalowane lokalnie",
		"S_STATE_INST_SUB": "modele na dysku, gotowe do przypisania",
		"S_PRESETS": "Który model dla którego języka",
		"S_PRESETS_HINT": "Model z wiersza «Auto» obsługuje wszystkie języki bez własnego modelu — na liście są szare. Kliknij język, aby nadać mu własny. Ten sam model rozpoznaje też, w jakim języku mówisz, gdy język źródłowy nie jest ustawiony.",
		"S_MFOLDER": "Własny model",
		"S_DICT_SAVE": "Zapisz",
		"S_OWNM_SUB": "Dodaj lokalny model rozpoznawania mowy",
		"S_OWNM_ONEFILE": "Jeden plik",
		"S_OWNM_FOLDERF": "Folder z plikami modelu",
		"S_OWNM_S1": "Otwórz folder modeli",
		"S_OWNM_S1S": "Folder docelowy:",
		"S_OWNM_S2": "Skopiuj model",
		"S_OWNM_S2S": "Wybierz jedną z obsługiwanych struktur",
		"S_OWNM_S3": "Uruchom aplikację ponownie",
		"S_OWNM_S3S": "Model pojawi się dla języków, które obsługuje",
		"S_AS_AUTO": "jak automatyczne wykrywanie",
		"S_REC_CHIP": "polecany",
		"S_BACK_AUTO": "Wróć do automatycznego wykrywania",
		"S_LANGS_COUNT": "języków: %d",
		"S_LANGS_UNKNOWN": "języki: nieznane",
		"S_TR_EN": "tłumaczy na %d języków",
		"S_TR_LIST": "tłumaczy: %s",
		"S_DL_GOING": "pobieranie:",
		"S_OPEN_FOLDER": "Otwórz folder",
		"S_UNLOAD": "Wyładuj z pamięci",
		"S_UNLOAD_SUB": "zwalnia pamięć; następne dyktowanie ponownie wczyta model",
		"S_UNLOAD_GO": "Wyładuj",
		"S_UNLOADED": "Wyładowano",
		"S_NOT_FOR_LANG": "%s nie rozpoznaje tego języka",
		"S_MANUAL_NOTE": "Nie można pobrać z aplikacji — licencja zabrania rozpowszechniania. Pobierz archiwum samodzielnie i rozpakuj do models/moonshine-uk.",
		"S_MANUAL_LINK": "Pobierz samodzielnie",
		"S_HF_FIT": "tylko pasujące do tego komputera",
		"S_HF_HIDDEN": "ukryte: %s",
		"S_WIZ_SKIP_DL": "Pobierz później",
		"S_WIZ_SKIP_NOTE": "Bez modelu dyktowanie nie zadziała. Pobierzesz go w sekcji „Języki i modele”.",
		"S_M_GIGAAM2": "poprzednia generacja rosyjskiego modelu: ta sama szybkość, ale bez interpunkcji",
		"S_M_MOONUK": "ukraiński model Moonshine: szybki i lekki, bez interpunkcji",
		"S_M_LOCAL": "znaleziony w folderze models; właściwości nieznane, dlatego bez pasków",
		"S_ALL_LANGS": "wszystkie języki",
		"S_OVPOS_SCHEME_SUB": "kliknij ekran — pasek stanie w tym miejscu",
		"S_OVDRAG": "Przeciągnij, gdzie chcesz",
		"S_OVMON": "Ekran",
		"S_OVMON_SUB": "na którym monitorze pokazywać pasek",
		"S_OVMON_CURSOR": "Ekran z kursorem",
		"S_M_NEMOTRON": "pisze w trakcie mowy: tekst pojawia się na pasku na żywo; 40 języków, sama stawia interpunkcję",
		"S_M_TINY": "najmniejszy i najszybszy, dla bardzo słabych komputerów; wyraźnie mniej dokładny",
		"S_STATE_LOADED": "Teraz w pamięci",
		"S_STATE_LOADED_SUB": "modele wyładowują się same po bezczynności",
		"S_STATE_WEEK": "W tym tygodniu",
		"S_ST_SUMMARY": "Podsumowanie", "S_ST_OVERLAY": "Pasek na ekranie", "S_ST_BEEP": "Sygnał dźwiękowy", "S_ST_AUTORUN": "Start z Windows", "S_ST_POST": "Przetwarzanie", "S_ST_LOCAL": "lokalnie", "S_ST_CHECKED": "sprawdzono %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "włączony", "S_ST_OFF_M": "wyłączony", "S_ST_MIC_OK": "sygnał w normie", "S_ST_MIC_BAD": "mikrofon milczy", "S_ST_CHECK": "Sprawdź", "S_ST_RECOG": "rozpoznaje %s", "S_ST_VER": "Wersja %s", "S_ST_LATEST": "najnowsza", "S_ST_OUTDATED": "nieaktualna", "S_ST_UPD_OK": "masz najnowszą wersję", "S_ST_UPD_DL": "Pobieram aktualizację…",
		"S_ST_QUICK": "Szybkie ustawienia",
		"S_ST_MODELS": "Modele",
		"S_ST_USAGE": "W tym tygodniu",
		"S_ST_READY": "Gotowe do dyktowania",
		"S_ST_OFF": "Wyłączone w zasobniku",
		"S_ST_OFF_SUB": "skrót nic nie robi, dopóki go nie włączysz",
		"S_ST_ENABLE": "Włącz",
		"S_ST_GOTO": "Otwórz to ustawienie na jego karcie",
		"S_ST_HOTKEY_GO": "Zmień skrót",
		"S_ST_UPD_LAST": "Wersja %s — najnowsza",
		"S_ST_UPD_HAVE": "Dostępna wersja %s",
		"S_ST_MEM": "Wolne %s GB z %s",
		"S_ST_MEM_SUB": "w pamięci: %s · na dysku: %d modeli, %s GB",
		"S_ST_MEM_NONE": "nic",
		"S_ST_LANG": "Język mowy",
		"S_ST_ASR": "Rozpoznawanie",
		"S_ST_ON": "włączone",
		"S_ST_OFF_W": "wyłączone",
		"S_ST_ON_F": "włączona",
		"S_ST_OFF_F": "wyłączona",
		"S_ST_ACTIVE": "aktywna",
		"S_ST_IDLE": "nie uruchamia się",
		"S_ST_DISK": "na dysku, %s",
		"S_ST_USAGE_SUB": "%d znaków · dziś %d · średnio %d znaków",
		"S_WEEK_OTHER": "pozostałe",
		"S_ST_NO_WEEK": "w tym tygodniu bez dyktowań",
		"S_ST_AUTORUN_SUB": "aplikacja nie uruchomi się sama",
		"S_ST_OVERLAY_SUB": "widoczny podczas nagrywania",
		"S_REPL_LANG": "Język reguły",
		"S_REPL_LANG_ALL": "wszystkie języki",
		"S_M_CANARY": "angielski, niemiecki, hiszpański, francuski — i sam tłumaczy między nimi",
		"S_M_QWEN3": "około 30 języków, sam stawia interpunkcję; najcięższy i najdokładniejszy w katalogu",
		"S_POSTAPI": "Serwer zewnętrzny",
		"S_POST_HINT": "Poprawia rozpoznany tekst według promptów: usuwa słowa-wypełniacze, naprawia interpunkcję, zmienia styl. Wyłączone — tekst wstawiany jest tak, jak rozpoznano.",
		"S_POST_MODEL": "Model",
		"S_SRC_LOCAL": "Lokalny",
		"S_SRC_USED": "w użyciu",
		"S_HF_GO": "Szukaj",
		"S_POSTAPI_HINT": "Domyślnie puste — cała obróbka działa lokalnie. Wpisz adres, a prompty wykona serwer zewnętrzny: OpenAI, Groq, własny vLLM — cokolwiek ze zgodnym API.",
		"S_POSTAPI_URL": "Adres",
		"S_POSTAPI_URL_SUB": "puste = model lokalny; przykład: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Model",
		"S_POSTAPI_KEY": "Klucz API",
		"S_POSTAPI_KEY_SET": "klucz zapisany (zaszyfrowany Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "brak klucza",
		"S_POSTAPI_SAVE": "Zapisz klucz",
		"S_POSTAPI_TIMEOUT": "Czas oczekiwania", "S_SEC_SHORT": "s",
		"S_POSTAPI_WARN": "⚠ Rozpoznany tekst dyktowań będzie trafiał pod ten adres. Dźwięk nigdy nie wychodzi. Klucz jest przechowywany zaszyfrowany.",
		"S_POSTAPI_ASK": "Wysyłać rozpoznany tekst do %s? Dźwięk zostaje na tym komputerze, ale tekst będzie go opuszczał.",
		"S_POSTAPI_BADGE": "serwer zewnętrzny",
		"S_NOT_INSTALLED": "niezainstalowany",
		"S_NAV_STATE": "Stan", "S_NAV_DICT": "Sterowanie i zachowanie", "S_NAV_MIC": "Mikrofon", "S_NAV_MODELS": "Języki i modele",
		"S_NAV_TEXT": "Reguły", "S_NAV_TR": "Tłumaczenie", "S_NAV_SYSTEM": "System", "S_NAV_ABOUT": "O programie",
		"S_STATE_HINT": "przytrzymaj i mów — tekst trafia tam, gdzie stoi kursor",
		"S_STATE_PROC": "Obróbka tekstu",
		"S_CHANGE_MODEL": "Zmień", "S_PICK_MODEL": "Dobierz", "S_STATE_GET": "Pobierz",
		"S_RETRY": "Spróbuj ponownie", "S_BERR_OPEN": "Otwórz ustawienia serwera",
		"S_STATE_LAST": "Ostatnie dyktowanie", "S_STATE_COPY": "Kopiuj", "S_STATE_MEM": "Pamięć",
		"S_STATE_MEM_SUB": "modele zostają w pamięci, pierwsze zdanie nie czeka",
		"S_HOTMODE":       "Tryb", "S_HOTMODE_HOLD": "przytrzymanie", "S_HOTMODE_TOGGLE": "przełącznik",
		"S_SUB_HOTMODE": "trzymaj klawisze albo naciśnij raz, by zacząć, i raz, by skończyć",
		"S_SILENCE":     "Zatrzymuj po ciszy", "S_SILENCE_MS": "Pauza przed zatrzymaniem, ms",
		"S_SUB_SILENCE": "tryb „przełącznik”: gdy już coś powiesz, pauza kończy nagranie, jakbyś znów nacisnął klawisze",
		"S_SUB_SILENCE_MS": "jak długo możesz milczeć; przed twoimi pierwszymi słowami nagranie nigdy się nie zatrzyma",
		"S_SUB_MINMS":   "pomija przypadkowe naciśnięcia",
		"S_SUB_ENTER":   "od razu wysyła wiadomość",
		"S_SUB_CLIP":    "obrazy i pliki wracają bez zmian",
		"S_SUB_TYPE":    "pomaga tam, gdzie pole nie przyjmuje wklejania",
		"S_SEC_OVERLAY": "Pasek na ekranie",
		"S_MIC_CHECK":   "Sprawdź mikrofon", "S_MIC_CHECK_SUB": "trzy sekundy nagrania i werdykt: poziom, przesterowanie, czy jest mowa", "S_MIC_CHECKING": "Sprawdzam…",
		"S_MCHECK": "Sprawdź zainstalowane modele", "S_MCHECK_SUB": "porównuje pliki modeli z wzorcowymi skrótami", "S_MCHECK_GO": "Sprawdź", "S_MCHECK_RUN": "Sprawdzam…",
		"S_HIST_INSERT": "Wklej",
		"S_MIC": "Mikrofon", "S_MIC_DEFAULT": "Domyślny systemowy", "S_MIC_REFRESH": "Odśwież listę",
		"S_MIC_LEVEL": "Poziom wejścia", "S_MIC_QUIET": "cisza",
		"S_SUB_THREADS": "więcej wątków nie zawsze znaczy szybciej — zmierz na swoim komputerze",
		"S_SEC_LLM":     "Model redaktora",
		"S_PUNCT":       "Znaki i wielkie litery", "S_SUB_PUNCT": "skąd biorą się znaki interpunkcyjne i wielkie litery",
		"S_PUNCT_MODEL": "z modelu", "S_PUNCT_LLM": "od modelu redaktora", "S_PUNCT_OFF": "usuwać",
		"S_SUB_DICT": "Słownik", "S_SUB_PROMPTS": "Prompty",
		"S_SUB_TRTARGET": "na niego tłumaczony jest tekst; dialog na pasku proponuje go pierwszy",
		"S_REMOTE_ABOUT": "Ustawiony jest serwer zdalny: dźwięk trafia do niego, a obietnica powyżej wtedy nie obowiązuje.",
		"S_UPD":          "Aktualizacje", "S_UPD_CHECK": "Sprawdź aktualizacje", "S_UPD_AUTO": "Sprawdzaj przy starcie",
		"S_SUB_UPD":  "jedyne zapytanie sieciowe poza pobieraniem modeli",
		"S_UPD_NONE": "Masz najnowszą wersję", "S_BADGE_MODELS": "Zainstalowane modele", "S_BADGE_MISS": "Model nie został pobrany", "S_BADGE_SYSTEM": "Ostrzeżenia wymagają uwagi", "S_BADGE_HIST": "Wpisy w historii", "S_LOG_OPEN": "Otwórz dziennik", "S_LOG": "Dziennik", "S_LOG_SUB": "wszystko, co aplikacja o sobie zapisuje", "S_UPD_AVAIL": "Dostępna jest wersja %s.",
		"S_UPD_GO": "Aktualizuj", "S_UPD_ERR": "Sprawdzenie nie powiodło się", "S_UPD_DL": "Pobieranie aktualizacji…",
		"S_SEC_SERVICE": "Usługa", "S_SUB_AUTOSTART": "wyłącz, jeśli sam uruchamiasz serwer",
		"S_SUB_PORT":    "rozpoznawanie samo się przeładuje",
		"S_MODEL_READY": "Model pobrany — wybierz go, żeby przełączyć",
		"S_FIT_OK":      "mieści się", "S_FIT_WARN": "na styk", "S_FIT_BAD": "za mało pamięci", "S_RAM": "Pamięć komputera:",
		"S_HF_PH":       "Nazwa modelu — np. qwen2.5 instruct",
		"S_NO_LLM":      "Nie ma jeszcze żadnego modelu — znajdź go w polu wyszukiwania poniżej.",
		"S_NO_LLM_PROF": "Prompty staną się dostępne, gdy będzie zainstalowany model — blok „Model” u góry tej karty.",
		"S_UPDATED":     "Ostatnia aktualizacja modelu", "S_PROF_EDIT": "Edytuj", "S_PROF_CLOSE": "Zwiń",
		"S_CONFIRM_DEL": "Usunąć model „%s”? Będzie można pobrać go ponownie.", "S_FREE": "wolne",
		"S_DEL_ACTIVE":     "Usunąć aktywny model „%s”? Rozpoznawanie zatrzyma się, dopóki nie wybierzesz innego — pobrać go można tutaj.",
		"S_REMOTE_WARN":    "Dźwięk będzie wysyłany na ten serwer. Tryb lokalny jest wyłączony.",
		"S_REMOTE_ASK":     "Dźwięk przestanie być przetwarzany na tym komputerze i będzie wysyłany na %s. Włączyć tryb zdalny?",
		"S_REMOTE_BADGE":   "ZDALNY",
		"S_OK":             "Tak", "S_CANCEL": "Anuluj", "S_DL_START": "Pobierz", "S_DL_CANCEL": "Przerwij pobieranie",
		"S_DL_ASK":    "Model „%s” nie jest pobrany (%s). Zacząć pobieranie?",
		"S_NOT_FOUND": "nic",   "S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autor i deweloper {app} — lokalnego narzędzia do dyktowania dla Windows: głos zamienia się w tekst dokładnie tam, gdzie stoi kursor, bez chmury i bez abonamentu.</p>" +
			"<p>Projekt jest otwarty: kod, budowanie i najnowsze wydania są na GitHubie.</p>" +
			"<ul>" +
			"<li>Repozytorium: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Profil autora: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Znalazłeś błąd albo masz pomysł — załóż zgłoszenie w repozytorium.</p>",
		"S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Co to jest</p>" +
			"<p>Mówisz, a słowa pojawiają się tam, gdzie miga kursor. Tak samo jak przy pisaniu na klawiaturze, tylko szybciej.</p>" +
			"<p class=\"hdim\">Wygląda to tak. Klikasz w pole, w którym zwykle piszesz: w e-mailu, na czacie, w wyszukiwarce, w dokumencie. Przytrzymujesz dwa klawisze i mówisz zdanie zwykłym głosem. Puszczasz. Po sekundzie zdanie jest w polu. Nie trzeba się nigdzie przełączać ani niczego kopiować: aplikacja działa w każdym programie i czeka w tle na swoje klawisze.</p>" +
			"<p class=\"hdim\">Całą pracę wykonuje twój komputer. Głos nie jest zapisywany: żyje kilka sekund w pamięci, zamienia się w tekst i znika. Aplikacja łączy się z internetem tylko w dwóch przypadkach: gdy pobierasz model rozpoznawania i gdy sprawdza, czy wyszła nowa wersja. Sprawdzanie można wyłączyć, a wtedy sieć nie jest jej w ogóle potrzebna.</p>" +
			"<p class=\"hdim\">Słowo „model” będzie się dalej często pojawiać. To po prostu plik, który umie rozpoznawać mowę. Pobierasz go raz, a potem leży w folderze obok aplikacji; im jest większy, tym dokładniej słyszy i tym wolniej działa.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Co się dzieje, gdy puścisz klawisze</p>" +
			"<div class=\"hflow\">" +
			"<span>nagrywanie</span><em>→</em><span>rozpoznawanie</span><em>→</em><span>zamiany</span><em>→</em><span>komendy</span><em>→</em><span>tłumaczenie</span><em>→</em><span>edycja</span><em>→</em><span>wstawianie</span>" +
			"</div>" +
			"<p>Zwykle trwa to sekundę lub dwie i nie trzeba znać poszczególnych kroków. Ale jeśli tekst kiedyś wyjdzie inny, niż się spodziewasz, ten łańcuch łatwo pokaże, który krok go zmienił.</p>" +
			"<ul>" +
			"<li><b>Nagrywanie.</b> Dopóki klawisze są wciśnięte, aplikacja słucha mikrofonu. Zbyt krótkie naciśnięcie pomija, żeby przypadkiem trącony klawisz nie zaczynał dyktowania.</li>" +
			"<li><b>Rozpoznawanie.</b> Głos zamienia się w litery. Zajmuje się tym model. Dla rosyjskiego zwykle działa GigaAM, dla polskiego Parakeet, dla większości pozostałych języków Nemotron, a Whisper zajmuje się resztą języków i tłumaczeniem; to tylko nazwy i nie musisz ich wybierać ręcznie.</li>" +
			"<li><b>Zamiany.</b> Twoja osobista lista „usłyszane → poprawne”. Jeśli model uparcie pisze „git hub” zamiast „GitHub”, poprawiasz to raz i od tej pory poprawia się samo.</li>" +
			"<li><b>Komendy.</b> Kilka fraz działa jak przyciski. Powiedz „nowy wiersz”, a tekst przejdzie do nowego wiersza, przy czym same te słowa nie trafią do pola.</li>" +
			"<li><b>Tłumaczenie.</b> Jeśli jest włączone, mówisz w jednym języku, a tekst trafia do pola w innym.</li>" +
			"<li><b>Edycja.</b> Krok opcjonalny. Drugi model porządkuje tekst: usuwa „yyy” i „no wiesz”, stawia przecinki, nadaje sformułowaniom służbowy ton.</li>" +
			"<li><b>Wstawianie.</b> Tekst trafia do pola. To, co przed dyktowaniem było w schowku, wraca na swoje miejsce.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Gdy to wszystko trwa, na dole ekranu widać wąski pasek z nazwą bieżącego kroku. Po nim od razu wiadomo, że aplikacja pracuje, a nie zawiesiła się.</p>" +
			"<p class=\"wh\" id=\"s-first\">Pierwsze uruchomienie</p>" +
			"<p>Przy pierwszym starcie aplikacja prowadzi cię przez trzy kroki. To zajmie minutę, a potem możesz dyktować.</p>" +
			"<ul>" +
			"<li><b>W jakim języku będziesz dyktować.</b> Najważniejsze pytanie: na podstawie odpowiedzi aplikacja sama dobierze odpowiedni model i pobierze go w tle, a postęp jest widoczny. Jeśli ten język różni się od języka okna, od razu pojawi się wiersz, w którym możesz zmienić także język interfejsu.</li>" +
			"<li><b>Próba.</b> Skrót klawiszowy (domyślnie Ctrl i klawisz Windows; żeby go zmienić, naciśnij to, co wygodnie trzymać jedną ręką), mikrofon ze wskaźnikiem poziomu na żywo i pole, w którym możesz podyktować zdanie na próbę. Jeśli wskaźnik nie rusza się, gdy mówisz, wybrany jest niewłaściwy mikrofon. Pole ożywa, gdy tylko model jest na miejscu.</li>" +
			"<li><b>Gotowe.</b> Przełącznik „Uruchamiaj z Windows” i kreator jest zakończony.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Każdy krok można pominąć i wrócić do niego później: wszystko to jest w ustawieniach. Model pobiera się dalej w tle, nawet jeśli przejdziesz dalej albo zakończysz kreator, zanim pobieranie się skończy. Jeśli zechcesz przejść kreator jeszcze raz, uruchom aplikację z parametrem <code>-wizard</code>. Po aktualizacji kreator już się nie pojawia.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">Pasek na dole ekranu</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Mów…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">tak aplikacja pokazuje, że cię słyszy</div>" +
			"</div>" +
			"<p>Podczas dyktowania to jedyna rzecz widoczna na ekranie. Pasek jest mały, nie przechwytuje myszy i nie przeszkadza w pracy: po prostu mówi, co się właśnie dzieje z twoim zdaniem.</p>" +
			"<ul>" +
			"<li><b>„Mów…”</b> Trwa nagrywanie. Słupki obok podskakują w rytm głosu. Jeśli leżą nieruchomo, nie słychać cię i nie ma sensu mówić dalej.</li>" +
			"<li><b>„Rozpoznawanie…”</b> Klawisze puszczone, pracuje model. Zwykle sekunda lub dwie. Pierwsze dyktowanie po uruchomieniu trwa dłużej: model jest wtedy wczytywany do pamięci.</li>" +
			"<li><b>„Tłumaczenie…”</b> Trwa tłumaczenie.</li>" +
			"<li><b>„Edycja: Cleanup (1/2)”</b> Tekst poprawia drugi model. Widać, które polecenie jest wykonywane i ile ich jest w sumie.</li>" +
			"<li><b>„Wstawiono: 128 znaków”</b> Gotowe, tekst jest już w polu. Pasek sam zgaśnie po kilku sekundach.</li>" +
			"<li><b>Krzyżyk.</b> Rozmyślisz się? Kliknij go na dowolnym etapie. Nagranie zostaje wyrzucone, do pola nic nie trafi.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Pasek można wyłączyć, ustawić w dowolnym rogu, przypiąć do kursora albo po prostu przeciągnąć myszą tam, gdzie nie przeszkadza. Wszystko to jest na zakładce „Sterowanie i zachowanie”.</p>" +
			"<p class=\"wh\" id=\"s-state\">Ekran „Stan”</p>" +
			"<p>Ten ekran odpowiada na jedno pytanie: czy można teraz dyktować. Jeśli nie, mówi dlaczego i daje przycisk, który prowadzi tam, gdzie to się naprawia.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Gotowe do dyktowania</b><span>przytrzymaj i mów — tekst trafia tam, gdzie stoi kursor</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">zielony kwadrat oznacza, że aplikacja nasłuchuje swoich klawiszy</div>" +
			"</div>" +
			"<p>Górna karta bywa trzech rodzajów. <b>Zielona</b>: wszystko w porządku, obok przypomnienie o klawiszach. <b>Żółta</b>: aplikację wyłączono przez ikonę w zasobniku obok zegara i klawisze nie działają, dopóki nie naciśniesz „Włącz” bezpośrednio na tej karcie. <b>Czerwona</b>: ustawienia sobie przeczą. Wybrano rozpoznawanie na innym komputerze, ale nie podano adresu.</p>" +
			"<p>Niżej jest to, na co zwykle tylko się zerka:</p>" +
			"<ul>" +
			"<li><b>Napisy.</b> Karta tuż pod górną: czy napisy są wyłączone, wczytują model, działają czy uruchamiają się ponownie, z modelem i tym, jak przychodzi tekst, ich skrót, jeśli jest ustawiony, przycisk, który je włącza lub wyłącza, i strzałka do zakładki „Napisy”.</li>" +
			"<li><b>W tym tygodniu.</b> Wykres kołowy: do jakich programów trafiały twoje dyktowania i ile znaków podyktowano.</li>" +
			"<li><b>Wersja.</b> Numer kompilacji i przycisk sprawdzania. Jeśli wyszła nowsza, aplikacja powie o tym właśnie tutaj.</li>" +
			"<li><b>Pamięć.</b> Ile jej jest wolnej, który model jest w niej teraz i ile wszystkie modele zajmują na dysku. Przycisk „Wyładuj” zwalnia pamięć od razu, na przykład przed wymagającą grą.</li>" +
			"<li><b>Szybkie ustawienia.</b> Pięć kart z tym, czego używa się najczęściej: mikrofon, język, pasek, dźwięk, start z Windows. Pokazują, jak wszystko jest teraz ustawione, a strzałka prowadzi na zakładkę, na której się to zmienia. Przełączników celowo tu nie ma, żeby niczego nie przestawić mimochodem.</li>" +
			"<li><b>Modele.</b> Co jest zainstalowane, co działa w tej chwili, za które języki odpowiada i co po prostu leży na dysku.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">Zakładka „System”</p>" +
			"<p>Wygląd aplikacji, jej współpraca z Windows i kilka przycisków serwisowych. Na co dzień zagląda się tu raz w roku.</p>" +
			"<p class=\"hsub\">Jak wygląda aplikacja</p>" +
			"<p>Trzy wiersze: język interfejsu, wygląd i kolor. Na działanie aplikacji nie wpływają wcale.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Język interfejsu</td><td>Język okna ustawień i napisów na pasku. Z językiem dyktowania nie ma nic wspólnego: okno może być po polsku, a dyktować będziesz po angielsku.</td></tr>" +
			"<tr><td>Wygląd</td><td>Ogólny styl. Jest ich siedem: „Edytor” (domyślny: ciemnoszary z niebieskim), „Fluent” (jak Ustawienia w Windows 11), „Studio” (grafit z pomarańczowym, prawie bez zaokrągleń), „Miękki” (jasny liliowy, zaokrąglony), „Dokument” (papier z niebieskim), „Neon” (fiolet z gradientem od różu do turkusu) i „Terminal” (zielone litery na czarnym tle, nawiasy wokół przycisków). Każdy z nich używa tej samej czcionki, IBM Plex Mono.</td></tr>" +
			"<tr><td>Kolor</td><td>Odcień tylko dla stylu „Terminal”: zielony, bursztynowy, niebieski, różowy. W tym samym kolorze będą pasek i ikona w zasobniku. Pozostałe style mają własne kolory i ten wiersz się chowa.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Uruchamianie i aktualizacje</p>" +
			"<p>Tutaj decydujesz, czy aplikacja startuje razem z systemem i czy śledzi nowe wersje.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Uruchamiaj z Windows</td><td>Aplikacja wystartuje przy logowaniu do systemu i będzie czekać na twoje klawisze. To zwykły wpis w autostarcie, taki sam, jaki można tam dodać samodzielnie.</td></tr>" +
			"<tr><td>Sprawdzaj przy starcie</td><td>Raz, przy starcie, aplikacja pyta, czy wyszła nowsza wersja. Poza pobieraniem modeli nie łączy się z internetem. Wyłącz to, a o nowych wersjach będziesz się dowiadywać na własną rękę.</td></tr>" +
			"<tr><td>Aktualizacje</td><td>Przycisk sprawdza od razu. Jeśli jest nowsza wersja, aplikacja pobierze ją i zainstaluje zamiast obecnej: ustawienia, modele i historia zostaną na miejscu.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Usługa</p>" +
			"<p>Trzy przyciski na wypadek, gdy coś poszło nie tak.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Dziennik</td><td>Otwiera plik, w którym aplikacja zapisuje, co robiła: który model się uruchomił, ile trwało rozpoznawanie, gdzie wystąpił błąd. Przyczyna jest zwykle opisana słowami, a nie kodami. Plik nie rośnie bez końca: stare wpisy są nadpisywane nowymi.</td></tr>" +
			"<tr><td>Wczytaj config.json ponownie</td><td>Wszystkie ustawienia są w jednym pliku tekstowym obok aplikacji. Jeśli zmieniasz go ręcznie albo przenosisz z innego komputera, ten przycisk wczyta zmiany bez ponownego uruchamiania.</td></tr>" +
			"<tr><td>Przywróć ustawienia</td><td>Ustawia wszystko tak jak po instalacji. Pobrane modele i historia zostają na miejscu: kasowane są tylko ustawienia.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Gdzie działa rozpoznawanie</p>" +
			"<p>Mowę rozpoznaje osobny mały program. Zwykle leży obok aplikacji, uruchamia się sam i nie musisz o nim myśleć. Można go jednak uruchomić na innym komputerze: na przykład na domowym, z mocną kartą graficzną, a dyktować z lekkiego laptopa. Wyboru dokonujesz przełącznikiem na kartach.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>Na tym komputerze</b><span>wątki CPU 4 · autostart włączony · port 8910 · whisper-server.exe obok aplikacji</span></div><button class=\"hbtn\">Edytuj</button></div>" +
			"<div class=\"hcap\">włączona karta jest podświetlona, a druga jest w tym czasie przygaszona</div>" +
			"</div>" +
			"<p>Przycisk „Edytuj” otwiera okno z czterema ustawieniami. Wszystkie mają wartości domyślne i prawie zawsze można ich nie ruszać.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Wątki CPU</td><td>Ile rdzeni procesora oddać rozpoznawaniu. Więcej rdzeni to szybciej, ale tylko do pewnej granicy: dalej nie ma zysku, a inne programy zaczynają zwalniać.</td></tr>" +
			"<tr><td>Uruchamiaj whisper-server automatycznie</td><td>Aplikacja sama uruchamia program rozpoznawania. Wyłączać to warto tylko w jednym przypadku: gdy uruchamiasz go ręcznie i nie chcesz, żeby aplikacja się wtrącała.</td></tr>" +
			"<tr><td>Port</td><td>Numer, pod którym aplikacja rozmawia z programem rozpoznawania wewnątrz twojego własnego komputera. Z zewnątrz nie jest widoczny. Zmieniać go trzeba tylko wtedy, gdy ten numer zajął inny program.</td></tr>" +
			"<tr><td>Ścieżka whisper-server</td><td>Gdzie leży program rozpoznawania. Wypełnia się sama. Własną warto wpisać tylko wtedy, gdy świadomie trzymasz go w innym miejscu: inaczej po przeniesieniu folderu dyktowanie przestanie się uruchamiać.</td></tr>" +
			"</table>" +
			"<p>Druga karta ma jedno ustawienie: adres komputera w sieci, na przykład <code>http://192.168.1.50:8080</code>. Pamiętaj: w tym wariancie nagrany głos trafia na tamtą maszynę, a pod zwykłym adresem idzie nieszyfrowany. W sieci domowej to normalne, ale przez internet lepiej tak nie pracować.</p>" +
			"<div class=\"hbad\">Jeśli wybrano „Na innym komputerze”, ale nie wpisano adresu, a twoim językiem zajmuje się Whisper, dyktowanie nie zadziała. Aplikacja powie o tym czerwoną kartą na ekranie „Stan”, czerwonym znacznikiem przy zakładce „System” i wierszem na dole okna.</div>" +
			"<p class=\"wh\" id=\"s-mic\">Zakładka „Mikrofon”</p>" +
			"<p>Ponad połowa nieudanych dyktowań to nie wina modelu ani ustawień, tylko mikrofonu: aplikacja słucha nie tego, co trzeba, albo prawie cię nie słychać. Na tej zakładce jest wszystko, żeby sprawdzić to w pół minuty.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Mikrofon</td><td>Którego urządzenia słuchać. Opcja „Domyślny systemowy” oznacza „to samo, co wybrano w Windows”. To wygodne, jeśli raz pracujesz w słuchawkach, a raz bez nich: aplikacja sama podąża za systemem.</td></tr>" +
			"<tr><td>Odśwież listę</td><td>Przydaje się, gdy zestaw słuchawkowy podłączono już po uruchomieniu aplikacji: lista zostanie wczytana od nowa.</td></tr>" +
			"<tr><td>Poziom wejścia</td><td>Wskaźnik na żywo. Powiedz coś, a słupki powinny wyraźnie się poruszać. Jeśli stoją, mikrofon jest wyciszony w Windows albo wybrano niewłaściwy.</td></tr>" +
			"<tr><td>Sprawdź mikrofon</td><td>Trzy sekundy nagrania i ocena prostymi słowami: czy jest wystarczająco głośno, czy nie ma przesterowania (gdy dźwięk „trzeszczy”, bo mikrofon jest za blisko) i czy w ogóle słychać mowę.</td></tr>" +
			"<tr><td>Sygnały dźwiękowe</td><td>Krótki dźwięk na początku i na końcu nagrania. Szczególnie przydatne, jeśli pasek jest wyłączony: wtedy przynajmniej słychać, że nagrywanie ruszyło.</td></tr>" +
			"<tr><td>Dźwięk sygnału</td><td>Jaki dokładnie dźwięk: głos Windows, dzwonek, marimba i kilka innych. Przycisk obok odtwarza wybrany, żeby nie zgadywać po nazwie.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">Zakładka „Historia”</p>" +
			"<p>Historia to lista tego, co podyktowano. Tylko tekst: głos nie jest zapisywany nigdzie i nigdy. Domyślnie jest wyłączona, i to celowo: do dyktowań często trafiają sprawy osobiste.</p>" +
			"<p>Warto ją włączyć, jeśli choć raz zdarzyła ci się taka sytuacja: tekst wkleił się do niewłaściwego okna i przepadł; trzeba powtórzyć długie zdanie; chcesz znaleźć to, co było dyktowane wczoraj, i wysłać jeszcze raz.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Jak długo przechowywać</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">dni</span></div>" +
			"<div class=\"hrow\"><label>Ile wpisów przechowywać</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">okres można podać w minutach, godzinach lub dniach, dowolną liczbą</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Przechowuj historię dyktowań</td><td>Główny wyłącznik. Po wyłączeniu nowe wpisy przestają się pojawiać, a stare są kasowane.</td></tr>" +
			"<tr><td>Jak długo przechowywać</td><td>Po jakim czasie wpis sam zniknie. Można ustawić pięć minut, jeśli historia jest potrzebna tylko na wypadek „wkleiło się nie tam”, a można trzydzieści dni.</td></tr>" +
			"<tr><td>Ile wpisów przechowywać</td><td>Drugi limit, według liczby. Działa razem z pierwszym: co nastąpi wcześniej, to zadziała. Przy „7 dni i nie więcej niż 50” pięćdziesiąty pierwszy wpis wypchnie najstarszy, nawet jeśli tydzień jeszcze nie minął.</td></tr>" +
			"<tr><td>Nigdy nie zapisuj z tych programów</td><td>Lista wyjątków. Wszystko, co dyktujesz do tych okien (bank, menedżer haseł, czat służbowy), nie trafi do historii, a samo dyktowanie działa jak zwykle. Program można wybrać z listy otwartych okien, bez szukania nazwy pliku.</td></tr>" +
			"</table>" +
			"<p>Z każdym wpisem można zrobić cztery rzeczy. <b>Kopiuj</b>: tekst trafi do schowka. <b>Wklej</b>: aplikacja poprosi o wskazanie miejsca. Kliknij właściwe pole w dowolnym programie, a tekst trafi właśnie tam. <b>Usuń</b>: usuwa jeden wpis. <b>Szukaj</b>: przeszukuje i tekst, i nazwę programu. Przycisk „Wyczyść” kasuje wszystko naraz, dlatego prosi o potwierdzenie.</p>" +
			"<p class=\"hdim\">Po prawej stronie każdego wpisu jest napisane, kiedy zniknie, na przykład „do 7 wrz”, a w ostatnim dniu ten napis robi się żółty.</p>" +
			"<p class=\"wh\" id=\"s-dict\">Zakładka „Sterowanie i zachowanie”</p>" +
			"<p>Jak uruchamia się dyktowanie i jak tekst trafia do pola. Zagląda się tu, gdy coś działa nie po twojej myśli: klawisze są niewygodne, tekst się nie wstawia, Enter naciska się sam.</p>" +
			"<p class=\"hsub\">Klawisze</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Skrót klawiszowy</td><td>Naciśnij „Zmień…” i przytrzymaj to, co ci wygodne. Aplikacja przechwytuje tę kombinację w całym systemie, więc nie wybieraj zajętej, na przykład Ctrl+C. Najrzadziej zajęte są kombinacje z klawiszem Windows. Język może mieć też własny skrót, przypisywany w zakładce „Języki i modele”.</td></tr>" +
			"<tr><td>Tryb</td><td>W trybie „przytrzymanie” nagrywa, dopóki klawisze są wciśnięte: puszczasz je i zdanie idzie do przetworzenia. W trybie „przełącznik” pierwsze naciśnięcie włącza nagrywanie, a drugie je wyłącza. Drugi tryb jest wygodny przy długich monologach, gdy ręka męczy się od trzymania klawiszy.</td></tr>" +
			"<tr><td>Zatrzymuj po ciszy</td><td>Tylko w trybie „przełącznik”. Gdy już coś powiesz, pauza tak długa jak ustawiona w „Pauza przed zatrzymaniem, ms” sama kończy nagranie, jakbyś znów nacisnął skrót, a pasek pokazuje „Zatrzymano po ciszy”. Przed twoimi pierwszymi słowami nagranie nigdy się nie zatrzyma, więc po naciśnięciu możesz spokojnie zebrać myśli; jeśli nic nie powiesz, i tak zakończy je limit „Maks. nagranie, s”. Gdy włączone jest „Naciśnij Enter po wstawieniu (auto-wysyłka)”, pauza wysyła wiadomość, więc wybierz dłuższą pauzę.</td></tr>" +
			"<tr><td>Min. nagranie, ms</td><td>Jak krótkie naciśnięcie uznać za przypadkowe. Trącisz klawisze na ułamek sekundy, a aplikacja nic nie zrobi, zamiast wstawiać skrawek słowa.</td></tr>" +
			"<tr><td>Maks. nagranie, s</td><td>Zabezpieczenie przed zapomnianym klawiszem: nagranie samo się urwie i nie będzie przez pół godziny nagrywać w próżnię.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Jak wstawiany jest tekst</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Naciśnij Enter po wstawieniu (auto-wysyłka)</td><td>Tekst nie tylko się wstawi, ale od razu zostanie wysłany: wiadomość pójdzie, wyszukiwanie się wykona. W komunikatorze oszczędza to jeden ruch, w e-mailu czy dokumencie prawie na pewno przeszkodzi.</td></tr>" +
			"<tr><td>Przywróć schowek po wstawieniu</td><td>Wstawianie idzie przez schowek, a aplikacja odkłada z powrotem to, co w nim było: obrazek, plik, skopiowany kawałek tekstu. Zostaw włączone, jeśli często coś kopiujesz i nie chcesz tego tracić.</td></tr>" +
			"<tr><td>Wpisywanie znak po znaku</td><td>Zamiast wkleić tekst w jednym kawałku, aplikacja wpisze go tak, jakby ktoś pisał na klawiaturze. To wolniejsze, ale ratuje tam, gdzie wklejanie ze schowka jest zablokowane: w niektórych formularzach bankowych, terminalach, starych programach.</td></tr>" +
			"<tr><td>Opóźnienie przed wstawieniem</td><td>Krótka pauza, zanim tekst zostanie przekazany oknu. Potrzebna rzadko: bywa, że program nie zdąży przyjąć tekstu i początek zdania przepada.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Pasek na ekranie</p>" +
			"<p>Tutaj decydujesz, czy w ogóle go pokazywać i gdzie dokładnie.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Pokazuj pasek</td><td>Wyłącz, jeśli przeszkadza albo wystarczają ci sygnały dźwiękowe.</td></tr>" +
			"<tr><td>Gdzie pokazywać pasek</td><td>Schemat ekranu z dziewięcioma miejscami: kliknij ten róg, w którym pasek nie będzie przeszkadzał.</td></tr>" +
			"<tr><td>Przy kursorze</td><td>Pasek pojawia się obok miejsca, w które dyktujesz. Jeśli program nie podaje, gdzie jest kursor, to obok wskaźnika myszy.</td></tr>" +
			"<tr><td>Własne miejsce</td><td>Pasek można po prostu przeciągnąć myszą, a on zapamięta to miejsce.</td></tr>" +
			"<tr><td>Ekran</td><td>Na którym monitorze go pokazywać, jeśli jest ich kilka.</td></tr>" +
			"<tr><td>Pokazywać rozpoznany tekst</td><td>Po wstawieniu na pasku widać samo zdanie, a nie liczbę znaków. Wygodne do sprawdzania, ale zdanie zobaczy każdy, kto patrzy na twój ekran.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">Zakładka „Języki i modele”</p>" +
			"<p>Tutaj mówisz, w jakim języku będziesz dyktować, a aplikacja decyduje, czym to rozpoznawać. Modeli jest kilka: jedne są szybkie, inne dokładne, jeszcze inne znają rzadkie języki.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Język mowy źródłowej</td><td>Język, w którym mówisz. Jeśli mówisz w kilku, wybierz „Auto”: model sam się zorientuje, choć trochę wolniej i trochę mniej dokładnie, niż gdy język jest podany wprost.</td></tr>" +
			"<tr><td>Który model dla którego języka</td><td>Kliknij język, a rozwinie się lista modeli, które go znają, z oceną szybkości i dokładności oraz rozmiarem pliku. Tak jeden język może zachować szybki model, a inny dostać dokładny. Osobny wiersz „Auto” jest wspólny: jego model obsługuje wszystkie języki, którym nie wybrano własnego, i na liście są one szare. Ten sam model rozpoznaje też, w jakim języku mówisz, gdy język mowy źródłowej nie jest ustawiony.</td></tr>" +
			"<tr><td>Własny skrót</td><td>W rozwiniętym języku jest miejsce na jego własną kombinację klawiszy. Przytrzymaj ją, a dyktowanie idzie w tym języku jego modelem, zamianami i regułami, niezależnie od ustawionego języka mowy źródłowej; główny skrót zachowuje język wybrany wyżej. Tak dyktuje się po angielsku jedną kombinacją, a we własnym języku inną, bez otwierania ustawień. Podczas nagrywania pasek pokazuje kod języka. Kombinacja nie może powtarzać głównej ani kombinacji innego języka.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Rosyjski<span class=\"hs\">GigaAM v3 · dokładność 5 z 5 · szybkość 5 z 5 · 232 MB</span></label><span class=\"hchip\"><i></i>aktywna</span></div>" +
			"<div class=\"hrow\"><label>Angielski<span class=\"hs\">jak automatyczne wykrywanie · Nemotron 3.5</span></label><button class=\"hbtn\">Dobierz</button></div>" +
			"<div class=\"hcap\">widać, który model obsługuje dany język i ile to kosztuje w szybkości i miejscu na dysku</div>" +
			"</div>" +
			"<p class=\"hsub\">Tłumaczenie</p>" +
			"<p>Niektóre modele potrafią nie tylko spisać mowę, ale od razu ją przetłumaczyć: mówisz w jednym języku, a w polu pojawia się tekst w innym. Tłumaczy Whisper: na angielski własnym wbudowanym zadaniem tłumaczenia, na pozostałe języki dzięki temu, że język wyniku podaje mu się jako język mowy; przy odległych parach językowych jakość bywa nierówna. Jeśli wybrany model nie umie tłumaczyć, obok przełącznika pojawi się pomarańczowy wiersz: przy zainstalowanym Whisperze pasek zaproponuje tłumaczenie za jego pomocą, a w przeciwnym razie tekst zostanie w języku mowy źródłowej.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Zmień język wyjściowego tekstu</td><td>Włącza tłumaczenie.</td></tr>" +
			"<tr><td>Domyślny język wyjściowego tekstu</td><td>Na jaki język tłumaczyć, jeśli nie wybierzesz innego w trakcie dyktowania.</td></tr>" +
			"<tr><td>Pytaj o język wyjściowego tekstu</td><td>Można nie pytać i zawsze tłumaczyć na język domyślny. Można pytać za każdym razem: po dyktowaniu na pasku pojawią się przyciski z językami. Można też pytać z limitem czasu: te same przyciski, ale jeśli nie zdążysz nacisnąć, zostanie użyty język domyślny.</td></tr>" +
			"<tr><td>Języki w oknie dialogowym</td><td>Które języki pokazywać jako te przyciski. Zostaw dwa lub trzy, których naprawdę używasz, inaczej wybieranie potrwa dłużej niż mówienie.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Modele na dysku</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Własny model</td><td>Jeśli masz własny model, umieść go w folderze modeli i uruchom aplikację ponownie: pojawi się na listach dla języków, które rozumie. Przycisk obok otwiera właściwy folder, żeby nie szukać go w Eksploratorze plików.</td></tr>" +
			"<tr><td>Wyładuj z pamięci</td><td>Gdy aplikacja działa, model zajmuje pamięć RAM. Ten przycisk zwalnia ją od razu, na przykład przed wymagającą grą. Następne dyktowanie wczyta model ponownie, a pierwsze zdanie potrwa dłużej niż zwykle.</td></tr>" +
			"<tr><td>Sprawdź zainstalowane modele</td><td>Porównuje pobrane pliki ze wzorcem. Przydaje się, jeśli połączenie zrywało się podczas pobierania i podejrzewasz, że plik się uszkodził.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">Zakładka „Reguły”</p>" +
			"<p>Trzy listy, które poprawiają tekst automatycznie, jeszcze zanim trafi do okna. Wszystkie trzy są opcjonalne: bez nich aplikacja świetnie działa.</p>" +
			"<p class=\"hsub\">Znaki i wielkie litery</p>" +
			"<p>Skąd w tekście biorą się kropki, przecinki i wielkie litery. Są trzy warianty: <b>z modelu</b> (model stawia je tak, jak usłyszał; szybko i zwykle wystarcza), <b>od modelu redaktora</b> (znaki dodaje drugi model; staranniej, ale każde dyktowanie wydłuża się o kilka sekund) i <b>usuwać</b> (tekst jednym ciągiem bez znaków; potrzebne rzadko, na przykład do pola wyszukiwania).</p>" +
			"<p class=\"hsub\">Słownik</p>" +
			"<p>Lista słów, których model zwykle nie zna: nazwiska, nazwy firm, terminy branżowe. To nie rozkaz „pisz tak”, tylko podpowiedź: gdy usłyszy coś podobnego, model chętniej wybierze twoją wersję.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">słowa dodaje się po jednym i każde ma własną kartę</div>" +
			"</div>" +
			"<p class=\"hdim\">Słownik czytają tylko modele Whisper. Nemotron, Parakeet, GigaAM i pozostałe modele w ogóle nie przyjmują podpowiedzi: dla nich narzędziem jest lista zamian poniżej. Domyślny zestaw słów zmienia się razem z językiem, dopóki nie wpiszesz własnego.</p>" +
			"<p class=\"hsub\">Zamiany po rozpoznaniu</p>" +
			"<p>To też lista słów, ale działa inaczej: to już nie podpowiedź, tylko gwarancja. To, co pojawi się po lewej, zmieni się w to, co po prawej. Każda zamiana ma trzy przełączniki: <b>Tylko całe słowa</b> (wtedy „kod” nie ruszy „kodowania”, i zwykle o to właśnie chodzi), <b>Uwzględniaj wielkość liter</b> (wtedy „Róża” i „róża” to różne słowa) i <b>Język reguły</b> (zamiana zadziała tylko wtedy, gdy dyktujesz w tym języku).</p>" +
			"<p class=\"hsub\">Komendy głosowe</p>" +
			"<p>Frazy, które nie trafiają do tekstu, tylko wykonują działanie. Złamanie wiersza, nowy akapit i własny tekst są rozpoznawane jako całe słowa w dowolnym miejscu zdania. Anulowanie i komendy formatowania działają tylko na końcu dyktowania, a usuwanie tylko wtedy, gdy powiesz je osobno. Komenda formatowania powiedziana osobno dotyczy następnego dyktowania, a pasek pokazuje „Następna fraza: …”.</p>" +
			"<ul>" +
			"<li><b>„nowy wiersz”</b>: przejść do nowego wiersza;</li>" +
			"<li><b>„nowy akapit”</b>: zostawić pustą linię między akapitami;</li>" +
			"<li><b>„anuluj”</b>: wyrzucić zdanie i niczego nie wstawiać;</li>" +
			"<li><b>„usuń to”</b>: usunąć to, co program wstawił ostatnio;</li>" +
			"<li><b>„jako nagłówek”</b>: postawić przed zdaniem znak, domyślnie #, a po nim pustą linię;</li>" +
			"<li><b>„jako lista”</b>: podzielić zdanie na przecinkach, średnikach i słowie „i” na wiersze zaczynające się od myślnika;</li>" +
			"<li><b>„wielką literą”</b>: zacząć zdanie wielką literą;</li>" +
			"<li><b>„w cudzysłowie”</b>: ująć zdanie w cudzysłów;</li>" +
			"<li><b>„bez kropki”</b>: usunąć kropkę na końcu;</li>" +
			"<li><b>własna komenda</b>: wstawić gotowy fragment, na przykład podpis, adres albo standardową odpowiedź.</li>" +
			"</ul>" +
			"<p>Usuwanie zabiera tylko to, co sam program wstawił ostatnio, i tylko wtedy, gdy od tamtej pory nic nie wpisano ani nie kliknięto, okno jest to samo, tekst nie miał złamań wiersza, nie został wysłany Enterem i minęło najwyżej pięć minut. Tekstu, który już był w oknie, program nigdy nie rusza.</p>" +
			"<p class=\"wh\" id=\"s-post\">Zakładka „Obróbka końcowa”</p>" +
			"<p>Obróbka końcowa to drugi model, który przepisuje rozpoznany tekst według twojego polecenia. Na przykład: „usuń słowa-wypełniacze i powtórzenia”, „zrób z tego oficjalny e-mail”, „zamień to na listę”. Polecenie pisze się zwykłymi słowami i nazywa się promptem.</p>" +
			"<p>Domyślnie jest wyłączona, i oto dlaczego: to jedyny krok, który wyraźnie spowalnia dyktowanie. Zamiast sekundy czy dwóch wstawianie zajmie kilka sekund, a na słabym komputerze nawet dziesięć.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Włącz przetwarzanie końcowe</td><td>Główny wyłącznik. Dopóki jest wyłączony, reszta zakładki nie działa, a tekst jest wstawiany tak, jak go rozpoznano.</td></tr>" +
			"<tr><td>Model · Lokalny</td><td>Redaktor, który działa na twoim komputerze. Nic nigdzie nie wychodzi, ale potrzebna jest pamięć: takie modele ważą od gigabajta wzwyż. Obok jest katalog, w którym model wybiera się i pobiera jednym kliknięciem.</td></tr>" +
			"<tr><td>Model · Serwer zewnętrzny</td><td>Zamiast własnego modelu tekst jest wysyłany do zewnętrznej usługi: OpenAI, Groq albo twojego własnego serwera. Szybciej i mądrzej, ale tekst opuszcza komputer. Klucz dostępu jest przechowywany w postaci zaszyfrowanej, a przycisk „Test połączenia” sprawdza, czy adres i klucz działają.</td></tr>" +
			"<tr><td>Prompty</td><td>Twoje polecenia dla redaktora. Zaznaczone są stosowane po kolei, z góry na dół: najpierw, powiedzmy, czyszczenie tekstu, potem styl służbowy. Jeśli nic nie jest zaznaczone, obróbka końcowa nic nie robi, a aplikacja ostrzeże o tym żółtym wierszem.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Głos nigdy nie trafia na serwer zewnętrzny, tylko rozpoznany tekst. Ale ten tekst jednak wychodzi, więc włączaj to świadomie.</div>" +
			"<p>Każdy prompt ma <b>nazwę</b> (pod tą nazwą widać go na pasku, gdy działa), <b>treść polecenia</b> zwykłymi słowami, opcjonalny <b>własny skrót klawiszowy</b> (dyktowanie z nim przepuści zdanie tylko przez ten prompt, z pominięciem pozostałych) i <b>test</b>: pole, w którym możesz wpisać zdanie i zobaczyć, co wyjdzie, niczego nie dyktując.</p>" +
			"<p class=\"wh\" id=\"s-subs\">Zakładka „Napisy”</p>" +
			"<p>Napisy zapisują to, co odtwarza komputer — rozmowę, film, wykład — i pokazują słowa w małym oknie, które zostaje nad innymi oknami. To samo rozpoznawanie co przy dyktowaniu, tylko dźwięk pochodzi z głośników, a nie z mikrofonu, a dyktowanie działa dalej, gdy napisy są włączone. Własne dźwięki aplikacji nie są zapisywane.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Włącz napisy</td><td>Włącza je i wyłącza. Wiersz obok mówi, co się dzieje: wyłączone, wczytywanie modelu, nasłuch albo ponowne uruchamianie rozpoznawania. To samo robią menu w zasobniku, własny skrót klawiszowy i krzyżyk na oknie.</td></tr>" +
			"<tr><td>Język mowy</td><td>Język tego, czego słuchasz, osobno od języka dyktowania. Model to ten, który zakładka „Języki i modele” daje temu językowi. Model strumieniowy pisze słowa w trakcie mówienia, pozostałe każde zdanie po jego końcu. Model Whisper na procesorze potrzebuje wielu sekund na zdanie, więc jego tekst mocno spóźnia się względem mowy; zaległe zdania trafiają do niego razem i wychodzą jako jeden wiersz.</td></tr>" +
			"<tr><td>Tłumaczenie</td><td>Pokazuje w oknie tłumaczenie zamiast rozpoznanego tekstu, w wybranym języku; tytuł okna brzmi wtedy „Tłumaczenie” z językami, na przykład EN → RU. Mowa jest cięta na frazy na bieżąco, a każda fraza od razu idzie do tłumaczenia. Fraza kończy się na końcu zdania albo przed słowem łączącym, a tłumacz dostaje razem z nią poprzednie frazy i własne ich tłumaczenia, więc zdanie pocięte na części jest tłumaczone spójnie. Tłumaczenie pojawia się równo, całymi słowami, w tempie czytania, a szybciej, gdy dużo się nazbierało, więc duży kawałek nie szarpie tekstem. Frazy jednego odcinka mowy płyną jednym akapitem, a po pauzie zaczyna się nowy wiersz. Fraza, której nie udało się przetłumaczyć, jest pokazywana tak, jak została usłyszana, bledszym kolorem. „Pokazuj też tekst oryginalny” dzieli okno cienką linią: u góry bledszym kolorem to, co powiedziano, u dołu tłumaczenie; linię można przeciągać myszą. Wiersz „Kto tłumaczy” decyduje, kto to robi. „Dokładny” to model tłumaczenia na karcie graficznej: do pobrania 2,3 GB, dla karty graficznej z co najmniej 4 GB własnej pamięci; ładuje się na kartę, gdy napisy startują, i schodzi z niej, gdy się kończą, a link obok otwiera warunki korzystania z niego. „Twój własny serwer” to serwer obróbki końcowej ustawiony na zakładce „Obróbka końcowa” — wtedy każda fraza trafia na ten serwer. „Wybiera program” bierze model tłumaczenia, jeśli jest pobrany, a karta graficzna się nadaje, inaczej twój serwer, jeśli jest ustawiony, inaczej model redaktor zainstalowany na zakładce „Obróbka końcowa”, który tłumaczy z grubsza: sens jest zrozumiały, ale słowa bywają nietrafione. Tłumacza wybranego ręcznie program nigdy nie zamienia na innego: jeśli nie może pracować, blok mówi, czego brakuje, a tekst jest pokazywany tak, jak został usłyszany. Tłumaczenie jest zapisywane w osobnym pliku obok transkrypcji: każda fraza z godziną, o której się zaczęła; fraza, której nie udało się przetłumaczyć, jest zapisywana tak, jak została usłyszana. Bez tłumacza nie ma tłumaczenia, a napisy działają jak dotąd.</td></tr>" +
			"<tr><td>Skrót klawiszowy</td><td>Włącza i wyłącza napisy z każdego okna. Nigdy nie zaczyna dyktowania i nie może powtarzać innego skrótu.</td></tr>" +
			"<tr><td>Zapisuj transkrypcję</td><td>Każde uruchomienie zapisuje plik tekstowy w folderze transcripts obok programu: pierwszy wiersz z datą, językiem i modelem, potem każdy gotowy wiersz z godziną, o której się zaczął. Dopóki pierwszy wiersz nie jest gotowy, nic nie jest zapisywane. Przy włączonym tłumaczeniu nazwa transkrypcji zawiera original i kod języka mowy, a obok pojawia się drugi plik z translation i kodem języka docelowego, na przykład 2026-11-12_10-18-29_original_EN.txt i 2026-11-12_10-18-29_translation_RU.txt.</td></tr>" +
			"<tr><td>Okno napisów</td><td>Gdzie stoi, jakie jest duże i jak duże są litery. Miejsce wybiera się na obrazku ekranu, tak jak dla paska dyktowania: kliknięcie jednego z dziewięciu miejsc stawia tam okno, a z „Własnym miejscem” może stać gdziekolwiek, gdy przeciągniesz znacznik na obrazku albo samo okno; przy kilku monitorach można wybrać ekran. Rozmiar zmienia się myszą za krawędzie i rogi okna; wysokość staje na całej liczbie wierszy, a szerokość i liczba wierszy są zapamiętywane dla każdej rozdzielczości ekranu. „Wiersze tekstu” obowiązują tam, gdzie rozmiaru jeszcze nie ustalono, a „Przywróć rozmiar” zapomina ustalone. Kółko myszy nad oknem przewija tekst wstecz o trzy wiersze na ząbek; gdy czytasz starszy tekst, nowy go nie przesuwa, a znacznik „do najnowszego ↓” w tytule wraca do bieżącej mowy.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Twoi rozmówcy powinni wiedzieć, że ich słowa są zapisywane.</div>" +
			"<p class=\"wh\" id=\"s-about\">„O programie” i „Kontakty”</p>" +
			"<p>Na zakładce „O programie” jest krótki opis tego, czym jest aplikacja i jak jest zbudowana, oraz lista wszystkich zewnętrznych bibliotek, z których się składa, wraz z licencjami. Nazwy są klikalne i prowadzą do stron tych projektów. Jest tam też link do repozytorium, w którym leżą kod źródłowy i wszystkie wersje.</p>" +
			"<p>„Kontakty” to trzy karty: e-mail, repozytorium i strona na błędy i pomysły. Kliknięcie adresu e-mail otworzy twój program pocztowy z wypełnionym już odbiorcą, a ikona obok po prostu skopiuje adres.</p>" +
			"<p class=\"wh\" id=\"s-files\">Co jest w folderze aplikacji</p>" +
			"<p>Aplikacja jest przenośna: wszystko, co do niej należy, leży w jednym folderze i nie trafia nigdzie indziej.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Wszystkie twoje ustawienia w jednym pliku tekstowym. Można go skopiować na inny komputer, żeby nie konfigurować wszystkiego od nowa.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>Dziennik. Pierwsze miejsce, do którego warto zajrzeć, jeśli aplikacja zachowuje się dziwnie.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>Historia dyktowań, jeśli jest włączona. Zwykły tekst: w razie potrzeby plik można po prostu usunąć.</td></tr>" +
			"<tr><td>models</td><td>Folder z pobranymi modelami, najcięższa część, kilka gigabajtów. Jego też można przenieść na inny komputer, żeby nie pobierać wszystkiego ponownie.</td></tr>" +
			"<tr><td>transcripts</td><td>Transkrypcje napisów, jeden plik tekstowy na uruchomienie, oraz ich tłumaczenia, gdy tłumaczenie jest włączone. To twoje dokumenty: deinstalator usuwa je tylko wtedy, gdy poprosisz także o usunięcie ustawień i danych.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>To właśnie te małe programy, które uruchamiają modele. Bez nich model jest tylko plikiem, którego nie ma czym odczytać.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Co trafia do internetu, a co nie</p>" +
			"<ul>" +
			"<li><b>Głos</b> nie jest zapisywany ani wysyłany, chyba że świadomie włączysz rozpoznawanie na innym komputerze.</li>" +
			"<li><b>Napisy</b> są rozpoznawane tutaj, tak samo jak dyktowanie, a transkrypcja zostaje w folderze transcripts; przy rozpoznawaniu na innym komputerze zdania trafiają na ten serwer; przy tłumaczeniu przez twój serwer obróbki końcowej każdy gotowy wiersz również tam trafia.</li>" +
			"<li><b>Tekst</b> nie opuszcza komputera, chyba że świadomie włączysz zewnętrzny serwer obróbki końcowej.</li>" +
			"<li><b>Historia</b> leży w pliku obok aplikacji i jest czyszczona według ustawionego przez ciebie okresu.</li>" +
			"<li><b>Sieć</b> jest potrzebna dwa razy: żeby pobrać model i sprawdzić aktualizację. Sprawdzanie aktualizacji można wyłączyć.</li>" +
			"<li><b>Klucz</b> do zewnętrznej usługi jest szyfrowany mechanizmami Windows i nigdy nie trafia do dziennika.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Jeśli coś poszło nie tak</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Tekst nie pojawia się w polu</td><td>Program nie przyjmuje wklejania ze schowka. Włącz „Wpisywanie znak po znaku” na zakładce „Sterowanie i zachowanie” albo dodaj krótkie opóźnienie przed wstawieniem.</td></tr>" +
			"<tr><td>Dyktowanie się nie uruchamia</td><td>Skrót klawiszowy jest zajęty przez inny program. Ustaw inny, najlepiej z klawiszem Windows, bo takie są zajęte rzadziej.</td></tr>" +
			"<tr><td>Nagrała się cisza</td><td>Mikrofon jest niewłaściwy albo wyciszony. Otwórz zakładkę „Mikrofon”, spójrz na wskaźnik poziomu i naciśnij „Sprawdź mikrofon”.</td></tr>" +
			"<tr><td>Pierwsze zdanie trwa bardzo długo</td><td>Tak ma być: model wczytuje się do pamięci przy pierwszym użyciu. Potem rozpoznawanie zajmuje sekundy.</td></tr>" +
			"<tr><td>Widać napis „Rozpoznawanie niedostępne”</td><td>Wybrano rozpoznawanie na innym komputerze, ale nie podano adresu. Wpisz adres albo wróć do wariantu „Na tym komputerze”.</td></tr>" +
			"<tr><td>Obróbka końcowa niczego nie zmienia</td><td>Nie zaznaczono żadnego promptu, nie wybrano modelu albo serwer zewnętrzny nie odpowiada. Przyczyna jest napisana na samej zakładce, obok wyłącznika.</td></tr>" +
			"<tr><td>Tekst jest dziwny, słowa są zamienione</td><td>Zajrzyj na zakładkę „Reguły”: najpewniej zadziałała zamiana albo komenda głosowa.</td></tr>" +
			"<tr><td>Nic z powyższego</td><td>„System” → „Otwórz dziennik”. Aplikacja zapisuje tam wszystko o sobie, a przyczyna jest zwykle opisana słowami.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">Nie ma tu odpowiedzi? Napisz e-mail z zakładki „Kontakty” albo załóż issue w repozytorium. Jedno i drugie czyta autor.</p>",
		"S_WIZ_HELLO":      "Dwie minuty na konfigurację",
		"S_WIZ_HELLO_TEXT": "{app} zamienia głos w tekst dokładnie w miejscu kursora: przytrzymaj skrót, powiedz zdanie, puść — tekst jest na miejscu. Wszystko liczy się na twoim komputerze, dźwięk nigdzie nie wychodzi.",
		"S_WIZ_LATER":      "Wszystko, co teraz wybierzemy, można później zmienić w ustawieniach.",
		"S_WIZ_MODEL_TEXT": "Model dla twojego języka zostanie dobrany automatycznie i pobierze się w tle.",
		"S_WIZ_PICKED":     "Dobrany dla języka: %s",
		"S_WIZ_DL_BG":      "Pobiera się w tle, możesz iść dalej",
		"S_WIZ_DL_FAIL":    "Pobieranie zostało przerwane: %s",
		"S_WIZ_RETRY":      "Spróbuj ponownie",
		"S_WIZ_TRY_DL":     "Model wciąż się pobiera (%d%%). Pole zacznie działać, gdy tylko będzie gotowy.",
		"S_WIZ_TRY_LOAD":   "Wczytuję model…",
		"S_WIZ_TRY_LIVE":   "Z modelem %s słowa pojawiają się na pasku, gdy jeszcze mówisz.",
		"S_WIZ_DONE_DL":    "Model wciąż się pobiera. Dyktowanie zacznie działać samo, gdy tylko model będzie na miejscu.",
		"S_WIZ_INPUT_TEXT": "Ten skrót będziesz trzymać podczas mówienia. Powiedz coś i sprawdź, czy pasek poziomu się rusza.",
		"S_WIZ_T_TRY":      "Próba",
		"S_WIZ_TRY_PH":     "tutaj pojawi się tekst",
		"S_WIZ_T_DONE":     "Gotowe",
		"S_WIZ_DONE_TEXT":  "{app} mieszka w zasobniku: lewy klik w ikonę — ustawienia, prawy — menu. Dyktować możesz w każdym oknie, w którym jest kursor tekstu.",
		"S_AUTORUN":        "Uruchamiaj z Windows",
		"S_AUTORUN_SUB":    "Wpis w autostarcie bieżącego użytkownika",
		"S_WIZ_SKIP":       "Pomiń",
		"S_WIZ_BACK":       "Wstecz",
		"S_WIZ_NEXT":       "Dalej",
		"S_WIZ_FINISH":     "Zakończ",
		"S_WIZ_WAIT":       "Czekam na pierwsze zdanie…",
		"S_WIZ_HEARD":      "Usłyszałem:",
		"S_WIZ_HAVE":       "Wszystko, co potrzebne, jest już pobrane",
		"S_WIZ_TRY_TEXT":   "Ustaw kursor w polu poniżej, przytrzymaj %s, powiedz zdanie i puść.",
	}
	msgs["uk"] = map[string]string{
		"menu.sub.start": "Запустити субтитри",
		"menu.sub.stop": "Зупинити субтитри",
		"sub.title": "Субтитри",
		"sub.title.tr": "Переклад",
		"sub.follow": "до останнього ↓",
		"sub.loading": "Завантажую модель…",
		"sub.listening": "Слухаю звук комп'ютера…",
		"sub.restarting": "Розпізнавання перезапускається…",
		"sub.err.nomodel": "Субтитри: для їхньої мови не встановлено модель",
		"sub.err.engine": "Субтитри: рушій розпізнавання не запустився",
		"sub.err.device": "Субтитри: немає пристрою відтворення, який можна слухати",
		"sub.err.lost": "Субтитри зупинено: рушій розпізнавання раз у раз падав",
		"app.name": "{app}", "already.running": "Застосунок уже запущено (іконка в треї).",
		"err.title": "{app} — помилка", "cfg.err.title": "{app} — помилка конфігурації",
		"err.details": "\n\nДеталі у {log}", "err.hook": "Не вдалося встановити хук клавіатури: %s",
		"err.mic": "Мікрофон: %s", "err.hotkey.cfg": "Сполучення клавіш у config.json: %s",
		"err.model.notfound": "Файл моделі не знайдено: %s\nПерезберіть проєкт (.\\build.ps1) або виправте \"model\" у config.json",
		"err.server.repeat":  "whisper-server постійно аварійно завершується — див. {log}",
		"err.server.dead":    "Сервер розпізнавання %s не відповідає (server_url/server_autostart у config.json)",
		"err.server.timeout": "whisper-server не відповів протягом %s",
		"err.server.start":   "whisper-server завершився під час запуску (див. лог)",
		"err.srv.noaddr":    "віддалений сервер розпізнавання не налаштовано — задайте адресу в налаштуваннях",
		"err.webview":        "Для вікна налаштувань потрібен Microsoft WebView2 Runtime (входить до Windows 11).\nЗараз відкриється сторінка завантаження — встановіть його та відкрийте налаштування знову.",
		"status.loading":     "Завантаження моделі…", "status.ready": "Готово — утримуйте %s і говоріть",
		"status.recording": "Йде запис…", "status.transcribing": "Розпізнаю…", "status.disabled": "Вимкнено",
		"status.server.restart": "Сервер розпізнавання впав, перезапускаю…", "status.cfg.err": "Помилка в config.json (див. лог)",
		"status.nomodel": "Модель розпізнавання не завантажена — оберіть її в налаштуваннях",
		"tf.title": "%s не перекладає. Перекласти?",
		"tf.yes": "Так, через %s",
		"tf.no": "Як є",
		"hk.lone": "%s — одиночна клавіша: доки вона призначена, набрати її звичайним способом не вийде",
		"state.loaded.none": "нічого не завантажено",
		"state.week": "%d диктувань · %d знаків",
		"snd.ok": "рівень у нормі",
		"snd.quiet": "тихо — говоріть ближче до мікрофона",
		"snd.clipped": "занадто гучно, звук обрізався",
		"snd.silent": "тиша в записі",
		"status.parked": "Рушій вивантажено — натисніть сполучення, щоб розбудити",
		"status.nomodel.lang": "Для мови %s не встановлено модель %s — відкрийте «Мови і моделі»",
		"menu.settings":  "Налаштування…", "menu.enable": "Увімкнути", "menu.disable": "Вимкнути",
		"menu.open.config": "Відкрити config.json", "menu.open.log": "Відкрити лог",
		"menu.about": "Про застосунок", "menu.quit": "Вихід",
		"menu.lastcopy": "Копіювати останній результат",
		"ov.copied":     "Скопійовано в буфер обміну", "ov.kept": "Скасовано — текст в «Останньому результаті»",
		"ov.llm.skipped": "Вставлено без профілю «%s»",
		"fd.title":       "Фокус змінився — вставити?", "fd.here": "Вставити сюди", "fd.copy": "Копіювати",
		"fd.keep":  "Залишити",
		"err.port.busy": "Порт %d зайнятий іншою програмою. Змініть порт у налаштуваннях.",
		"ov.speak": "Говоріть…", "ov.transcribing": "Розпізнаю", "ov.asking": "Чекаю на відповідь", "ov.inserted": "Вставлено: %d символів", "ov.left": "лишилося %d с", "ov.esc": "1…9 · Enter · Esc — скасувати", "err.net.dns": "Немає зв'язку з %s — перевірте інтернет", "err.net.timeout": "Сервер не відповів вчасно — спробуйте ще раз", "err.net.down": "Не вдалося з'єднатися — перевірте інтернет", "err.net.cert": "З'єднання незахищене — перевірте дату та антивірус", "err.answer": "Сервер відповів незрозуміло — спробуйте пізніше", "err.file.missing": "Файл не знайдено", "err.file.denied": "Немає доступу до файлу — закрийте програму, яка його тримає", "err.disk.full": "На диску немає місця", "err.cancelled": "Скасовано", "err.generic": "Не вийшло — подробиці в журналі", "err.server.launch": "Не вдалося запустити %s — перевірте шлях до сервера в розділі «Система»",
		"ov.err.mic":       "Мікрофон недоступний — перевірте пристрій у налаштуваннях",
		"ov.err.recognize": "Помилка розпізнавання (див. лог)", "ov.err.paste": "Не вставилося — текст у «Останньому результаті»",
		"ov.moved":  "Вікно змінилося — текст у буфері обміну",
		"copy.ok":   "Скопійовано",
		"copy.none": "Нема чого копіювати",
		"copy.fail": "Не вдалося скопіювати: %s",
		"mic.busy":  "Триває диктування, зараз перевірити не можна", "mic.check.ok": "Чути добре: пік %.0f дБ, мовлення на %.0f%% запису",
		"mic.check.quiet": "Надто тихо: пік %.0f дБ — додайте гучності мікрофона у Windows або сядьте ближче", "mic.check.clipped": "Перевантаження: обрізано %.1f%% відліків — зменште гучність мікрофона", "mic.check.silent": "Мовлення не чути — перевірте, чи вибрано той мікрофон і чи не вимкнений він",
		"ov.quiet": "Надто тихо, майже нічого не чути", "ov.clipped": "Перевантаження — звук обрізано",
		"ov.cmd.cancelled": "Скасовано голосом",
		"ov.undo.done": "Стерто: %d символів", "ov.undo.none": "Нічого стирати", "ov.undo.typed": "Не стерто: ви друкували або клацали після вставки",
		"ov.undo.window": "Не стерто: інше вікно", "ov.undo.sent": "Не стерто: текст уже надіслано через Enter", "ov.undo.old": "Не стерто: вставлено надто давно",
		"ov.undo.breaks": "Не стерто: у тексті були перенесення рядків", "ov.undo.long": "Не стерто: текст надто довгий", "ov.undo.moved": "Стирання перервано: вікно змінилося",
		"ov.shape.next": "Наступна фраза: %s", "cmd.heading": "заголовок", "cmd.list": "список", "cmd.capital": "з великої", "cmd.quotes": "в лапках", "cmd.nostop": "без крапки",
		"ov.silence":       "Тиша — нічого не розпізнано", "ov.server.loading": "Сервер ще завантажується",
		"ov.stop.silence":  "Зупинено через тишу",
		"ov.tooshort":  "Занадто коротко — тримайте клавіші довше",
		"ov.cancelled": "Скасовано", "ov.editing": "Редагую: %s", "ov.translating": "Перекладаю",
		"ov.llm.needed": "Ця мова потребує LLM-модуль", "td.title": "Перекласти на:", "td.plain": "Перекласти мовою вихідного мовлення",
		"cap.title": "{app} — сполучення клавіш", "cap.prompt": "Натисніть нове сполучення клавіш\n\nзараз: %s   ·   Esc — скасувати",
		"cap.selected": "Обрано: %s", "cap.cancelled": "Скасовано", "hk.taken": "Сполучення %s зайняте Windows: %s. Диктування може не початися", "hk.lock": "блокування комп’ютера", "hk.desktop": "показати робочий стіл", "hk.explorer": "провідник", "hk.run": "вікно «Виконати»", "hk.settings": "параметри Windows", "hk.search": "пошук", "hk.center": "центр сповіщень", "hk.menu": "меню досвідченого користувача", "hk.clipboard": "журнал буфера обміну", "hk.gamebar": "ігрова панель", "hk.voice": "голосове введення Windows", "hk.project": "проєціювання на екран", "hk.tasks": "подання завдань", "hk.layout": "зміна розкладки", "hk.newdesktop": "новий робочий стіл", "hk.closedesktop": "закрити робочий стіл", "hk.snip": "знімок екрана", "hk.switch": "перемикання вікон", "hk.close": "закрити вікно", "hk.cycle": "перебір вікон", "hk.start": "меню «Пуск»", "hk.taskmgr": "диспетчер завдань", "hk.secure": "екран безпеки",
		"err.hotkey.dup":    "Сполучення «%s» призначено двічі — хоткеї не повинні збігатися",
		"cfg.err.recovered": "config.json пошкоджено (%s).\nФайл збережено як %s, налаштування скинуто до типових.",
		"err.disk.space":    "мало місця на диску: вільно %d МБ, потрібно ~%d МБ",
		"err.model.dir":    "не вдалося створити теку %s — перевірте доступ до неї",
		"err.save":          "не вдалося зберегти налаштування: %s — лишив попередні",
		"err.port":          "порт %d не підходить: потрібен номер від 1024 до 65535",
		"err.nolangs":       "залиште хоча б одну мову в списку для питання про переклад",
		"ov.mic.lost":       "Мікрофон відключено — запис перервано",
		"err.hash":          "завантажений файл пошкоджено — спробуйте ще раз",
		"models.check.ok":   "Перевірено моделей: %d — усі файли цілі",
		"models.check.none": "Немає що перевіряти — жодна встановлена модель не має еталонного хешу",
		"models.check.bad":  "Пошкоджені файли: %s — завантажте модель ще раз",
		"hist.insert.gone":  "запис не знайдено",
		"ov.aim": "Клацніть поле, куди вставити · Esc — скасувати",
		"hist.aim.armed": "клацніть поле, куди вставити",
		"hist.aim.busy": "вже чекаю на клацання",
		"hist.aim.off": "вставку скасовано",
		"hist.insert.nowin": "нікуди вставляти — текст скопійовано в буфер",
		"hist.insert.ok":    "вставлено в «%s»",
		"lists.bad":         "файл не підходить",
		"lists.saved":       "збережено в %s",
		"lists.added":       "додано: %d, пропущено: %d",
		"lists.save.title":  "Куди зберегти списки",
		"lists.open.title":  "Звідки завантажити списки",
		"un.title":          "{app} — видалення", "un.confirm": "Видалити {app} з цього комп'ютера?",
		"un.data": "Видалити також налаштування, завантажені моделі та розшифровки субтитрів?", "un.done": "{app} видалено.",
		"model.switching": "Перемикаю модель — розпізнавач перезапускається…", "model.del.active": "Не можна видалити активну модель",
		"model.del.ok":       "Модель видалено",
		"about.text":         "{app} %s\n\nГолос → текст у позицію курсора.\nПоставте курсор у поле введення, утримуйте %s, скажіть фразу, відпустіть — текст вставиться сам.\n\nРозпізнавання: whisper.cpp, повністю локально й офлайн.\nМодель: %s (мова: %s)\n\nНалаштування: клік по іконці в треї або config.json.\nЛоги: {log} (макс. ~2 МБ).",
		"ov.notranslate":     "Активна модель не перекладає — вставлено як розпізнано",
		"ov.engine.fallback": "Другий рушій не запустився — лишаємось на поточному",
		"route.speech":       "Мовлення %s", "route.other": "Інші мови", "route.translate": "Переклад",
		"route.lang.auto":    "будь-яка мова",
		"route.why.language": "тут точніше, з розділовими", "route.why.otherlang": "99 мов",
		"route.why.translate": "перекладає лише Whisper", "route.why.notinstalled": "російську модель не встановлено",
		"route.why.unknownlang": "мову не задано — розпізнає лише Whisper", "route.why.forced": "примусово в config.json",
		"status.line": "Готово · %s · %.1f ГБ вільно", "state.ram.free": "%d МБ вільно",
		"ago.now": "щойно", "ago.min": "%d хв тому", "ago.hour": "%d год тому",
		"chars": "%d символів", "inserted.into": "вставлено в %s",
		"punct.prompt":         "Додай розділові знаки та великі літери. Не змінюй слова, не перекладай, нічого не додавай. Поверни лише виправлений текст.",
		"err.sherpa.notfound":  "розпізнавач sherpa не знайдено: %s",
		"err.sherpa.start":     "sherpa-server завершився під час запуску (див. журнал)",
		"err.sherpa.translate": "ця модель не вміє перекладати",
		"err.sherpa.model":     "Файл моделі не знайдено: %s — завантажте його в налаштуваннях або виправте sherpa_model у config.json",
		"srv.restarting":       "Перезапускаю розпізнавач із новими налаштуваннями…",
	}
	settingsStrings["uk"] = map[string]string{
		"S_SUBS_ON": "Субтитри йдуть",
		"S_SUBS_HERO_HINT": "запустіть — дзвінок чи відео стануть текстом на екрані",
		"S_SUBS_HERO_START": "Запустити вручну",
		"S_SUBS_HERO_STOP": "Зупинити",
		"S_UPDATED_TO": "Програму оновлено до версії %s",
		"S_SUBS_HERO_LIVE": "%s · слова під час мовлення",
		"S_SUBS_HERO_PHRASES": "%s · текст після кожної фрази",
		"S_NAV_SUBS": "Субтитри",
		"S_SUBS_HINT": "Застосунок записує те, що звучить на комп'ютері, — дзвінок, відео — і показує слова у вікні поверх інших.",
		"S_SUBS_NOTICE": "Співрозмовники мають знати, що їхнє мовлення записується текстом.",
		"S_SUBS_START": "Запустити субтитри",
		"S_SUBS_STOP": "Зупинити субтитри",
		"S_SUBS_OFF": "Субтитри вимкнено",
		"S_SUBS_LOADING": "Завантажую %s…",
		"S_SUBS_LISTEN_LIVE": "Слухаю звук комп'ютера · %s, слова під час мовлення",
		"S_SUBS_LISTEN_PHRASES": "Слухаю звук комп'ютера · %s, текст після кожної фрази",
		"S_SUBS_RESTARTING": "Розпізнавання перезапускається…",
		"S_SUBS_LANG": "Мова мовлення",
		"S_SUBS_LANG_SUB": "Мова того, що ви слухаєте, окремо від мови диктування",
		"S_SUBS_M_LIVE": "Писатиме %s: слова з'являються під час мовлення.",
		"S_SUBS_M_PHRASES": "Писатиме %s: текст з'являється після того, як закінчиться кожна фраза.",
		"S_SUBS_M_ABSENT": "Для цієї мови вибрано модель %s, але її не встановлено.",
		"S_SUBS_M_WHISPER": "%s — модель Whisper: на процесорі їй потрібно багато секунд на кожну фразу, і текст сильно відставатиме від мовлення. Для субтитрів краще підходить потокова або швидка модель.",
		"S_SUBS_GO_MODELS": "Відкрити «Мови і моделі»",
		"S_SUBS_REMOTE": "Розпізнавання йде на іншому комп'ютері: кожна фраза надсилається на цей сервер.",
		"S_SUBS_TR": "Переклад",
		"S_SUBS_TR_SUB": "Вікно показує переклад замість розпізнаного тексту",
		"S_SUBS_TR_OFF": "Вимкнено",
		"S_SUBS_TR_LOCAL": "Перекладає %s: переклад іде фраза за фразою, за мить після мовлення.",
		"S_SUBS_TR_API": "Перекладає ваш сервер постобробки (%s): кожна фраза надсилається на нього, і переклад іде фраза за фразою, за мить після мовлення.",
		"S_SUBS_TR_NONE": "Немає ні моделі перекладу, ні сервера постобробки, ні моделі-редактора, тому перекладу не буде.",
		"S_SUBS_TR_GO": "Відкрити «Постобробку»",
		"S_SUBS_TR_SAME": "Мова перекладу збігається з мовою мовлення: перекладати нічого.",
		"S_SUBS_WHO": "Хто перекладає",
		"S_SUBS_WHO_SUB": "Програма бере найкращого з доступних, або ви закріплюєте одного",
		"S_SUBS_WHO_AUTO": "Обирає програма",
		"S_SUBS_WHO_QUALITY": "Якісний: модель на відеокарті",
		"S_SUBS_WHO_SERVER": "Ваш сервер",
		"S_SUBS_TRM": "Модель перекладу",
		"S_SUBS_TRM_SUB": "%s, працює на відеокарті, що має щонайменше 4 ГБ власної пам'яті",
		"S_SUBS_TRM_ABSENT": "Не завантажена",
		"S_SUBS_TRM_GETTING": "Завантажується… %s%",
		"S_SUBS_TRM_CHECKING": "Перевіряється файл…",
		"S_SUBS_TRM_READY": "Готова",
		"S_SUBS_TRM_RUNNING": "Готова, завантажена на відеокарту",
		"S_SUBS_TRM_BAD": "Файл у теці models пошкоджений або неповний",
		"S_SUBS_TRM_ERR": "Не завантажена: %s",
		"S_SUBS_TRM_SMALL": "На цьому комп'ютері вона не запрацює: відеокарта має %s МБ власної пам'яті, а потрібно 4 ГБ.",
		"S_SUBS_TRM_INTEGRATED": "На цьому комп'ютері вона не запрацює: тут лише вбудована графіка, яка ділить пам'ять із комп'ютером.",
		"S_SUBS_TRM_NOCARD": "На цьому комп'ютері вона не запрацює: Windows не повідомляє про жодну відеокарту.",
		"S_SUBS_TRM_UNKNOWN": "На цьому комп'ютері вона не запрацює: ця версія Windows не описує свою відеокарту.",
		"S_SUBS_TRM_DOWN": "Не запустилася: %s",
		"S_SUBS_TRM_DOWN_NEXT": "До наступного запуску субтитрів перекладає наступний у черзі.",
		"S_SUBS_TRM_DOWN_HEARD": "До наступного запуску субтитрів текст показується так, як його почуто.",
		"S_SUBS_TRM_DEL_SURE": "Видалити?",
		"S_SUBS_TRM_W_VULKAN": "драйвер відеокарти не підтримує Vulkan",
		"S_SUBS_TRM_W_DEVICE": "рушій не бачить цю відеокарту",
		"S_SUBS_TRM_W_MEMORY": "на відеокарті бракує вільної пам'яті",
		"S_SUBS_TRM_W_CRASH": "рушій зупинився під час завантаження",
		"S_SUBS_TRM_W_TIMEOUT": "рушій не відповів вчасно",
		"S_SUBS_TRM_W_PORT": "його порт зайнятий іншою програмою",
		"S_SUBS_TRM_W_EXE": "поруч із програмою немає llama-server-gpu.exe",
		"S_SUBS_TRM_DL": "Завантажити",
		"S_SUBS_TRM_CANCEL": "Скасувати",
		"S_SUBS_TRM_DEL": "Видалити",
		"S_SUBS_TRM_TERMS": "Умови використання моделі",
		"S_SUBS_TR_QUALITY": "Перекладає модель перекладу на відеокарті (%s): переклад іде фраза за фразою, за мить після мовлення.",
		"S_SUBS_TR_MISS_CARD": "Обрано «Якісний», але відеокарта цього комп'ютера не підходить, тому текст показується так, як його почуто.",
		"S_SUBS_TR_MISS_MODEL": "Обрано «Якісний», але модель перекладу не завантажена, тому текст показується так, як його почуто.",
		"S_SUBS_TR_MISS_START": "Обрано «Якісний», але модель перекладу не запустилася на відеокарті, тому до наступного запуску субтитрів текст показується так, як його почуто.",
		"S_SUBS_TR_MISS_SERVER": "Обрано ваш сервер, але сервер постобробки не налаштований, тому текст показується так, як його почуто.",
		"S_SUBS_ORIG": "Показувати й початковий текст",
		"S_SUBS_ORIG_SUB": "вікно показує над перекладом те, що було сказано; межу між ними можна тягнути мишею",
		"S_SUBS_KEY": "Сполучення клавіш",
		"S_SUBS_KEY_SUB": "Запускає і зупиняє субтитри з будь-якого вікна",
		"S_SUBS_SAVE": "Зберігати розшифровку",
		"S_SUBS_SAVE_SUB": "Один текстовий файл на кожен запуск у теці transcripts поруч із програмою, а з увімкненим перекладом — другий, з перекладом",
		"S_SUBS_FOLDER": "Тека розшифровок",
		"S_SUBS_WINDOW": "Вікно субтитрів",
		"S_SUBS_ROWS": "Рядків тексту",
		"S_SUBS_SIZE": "Розмір тексту",
		"S_SUBS_SIZE_S": "Дрібний",
		"S_SUBS_SIZE_M": "Середній",
		"S_SUBS_SIZE_L": "Великий",
		"S_SUBS_BOX": "Розмір вікна",
		"S_SUBS_BOX_SUB": "тягніть краї та кути вікна мишею; розмір запам'ятовується для кожної роздільності екрана, а «Рядків тексту» діє, доки його немає",
		"S_SUBS_BOX_BTN": "Повернути розмір",
		"S_SUBS_BOX_DONE": "Вікно повернеться до розміру за замовчуванням",
		"S_SUBS_FREE": "Своє місце",
		"S_SUBS_FREE_SUB": "вікно можна перетягнути будь-куди, місце запам'ятовується для кожної роздільності екрана",
		"S_SUBS_POS": "Де показувати вікно",
		"S_SUBS_POS_SCHEME_SUB": "клацніть по екрану — вікно стане туди",
		"S_SUBS_POS_DRAG_SUB": "тягніть позначку або саме вікно мишею — воно стане будь-де",
		"S_SUBS_MON": "Екран",
		"S_SUBS_MON_SUB": "на якому моніторі показувати вікно субтитрів",
		"S_TITLE": "{app} — налаштування", "S_DICT_HINT": "Терміни, імена та абревіатури через кому — підказка слуху, не команди. Працює для Whisper; російська мова через GigaAM це ігнорує. Типовий набір іде за мовою розпізнавання, доки ви не впишете своє.",
		"S_TR_DEFAULT": "Змінити мову текстового виводу", "S_TR_TARGET": "Мова текстового виводу за замовчуванням", "S_TR_ASK": "Питати мову текстового виводу", "S_TR_ASK_NEVER": "Не питати — перекладати одразу",
		"S_SRCLANG_SUB": "нею ви говорите; вона визначає модель розпізнавання",
		"S_TR_LANGS_SUB": "ці мови будуть кнопками на смужці під час вставлення",
		"S_TR_UNAVAIL": "%s не вміє перекладати — текст залишиться мовою вихідного мовлення",
		"S_TR_GET_WHISPER": "Завантажити %s для перекладу (%d МБ)",
		"S_TR_STANDIN": "%s не вміє перекладати — на смужці запропонує %s",
		"S_TR_LOCK": "%s не можна прибрати зі списку — це мова текстового виводу за замовчуванням. Виберіть іншу мову за замовчуванням, і тоді %s можна буде виключити.",
		"S_TR_LOCK_OK": "Зрозуміло",
		"S_TR_ONE": "Позначено кілька мов, але без запитання переклад завжди йтиме в одну — %s (мова виводу за замовчуванням). Решта залишаться позначеними, але вимкненими.",
		"S_TR_NOMODEL": "%s не вміє перекладати. Якщо продовжити, переклад буде вимкнено й недоступний, поки працює ця модель.",
		"S_TR_CONFIRM": "Підтвердити",
		"S_TR_ASK_ALWAYS": "Питати щоразу", "S_TR_ASK_TIMEOUT": "Питати, з таймаутом", "S_TR_SECONDS": "Таймаут, с",
		"S_TR_LANGS": "Мови в діалозі",
		"S_LLM_HINT": "Позначені профілі застосовуються по черзі, зверху вниз, під час звичайного диктування. Нічого не позначено — текст вставляється як є.",
		"S_PROF_ADD": "Додати", "S_PROF_NAME": "Ім'я", "S_PROF_PROMPT": "Промпт", "S_PROF_TEST": "Перевірка",
		"S_PROF_EDIT": "Редагувати", "S_PROF_CLOSE": "Згорнути",
		"S_CONFIRM_DEL": "Видалити модель «%s»? Її можна буде завантажити знову.", "S_FREE": "вільно",
		"S_DEL_ACTIVE":     "Видалити активну модель «%s»? Розпізнавання зупиниться, доки ви не виберете іншу — завантажити її можна тут же.",
		"S_SUB_PROMPTS":    "Промпти",
		"S_SUB_DICT":       "Словник", "S_UPD": "Оновлення", "S_UPD_CHECK": "Перевірити оновлення", "S_UPD_AUTO": "Перевіряти під час запуску",
		"S_UPD_NONE": "Встановлено останню версію", "S_BADGE_MODELS": "Встановлені моделі", "S_BADGE_MISS": "Модель не завантажена", "S_BADGE_SYSTEM": "Попередження — потрібна увага", "S_BADGE_HIST": "Записів в історії", "S_LOG_OPEN": "Відкрити лог", "S_LOG": "Журнал роботи", "S_LOG_SUB": "усе, що застосунок пише про себе", "S_UPD_AVAIL": "Доступна версія %s.",
		"S_UPD_GO": "Оновити", "S_UPD_ERR": "Не вдалося перевірити оновлення", "S_UPD_DL": "Завантажую оновлення…",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Автор і розробник {app} — локального диктувальника для Windows: голос перетворюється на текст просто в позиції курсора, без хмар і підписок.</p>" +
			"<p>Проєкт відкритий: вихідний код, збірка та свіжі версії — на GitHub.</p>" +
			"<ul>" +
			"<li>Репозиторій: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Профіль автора: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Знайшли помилку або маєте ідею — створіть issue в репозиторії.</p>",
		"S_RAM": "Пам'ять комп'ютера:", "S_HF_PH": "Назва моделі — наприклад, qwen2.5 instruct",
		"S_FIT_OK": "поміститься", "S_FIT_WARN": "впритул", "S_FIT_BAD": "не вистачить пам'яті",
		"S_NO_LLM":      "Ще не встановлено жодної моделі — знайдіть і завантажте в полі пошуку нижче.",
		"S_NO_LLM_PROF": "Промпти стануть доступні після встановлення моделі — блок «Модель» вище на цій вкладці.",
		"S_UPDATED":     "Дата останнього оновлення моделі",
		"S_HOTKEY":      "Сполучення клавіш", "S_CHANGE": "Змінити…", "S_UILANG": "Мова інтерфейсу", "S_AUTO": "Як у системі",
		"S_SEC_SOUND": "Звук", "S_BEEP": "Звукові сигнали запису", "S_SOUND": "Сигнал",
		"S_SND_SPEECH": "Голос Windows", "S_SND_CHIME": "Дзвіночок", "S_SND_SOFT": "М'який", "S_SND_MARIMBA": "Марімба",
		"S_SND_BLIP": "Бліп", "S_SND_POP": "Поп",
		"S_AUTOENTER": "Enter після вставлення (автовідправлення)", "S_RESTORE": "Відновлювати буфер обміну після вставлення",
		"S_NAV_HISTORY": "Історія", "S_HIST_ON": "Зберігати історію диктувань", "S_HIST_ON_SUB": "лише текст, на цьому комп'ютері; звук не зберігається ніколи",
		"S_HIST_DAYS": "Скільки днів зберігати", "S_HIST_MAX": "Скільки записів зберігати",
		"S_HIST_SKIP": "Не записувати з цих програм", "S_HIST_SKIP_SUB": "через кому: keepass.exe, 1password.exe", "S_SKIP_ADD_DLG": "Додавання програми", "S_SKIP_EDIT_DLG": "Зміна програми", "S_SKIP_NAME": "Назва програми", "S_SKIP_NAME_SUB": "Файл без шляху: keepass.exe. Зірочка в кінці ловить усі версії: 1password*", "S_SKIP_OPEN": "Відкриті зараз програми", "S_SKIP_REFRESH": "Оновити список", "S_SKIP_PICKED": "Вибрано %d з %d", "S_SKIP_NONE": "Нічого не вибрано", "S_SKIP_EMPTY": "Список порожній — історія пишеться з усіх програм", "S_SKIP_ADD_BTN": "Додати програму", "S_SKIP_HINT": "Те, що надиктовано в ці програми, до історії не потрапляє. Сама вставка працює як завжди.",
		"S_HIST_LIST": "Записи", "S_HIST_CLEAR": "Очистити", "S_HIST_TILL": "до %s", "S_HIST_TILL1": "до завтра", "S_HIST_TILL_FULL": "Буде видалено %s — строк зберігання %s", "S_HIST_LIST_HINT": "Те, що надиктовано: скопіювати, вставити в будь-яке вікно або видалити.", "S_HIST_COPY": "Копіювати",
		"S_HIST_KEEP": "Скільки зберігати",
		"S_UNIT_MIN": "хвилин",
		"S_UNIT_HOUR": "годин",
		"S_UNIT_DAY": "днів",
		"S_HIST_FIND": "Знайти в історії…", "S_HIST_EMPTY": "Історії поки немає", "S_HIST_ASK": "Видалити всю історію диктувань?",
		"S_SEC_CMD": "Голосові команди", "S_CMD_HINT": "Сказане перетворюється на перенесення рядка, знак, скасування, стирання останньої вставки або оформлення фрази замість того, щоб потрапити в текст. Шукаються цілими словами, застосовуються згори вниз, уже після замін; команди оформлення працюють у кінці диктування.",
		"S_CMD_ADD": "Додати команду", "S_CMD_PRESET": "Додати звичайні", "S_CMD_PH": "фраза, яку ви промовите",
		"S_CMD_NEWLINE": "перенесення рядка", "S_CMD_PARAGRAPH": "новий абзац", "S_CMD_TEXT": "підставити текст", "S_CMD_CANCEL": "скасувати диктування",
		"S_CMD_TEXT_PH": "що підставити", "S_CMD_EMPTY": "Команд поки немає", "S_CMD_DEL": "Видалити команду",
		"S_CMD_P_NEWLINE": "новий рядок", "S_CMD_P_PARAGRAPH": "новий абзац", "S_CMD_P_CANCEL": "скасувати",
		"S_CMD_UNDO": "стерти останню вставку", "S_CMD_HEADING": "зробити заголовком", "S_CMD_LIST": "зробити списком",
		"S_CMD_CAPITAL": "перша літера велика", "S_CMD_QUOTES": "взяти в лапки", "S_CMD_NOSTOP": "прибрати крапку в кінці",
		"S_CMD_HEADING_PH": "знак перед заголовком, наприклад #", "S_CMD_LIST_PH": "слова перед останнім пунктом, наприклад і", "S_CMD_QUOTES_PH": "відкривальні й закривальні лапки, наприклад «»",
		"S_CMD_END_NOTE": "Скажіть у кінці диктування; сказана окремо, команда оформить наступне.",
		"S_CMD_UNDO_NOTE": "Скажіть окремо. Стирає останню вставку, якщо ви відтоді нічого не друкували й не клацали, а вікно те саме; текст, що вже був у вікні, не зачіпається ніколи.",
		"S_CMD_P_UNDO": "видали це", "S_CMD_P_HEADING": "як заголовок", "S_CMD_P_LIST": "як список", "S_CMD_P_CAPITAL": "з великої літери",
		"S_CMD_P_QUOTES": "в лапках", "S_CMD_P_NOSTOP": "без крапки", "S_CMD_P_AND": "і, й, та", "S_CMD_P_MARKS": "«»",
		"S_SEC_REPLACE": "Заміни після розпізнавання", "S_REPLACE_HINT": "Те, що почулося неправильно, стає тим, що ви мали на увазі — одразу після розпізнавання, до промптів. Застосовуються згори вниз.",
		"S_REPL_WHOLE_FULL": "Лише цілі слова", "S_REPL_CASE_FULL": "Враховувати регістр", "S_CMD_ACTION": "Дія",
		"S_FM_ADD": "Додати",
		"S_TIP_REPL_LANG": "Правило діє, лише коли диктування йде вибраною мовою. «всі мови» — діє завжди.",
		"S_TIP_REPL_CASE": "Великі й малі літери різняться: «гіт» і «Гіт» — різні слова. Вимкнено — регістр не важливий.",
		"S_TIP_REPL_WHOLE": "Заміна спрацьовує, лише якщо текст стоїть окремим словом. Вимкнено — шукається й усередині інших слів.",
		"S_TIP_CMD_ACTION": "Що станеться, коли ви промовите фразу: перенесення рядка, новий абзац, підстановка свого тексту, скасування диктування, стирання останньої вставки або оформлення фрази — заголовком, списком, з великої літери, в лапках чи без крапки.",
		"S_LIST_FILTER_PH": "знайти…",
		"S_REPL_DEL": "Видалити заміну",
		"S_LIST_NOTHING": "Нічого не знайдено: «%s»",
		"S_FM_T_REPL_ADD": "Додавання заміни", "S_FM_T_REPL_EDIT": "Зміна заміни",
		"S_FM_T_CMD_ADD": "Додавання команди", "S_FM_T_CMD_EDIT": "Зміна команди",
		"S_MT_DEL": "Видалення моделі", "S_MT_DEL_PROMPT": "Видалення промпту", "S_MT_DL": "Завантаження моделі",
		"S_MT_TR_OFF": "Вимкнення перекладу", "S_MT_TR_ONE": "Переклад без запитання", "S_MT_TR_LOCK": "Мова виводу за замовчуванням",
		"S_MT_REMOTE": "Віддалений сервер", "S_MT_POST": "Зовнішній сервер", "S_MT_HIST": "Очищення історії",
		"S_MT_RESET": "Скидання налаштувань", "S_MT_EXE": "Шлях до сервера",
		"S_DICT_ADD": "Додати слово", "S_FM_T_DICT_ADD": "Додавання слова", "S_DICT_EMPTY": "Слів поки немає",
		"S_DICT_ADD_PH": "слово або кілька через кому",
		"S_DICT_NOMODEL": "Поточна модель %s не підтримує словник — його читають лише моделі Whisper.",
		"S_OV_FREE": "Своє місце", "S_OV_FREE_SUB": "смужку можна перетягнути будь-куди",
		"S_OVPOS_DRAG_SUB": "тягніть смужку мишею — вона стане будь-де",
		"S_OVMON_N": "Екран %d",
		"S_POST_ENABLE": "Увімкнути постобробку",
		"S_API_SUM_URL": "адреса", "S_API_SUM_MODEL": "модель", "S_API_SUM_KEY": "ключ", "S_API_SUM_TIMEOUT": "очікування",
		"S_API_SUM_STATE": "стан", "S_API_NO_MODEL": "не вказано",
		"S_API_NONE": "не налаштований — постобробка йде локально",
		"S_POSTAPI_SETUP": "Налаштувати", "S_API_EDIT": "Змінити", "S_API_KEY_DEL": "Видалити ключ", "S_API_DLG": "Зовнішній сервер",
		"S_LLM_CATALOG": "Каталог моделей", "S_LLM_BLOCK": "Встановлені моделі", "S_LLM_NONE_HINT": "Жодної моделі не встановлено — завантажте знайдену стрілкою, і вона зʼявиться тут", "S_LLM_IN_MEM": "у памʼяті", "S_LLM_ON_DISK": "на диску", "S_LLM_EJECT": "Вивантажити з памʼяті", "S_LLM_FOUND": "знайдено %d", "S_LLM_NOSEARCH": "пошук не запускали", "S_LLM_SEARCH_HINT": "Введіть назву моделі та натисніть «Знайти»", "S_LLM_PICK_WAIT": "Буде доступний, коли модель завантажиться", "S_LLM_INSTALLED": "встановлені",
		"S_LLM_SUM_MODEL": "модель", "S_LLM_SUM_SIZE": "розмір", "S_LLM_SUM_COUNT": "встановлено", "S_LLM_SUM_RAM": "пам'ять",
		"S_DLG_CLOSE": "Закрити", "S_LLM_NOPICK": "не вибрана", "S_NO_PROMPTS": "Промптів поки немає", "S_PROF_DRAG": "перетягніть, щоб змінити порядок",
		"S_PROF_NAME_PH": "як назвати промпт", "S_PROF_TEST_PH": "напишіть фразу для перевірки",
		"S_PF_NEW": "Новий промпт", "S_PF_EDIT": "Зміна промпту",
		"S_POST_NO_MODEL": "увімкнено, але модель не вибрано", "S_POST_NO_API": "увімкнено, але сервер не налаштовано", "S_POST_BAD": "сервер не відповів: %s", "S_POST_NO_PROMPT": "увімкнено, але не позначено жодного промпту", "S_API_TEST": "Тест з'єднання", "S_API_TEST_RUN": "Перевіряю…", "S_API_TEST_OK": "Сервер відповів", "S_API_CLEAR": "Очистити", "S_API_CLEAR_ASK": "Видалити адресу, модель і ключ зовнішнього сервера? Постобробка повернеться до локальної моделі.", "S_RAM_AVAIL": "Доступно памʼяті %s ГБ із %s ГБ", "S_RAM_OF": "%s ГБ із %s ГБ",
		"S_REPL_ADD": "Додати заміну", "S_REPL_FROM_PH": "гіт хаб", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "цілі слова", "S_REPL_CASE": "регістр", "S_REPL_EMPTY": "Замін поки немає",
		"S_PASTE_DELAY": "Затримка перед вставкою", "S_PASTE_DELAY_SUB": "коли програма не встигає прийняти текст",
		"S_OVPOS": "Де показувати смужку", "S_OVPOS_SUB": "біля курсора — поряд із місцем введення; якщо застосунок його не показує, поряд із вказівником миші",
		"S_OVPOS_CARET": "Біля курсора",
		"S_OVTEXT": "Показувати розпізнаний текст", "S_OVTEXT_SUB": "у смужці після вставки, замість кількості символів",
		"S_OVERLAY": "Показувати смужку", "S_OVERLAY_SUB": "під час диктування на екрані видно, що триває запис", "S_TYPEMODE": "Посимвольне введення (для полів без вставлення)",
		"S_RECLANG": "Мова вихідного мовлення", "S_RECAUTO": "Автовизначення",
		"S_LANG_KEY": "Своє сполучення", "S_LANG_KEY_SET": "Призначити…", "S_LANG_KEY_CLEAR": "Прибрати сполучення", "S_LANG_KEY_HINT": "Утримуйте його, щоб диктувати цією мовою її моделлю. Основне сполучення й далі використовує мову вихідного мовлення, вибрану вище.",
		"S_DL": "Завантажити", "S_DEL": "Видалити",
		"S_M_BASE": "швидка, для слабких ПК", "S_M_SMALL": "баланс швидкості й точності", "S_M_MED": "Whisper для будь-якої мови; та, що перекладає", "S_M_TURBO": "максимум точності на CPU", "S_M_GIGAAM": "точніша російською, сама ставить розділові знаки", "S_M_PARAKEET": "25 європейських мов, сама ставить розділові знаки",
		"S_MIC_CHECK": "Перевірити мікрофон", "S_MIC_CHECK_SUB": "три секунди запису та розбір: гучність, перевантаження, чи є мовлення", "S_MIC_CHECKING": "Перевіряю…",
		"S_MCHECK": "Перевірити встановлені моделі", "S_MCHECK_SUB": "звіряє файли моделей з еталонними хешами", "S_MCHECK_GO": "Перевірити", "S_MCHECK_RUN": "Перевіряю…",
		"S_HIST_INSERT": "Вставити",
		"S_MIC": "Мікрофон", "S_MIC_DEFAULT": "Системний за замовчуванням", "S_MIC_REFRESH": "Оновити список",
		"S_MIC_LEVEL": "Рівень сигналу", "S_MIC_QUIET": "тихо",
		"S_THREADS": "Потоки CPU", "S_MINMS": "Мін. запис, мс", "S_MAXSEC": "Макс. запис, с",
		"S_AUTOSTART": "Запускати whisper-server автоматично", "S_PORT": "Порт", "S_SERVEREXE": "Шлях до whisper-server", "S_SERVEREXE_SUB": "заповнюється самотужки; змінюйте, лише якщо сервер лежить в іншому місці", "S_EXE_RESET": "Скинути", "S_EXE_WARN": "Застосунок знаходить whisper-server поруч із собою. Зі шляхом, вписаним вручну, після перенесення теки розпізнавання перестане запускатися. Змінити?", "S_RESET_ALL": "Скинути налаштування", "S_RESET_ALL_SUB": "усе, крім моделей та історії, повертається до заводського", "S_RESET_ALL_BTN": "Скинути", "S_RESET_ALL_ASK": "Повернути всі налаштування до заводських? Моделі, історія та промпти залишаться.", "S_RELOAD_CFG": "Перечитати config.json", "S_RELOAD_CFG_SUB": "якщо файл правили руками", "S_RELOAD_CFG_BTN": "Перечитати", "S_UPD_FOUND": "Є версія %s", "S_THEME": "Колір", "S_THEME_SUB": "колір вікна, смужки та значка в лотку", "S_THEME_GREEN": "Зелений", "S_THEME_AMBER": "Бурштиновий", "S_THEME_BLUE": "Синій", "S_THEME_PINK": "Рожевий", "S_THEME_EDITOR": "Редактор", "S_THEME_NEON": "Неон", "S_WND_MAX": "Розгорнути на весь екран", "S_WND_RESTORE": "Повернути попередній розмір", "S_WND_MIN": "Згорнути в лоток", "S_WND_CLOSE": "Закрити вікно", "S_SKIN": "Дизайн", "S_SKIN_SUB": "шрифт, форма, ефекти та анімація", "S_SKIN_TERMINAL": "Термінал", "S_SKIN_SOFT": "М'який", "S_SKIN_PAPER": "Документ", "S_SKIN_FLUENT": "Fluent", "S_SKIN_STUDIO": "Студія",
		"S_SERVERURL": "Зовнішній сервер (URL)", "S_URLHINT": "Якщо задано — свій сервер не запускається",
		"S_STT_SRV": "Сервер розпізнавання",
		"S_STT_SRV_HINT": "Whisper-моделі запускає окрема програма. Вона може працювати на цьому комп’ютері або на іншому — виберіть, який використати.",
		"S_SRV_LOCAL": "На цьому комп’ютері",
		"S_SRV_REMOTE": "На іншому комп’ютері",
		"S_SRV_REMOTE_HINT": "Той самий whisper-server, запущений деінде: домашній сервер, машина з відеокартою, сусідній комп’ютер.",
		"S_SRV_K_AUTO": "автозапуск",
		"S_SRV_K_FILE": "файл",
		"S_SRV_K_ADDR": "адреса",
		"S_SRV_K_CHECK": "перевірка",
		"S_SRV_NEAR": "whisper-server.exe поряд із застосунком",
		"S_SRV_NOADDR": "не задано",
		"S_SRV_NOCHECK": "не перевіряли",
		"S_SRV_LOCAL_DLG": "Локальний сервер розпізнавання",
		"S_SRV_ADDR": "Адреса сервера",
		"S_SRV_ADDR_SUB": "адреса whisper-server на іншій машині, разом із портом",
		"S_SRV_ON": "увімкнено",
		"S_SRV_OFF": "вимкнено",
		"S_SRV_K_THREADS": "потоки CPU",
		"S_SRV_K_PORT": "порт",
		"S_SRV_DOWN": "Розпізнавання недоступне",
		"S_SRV_DOWN_WHY": "віддалений сервер розпізнавання не налаштовано — задайте адресу в налаштуваннях",
		"S_SRV_DOWN_GO": "Відкрити налаштування сервера",
		"S_SRV_WARN_NOW": "Диктування зараз не працює: вибрано віддалений сервер, а його адреси немає.",
		"S_SRV_WARN_LATER": "Щойно буде вибрано модель Whisper, розпізнавання не працюватиме: адреси віддаленого сервера немає.",
		"S_SAVED":      "Збережено",
		"S_ABOUT_HTML": "<p><b>Голос → текст у позицію курсора.</b></p><p>Поставте курсор у поле введення, утримуйте сполучення клавіш, скажіть фразу, відпустіть — текст вставиться сам.</p><p>Повністю локально й офлайн. Технології: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; моделі з Hugging Face.</p><p>Логи не перевищують ~2 МБ.</p>",
		  "S_SEARCH": "Знайти налаштування…",
		"S_GRP_GENERAL": "Загальне", "S_GRP_SPEECH": "Обробка мовлення", "S_GRP_INFO": "Відомості", "S_NAV_POST": "Постобробка", "S_NAV_HELP": "Довідка", "S_NAV_CONTACTS": "Контакти", "S_HIST_ADD": "Додати", "S_CONTACT_MAIL": "Пошта", "S_DICT_MODEL": "Модель розпізнавання", "S_LIB_ACC": "точність", "S_LIB_SPD": "швидкість",
		"S_HELP_TOC": "На цій сторінці",
		"S_HELP_TOC_SHOW": "Показати зміст — вікно стане ширшим",
		"S_HELP_TOC_HIDE": "Сховати зміст і повернути ширину вікна",
		"S_CONTACT_TITLE": "Зв’язатися",
		"S_ABOUT_DEPS": "Зовнішні модулі",
		"S_ABOUT_DEPS_HINT": "Чужий код, вбудований у застосунок, та його ліцензії. Клац по назві відкриває сторінку проєкту.",
		"S_DEP_WHISPER": "запускає моделі Whisper",
		"S_DEP_LLAMA": "постобробка тексту, моделі GGUF, переклад субтитрів на відеокарті (llama-server-gpu.exe)",
		"S_DEP_VULKAN": "заголовки й бібліотека імпорту, з якими зібрано рушій для відеокарти; саме середовище Vulkan приходить із драйвером відеокарти",
		"S_DEP_SHERPA": "рушій GigaAM, Parakeet, Canary, Qwen, Moonshine",
		"S_DEP_GGML": "тензорна бібліотека всередині whisper.cpp і llama.cpp",
		"S_DEP_ONNX": "запускає моделі всередині sherpa-onnx",
		"S_DEP_WEBVIEW": "вікно налаштувань на WebView2",
		"S_DEP_WV2RT": "компонент Windows, який малює це вікно",
		"S_DEP_MALGO": "захоплення звуку з мікрофона",
		"S_DEP_MINIAUDIO": "звуковий шар усередині malgo",
		"S_DEP_WS": "зв’язок із sherpa-server",
		"S_DEP_XSYS": "виклики WinAPI з Go",
		"S_DEP_WINLOADER": "завантаження DLL усередині go-webview2",
		"S_DEP_PLEX": "шрифт інтерфейсу",
		"S_DEP_HF": "каталог, звідки завантажуються моделі",
		"S_CONTACT_HINT": "Помилка, ідея, питання про налаштування — пишіть на пошту, якщо розмова особиста, або створюйте issue, якщо це помилка.",
		"S_CONTACT_REPO": "Репозиторій",
		"S_CONTACT_ISSUES": "Помилки та ідеї",
		"S_CONTACT_WRITE": "Написати листа",
		"S_CONTACT_OPEN": "Відкрити",
		"S_STATE_ACTIVE": "Розпізнає",
		"S_STATE_USED": "Задіяні моделі",
		"S_STATE_INST": "Встановлені локально",
		"S_STATE_INST_SUB": "моделі на диску, готові до призначення",
		"S_PRESETS": "Яка модель якій мові",
		"S_PRESETS_HINT": "Модель із рядка «Автовизначення» працює для всіх мов, яким не обрано власну — у списку вони сірі. Клацніть мову, щоб дати їй окрему модель. Вона ж визначає, якою мовою говорять, коли мову вихідного мовлення не задано.",
		"S_MFOLDER": "Своя модель",
		"S_DICT_SAVE": "Зберегти",
		"S_OWNM_SUB": "Додайте локальну модель розпізнавання мовлення",
		"S_OWNM_ONEFILE": "Один файл",
		"S_OWNM_FOLDERF": "Тека з файлами моделі",
		"S_OWNM_S1": "Відкрийте теку моделей",
		"S_OWNM_S1S": "Тека призначення:",
		"S_OWNM_S2": "Скопіюйте модель",
		"S_OWNM_S2S": "Виберіть одну з підтримуваних структур",
		"S_OWNM_S3": "Перезапустіть застосунок",
		"S_OWNM_S3S": "Модель з’явиться для мов, які вона підтримує",
		"S_AS_AUTO": "як автовизначення",
		"S_REC_CHIP": "рекомендована",
		"S_BACK_AUTO": "Повернути як автовизначення",
		"S_LANGS_COUNT": "мов: %d",
		"S_LANGS_UNKNOWN": "мови: невідомі",
		"S_TR_EN": "перекладає на %d мов",
		"S_TR_LIST": "перекладає: %s",
		"S_DL_GOING": "завантажується:",
		"S_OPEN_FOLDER": "Відкрити теку",
		"S_UNLOAD": "Вивантажити з пам’яті",
		"S_UNLOAD_SUB": "пам’ять звільниться; наступне диктування завантажить модель знову",
		"S_UNLOAD_GO": "Вивантажити",
		"S_UNLOADED": "Вивантажено",
		"S_NOT_FOR_LANG": "%s не розпізнає цю мову",
		"S_MANUAL_NOTE": "Завантажити з програми не можна — ліцензія забороняє поширення. Скачайте архів за посиланням і розпакуйте в models/moonshine-uk.",
		"S_MANUAL_LINK": "Скачати самому",
		"S_HF_FIT": "лише ті, що підходять цьому комп’ютеру",
		"S_HF_HIDDEN": "приховано: %s",
		"S_WIZ_SKIP_DL": "Скачати пізніше",
		"S_WIZ_SKIP_NOTE": "Без моделі диктування не запрацює. Скачати можна в розділі «Мови і моделі».",
		"S_M_GIGAAM2": "попереднє покоління російської моделі: та сама швидкість, але без розділових знаків",
		"S_M_MOONUK": "українська модель Moonshine: швидка й легка, без розділових знаків",
		"S_M_LOCAL": "знайдена в теці models; характеристики невідомі, тому смужок немає",
		"S_ALL_LANGS": "усі мови",
		"S_OVPOS_SCHEME_SUB": "клацніть по екрану — смужка стане туди",
		"S_OVDRAG": "Перетягніть, куди потрібно",
		"S_OVMON": "Екран",
		"S_OVMON_SUB": "на якому моніторі показувати смужку",
		"S_OVMON_CURSOR": "Екран із курсором",
		"S_M_NEMOTRON": "друкує під час мовлення: текст з’являється на смужці наживо; 40 мов, знаки сама",
		"S_M_TINY": "найменша і найшвидша, для дуже слабких машин; помітно менш точна",
		"S_STATE_LOADED": "Зараз у пам’яті",
		"S_STATE_LOADED_SUB": "моделі вивантажуються самі після простою",
		"S_STATE_WEEK": "За тиждень",
		"S_ST_SUMMARY": "Зведення", "S_ST_OVERLAY": "Смужка на екрані", "S_ST_BEEP": "Звуковий сигнал", "S_ST_AUTORUN": "Запуск із Windows", "S_ST_POST": "Постобробка", "S_ST_LOCAL": "локально", "S_ST_CHECKED": "перевірено %s", "S_ST_GB": "%s ГБ", "S_ST_ON_M": "увімкнений", "S_ST_OFF_M": "вимкнений", "S_ST_MIC_OK": "сигнал у нормі", "S_ST_MIC_BAD": "мікрофон не відповідає", "S_ST_CHECK": "Перевірити", "S_ST_RECOG": "розпізнає %s", "S_ST_VER": "Версія %s", "S_ST_LATEST": "остання", "S_ST_OUTDATED": "не остання", "S_ST_UPD_OK": "у вас остання версія", "S_ST_UPD_DL": "Завантажую оновлення…",
		"S_ST_QUICK": "Швидкі налаштування",
		"S_ST_MODELS": "Моделі",
		"S_ST_USAGE": "Використання за тиждень",
		"S_ST_READY": "Готово до диктування",
		"S_ST_OFF": "Вимкнено в треї",
		"S_ST_OFF_SUB": "сполучення не працює, доки не увімкнете назад",
		"S_ST_ENABLE": "Увімкнути",
		"S_ST_GOTO": "Відкрити налаштування на його вкладці",
		"S_ST_HOTKEY_GO": "Змінити сполучення",
		"S_ST_UPD_LAST": "Версія %s — остання",
		"S_ST_UPD_HAVE": "Доступна версія %s",
		"S_ST_MEM": "Вільно %s ГБ з %s",
		"S_ST_MEM_SUB": "у памʼяті: %s · на диску: %d моделей, %s ГБ",
		"S_ST_MEM_NONE": "нічого",
		"S_ST_LANG": "Мова мовлення",
		"S_ST_ASR": "Розпізнавання",
		"S_ST_ON": "увімкнено",
		"S_ST_OFF_W": "вимкнено",
		"S_ST_ON_F": "увімкнена",
		"S_ST_OFF_F": "вимкнена",
		"S_ST_ACTIVE": "активна",
		"S_ST_IDLE": "не запускається",
		"S_ST_DISK": "лежить на диску, %s",
		"S_ST_USAGE_SUB": "%d знаків · сьогодні %d · у середньому %d знаків",
		"S_WEEK_OTHER": "інші",
		"S_ST_NO_WEEK": "цього тижня диктувань не було",
		"S_ST_AUTORUN_SUB": "застосунок не підніметься сам",
		"S_ST_OVERLAY_SUB": "видно під час запису",
		"S_REPL_LANG": "Мова правила",
		"S_REPL_LANG_ALL": "усі мови",
		"S_M_CANARY": "англійська, німецька, іспанська, французька — і сама перекладає між ними",
		"S_M_QWEN3": "близько 30 мов, сама ставить знаки; найважча і найточніша в каталозі",
		"S_POSTAPI": "Зовнішній сервер",
		"S_POST_HINT": "Править розпізнаний текст за промптами: прибирає слова-паразити, лагодить пунктуацію, змінює стиль. Вимкнено — текст вставляється як розпізнано.",
		"S_POST_MODEL": "Модель",
		"S_SRC_LOCAL": "Локальна",
		"S_SRC_USED": "використовується",
		"S_HF_GO": "Знайти",
		"S_POSTAPI_HINT": "Типово порожньо — вся постобробка локальна. Впишіть адресу, і промпти виконуватиме зовнішній сервер: OpenAI, Groq, власний vLLM — будь-що із сумісним API.",
		"S_POSTAPI_URL": "Адреса",
		"S_POSTAPI_URL_SUB": "порожньо = локальна модель; приклад: https://api.openai.com/v1",
		"S_POSTAPI_MODEL": "Модель",
		"S_POSTAPI_KEY": "Ключ API",
		"S_POSTAPI_KEY_SET": "ключ збережено (зашифровано Windows DPAPI)",
		"S_POSTAPI_KEY_NONE": "ключа немає",
		"S_POSTAPI_SAVE": "Зберегти ключ",
		"S_POSTAPI_TIMEOUT": "Очікування відповіді", "S_SEC_SHORT": "сек",
		"S_POSTAPI_WARN": "⚠ Розпізнаний текст диктувань ітиме на цю адресу. Звук не виходить ніколи. Ключ зберігається зашифрованим.",
		"S_POSTAPI_ASK": "Надсилати розпізнаний текст на %s? Звук залишиться на компʼютері, але текст виходитиме назовні.",
		"S_POSTAPI_BADGE": "зовнішній сервер",
		"S_NOT_INSTALLED": "не встановлено",
		"S_NAV_STATE": "Стан", "S_NAV_DICT": "Керування і поведінка", "S_NAV_MIC": "Мікрофон", "S_NAV_MODELS": "Мови і моделі",
		"S_NAV_TEXT": "Правила", "S_NAV_TR": "Переклад", "S_NAV_SYSTEM": "Система", "S_NAV_ABOUT": "Про програму",
		"S_STATE_HINT": "тримайте й говоріть — текст з'явиться там, де стоїть курсор",
		"S_STATE_PROC": "Постобробка",
		"S_CHANGE_MODEL": "Змінити", "S_PICK_MODEL": "Підібрати", "S_STATE_GET": "Завантажити",
		"S_RETRY": "Повторити", "S_BERR_OPEN": "Відкрити налаштування сервера",
		"S_STATE_LAST": "Останнє диктування", "S_STATE_COPY": "Копіювати", "S_STATE_MEM": "Пам'ять",
		"S_STATE_MEM_SUB": "моделі лишаються в пам'яті, перша фраза йде без затримки",
		"S_HOTMODE":       "Режим", "S_HOTMODE_HOLD": "утримання", "S_HOTMODE_TOGGLE": "перемикач",
		"S_SUB_HOTMODE": "тримайте клавіші або натисніть раз, щоб почати, і раз, щоб зупинити",
		"S_SILENCE":     "Зупинка через тишу", "S_SILENCE_MS": "Пауза до зупинки, мс",
		"S_SUB_SILENCE": "режим «перемикач»: коли ви вже щось сказали, пауза завершує запис, ніби ви знову натиснули клавіші",
		"S_SUB_SILENCE_MS": "скільки можна мовчати; до ваших перших слів запис не зупиниться",
		"S_SUB_MINMS":   "відсікає випадкові натискання",
		"S_SUB_ENTER":   "надсилає повідомлення одразу",
		"S_SUB_CLIP":    "зображення й файли повертаються як були",
		"S_SUB_TYPE":    "допомагає там, де поле не приймає вставку",
		"S_SEC_OVERLAY": "Смужка на екрані",
		"S_SUB_THREADS": "більше потоків не завжди швидше — виміряйте на своїй машині",
		"S_SEC_LLM":     "Модель-редактор",
		"S_PUNCT":       "Розділові знаки й великі літери", "S_SUB_PUNCT": "звідки беруться розділові знаки й великі літери",
		"S_PUNCT_MODEL": "від моделі", "S_PUNCT_LLM": "від моделі-редактора", "S_PUNCT_OFF": "прибирати",
		"S_SUB_TRTARGET": "у неї перекладається текст; у діалозі на смужці вона запропонована першою",
		"S_REMOTE_ABOUT": "Задано віддалений сервер: звук іде на нього, і обіцянка вище не діє, поки він увімкнений.",
		"S_SUB_UPD":      "єдиний мережевий запит, окрім завантаження моделей",
		"S_SEC_SERVICE":  "Службове", "S_SUB_AUTOSTART": "вимкніть, якщо запускаєте сервер самі",
		"S_SUB_PORT":     "розпізнавач перезапуститься сам",
		"S_MODEL_READY":  "Модель завантажено — оберіть її, щоб перемкнутися",
		"S_REMOTE_WARN":  "Звук ітиме на цей сервер. Локальний режим вимкнено.",
		"S_REMOTE_ASK":   "Аудіо перестане оброблятися на цьому комп'ютері й надсилатиметься на %s. Увімкнути віддалений режим?",
		"S_REMOTE_BADGE": "ВІДДАЛЕНО",
		"S_OK":           "Так", "S_CANCEL": "Скасувати", "S_DL_START": "Завантажити", "S_DL_CANCEL": "Скасувати завантаження",
		"S_DL_ASK":    "Модель «%s» не завантажена (%s). Почати завантаження?",
		"S_NOT_FOUND": "нічого",   "S_HELP_HTML": "<p class=\"wh\" id=\"s-what\">Що це таке</p>" +
			"<p>Ви говорите, а слова з'являються там, де блимає курсор. Наче ви їх надрукували, тільки швидше.</p>" +
			"<p class=\"hdim\">Виглядає це так. Ви клацаєте в поле, де зазвичай друкуєте: у листі, у чаті, у рядку пошуку, у документі. Затискаєте дві клавіші й промовляєте фразу звичайним голосом. Відпускаєте. За секунду фраза вже стоїть у полі. Нікуди не треба перемикатися й нічого не треба копіювати: застосунок працює поверх будь-яких програм і чекає на свої клавіші у фоні.</p>" +
			"<p class=\"hdim\">Усе обчислює ваш комп'ютер. Голос не зберігається: він кілька секунд живе в пам'яті, перетворюється на текст і зникає. В інтернет застосунок виходить лише у двох випадках: коли ви завантажуєте модель розпізнавання і коли він перевіряє, чи не вийшла нова версія. Перевірку можна вимкнути, і тоді мережа йому зовсім не потрібна.</p>" +
			"<p class=\"hdim\">Слово «модель» далі траплятиметься часто. Це просто файл, який уміє розбирати мовлення. Він завантажується один раз і лежить у теці поруч із застосунком; що більший файл, то точніше він чує і то повільніше працює.</p>" +
			"<p class=\"wh\" id=\"s-flow\">Що відбувається, коли ви відпустили клавіші</p>" +
			"<div class=\"hflow\">" +
			"<span>запис</span><em>→</em><span>розпізнавання</span><em>→</em><span>заміни</span><em>→</em><span>команди</span><em>→</em><span>переклад</span><em>→</em><span>редагування</span><em>→</em><span>вставлення</span>" +
			"</div>" +
			"<p>Зазвичай це секунда-дві, і розбиратися в кроках немає потреби. Але якщо текст колись вийде не таким, як ви чекали, за цим ланцюжком легко зрозуміти, що саме його змінило.</p>" +
			"<ul>" +
			"<li><b>Запис.</b> Поки клавіші затиснуті, застосунок слухає мікрофон. Надто коротке натискання він пропускає повз вуха, щоб випадково зачеплена клавіша не починала диктування.</li>" +
			"<li><b>Розпізнавання.</b> Голос перетворюється на літери. Цим займається модель. Для російської зазвичай працює GigaAM, для польської — Parakeet, для більшості інших мов — Nemotron, а Whisper бере на себе решту мов і переклад; це просто назви, і вибирати їх вручну не обов'язково.</li>" +
			"<li><b>Заміни.</b> Ваш особистий список «почуто → треба». Якщо модель уперто пише «гіт хаб» замість «GitHub», ви виправляєте це один раз, і далі це виправляється автоматично.</li>" +
			"<li><b>Команди.</b> Кілька фраз працюють як кнопки. Скажіть «новий рядок», і текст перейде на новий рядок, а самі ці слова в поле не потраплять.</li>" +
			"<li><b>Переклад.</b> Якщо ви його ввімкнули, ви говорите однією мовою, а в поле потрапляє текст іншою.</li>" +
			"<li><b>Редагування.</b> Необов'язковий крок. Друга модель причісує текст: прибирає «еее» і «ну типу», розставляє коми, робить формулювання діловим.</li>" +
			"<li><b>Вставлення.</b> Текст потрапляє в поле. Те, що лежало у вашому буфері обміну до диктування, повертається на місце.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Поки все це триває, унизу екрана видно вузьку смужку з назвою поточного кроку. З неї одразу зрозуміло, що застосунок зайнятий, а не завис.</p>" +
			"<p class=\"wh\" id=\"s-first\">Перший запуск</p>" +
			"<p>Під час першого запуску застосунок проводить вас через три кроки. Це забере хвилину, а після них можна диктувати.</p>" +
			"<ul>" +
			"<li><b>Якою мовою ви диктуватимете.</b> Найважливіше запитання: за вашою відповіддю застосунок сам підбере потрібну модель і завантажить її у фоні, і видно, скільки вже готово. Якщо ця мова відрізняється від мови вікна, тут же з'явиться рядок, де можна змінити й мову інтерфейсу.</li>" +
			"<li><b>Проба.</b> Сполучення клавіш (за замовчуванням Ctrl і клавіша Windows; щоб змінити, натисніть те, що вам зручно тримати однією рукою), мікрофон із живим індикатором рівня і поле, куди можна надиктувати пробну фразу. Якщо під час мовлення індикатор не ворушиться, вибрано не той мікрофон. Поле оживає, щойно модель на місці.</li>" +
			"<li><b>Готово.</b> Перемикач «Запускати разом із Windows» — і майстер завершено.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Будь-який крок можна пропустити й повернутися до нього пізніше: усе це є в налаштуваннях. Модель докачується у фоні, навіть якщо ви пішли далі або завершили майстер раніше, ніж вона завантажилася. Якщо захочете пройти майстер наново, запустіть застосунок із ключем <code>-wizard</code>. Під час оновлення він більше не з'являється.</p>" +
			"<p class=\"wh\" id=\"s-overlay\">Смужка внизу екрана</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hpill\"><span class=\"hdot\"></span><span>Говоріть…</span><span class=\"hbars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i></span><span>✕</span></div>" +
			"<div class=\"hcap\">так застосунок показує, що чує вас</div>" +
			"</div>" +
			"<p>Під час диктування це єдине, що видно на екрані. Смужка маленька, мишу не перехоплює й працювати не заважає: вона просто розповідає, що зараз відбувається з вашою фразою.</p>" +
			"<ul>" +
			"<li><b>«Говоріть…»</b> Іде запис. Стовпчики поруч стрибають у такт голосу. Якщо вони лежать нерухомо, вас не чути, і говорити далі немає сенсу.</li>" +
			"<li><b>«Розпізнаю…»</b> Ви відпустили клавіші, працює модель. Зазвичай секунда-дві. Найперше диктування після запуску триває довше: саме тоді модель завантажується в пам'ять.</li>" +
			"<li><b>«Перекладаю…»</b> Іде переклад.</li>" +
			"<li><b>«Редагую: Cleanup (1/2)»</b> Текст править друга модель. Видно, яка саме вказівка виконується і скільки їх усього.</li>" +
			"<li><b>«Вставлено: 128 символів»</b> Готово, текст уже в полі. Смужка згасне сама за кілька секунд.</li>" +
			"<li><b>Хрестик.</b> Якщо передумали, натисніть його на будь-якому кроці. Запис викидається, і в поле нічого не потрапить.</li>" +
			"</ul>" +
			"<p class=\"hdim\">Смужку можна вимкнути, поставити в будь-який кут, прив'язати до курсора або просто перетягнути мишею туди, де вона не заважає. Усе це на вкладці «Керування і поведінка».</p>" +
			"<p class=\"wh\" id=\"s-state\">Екран «Стан»</p>" +
			"<p>Цей екран відповідає на одне запитання: чи можна диктувати просто зараз. Якщо ні, він пояснює чому й дає кнопку, яка веде туди, де це можна виправити.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\"><span class=\"hled\"></span><div class=\"ht\"><b>Готово до диктування</b><span>тримайте й говоріть — текст з'явиться там, де стоїть курсор</span></div><span class=\"hkey\">ctrl</span><span style=\"color:var(--dim)\">+</span><span class=\"hkey\">win</span></div>" +
			"<div class=\"hcap\">зелений квадрат означає, що застосунок чекає на свої клавіші</div>" +
			"</div>" +
			"<p>Верхня картка буває трьох видів. <b>Зелена</b>: усе гаразд, поруч нагадування про клавіші. <b>Жовта</b>: ви вимкнули застосунок через значок біля годинника, і клавіші не працюють, доки ви не натиснете «Увімкнути» тут же, на картці. <b>Червона</b>: налаштування суперечать одне одному. Вибрано розпізнавання на іншому комп'ютері, а адресу не вказано.</p>" +
			"<p>Нижче те, на що зазвичай лише кидають оком:</p>" +
			"<ul>" +
			"<li><b>Субтитри.</b> Картка одразу під верхньою: чи вимкнені субтитри, завантажують модель, ідуть чи перезапускаються, з моделлю і тим, як приходить текст, їхнє сполучення клавіш, якщо його задано, кнопка, яка їх запускає або зупиняє, і стрілка на вкладку «Субтитри».</li>" +
			"<li><b>Використання за тиждень.</b> Кругова діаграма: у які програми йшли ваші диктування і скільки знаків ви наговорили.</li>" +
			"<li><b>Версія.</b> Номер збірки й кнопка перевірки. Якщо вийшла новіша, застосунок скаже про це тут же.</li>" +
			"<li><b>Пам'ять.</b> Скільки її вільно, яка модель зараз у ній сидить і скільки всі моделі важать на диску. Кнопка «Вивантажити» звільняє пам'ять негайно, наприклад перед важкою грою.</li>" +
			"<li><b>Швидкі налаштування.</b> П'ять карток із найуживанішим: мікрофон, мова, смужка, звук, запуск із Windows. Вони показують, як усе налаштовано зараз, а стрілка веде на вкладку, де це змінюється. Перемикачів тут навмисно немає, щоб нічого не зачепити мимохідь.</li>" +
			"<li><b>Моделі.</b> Що встановлено, що працює саме зараз, за які мови відповідає і що просто лежить на диску.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-system\">Вкладка «Система»</p>" +
			"<p>Зовнішній вигляд застосунку, його взаємодія з Windows і кілька службових кнопок. У звичайному житті сюди заходять раз на рік.</p>" +
			"<p class=\"hsub\">Як виглядає застосунок</p>" +
			"<p>Три рядки: мова інтерфейсу, дизайн і колір. На роботу вони зовсім не впливають.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Мова інтерфейсу</td><td>Мова вікна налаштувань і написів на смужці. До мови диктування стосунку не має: вікно може бути українською, а диктувати ви будете англійською.</td></tr>" +
			"<tr><td>Дизайн</td><td>Загальний вигляд. Їх сім: «Редактор» (за замовчуванням: темно-сірий із синім), «Fluent» (як у Параметрах Windows 11), «Студія» (графіт із помаранчевим, майже без заокруглень), «М'який» (світло-бузковий, заокруглений), «Документ» (папір із синім), «Неон» (фіолетовий із градієнтом від рожевого до бірюзового) і «Термінал» (зелені літери на чорному, дужки навколо кнопок). Усі сім набрано одним шрифтом, IBM Plex Mono.</td></tr>" +
			"<tr><td>Колір</td><td>Відтінок лише для дизайну «Термінал»: зелений, бурштиновий, синій, рожевий. Того самого кольору будуть смужка і значок біля годинника. В інших дизайнів кольори свої, і цей рядок ховається.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Запуск і оновлення</p>" +
			"<p>Тут вирішується, чи стартує застосунок разом із системою і чи стежить він за новими версіями.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Запускати разом із Windows</td><td>Застосунок запуститься, щойно ви ввійдете в систему, і чекатиме на ваші клавіші. Це звичайний ярлик в автозавантаженні, такий самий, який ви могли б покласти туди самі.</td></tr>" +
			"<tr><td>Перевіряти під час запуску</td><td>Один раз під час старту застосунок питає, чи не вийшла новіша версія. Більше в інтернет він не ходить, якщо не рахувати завантаження моделей. Вимкніть це, і про нові версії дізнаватиметеся самі.</td></tr>" +
			"<tr><td>Оновлення</td><td>Кнопка перевіряє просто зараз. Якщо є новіша версія, застосунок завантажить її і встановить поверх поточної: налаштування, моделі та історія залишаться на місці.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Службове</p>" +
			"<p>Три кнопки на випадок, коли щось пішло не так.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Журнал роботи</td><td>Відкриває файл, куди застосунок записує, що він робив: яка модель запустилася, скільки тривало розпізнавання, де сталася помилка. Причину зазвичай описано словами, а не кодами. Файл не розростається: старе затирається новим.</td></tr>" +
			"<tr><td>Перечитати config.json</td><td>Усі налаштування лежать в одному текстовому файлі поруч із застосунком. Якщо ви правили його вручну або принесли з іншого комп'ютера, ця кнопка підхопить зміни без перезапуску.</td></tr>" +
			"<tr><td>Скинути налаштування</td><td>Повертає все так, як було після встановлення. Завантажені моделі та історія залишаються на місці: стираються лише налаштування.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Де працює розпізнавання</p>" +
			"<p>Мовлення розбирає окрема маленька програма. Зазвичай вона лежить поруч із застосунком, запускається сама, і думати про неї вам не потрібно. Але її можна поставити й на інший комп'ютер: наприклад, на домашній, де потужна відеокарта, а диктувати з легкого ноутбука. Вибір робиться перемикачем на картках.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hcard\" style=\"border-color:var(--hi);box-shadow:var(--higlow)\"><span class=\"hsw on\"><i></i></span><div class=\"ht\"><b>На цьому комп’ютері</b><span>потоки CPU 4 · автозапуск увімкнено · порт 8910 · whisper-server.exe поряд із застосунком</span></div><button class=\"hbtn\">Редагувати</button></div>" +
			"<div class=\"hcap\">увімкнену картку підсвічено, а другу тим часом пригашено</div>" +
			"</div>" +
			"<p>Кнопка «Редагувати» відкриває вікно з чотирма налаштуваннями. Усі вони мають значення за замовчуванням, і майже завжди їх можна не чіпати.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Потоки CPU</td><td>Скільки ядер процесора віддати розпізнаванню. Що більше ядер, то швидше, але тільки до певної межі: далі виграшу немає, а інші програми починають пригальмовувати.</td></tr>" +
			"<tr><td>Запускати whisper-server автоматично</td><td>Застосунок сам запускає розпізнавач. Вимикати це варто лише в одному випадку: якщо ви запускаєте його вручну й не хочете, щоб застосунок утручався.</td></tr>" +
			"<tr><td>Порт</td><td>Номер, за яким застосунок спілкується з розпізнавачем усередині вашого ж комп'ютера. Ззовні його не видно. Змінювати його треба, лише якщо цей номер зайняла інша програма.</td></tr>" +
			"<tr><td>Шлях до whisper-server</td><td>Де лежить розпізнавач. Заповнюється сам. Вписувати свій шлях варто, лише якщо ви свідомо поклали його в інше місце: інакше після перенесення теки диктування перестане запускатися.</td></tr>" +
			"</table>" +
			"<p>Друга картка має одне налаштування: адресу комп'ютера в мережі, наприклад <code>http://192.168.1.50:8080</code>. Майте на увазі: у цьому варіанті записаний голос іде на ту машину, і за звичайною адресою він передається незашифрованим. Для домашньої мережі це нормально, а через інтернет так краще не працювати.</p>" +
			"<div class=\"hbad\">Якщо вибрано «На іншому комп’ютері», а адресу не вписано, і при цьому вашу мову веде Whisper, диктування не працюватиме. Застосунок скаже про це червоною карткою на екрані «Стан», червоною позначкою біля вкладки «Система» і рядком унизу вікна.</div>" +
			"<p class=\"wh\" id=\"s-mic\">Вкладка «Мікрофон»</p>" +
			"<p>Понад половина невдалих диктувань пов'язана не з моделлю й не з налаштуваннями, а з мікрофоном: застосунок слухає не той мікрофон, або вас майже не чути. На цій вкладці є все, щоб перевірити це за пів хвилини.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Мікрофон</td><td>Який пристрій слухати. Варіант «Системний за замовчуванням» означає «те саме, що вибрано в Windows». Це зручно, якщо ви то в навушниках, то без них: застосунок сам підлаштується під систему.</td></tr>" +
			"<tr><td>Оновити список</td><td>Знадобиться, якщо гарнітуру підключили вже після запуску застосунку: список перечитається.</td></tr>" +
			"<tr><td>Рівень сигналу</td><td>Живий індикатор. Скажіть щось, і стовпчики мають помітно рухатися. Не рухаються: мікрофон приглушено в Windows або вибрано не той.</td></tr>" +
			"<tr><td>Перевірити мікрофон</td><td>Три секунди запису й розбір простими словами: чи досить гучно, чи немає перевантаження (коли звук «хрипить» через надто близький мікрофон) і чи взагалі чути мовлення.</td></tr>" +
			"<tr><td>Звукові сигнали запису</td><td>Короткий звук на початку й наприкінці запису. Особливо корисно, якщо смужку ви вимкнули: тоді принаймні чути, що запис почався.</td></tr>" +
			"<tr><td>Сигнал</td><td>Який саме звук: голос Windows, дзвіночок, марімба та ще кілька. Кнопка поруч відтворює вибраний, щоб не вгадувати за назвою.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-hist\">Вкладка «Історія»</p>" +
			"<p>В історії зберігається список того, що ви наговорили. Лише текст: голос не зберігається ніде й ніколи. За замовчуванням історію вимкнено, і це свідомо: у диктування часто потрапляє особисте.</p>" +
			"<p>Увімкнути її варто, якщо ви хоч раз потрапляли в таку ситуацію: вставили текст не в те вікно, і він зник; треба повторити довгу фразу; хочеться знайти, що диктували вчора, і надіслати ще раз.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Скільки зберігати</label><span class=\"hnum\"><span class=\"hfld\">7</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span><span class=\"hfld\">днів</span></div>" +
			"<div class=\"hrow\"><label>Скільки записів зберігати</label><span class=\"hnum\"><span class=\"hfld\">200</span><span class=\"hspin\"><i>▲</i><i>▼</i></span></span></div>" +
			"<div class=\"hcap\">строк можна задати у хвилинах, годинах або днях, будь-яким числом</div>" +
			"</div>" +
			"<table class=\"hopt\">" +
			"<tr><td>Зберігати історію диктувань</td><td>Головний вимикач. Вимкнули, і нові записи більше не з'являються, а старі стираються.</td></tr>" +
			"<tr><td>Скільки зберігати</td><td>Через скільки часу запис зникне сам. Можна поставити п'ять хвилин, якщо історія потрібна лише на випадок «вставилося не туди», а можна тридцять днів.</td></tr>" +
			"<tr><td>Скільки записів зберігати</td><td>Друга межа, за кількістю. Працює разом із першою: що настане раніше, те й спрацює. Якщо стоїть «7 днів і не більше 50», п'ятдесят перший запис витіснить найстаріший, навіть якщо тиждень ще не минув.</td></tr>" +
			"<tr><td>Не записувати з цих програм</td><td>Список винятків. Усе, що ви диктуєте в ці вікна (банк, менеджер паролів, робочий чат), до історії не потрапить, а саме диктування працюватиме як звичайно. Програму можна вибрати зі списку відкритих вікон, не шукаючи назву файлу.</td></tr>" +
			"</table>" +
			"<p>З кожним записом можна зробити чотири речі. <b>Копіювати</b>: текст ляже в буфер обміну. <b>Вставити</b>: застосунок попросить показати, куди саме. Клацніть потрібне поле в будь-якій програмі, і текст з'явиться там. <b>Видалити</b>: прибирає один запис. <b>Знайти</b>: пошук іде і за текстом, і за назвою програми. Кнопка «Очистити» стирає все одразу, тому просить підтвердження.</p>" +
			"<p class=\"hdim\">Праворуч біля кожного запису написано, коли він зникне: «до 7 вер.», а в останній день цей напис стає жовтим.</p>" +
			"<p class=\"wh\" id=\"s-dict\">Вкладка «Керування і поведінка»</p>" +
			"<p>Як запускається диктування і як текст потрапляє в поле. Сюди заходять, коли щось працює не так, як хочеться: клавіші незручні, текст не вставляється, Enter натискається сам.</p>" +
			"<p class=\"hsub\">Клавіші</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Сполучення клавіш</td><td>Натисніть «Змінити…» і затисніть те, що вам зручно. Застосунок перехоплює це сполучення в усій системі, тому не беріть зайняте, наприклад Ctrl+C. Найрідше зайняті сполучення з клавішею Windows. Мова може мати і своє сполучення — його призначають на вкладці «Мови і моделі».</td></tr>" +
			"<tr><td>Режим</td><td>У режимі «утримання» запис іде, поки клавіші затиснуті: відпустили, і фраза пішла в роботу. У режимі «перемикач» перше натискання вмикає запис, а друге вимикає. Другий режим зручний для довгих монологів, коли рука втомлюється тримати клавіші.</td></tr>" +
			"<tr><td>Зупинка через тишу</td><td>Лише в режимі «перемикач». Коли ви вже щось сказали, пауза тієї довжини, що задана в «Пауза до зупинки, мс», сама завершує запис, ніби ви знову натиснули сполучення клавіш, а на смужці з'являється «Зупинено через тишу». До ваших перших слів запис не зупиниться, тож після натискання можна спокійно зібратися з думками; якщо ви нічого не скажете, запис усе одно завершить межа «Макс. запис, с». Якщо увімкнено «Enter після вставлення (автовідправлення)», пауза надішле повідомлення, тому виберіть довшу паузу.</td></tr>" +
			"<tr><td>Мін. запис, мс</td><td>Наскільки коротке натискання вважати випадковим. Зачепили клавіші на частку секунди, і застосунок промовчить, замість того щоб вставляти уривок слова.</td></tr>" +
			"<tr><td>Макс. запис, с</td><td>Страховка від забутої клавіші: запис обірветься сам і не писатиме в порожнечу пів години.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Як вставляється текст</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Enter після вставлення (автовідправлення)</td><td>Текст не просто вставиться, а одразу відправиться: повідомлення піде, пошук виконається. У месенджері це заощаджує рух, а в листі чи документі майже напевно заважатиме.</td></tr>" +
			"<tr><td>Відновлювати буфер обміну після вставлення</td><td>Вставлення йде через буфер обміну, і застосунок кладе назад те, що там лежало: картинку, файл, скопійований шматок тексту. Залиште ввімкненим, якщо часто щось копіюєте й не хочете це втрачати.</td></tr>" +
			"<tr><td>Посимвольне введення (для полів без вставлення)</td><td>Замість вставлення одним шматком застосунок набере текст, ніби ви друкуєте. Це повільніше, зате рятує там, де вставлення з буфера заборонене: у деяких банківських формах, терміналах, старих програмах.</td></tr>" +
			"<tr><td>Затримка перед вставкою</td><td>Невелика пауза, перш ніж віддати текст вікну. Потрібна рідко: буває, що програма не встигає прийняти текст, і початок фрази губиться.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Смужка на екрані</p>" +
			"<p>Тут вирішується, чи показувати її взагалі і де саме.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Показувати смужку</td><td>Вимкніть, якщо вона заважає або вам досить звукових сигналів.</td></tr>" +
			"<tr><td>Де показувати смужку</td><td>Схема екрана з дев'ятьма місцями: клацніть той кут, де смужка не заважатиме.</td></tr>" +
			"<tr><td>Біля курсора</td><td>Смужка з'являється поруч із місцем, куди ви диктуєте. Якщо програма не повідомляє, де курсор, то поруч із вказівником миші.</td></tr>" +
			"<tr><td>Своє місце</td><td>Смужку можна просто перетягнути мишею, і вона запам'ятає це місце.</td></tr>" +
			"<tr><td>Екран</td><td>На якому моніторі її показувати, якщо їх кілька.</td></tr>" +
			"<tr><td>Показувати розпізнаний текст</td><td>Після вставлення на смужці видно саму фразу, а не кількість символів. Зручно для перевірки, але фразу побачить кожен, хто дивиться на ваш екран.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-models\">Вкладка «Мови і моделі»</p>" +
			"<p>Тут ви кажете, якою мовою диктуватимете, а застосунок вирішує, чим це розпізнавати. Моделей кілька: одні швидкі, інші точні, а треті знають рідкісні мови.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Мова вихідного мовлення</td><td>Мова, якою ви говорите. Якщо мов кілька, поставте «Автовизначення»: модель розбереться сама, хоча трохи повільніше й трохи менш точно, ніж коли мову названо прямо.</td></tr>" +
			"<tr><td>Яка модель якій мові</td><td>Клацніть мову, і розгорнеться список моделей, які її знають, з оцінкою швидкості й точності та розміром файлу. Так одна мова може лишити собі швидку модель, а інша отримати точну. Окремий рядок «Автовизначення» спільний: його модель обслуговує всі мови, яким не вибрано власної, і в списку вони сірі. Вона ж визначає, якою мовою говорять, коли мову вихідного мовлення не задано.</td></tr>" +
			"<tr><td>Своє сполучення</td><td>Усередині розгорнутої мови є місце для її власного сполучення клавіш. Утримуйте його — і диктування йде цією мовою з її моделлю, замінами й правилами, хоч би що стояло в мові вихідного мовлення; основне сполучення й далі використовує мову, вибрану вище. Так англійську диктують одним сполученням, а рідну мову іншим, не заходячи в налаштування. Поки триває запис, смужка показує код мови. Сполучення не може повторювати основне або сполучення іншої мови.</td></tr>" +
			"</table>" +
			"<div class=\"hshot\">" +
			"<div class=\"hrow\"><label>Російська<span class=\"hs\">GigaAM v3 · точність 5 з 5 · швидкість 5 з 5 · 232 МБ</span></label><span class=\"hchip\"><i></i>активна</span></div>" +
			"<div class=\"hrow\"><label>Англійська<span class=\"hs\">як автовизначення · Nemotron 3.5</span></label><button class=\"hbtn\">Підібрати</button></div>" +
			"<div class=\"hcap\">видно, яка модель веде мову і чого це коштує за швидкістю та місцем на диску</div>" +
			"</div>" +
			"<p class=\"hsub\">Переклад</p>" +
			"<p>Деякі моделі вміють не просто розшифрувати мовлення, а одразу його перекласти: ви говорите однією мовою, а в поле потрапляє текст іншою. Перекладає Whisper: на англійську своїм штатним завданням перекладу, на інші мови завдяки тому, що мову результату йому називають як мову мовлення; на далеких парах мов якість нерівна. Якщо вибрана модель перекладати не вміє, поруч із перемикачем з'явиться помаранчевий рядок: якщо Whisper встановлено, смужка запропонує перекласти ним, інакше текст залишиться мовою вихідного мовлення.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Змінити мову текстового виводу</td><td>Вмикає переклад.</td></tr>" +
			"<tr><td>Мова текстового виводу за замовчуванням</td><td>На яку мову перекладати, якщо ви не вибрали іншу просто під час диктування.</td></tr>" +
			"<tr><td>Питати мову текстового виводу</td><td>Можна не питати й завжди перекладати на мову за замовчуванням. Можна питати щоразу: після диктування на смужці з'являться кнопки з мовами. А можна питати з таймаутом: ті самі кнопки, але якщо ви не встигли натиснути, буде взято мову за замовчуванням.</td></tr>" +
			"<tr><td>Мови в діалозі</td><td>Які мови показувати цими кнопками. Залиште дві-три, якими справді користуєтеся, інакше вибирати буде довше, ніж говорити.</td></tr>" +
			"</table>" +
			"<p class=\"hsub\">Моделі на диску</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Своя модель</td><td>Якщо у вас є своя модель, покладіть її в теку моделей і перезапустіть застосунок: вона з'явиться в списках для тих мов, які розуміє. Кнопка поруч відкриває потрібну теку, щоб не шукати її в Провіднику.</td></tr>" +
			"<tr><td>Вивантажити з пам’яті</td><td>Поки застосунок працює, модель займає оперативну пам'ять. Ця кнопка звільняє її негайно, наприклад перед важкою грою. Наступне диктування завантажить модель знову, і перша фраза триватиме довше, ніж зазвичай.</td></tr>" +
			"<tr><td>Перевірити встановлені моделі</td><td>Звіряє завантажені файли з еталоном. Знадобиться, якщо зв'язок обривався під час завантаження і є підозра, що файл пошкодився.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-rules\">Вкладка «Правила»</p>" +
			"<p>Три списки, які правлять текст автоматично, ще до того, як він потрапить у вікно. Усі три необов'язкові: без них застосунок чудово працює.</p>" +
			"<p class=\"hsub\">Розділові знаки й великі літери</p>" +
			"<p>Звідки в тексті беруться крапки, коми й великі літери. Варіантів три: <b>від моделі</b> (як модель почула, так і розставила; швидко і зазвичай достатньо), <b>від моделі-редактора</b> (знаки ставить друга модель; акуратніше, але кожне диктування стає на кілька секунд довшим) і <b>прибирати</b> (текст іде суцільним рядком без знаків; це потрібно рідко, наприклад для рядка пошуку).</p>" +
			"<p class=\"hsub\">Словник</p>" +
			"<p>Список слів, яких модель зазвичай не знає: прізвища, назви компаній, професійні терміни. Це не наказ «пиши так», а підказка: почувши щось схоже, модель охочіше вибере ваш варіант.</p>" +
			"<div class=\"hshot\">" +
			"<div class=\"hchips\"><span class=\"hchip\">GitHub</span><span class=\"hchip\">Docker</span><span class=\"hchip\">Kubernetes</span><span class=\"hchip\">HoldToType</span></div>" +
			"<div class=\"hcap\">слова додаються по одному й живуть окремими картками</div>" +
			"</div>" +
			"<p class=\"hdim\">Словник читають лише моделі Whisper. Nemotron, Parakeet, GigaAM та інші моделі підказок не приймають узагалі: для них інструмент — список замін нижче. Стандартний набір слів змінюється разом із мовою, доки ви не впишете свій.</p>" +
			"<p class=\"hsub\">Заміни після розпізнавання</p>" +
			"<p>Теж список слів, але працює він інакше: це вже не підказка, а гарантія. Те, що трапилося ліворуч, стане тим, що праворуч. Кожна заміна має три перемикачі: <b>Лише цілі слова</b> (тоді «код» не зачепить «кодування», і зазвичай саме це й потрібно), <b>Враховувати регістр</b> (тоді «Надія» і «надія» вважаються різними словами) і <b>Мова правила</b> (заміна спрацює лише тоді, коли ви диктуєте цією мовою).</p>" +
			"<p class=\"hsub\">Голосові команди</p>" +
			"<p>Фрази, які не потрапляють у текст, а виконують дію. Перенесення рядка, новий абзац і своя команда знаходяться цілими словами в будь-якому місці фрази. Скасування й команди оформлення працюють лише в кінці диктування, стирання — лише якщо сказати його окремо. Команда оформлення, сказана окремо, застосовується до наступного диктування, і смужка показує «Наступна фраза: …».</p>" +
			"<ul>" +
			"<li><b>«новий рядок»</b>: перейти на новий рядок;</li>" +
			"<li><b>«новий абзац»</b>: залишити порожній рядок між абзацами;</li>" +
			"<li><b>«скасувати»</b>: викинути фразу й нічого не вставляти;</li>" +
			"<li><b>«видали це»</b>: стерти те, що програма вставила останнім;</li>" +
			"<li><b>«як заголовок»</b>: поставити перед фразою знак, типово #, а після неї порожній рядок;</li>" +
			"<li><b>«як список»</b>: розбити фразу за комами, крапками з комою й сполучником «і» на рядки з дефісом на початку;</li>" +
			"<li><b>«з великої літери»</b>: почати фразу з великої літери;</li>" +
			"<li><b>«в лапках»</b>: взяти фразу в лапки;</li>" +
			"<li><b>«без крапки»</b>: прибрати крапку в кінці;</li>" +
			"<li><b>своя команда</b>: вставити заготовку, наприклад підпис, адресу чи типову відповідь.</li>" +
			"</ul>" +
			"<p>Стирання прибирає лише те, що останнім вставила сама програма, і лише якщо відтоді нічого не друкували й не клацали, вікно те саме, у тексті не було перенесень рядків, його не надіслано через Enter і минуло не більше п'яти хвилин. Текст, який уже був у вікні, не зачіпається ніколи.</p>" +
			"<p class=\"wh\" id=\"s-post\">Вкладка «Постобробка»</p>" +
			"<p>Постобробку виконує друга модель: вона переписує розпізнаний текст за вашою вказівкою. Наприклад: «прибери слова-паразити й повтори», «зроби з цього діловий лист», «оформи це списком». Вказівку пишуть звичайними словами, і називається вона промптом.</p>" +
			"<p>За замовчуванням постобробку вимкнено, і ось чому: це єдиний крок, який помітно сповільнює диктування. Замість секунди-двох вставлення забере кілька секунд, а на слабкій машині й цілий десяток.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Увімкнути постобробку</td><td>Головний вимикач. Поки його вимкнено, решта на вкладці не працює, і текст вставляється таким, яким його розпізнали.</td></tr>" +
			"<tr><td>Модель · Локальна</td><td>Редактор, який працює на вашому комп'ютері. Нікуди нічого не йде, але потрібна пам'ять: такі моделі важать від гігабайта. Поруч є каталог, де модель вибирають і завантажують одним клацанням.</td></tr>" +
			"<tr><td>Модель · Зовнішній сервер</td><td>Замість власної моделі текст надсилається в чужий сервіс: OpenAI, Groq або ваш власний сервер. Це швидше й розумніше, але текст покидає комп'ютер. Ключ доступу зберігається зашифрованим, а кнопка «Тест з'єднання» перевіряє, що адреса й ключ робочі.</td></tr>" +
			"<tr><td>Промпти</td><td>Ваші вказівки редакторові. Позначені застосовуються по черзі, згори вниз: спершу, скажімо, очищення тексту, потім діловий стиль. Якщо нічого не позначено, постобробка нічого не робить, і застосунок попередить про це жовтим рядком.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Голос на зовнішній сервер не йде ніколи, тільки розпізнаний текст. Але текст таки йде, тож вмикайте це свідомо.</div>" +
			"<p>Кожен промпт має <b>назву</b> (під нею його видно на смужці під час роботи), <b>текст вказівки</b> звичайними словами, необов'язкове <b>власне сполучення клавіш</b> (диктування з ним прожене фразу лише через цей промпт, оминаючи решту) і <b>перевірку</b>: поле, де можна вписати фразу й подивитися, що вийде, нічого не диктуючи.</p>" +
			"<p class=\"wh\" id=\"s-subs\">Вкладка «Субтитри»</p>" +
			"<p>Субтитри записують те, що звучить на комп'ютері, — дзвінок, відео, лекцію — і показують слова в невеликому вікні поверх інших вікон. Це те саме розпізнавання, що й для диктування, тільки звук береться з динаміків, а не з мікрофона, і диктування працює, поки йдуть субтитри. Звуки самого застосунку не записуються.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Запустити субтитри</td><td>Запускає і зупиняє їх. Рядок поруч каже, що відбувається: вимкнено, завантажується модель, триває прослуховування або розпізнавання перезапускається. Те саме роблять меню в треї, власне сполучення клавіш і хрестик на вікні.</td></tr>" +
			"<tr><td>Мова мовлення</td><td>Мова того, що ви слухаєте, окремо від мови диктування. Модель — та, яку цій мові дано на вкладці «Мови і моделі». Потокова модель пише слова під час мовлення, решта — кожну фразу після її кінця. Моделі Whisper на процесорі потрібно багато секунд на фразу, і її текст сильно відстає від мовлення; фрази, що накопичилися, йдуть до неї разом і виходять одним рядком.</td></tr>" +
			"<tr><td>Переклад</td><td>Показує у вікні переклад замість розпізнаного тексту, вибраною мовою; заголовок вікна тоді каже «Переклад» і мови, наприклад EN → RU. Мовлення ріжеться на фрази просто по ходу, і кожна фраза одразу йде на переклад. Фраза закінчується на кінці речення або перед сполучним словом, а перекладач отримує разом із нею попередні фрази та свої переклади, тож речення, розрізане на частини, перекладається зв'язно. Готовий переклад з'являється рівно, цілими словами, зі швидкістю читання, а якщо його накопичилося багато — швидше, тож великий шматок не зсуває текст ривком. Фрази одного відрізка мовлення йдуть одним абзацом, а після паузи починається новий рядок. Фраза, яку не вдалося перекласти, показується так, як її почуто, тьмянішим кольором. «Показувати й початковий текст» ділить вікно тонкою лінією: згори тьмяніше те, що було сказано, знизу переклад; лінію можна тягнути мишею. Хто перекладає, вирішує рядок «Хто перекладає». «Якісний» — модель перекладу на відеокарті: завантажується один раз, 2,3 ГБ, і потрібна відеокарта, що має щонайменше 4 ГБ власної пам'яті; модель завантажується на відеокарту, коли субтитри починаються, і йде з неї, коли вони закінчуються, а посилання поруч відкриває умови її використання. «Ваш сервер» — сервер постобробки, заданий на вкладці «Постобробка», і тоді кожна фраза надсилається на цей сервер. «Обирає програма» бере модель перекладу, якщо її завантажено і відеокарта підходить, інакше ваш сервер, якщо його задано, інакше модель-редактор із вкладки «Постобробка», яка перекладає грубо: зміст зрозумілий, але слова бувають неточними. Перекладача, вибраного вручну, програма ніколи не підміняє іншим: якщо він не може працювати, блок каже, чого бракує, а текст показується так, як його почуто. Переклад зберігається окремим файлом поруч із розшифровкою: кожна фраза з часом, коли вона почалася; фраза, яку не вдалося перекласти, записується так, як її почуто. Без перекладача перекладу немає, і субтитри працюють як раніше.</td></tr>" +
			"<tr><td>Сполучення клавіш</td><td>Запускає і зупиняє субтитри з будь-якого вікна. Диктування воно не починає і не може повторювати інше сполучення.</td></tr>" +
			"<tr><td>Зберігати розшифровку</td><td>Кожен запуск пише текстовий файл у теку transcripts поруч із програмою: перший рядок із датою, мовою і моделлю, далі кожен готовий рядок із часом, коли він почався. Поки перший рядок не готовий, нічого не пишеться. Коли переклад увімкнено, у назві розшифровки стоять original і код мови мовлення, а поруч з'являється другий файл, із translation і кодом мови перекладу, наприклад 2026-11-12_10-18-29_original_EN.txt і 2026-11-12_10-18-29_translation_RU.txt.</td></tr>" +
			"<tr><td>Вікно субтитрів</td><td>Де воно стоїть, якого воно розміру і наскільки великі літери. Місце вибирається на картинці екрана так само, як для смужки диктування: клацання по одному з дев'яти місць ставить вікно туди, а зі «Своїм місцем» його можна поставити будь-де, перетягнувши позначку на картинці або саме вікно; за кількох моніторів можна вибрати екран. Розмір змінюється мишею за краї й кути вікна; висота стає на цілу кількість рядків, а ширина й кількість рядків запам'ятовуються для кожної роздільності екрана. «Рядків тексту» діє там, де розмір ще не задано, а «Повернути розмір» забуває задані. Коліщатко миші над вікном гортає текст назад по три рядки; поки ви читаєте старіше, новий текст його не зсуває, а позначка «до останнього ↓» у заголовку повертає до поточного мовлення.</td></tr>" +
			"</table>" +
			"<div class=\"hwarn\">Співрозмовники мають знати, що їхнє мовлення записується текстом.</div>" +
			"<p class=\"wh\" id=\"s-about\">«Про програму» і «Контакти»</p>" +
			"<p>На вкладці «Про програму» є коротка розповідь про те, що це таке і як воно влаштоване, а також список усіх сторонніх бібліотек, з яких зібрано застосунок, з їхніми ліцензіями. Назви клікабельні й ведуть на сторінки цих проєктів. Там же є посилання на репозиторій, де лежать вихідний код і всі версії.</p>" +
			"<p>Вкладка «Контакти» містить три картки: пошта, репозиторій і сторінка для помилок та ідей. Клацання по адресі пошти відкриє вашу поштову програму з уже заповненим одержувачем, а значок поруч просто скопіює адресу.</p>" +
			"<p class=\"wh\" id=\"s-files\">Що лежить у теці застосунку</p>" +
			"<p>Застосунок портативний: усе його господарство лежить в одній теці й більше нікуди не розповзається.</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>config.json</td><td>Усі ваші налаштування одним текстовим файлом. Його можна скопіювати на інший комп'ютер, щоб не налаштовувати все наново.</td></tr>" +
			"<tr><td>holdtotype.log</td><td>Журнал роботи. Перше місце, куди варто зазирнути, якщо застосунок поводиться дивно.</td></tr>" +
			"<tr><td>holdtotype-history.json</td><td>Історія диктувань, якщо ви її ввімкнули. Звичайний текст: за бажання файл можна просто видалити.</td></tr>" +
			"<tr><td>models</td><td>Тека із завантаженими моделями, найважча частина, кілька гігабайтів. Її теж можна перенести на інший комп'ютер, щоб не завантажувати все заново.</td></tr>" +
			"<tr><td>transcripts</td><td>Розшифровки субтитрів, по текстовому файлу на запуск, і переклади до них, коли переклад увімкнено. Це ваші документи: програма видалення стирає їх, тільки якщо ви попросите видалити й налаштування з даними.</td></tr>" +
			"<tr><td>whisper-server, sherpa-server, llama-server</td><td>Ті самі маленькі програми, які запускають моделі. Без них модель лише файл, який нема чим прочитати.</td></tr>" +
			"</table>" +
			"<p class=\"wh\" id=\"s-privacy\">Що йде в інтернет, а що ні</p>" +
			"<ul>" +
			"<li><b>Голос</b> не зберігається й не надсилається, якщо ви самі не ввімкнули розпізнавання на іншому комп'ютері.</li>" +
			"<li><b>Субтитри</b> розпізнаються тут же, як диктування, а розшифровка лишається в теці transcripts; при розпізнаванні на іншому комп'ютері фрази йдуть на той сервер; під час перекладу через ваш сервер постобробки туди ж іде кожен завершений рядок.</li>" +
			"<li><b>Текст</b> не покидає машину, якщо ви самі не ввімкнули зовнішній сервер постобробки.</li>" +
			"<li><b>Історія</b> лежить файлом поруч із застосунком і очищується за строком, який ви задали.</li>" +
			"<li><b>Мережа</b> потрібна двічі: щоб завантажити модель і перевірити оновлення. Перевірку оновлень можна вимкнути.</li>" +
			"<li><b>Ключ</b> від зовнішнього сервісу зберігається зашифрованим засобами Windows і ніколи не потрапляє в журнал.</li>" +
			"</ul>" +
			"<p class=\"wh\" id=\"s-trouble\">Якщо щось пішло не так</p>" +
			"<table class=\"hopt\">" +
			"<tr><td>Текст не з'являється в полі</td><td>Програма не приймає вставлення з буфера. Увімкніть «Посимвольне введення (для полів без вставлення)» на вкладці «Керування і поведінка» або додайте невелику затримку перед вставкою.</td></tr>" +
			"<tr><td>Диктування не запускається</td><td>Сполучення клавіш зайняте іншою програмою. Задайте інше, краще з клавішею Windows: такі займають рідше.</td></tr>" +
			"<tr><td>Записалася тиша</td><td>Мікрофон не той або приглушений. Відкрийте вкладку «Мікрофон», подивіться на індикатор рівня й натисніть «Перевірити мікрофон».</td></tr>" +
			"<tr><td>Перша фраза обробляється дуже довго</td><td>Так і задумано: модель завантажується в пам'ять під час першого звернення. Далі розпізнавання займає секунди.</td></tr>" +
			"<tr><td>Написано «Розпізнавання недоступне»</td><td>Вибрано розпізнавання на іншому комп'ютері, а адресу не задано. Впишіть адресу або поверніться до варіанта «На цьому комп’ютері».</td></tr>" +
			"<tr><td>Постобробка нічого не змінює</td><td>Не позначено жодного промпту, не вибрано модель або зовнішній сервер не відповідає. Причину написано на самій вкладці, поруч із вимикачем.</td></tr>" +
			"<tr><td>Текст дивний, слова замінено</td><td>Зазирніть на вкладку «Правила»: найімовірніше, спрацювала заміна або голосова команда.</td></tr>" +
			"<tr><td>Нічого з переліченого</td><td>«Система» → «Відкрити лог». Застосунок пише про себе все, і причину зазвичай описано словами.</td></tr>" +
			"</table>" +
			"<p class=\"hdim\" style=\"margin-top:16px\">Не знайшли відповіді? Напишіть на пошту з вкладки «Контакти» або створіть issue в репозиторії. І те, й інше читає автор.</p>",
		"S_WIZ_HELLO":      "Дві хвилини на налаштування",
		"S_WIZ_HELLO_TEXT": "{app} перетворює голос на текст просто в позиції курсора: затиснули сполучення клавіш, сказали фразу, відпустили — текст на місці. Усе рахується на вашому комп'ютері, звук нікуди не йде.",
		"S_WIZ_LATER":      "Усе, що ми зараз оберемо, потім змінюється в налаштуваннях.",
		"S_WIZ_MODEL_TEXT": "Модель під вашу мову добирається сама й завантажується у фоні.",
		"S_WIZ_PICKED":     "Дібрана для мови: %s",
		"S_WIZ_DL_BG":      "Завантажується у фоні, можна йти далі",
		"S_WIZ_DL_FAIL":    "Завантаження перервалося: %s",
		"S_WIZ_RETRY":      "Спробувати знову",
		"S_WIZ_TRY_DL":     "Модель ще завантажується (%d%%). Поле оживе, щойно вона буде готова.",
		"S_WIZ_TRY_LOAD":   "Завантажую модель…",
		"S_WIZ_TRY_LIVE":   "З моделлю %s слова з'являються на смужці, поки ви ще говорите.",
		"S_WIZ_DONE_DL":    "Модель ще завантажується. Диктування запрацює саме, щойно модель буде на місці.",
		"S_WIZ_INPUT_TEXT": "Це сполучення ви триматимете під час мовлення. Скажіть щось і перевірте, що смужка рівня рухається.",
		"S_WIZ_T_TRY":      "Проба",
		"S_WIZ_TRY_PH":     "текст з'явиться тут",
		"S_WIZ_T_DONE":     "Готово",
		"S_WIZ_DONE_TEXT":  "{app} живе у треї: лівий клік по значку — налаштування, правий — меню. Диктувати можна в будь-якому вікні, де є курсор введення.",
		"S_AUTORUN":        "Запускати разом із Windows",
		"S_AUTORUN_SUB":    "Запис в автозапуску поточного користувача",
		"S_WIZ_SKIP":       "Пропустити",
		"S_WIZ_BACK":       "Назад",
		"S_WIZ_NEXT":       "Далі",
		"S_WIZ_FINISH":     "Завершити",
		"S_WIZ_WAIT":       "Чекаю на першу фразу…",
		"S_WIZ_HEARD":      "Почув:",
		"S_WIZ_HAVE":       "Усе потрібне вже завантажено",
		"S_WIZ_TRY_TEXT":   "Поставте курсор у поле нижче, затисніть %s, скажіть фразу й відпустіть.",
	}
}
