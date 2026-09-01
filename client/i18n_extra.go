package main

func init() {
	msgs["de"] = map[string]string{
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
		"ov.silence":       "Stille — nichts erkannt", "ov.server.loading": "Server lädt noch",
		"ov.tooshort":  "Zu kurz — Tasten länger halten",
		"ov.cancelled": "Abgebrochen", "ov.editing": "Bearbeite: %s", "ov.translating": "Übersetze",
		"ov.llm.needed": "Diese Sprache benötigt das LLM-Modul", "td.title": "Übersetzen nach:", "td.plain": "Ohne Übersetzung",
		"cap.title": "{app} — Tastenkürzel", "cap.prompt": "Neue Tastenkombination drücken\n\njetzt: %s   ·   Esc — Abbrechen",
		"cap.selected": "Gewählt: %s", "cap.cancelled": "Abgebrochen", "hk.taken": "%s ist von Windows belegt: %s. Das Diktat startet womöglich nie", "hk.lock": "Computer sperren", "hk.desktop": "Desktop anzeigen", "hk.explorer": "Explorer", "hk.run": "Ausführen-Fenster", "hk.settings": "Windows-Einstellungen", "hk.search": "Suche", "hk.center": "Infocenter", "hk.menu": "Power-User-Menü", "hk.clipboard": "Zwischenablageverlauf", "hk.gamebar": "Gamebar", "hk.voice": "Windows-Spracheingabe", "hk.project": "Projizieren", "hk.tasks": "Taskansicht", "hk.layout": "Tastaturlayout wechseln", "hk.newdesktop": "neuer Desktop", "hk.closedesktop": "Desktop schließen", "hk.snip": "Screenshot-Werkzeug", "hk.switch": "Fenster wechseln", "hk.close": "Fenster schließen", "hk.cycle": "Fenster durchblättern", "hk.start": "Startmenü", "hk.taskmgr": "Task-Manager", "hk.secure": "Sicherheitsbildschirm",
		"model.switching": "Modellwechsel — der Erkenner startet neu…", "model.del.active": "Aktives Modell kann nicht gelöscht werden",
		"model.del.ok":   "Modell gelöscht",
		"about.text":     "{app} %s\n\nSprache → Text an der Cursorposition.\nCursor in ein Eingabefeld setzen, %s halten, sprechen, loslassen — der Text wird eingefügt.\n\nErkennung: whisper.cpp, vollständig lokal und offline.\nModell: %s (Sprache: %s)\n\nEinstellungen: Klick auf das Tray-Symbol oder config.json.\nLogs: {log} (max. ~2 MB).",
		"status.nomodel": "Kein Erkennungsmodell geladen — wählen Sie eines in den Einstellungen",
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
		"un.data": "Auch Einstellungen und geladene Modelle löschen?", "un.done": "{app} wurde entfernt.",
		"srv.restarting": "Erkennung wird mit den neuen Einstellungen neu gestartet…",
	}
	msgs["fr"] = map[string]string{
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
		"ov.silence":       "Silence — rien reconnu", "ov.server.loading": "Le serveur charge encore",
		"ov.tooshort":  "Trop court — maintenez les touches plus longtemps",
		"ov.cancelled": "Annulé", "ov.editing": "Édition : %s", "ov.translating": "Traduction",
		"ov.llm.needed": "Cette langue nécessite le module LLM", "td.title": "Traduire vers :", "td.plain": "Sans traduction",
		"cap.title": "{app} — raccourci", "cap.prompt": "Appuyez une nouvelle combinaison\n\nactuel : %s   ·   Échap — annuler",
		"cap.selected": "Choisi : %s", "cap.cancelled": "Annulé", "hk.taken": "%s est pris par Windows : %s. La dictée risque de ne jamais démarrer", "hk.lock": "verrouiller l'ordinateur", "hk.desktop": "afficher le bureau", "hk.explorer": "explorateur de fichiers", "hk.run": "fenêtre Exécuter", "hk.settings": "paramètres Windows", "hk.search": "recherche", "hk.center": "centre de notifications", "hk.menu": "menu avancé", "hk.clipboard": "historique du presse-papiers", "hk.gamebar": "barre de jeu", "hk.voice": "saisie vocale Windows", "hk.project": "projeter sur un écran", "hk.tasks": "affichage des tâches", "hk.layout": "changer de disposition", "hk.newdesktop": "nouveau bureau", "hk.closedesktop": "fermer le bureau", "hk.snip": "outil capture d'écran", "hk.switch": "changer de fenêtre", "hk.close": "fermer la fenêtre", "hk.cycle": "parcourir les fenêtres", "hk.start": "menu Démarrer", "hk.taskmgr": "gestionnaire des tâches", "hk.secure": "écran de sécurité",
		"model.switching": "Changement de modèle — redémarrage…", "model.del.active": "Impossible de supprimer le modèle actif",
		"model.del.ok":   "Modèle supprimé",
		"about.text":     "{app} %s\n\nVoix → texte à la position du curseur.\nPlacez le curseur, maintenez %s, parlez, relâchez — le texte s'insère.\n\nReconnaissance : whisper.cpp, entièrement locale et hors ligne.\nModèle : %s (langue : %s)\n\nRéglages : clic sur l'icône ou config.json.\nLogs : {log} (max ~2 Mo).",
		"status.nomodel": "Aucun modèle de reconnaissance téléchargé — choisissez-en un dans les réglages",
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
		"un.data": "Supprimer aussi les réglages et les modèles téléchargés ?", "un.done": "{app} a été supprimé.",
		"srv.restarting": "Redémarrage de la reconnaissance avec les nouveaux réglages…",
	}
	msgs["es"] = map[string]string{
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
		"ov.silence":       "Silencio — nada reconocido", "ov.server.loading": "El servidor aún carga",
		"ov.tooshort":  "Demasiado corto: mantén las teclas más tiempo",
		"ov.cancelled": "Cancelado", "ov.editing": "Editando: %s", "ov.translating": "Traduciendo",
		"ov.llm.needed": "Este idioma requiere el módulo LLM", "td.title": "Traducir a:", "td.plain": "Sin traducción",
		"cap.title": "{app} — atajo", "cap.prompt": "Pulse una nueva combinación\n\nactual: %s   ·   Esc — cancelar",
		"cap.selected": "Elegido: %s", "cap.cancelled": "Cancelado", "hk.taken": "%s lo usa Windows: %s. Puede que el dictado no empiece nunca", "hk.lock": "bloquear el equipo", "hk.desktop": "mostrar el escritorio", "hk.explorer": "explorador de archivos", "hk.run": "ventana Ejecutar", "hk.settings": "configuración de Windows", "hk.search": "búsqueda", "hk.center": "centro de notificaciones", "hk.menu": "menú avanzado", "hk.clipboard": "historial del portapapeles", "hk.gamebar": "barra de juego", "hk.voice": "escritura por voz de Windows", "hk.project": "proyectar en pantalla", "hk.tasks": "vista de tareas", "hk.layout": "cambiar la distribución", "hk.newdesktop": "escritorio nuevo", "hk.closedesktop": "cerrar el escritorio", "hk.snip": "herramienta de recortes", "hk.switch": "cambiar de ventana", "hk.close": "cerrar la ventana", "hk.cycle": "recorrer ventanas", "hk.start": "menú Inicio", "hk.taskmgr": "administrador de tareas", "hk.secure": "pantalla de seguridad",
		"model.switching": "Cambiando modelo — reiniciando…", "model.del.active": "No se puede borrar el modelo activo",
		"model.del.ok":   "Modelo borrado",
		"about.text":     "{app} %s\n\nVoz → texto en la posición del cursor.\nColoque el cursor, mantenga %s, hable, suelte — el texto se inserta.\n\nReconocimiento: whisper.cpp, totalmente local y sin conexión.\nModelo: %s (idioma: %s)\n\nAjustes: clic en el icono o config.json.\nLogs: {log} (máx ~2 MB).",
		"status.nomodel": "No hay ningún modelo descargado — elige uno en los ajustes",
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
		"un.data": "¿Borrar también los ajustes y los modelos descargados?", "un.done": "{app} se ha quitado.",
		"srv.restarting": "Reiniciando el reconocedor con los nuevos ajustes…",
	}
	msgs["it"] = map[string]string{
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
		"ov.silence":       "Silenzio — nulla riconosciuto", "ov.server.loading": "Il server sta ancora caricando",
		"ov.tooshort":  "Troppo breve — tieni premuti i tasti più a lungo",
		"ov.cancelled": "Annullato", "ov.editing": "Modifica: %s", "ov.translating": "Traduzione",
		"ov.llm.needed": "Questa lingua richiede il modulo LLM", "td.title": "Traduci in:", "td.plain": "Senza traduzione",
		"cap.title": "{app} — scorciatoia", "cap.prompt": "Premi una nuova combinazione\n\nattuale: %s   ·   Esc — annulla",
		"cap.selected": "Scelto: %s", "cap.cancelled": "Annullato", "hk.taken": "%s è occupata da Windows: %s. La dettatura potrebbe non partire", "hk.lock": "bloccare il computer", "hk.desktop": "mostrare il desktop", "hk.explorer": "esplora file", "hk.run": "finestra Esegui", "hk.settings": "impostazioni di Windows", "hk.search": "ricerca", "hk.center": "centro notifiche", "hk.menu": "menu avanzato", "hk.clipboard": "cronologia appunti", "hk.gamebar": "barra di gioco", "hk.voice": "dettatura di Windows", "hk.project": "proiettare su schermo", "hk.tasks": "visualizzazione attività", "hk.layout": "cambio layout", "hk.newdesktop": "nuovo desktop", "hk.closedesktop": "chiudere il desktop", "hk.snip": "strumento di cattura", "hk.switch": "cambiare finestra", "hk.close": "chiudere la finestra", "hk.cycle": "scorrere le finestre", "hk.start": "menu Start", "hk.taskmgr": "gestione attività", "hk.secure": "schermata di sicurezza",
		"model.switching": "Cambio modello — riavvio…", "model.del.active": "Impossibile eliminare il modello attivo",
		"model.del.ok":   "Modello eliminato",
		"about.text":     "{app} %s\n\nVoce → testo alla posizione del cursore.\nPosiziona il cursore, tieni %s, parla, rilascia — il testo viene inserito.\n\nRiconoscimento: whisper.cpp, completamente locale e offline.\nModello: %s (lingua: %s)\n\nImpostazioni: clic sull'icona o config.json.\nLog: {log} (max ~2 MB).",
		"status.nomodel": "Nessun modello di riconoscimento scaricato — scegline uno nelle impostazioni",
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
		"un.data": "Eliminare anche impostazioni e modelli scaricati?", "un.done": "{app} è stato rimosso.",
		"srv.restarting": "Riavvio del riconoscitore con le nuove impostazioni…",
	}
	msgs["pl"] = map[string]string{
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
		"ov.silence":       "Cisza — nic nie rozpoznano", "ov.server.loading": "Serwer wciąż się ładuje",
		"ov.tooshort":  "Za krótko — przytrzymaj klawisze dłużej",
		"ov.cancelled": "Anulowano", "ov.editing": "Edycja: %s", "ov.translating": "Tłumaczenie",
		"ov.llm.needed": "Ten język wymaga modułu LLM", "td.title": "Tłumacz na:", "td.plain": "Bez tłumaczenia",
		"cap.title": "{app} — skrót", "cap.prompt": "Naciśnij nową kombinację\n\nobecna: %s   ·   Esc — anuluj",
		"cap.selected": "Wybrano: %s", "cap.cancelled": "Anulowano", "hk.taken": "%s jest zajęty przez Windows: %s. Dyktowanie może się nie zacząć", "hk.lock": "blokada komputera", "hk.desktop": "pokaż pulpit", "hk.explorer": "eksplorator plików", "hk.run": "okno Uruchom", "hk.settings": "ustawienia Windows", "hk.search": "wyszukiwanie", "hk.center": "centrum powiadomień", "hk.menu": "menu zaawansowane", "hk.clipboard": "historia schowka", "hk.gamebar": "pasek gier", "hk.voice": "pisanie głosem Windows", "hk.project": "projektowanie obrazu", "hk.tasks": "widok zadań", "hk.layout": "zmiana układu klawiatury", "hk.newdesktop": "nowy pulpit", "hk.closedesktop": "zamknięcie pulpitu", "hk.snip": "narzędzie wycinania", "hk.switch": "przełączanie okien", "hk.close": "zamknięcie okna", "hk.cycle": "przeglądanie okien", "hk.start": "menu Start", "hk.taskmgr": "menedżer zadań", "hk.secure": "ekran zabezpieczeń",
		"model.switching": "Zmiana modelu — restart…", "model.del.active": "Nie można usunąć aktywnego modelu",
		"model.del.ok":   "Model usunięty",
		"about.text":     "{app} %s\n\nGłos → tekst w pozycji kursora.\nUstaw kursor, przytrzymaj %s, mów, puść — tekst zostanie wstawiony.\n\nRozpoznawanie: whisper.cpp, w pełni lokalnie i offline.\nModel: %s (język: %s)\n\nUstawienia: kliknij ikonę lub config.json.\nLogi: {log} (maks. ~2 MB).",
		"status.nomodel": "Nie pobrano żadnego modelu — wybierz go w ustawieniach",
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
		"un.data": "Usunąć także ustawienia i pobrane modele?", "un.done": "{app} został usunięty.",
		"srv.restarting": "Ponowne uruchamianie rozpoznawania z nowymi ustawieniami…",
	}

	settingsStrings["de"] = map[string]string{
		"S_TITLE": "{app} — Einstellungen", "S_DICT_HINT": "Begriffe, Namen und Abkürzungen, durch Kommas getrennt — ein Hinweis fürs Gehör, keine Befehle. Gilt für Whisper; russische Sprache über GigaAM ignoriert das. Der Standardsatz folgt der Erkennungssprache, bis Sie eigenes eintragen.",
		"S_TR_DEFAULT": "Sprache der Textausgabe ändern", "S_TR_TARGET": "Standardsprache der Textausgabe", "S_TR_ASK": "Nach der Ausgabesprache fragen", "S_TR_ASK_NEVER": "Nicht fragen — sofort übersetzen",
		"S_SRCLANG_SUB": "Sie sprechen sie; sie bestimmt das Erkennungsmodell",
		"S_TR_LANGS_SUB": "diese Sprachen werden beim Einfügen zu Schaltflächen auf der Leiste",
		"S_TR_UNAVAIL": "nicht verfügbar — %s kann nicht übersetzen",
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
		"S_SEC_CMD": "Sprachbefehle", "S_CMD_HINT": "Gesagtes wird zu einem Zeilenumbruch, einem Zeichen oder einem Abbruch, statt im Text zu landen. Als ganze Wörter erkannt, von oben nach unten angewendet, nach den Ersetzungen.",
		"S_CMD_ADD": "Befehl hinzufügen", "S_CMD_PRESET": "Übliche hinzufügen", "S_CMD_PH": "die Phrase, die Sie sagen werden",
		"S_CMD_NEWLINE": "Zeilenumbruch", "S_CMD_PARAGRAPH": "neuer Absatz", "S_CMD_TEXT": "Text einfügen", "S_CMD_CANCEL": "Diktat abbrechen",
		"S_CMD_TEXT_PH": "was einfügen", "S_CMD_EMPTY": "Noch keine Befehle", "S_CMD_DEL": "Befehl löschen",
		"S_CMD_P_NEWLINE": "neue Zeile", "S_CMD_P_PARAGRAPH": "neuer Absatz", "S_CMD_P_CANCEL": "abbrechen",
		"S_SEC_REPLACE": "Ersetzungen nach der Erkennung", "S_REPLACE_HINT": "Falsch Gehörtes wird zu dem, was gemeint war — direkt nach der Erkennung, vor den Prompts. Von oben nach unten angewendet.",
		"S_REPL_WHOLE_FULL": "Nur ganze Wörter", "S_REPL_CASE_FULL": "Groß-/Kleinschreibung beachten", "S_CMD_ACTION": "Aktion",
		"S_FM_ADD": "Hinzufügen",
		"S_TIP_REPL_LANG": "Die Regel greift nur, wenn Sie in der gewählten Sprache diktieren. „alle Sprachen“ — sie greift immer.",
		"S_TIP_REPL_CASE": "Groß- und Kleinbuchstaben zählen: „git“ und „Git“ sind verschiedene Wörter. Aus — die Schreibung ist egal.",
		"S_TIP_REPL_WHOLE": "Die Ersetzung greift nur, wenn der Text als eigenes Wort steht. Aus — sie trifft auch innerhalb anderer Wörter.",
		"S_TIP_CMD_ACTION": "Was passiert, wenn Sie die Phrase sagen: Zeilenumbruch, neuer Absatz, eigener Text oder Abbruch des Diktats.",
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
		"S_DL": "Laden", "S_DEL": "Löschen",
		"S_M_BASE": "schnell, für schwache PCs", "S_M_SMALL": "ausgewogen", "S_M_MED": "genauer, empfohlen", "S_M_TURBO": "beste Genauigkeit auf CPU", "S_M_PARAKEET": "25 europäische Sprachen, setzt Satzzeichen selbst",
		"S_THREADS":  "CPU-Threads", "S_MINMS": "Min. Aufnahme, ms", "S_MAXSEC": "Max. Aufnahme, s",
		"S_AUTOSTART": "whisper-server automatisch starten", "S_PORT": "Port", "S_SERVEREXE": "Pfad zu whisper-server", "S_SERVEREXE_SUB": "wird automatisch ausgefüllt; nur ändern, wenn der Server woanders liegt", "S_EXE_RESET": "Zurücksetzen", "S_EXE_WARN": "Die App findet whisper-server neben sich. Mit einem handgeschriebenen Pfad startet die Erkennung nach dem Verschieben des Ordners nicht mehr. Ändern?", "S_RESET_ALL": "Einstellungen zurücksetzen", "S_RESET_ALL_SUB": "alles außer Modellen und Verlauf auf Werkszustand", "S_RESET_ALL_BTN": "Zurücksetzen", "S_RESET_ALL_ASK": "Alle Einstellungen auf Werkszustand? Modelle, Verlauf und Prompts bleiben erhalten.", "S_RELOAD_CFG": "config.json neu einlesen", "S_RELOAD_CFG_SUB": "falls Sie die Datei von Hand bearbeitet haben", "S_RELOAD_CFG_BTN": "Neu einlesen", "S_UPD_FOUND": "Version %s ist da", "S_THEME": "Farbe", "S_THEME_SUB": "Farbe von Fenster, Leiste und Symbol im Infobereich", "S_THEME_GREEN": "Grün", "S_THEME_AMBER": "Bernstein", "S_THEME_BLUE": "Blau", "S_THEME_PINK": "Rosa", "S_THEME_EDITOR": "Editor", "S_THEME_NEON": "Neon", "S_WND_MAX": "Bildschirm füllen", "S_WND_RESTORE": "Vorherige Größe", "S_WND_MIN": "In den Infobereich", "S_WND_CLOSE": "Fenster schließen", "S_SKIN": "Design", "S_SKIN_SUB": "Schrift, Form, Effekte und Bewegung", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Sanft", "S_SKIN_PAPER": "Dokument",
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
		"S_DEP_LLAMA": "Nachbearbeitung, GGUF-Modelle",
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
		"S_PRESETS_HINT": "Klicken Sie eine Sprache an — darunter öffnet sich die Modellwahl für sie. Sprachen ohne eigenes Modell nutzen das der automatischen Erkennung.",
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
		"S_TR_EN": "übersetzt ins Englische",
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
		"S_ST_SUMMARY": "Übersicht", "S_ST_OVERLAY": "Anzeige auf dem Bildschirm", "S_ST_BEEP": "Tonsignal", "S_ST_AUTORUN": "Start mit Windows", "S_ST_POST": "Nachbearbeitung", "S_ST_LOCAL": "lokal", "S_ST_CHECKED": "geprüft %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "an", "S_ST_OFF_M": "aus", "S_ST_MIC_OK": "Signal ist in Ordnung", "S_ST_MIC_BAD": "das Mikrofon schweigt", "S_ST_CHECK": "Prüfen", "S_ST_RECOG": "erkannt von %s", "S_ST_VER": "Version %s", "S_ST_LATEST": "aktuell", "S_ST_OUTDATED": "veraltet", "S_ST_UPD_OK": "Sie haben die neueste Version", "S_ST_UPD_DL": "Update wird geladen…",
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
		"S_WIZ_NEED_MODEL": "Laden Sie zuerst ein Modell — ohne Modell gibt es nichts zu erkennen",
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
		"S_HELP_HTML": "<p class=\"wh\">Wie es funktioniert</p>" +
			"<p>Kürzel halten — die Aufnahme läuft (die Leiste am unteren Bildschirmrand zeigt Ihren Pegel). Loslassen — der Ton wird erkannt, bei Bedarf übersetzt und durch die Prompts geschickt, und der fertige Text landet an der Cursorposition. Das ✕ auf der Leiste bricht in jedem Schritt ab.</p>" +
			"<p>Der ganze Weg: <b>Aufnahme → Erkennung (GigaAM für Russisch, sonst Whisper) → Übersetzung (falls aktiv) → Prompts (LLM) → Einfügen</b>. Jeder Schritt ist auf der Leiste zu sehen.</p>" +
			"<p class=\"wh\">Erster Start</p>" +
			"<p>Beim allerersten Start öffnet sich ein Assistent mit fünf Schritten: Sprache der Oberfläche, Sprache des Diktats (das Modell sucht er aus und lädt es), Tastenkürzel und Mikrofon mit Pegelbalken, ein Feld zum Ausprobieren und zuletzt der Start mit Windows. Überspringen geht — es läuft auch so; zurückholen mit <b>{exe} -wizard</b>. Bei einem Update erscheint er nicht.</p>" +
			"<p class=\"wh\">Die Leiste</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Sprechen…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Sprechen…</b> — Aufnahme: ein roter Punkt und die Pegelbalken.</li>" +
			"<li><b>Erkenne…</b> — Whisper arbeitet; beim Übersetzen — „Übersetze“, beim Bearbeiten — „Bearbeite: Name (1/2)“.</li>" +
			"<li><b>Eingefügt: N Zeichen</b> — fertig; bei Fehlern oder Stille steht dort kurz der Grund.</li>" +
			"<li>Das ✕ rechts bricht in jedem Schritt ab; die Leiste nimmt niemals den Eingabefokus. Leiste und Animation lassen sich unter „Diktat“ abschalten.</li>" +
			"<li>Wo die Leiste erscheint — unten, oben oder am Cursor — und ob sie den erkannten Text statt der Zeichenzahl zeigt, stellen Sie unter „Diktat“ ein.</li>" +
			"<li>Solange die Leiste fragt, sagt ihre obere Zeile genau das — „Warte auf Ihre Antwort“ — und der Punkt hört auf zu pulsieren. Jede Antwort trägt ihre Nummer: 1…9 wählen eine aus, Enter nimmt die hervorgehobene, Esc bricht alles ab; die Tasten stehen rechts in derselben Zeile. Zehn Sekunden vor dem Aufnahmelimit läuft auf der Leiste ein bernsteinfarbener Countdown.</li>" +
			"<li>In der Titelzeile stehen drei Schaltflächen: in den Infobereich, Bildschirm füllen und schließen. Dieselbe Schaltfläche holt das gefüllte Fenster auf die vorherige Größe zurück, und die mit der Maus gesetzte Größe bleibt erhalten. Unter 760×500 geht das Fenster nicht.</li>" +
			"<li>Lange Namen — Gerät, Modell, Datei — werden auf den Status-Karten mit Auslassungspunkten gekürzt, damit die Karten in einer Linie stehen; der ganze Name erscheint als Hinweis, wenn der Zeiger auf der Karte ruht. Die Hinweise sind in den Farben des aktuellen Aussehens gezeichnet, nicht in denen des Systems.</li>" +
			"<li>Das Aussehen kommt aus zwei Listen im Bereich „System“. „Design“ bestimmt Schrift, Form, Rahmenstärke, Schein und Art der Bewegung; es gibt drei: „Terminal“ (grün, Standard), „Editor“ (flaches Grau, ohne Schein) und „Neon“ (violett, abgerundet). „Farbe“ wird nur beim Terminal angeboten und ändert allein die Farbe von Fenster, Leiste und Symbol im Infobereich: Grün, Bernstein, Blau, Rosa. Die anderen Designs bringen ihre eigenen Farben mit. Die Wahl gilt sofort, ohne Neustart.</li>" +
			"</ul>" +
			"<p class=\"wh\">Die Frage nach der Übersetzung</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Erkenne…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Übersetzen nach:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Ohne Übersetzung</span></div></div>" +
			"<p>Die Leiste fragt selbst, auf einer zweiten Zeile, sobald Sie das Kürzel loslassen — in den Modi „jedes Mal fragen“ und „mit Timeout fragen“. Die Schaltflächen kommen aus „Sprachen im Dialog“; die Zielsprache ist hervorgehoben. Mit Timeout schrumpft unter dieser Schaltfläche ein Strich: läuft er aus, gilt die hervorgehobene Sprache. <b>Ohne Übersetzung</b> fügt den Text so ein, wie er gehört wurde; das ✕ auf der Leiste bricht alles ab. Auch die Tastatur geht: Enter nimmt die hervorgehobene Antwort, 1…9 wählen eine Schaltfläche, Esc bricht ab.</p>" +
			"<p class=\"wh\">Sicheres Einfügen</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Erkenne…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Fokus gewechselt — einfügen?</span><span class=\"mock-btn on mock-cd\">Hier einfügen</span><span class=\"mock-btn\">Kopieren</span></div></div>" +
			"<ul>" +
			"<li>Das Zielfenster wird in dem Moment gemerkt, in dem Sie das Kürzel drücken. Hat sich der Fokus während der Verarbeitung geändert, wird nichts eingefügt — die Leiste fragt auf ihrer zweiten Zeile: <b>Hier einfügen</b> (ins aktuelle Fenster), <b>Kopieren</b> (in die Zwischenablage) oder das ✕. Läuft die Zeit ab, wird nicht eingefügt und der Text bleibt im letzten Ergebnis.</li>" +
			"<li>Enter nach dem Einfügen wird nur gedrückt, wenn das Zielfenster dasselbe geblieben ist.</li>" +
			"<li><b>Letztes Ergebnis</b> — der fertige Text jedes Diktats bleibt bis zum nächsten im Speicher; im Menü im Infobereich gibt es „Letztes Ergebnis kopieren“. Ein misslungenes Einfügen oder ein Fokuswechsel kostet nie ein Diktat.</li>" +
			"</ul>" +
			"<p class=\"wh\">Mikrofon prüfen</p>" +
			"<p>Die Schaltfläche „Test“ unter „Mikrofon“ nimmt drei Sekunden auf und zerlegt sie: Spitzenpegel in Dezibel, wie viel der Aufnahme wirklich Sprache enthält und wie viele Abtastwerte abgeschnitten wurden. Die Antwort kommt in Worten: klingt gut, zu leise — Pegel in Windows anheben, übersteuert — Pegel senken, keine Sprache gehört — ist das richtige Mikrofon gewählt. Dasselbe wird nach jedem Diktat gemessen und ins Log geschrieben; kommt die Erkennung leer zurück, nennt die Leiste den Grund — zu leise, übersteuert oder Stille — statt nur zu sagen, sie habe nichts gehört.</p>" +
			"<p class=\"wh\">Aus dem Verlauf einfügen</p>" +
			"<p>Jeder Eintrag im Verlauf hat die Schaltfläche „Einfügen“: sie holt das Fenster zurück, aus dem Sie die Einstellungen geöffnet haben, und fügt den Text dort ein wie ein gewöhnliches Diktat. Gibt es kein solches Fenster, landet der Text einfach in der Zwischenablage, und das Programm sagt es.</p>" +
			"<p class=\"wh\">Die Listen in einer Datei</p>" +
			"<p>Ersetzungen und Sprachbefehle lassen sich in eine .json-Datei speichern und auf einem anderen Rechner laden — die Schaltflächen unter der Befehlsliste im Bereich „Text“. Beim Laden wird nichts überschrieben: nur die Zeilen, die noch fehlen, kommen hinzu, und das Programm nennt die Zahl der hinzugefügten und der übersprungenen.</p>" +
			"<p class=\"wh\">Unversehrtheit der Dateien</p>" +
			"<p>Zu jedem Modell aus dem Katalog gehört ein bekannter SHA-256-Referenz-Hash. Nach dem Herunterladen wird die Datei damit verglichen: passt sie nicht, wird sie gelöscht und der Download kann wiederholt werden. Die Schaltfläche „Prüfen“ im Bereich „Modelle“ vergleicht die bereits installierten Modelle genauso, und beim Update wird auch das heruntergeladene Installationsprogramm geprüft — eine fremde Datei wird nicht gestartet.</p>" +
			"<p class=\"wh\">Verlauf der Diktate</p>" +
			"<p>Der Abschnitt „Verlauf“ in der linken Spalte bewahrt auf, was Sie diktiert haben: nur Text, nur auf diesem Rechner, Ton wird nie gespeichert. Standardmäßig aus, eingeschaltet wird er mit einem Schalter an derselben Stelle. Einträge bleiben eine eingestellte Zahl von Tagen und bis zu einer eingestellten Anzahl, Älteres fällt von selbst heraus; „Aus diesen Programmen nie aufzeichnen“ listet durch Komma getrennt jene, aus denen nichts gespeichert werden soll — Passwortmanager, Banking. Die Suche greift auf Text und Programmnamen, die Schaltfläche neben einem Eintrag legt ihn in die Zwischenablage, und „Leeren“ entfernt alles samt der Datei <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Sprachbefehle</p>" +
			"<p>Unter den Ersetzungen im Bereich „Text“ steht eine Liste von Befehlen: Gesagtes wird zur Handlung statt zu Wörtern. „Neue Zeile“ und „neuer Absatz“ setzen einen Umbruch — Modelle tun das nie; „abbrechen“ verwirft das ganze Diktat und fügt nichts ein; „Text einfügen“ setzt beliebiges ein, auch ein Smiley. Die Schaltfläche neben der Liste füllt sie mit den üblichen Wendungen in der Sprache der Oberfläche. Befehle werden als ganze Wörter erkannt und laufen nach den Ersetzungen, damit Prompts und Übersetzung schon den fertigen Text bekommen. Überflüssige Leerzeichen um die Umbrüche verschwinden von selbst. Das Feld darunter probiert Ersetzungen und Befehle an jedem Satz: ein Umbruch erscheint als ⏎.</p>" +
			"<p class=\"wh\">Ersetzungen nach der Erkennung</p>" +
			"<p>Unter „Text“ lässt sich auflisten, was das Modell falsch hört und was daraus werden soll: „git hub“ → GitHub, Nachnamen, hauseigene Begriffe. Ersetzungen laufen direkt nach der Erkennung, vor den Prompts, damit der Editor bereits die richtigen Wörter bekommt. Die Übersetzung ins Englische passiert innerhalb der Erkennung, die Ersetzungen sehen also bereits den übersetzten Text. Standardmäßig gelten ganze Wörter und keine Groß-/Kleinschreibung; die beiden Schalter daneben ändern das. Regeln greifen von oben nach unten. Das Feld darunter probiert sie an jedem Satz aus, ganz ohne Diktat.</p>" +
			"<p class=\"wh\">Regeln pro Programm</p>" +
			"<p>Unter „Diktat“ lassen sich Regeln für einzelne Programme festlegen: womit eingefügt wird (Zwischenablage oder zeichenweise), ob Enter gedrückt wird, wie lange vor dem Einfügen gewartet wird und welche Prompts gelten. Ein Programm wird über seinen Dateinamen angegeben — <b>chrome.exe</b>; in einer Regel dürfen mehrere durch Komma stehen, ein Stern am Ende fängt alle Namen mit diesem Anfang. Die erste passende Regel gewinnt; ohne Regeln oder ohne Treffer gilt alles wie in den allgemeinen Einstellungen. Die Schaltfläche neben der Liste trägt das zuletzt beschriebene Programm ein.</p>" +
			"<p class=\"wh\">Diktat</p>" +
			"<ul>" +
			"<li><b>Tastenkürzel</b> — das Hauptkürzel fürs Diktieren. Jede Kombination lässt sich aufnehmen; linke und rechte Modifikatoren werden unterschieden. Kürzel für Diktat, Übersetzung und Profile müssen eindeutig sein — eine Dopplung verhindert das Speichern.</li>" +
			"<li><b>Modus</b> — Tasten halten oder einmal drücken zum Starten und noch einmal zum Beenden.</li>" +
			"<li><b>Oberflächensprache</b> — wechselt sofort; „Wie im System“ folgt Windows.</li>" +
			"<li><b>Erkennungssprache</b> — ein Hinweis für Whisper; „auto“ erkennt die Sprache am Klang.</li>" +
			"<li><b>Ton</b> — Signale für Start und Ende: mehrere Sätze plus Windows-Systemklänge, ▶ spielt sie vor.</li>" +
			"<li><b>Enter nach dem Einfügen</b> — schickt den diktierten Text sofort ab (praktisch in Messengern).</li>" +
			"<li><b>Zwischenablage wiederherstellen</b> — legt den vorherigen Inhalt vollständig zurück, auch Bilder, Dateien und formatierten Text. Lässt sich der Inhalt nicht sichern, bleibt die Zwischenablage unangetastet und der Text wird Zeichen für Zeichen getippt.</li>" +
			"<li><b>Leiste und Animation</b> — die Statusanzeige am unteren Rand; die Animation lässt sich abschalten.</li>" +
			"<li><b>Zeichenweise einfügen</b> — statt Strg+V werden Tastendrücke simuliert, für Felder, die das Einfügen verweigern.</li>" +
			"</ul>" +
			"<p class=\"wh\">Erkennung</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Automatische Erkennung</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Russisch</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · wie automatische Erkennung — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Welches Modell für welche Sprache</b> — die Registerkarte „Sprachen & Modelle“ ist eine Liste der Sprachen. Klicken Sie eine an — darunter öffnen sich die Modelle, die sie beherrschen: das zugewiesene und das empfohlene zuerst, fehlende mit Größe und Download-Pfeil. Ein Klick auf die Karte ist die Wahl; ein fehlendes Modell lädt sich selbst und übernimmt, sobald es bereit ist. Sprachen ohne eigenes Modell erben das der automatischen Erkennung und stehen gedimmt da.</li>" +
			"<li><b>Der Katalog</b> — Whisper: Base (schnell, für schwache PCs), Small (die Balance), Medium und Turbo (genauer und langsamer; „q5“ ist die quantisierte Fassung: etwas kleiner und schneller, fast ohne Qualitätsverlust), sie übersetzen auch ins Englische; GigaAM v3 ist auf Russisch genauer und setzt selbst Satzzeichen; Parakeet v3 deckt 25 europäische Sprachen ab; Nemotron 3.5 tippt, während Sie sprechen. Geladen wird aus den offiziellen Hugging-Face-Repositorien, jede Datei gegen ihren Referenz-Hash geprüft.</li>" +
			"<li><b>Eigenes Modell</b> — geeignet ist Whisper als einzelne ggml-*.bin-Datei oder ein sherpa-onnx-Modellordner (encoder.onnx, decoder.onnx, tokens.txt). Legen Sie es in den Ordner models neben der Anwendung und starten Sie sie neu — das Modell erscheint in der Auswahl passender Sprachen; seine Werte sind unbekannt, darum ehrlich ohne Balken.</li>" +
			"<li>Der Erkennungsserver hält das Modell zwischen den Sätzen im Speicher — das erste Diktat nach dem Start dauert länger (Laden), danach braucht die Erkennung ein bis drei Sekunden.</li>" +
			"<li><b>Wörterbuch</b> — Begriffe, Namen und Abkürzungen, durch Kommas getrennt. Ein Hinweis für Whispers „Gehör“, damit seltene Wörter richtig ankommen; keine Befehle.</li>" +
			"<li><b>Mikrofon</b> — Gerätewahl mit Pegelanzeige (sprechen Sie, und der Balken bewegt sich, dann wird das Gerät gehört). Wird das gewählte Gerät abgezogen, greift das Systemgerät; eine Aufnahme ohne Sprache wird gar nicht erst zur Erkennung geschickt — stattdessen meldet die Leiste „Stille“.</li>" +
			"<li><b>Dienst</b> — der Erkennungsserver startet selbst und läuft lokal. Port, Pfad oder ein entfernter Server lassen sich ändern; die Erkennung startet danach von allein neu.</li>" +
			"<li><b>Übersetzung</b> — übersetzt wird ausschließlich mit Whisper: ins Englische im nativen Modus, in andere Sprachen <b>experimentell</b> über die erzwungene Ausgabesprache (die Qualität hängt vom Sprachpaar ab; große Sprachen gelingen am besten). Das Turbo-Modell ist dafür nicht trainiert — die Einstellungen warnen, solange es aktiv ist. „Immer in die Zielsprache übersetzen“ übersetzt jedes Diktat ohne Rückfrage. Ohne dieses Häkchen gilt der Fragemodus: immer oder mit Timeout — vor der Erkennung erscheint der Sprachdialog, und nach Ablauf gilt die Zielsprache. Das eigene Übersetzungskürzel übersetzt einmalig, ohne das normale Diktat zu verändern.</li>" +
			"</ul>" +
			"<p class=\"wh\">Nachbearbeitung (LLM)</p>" +
			"<p>Eine freiwillige zweite Schicht: ein lokales Sprachmodell (llama.cpp) bearbeitet den erkannten Text nach Ihren Prompts — entfernt Füllwörter, ändert den Stil, formatiert. Vollständig offline, nur auf der CPU.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modelle</b> — die installierten Bearbeitungsmodelle; der Radioknopf wählt das aktive (gilt sofort), das ✕ löscht (auch das aktive — dann ist die Nachbearbeitung aus). Der Fortschritt beim Laden steht ebenfalls hier.</li>" +
			"<li><b>Suche</b> — GGUF-Modelle auf Hugging Face nach Namen (etwa „qwen2.5 instruct“). Jedes Repository zeigt Datum der letzten Änderung, Zahl der Downloads und ein ↗ zur Modellseite; ein Klick auf die Zeile öffnet die Quant-Dateien. Die Anzeige ● ≈N GB wird mit dem <b>freien</b> Arbeitsspeicher verglichen (er steht über der Liste).</li>" +
			"<li><b>Welches Quant:</b> die Zahl sind Bits pro Gewicht (Q4 — der gute Mittelweg, Q8 — fast unkomprimiert, Q3 — spart Speicher auf Kosten der Qualität); K_M ist besser als K_S; IQ4 ist die neuere Generation und bei gleicher Größe besser als die klassischen. Die Anzeige ● ≈N GB schätzt den nötigen Speicher (Datei plus Reserve für den Kontext): grün passt, gelb ist knapp, rot passt nicht.</li>" +
			"<li>Ein Modell mit 1,5–3B bearbeitet schnell; 7–9B ist merklich klüger, braucht auf der CPU aber Sekunden pro Durchgang. Der LLM-Server startet beim ersten Einsatz und hält das Modell warm.</li>" +
			"</ul>" +
			"<p class=\"wh\">Prompts</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Aufräumen</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Geschäftsstil</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Ein Prompt ist eine Anweisung an das Bearbeitungsmodell. Zwei sind von Anfang an dabei: „Aufräumen“ (entfernt Füllwörter, Wiederholungen und Fehlstarts, richtet die Satzzeichen) und „Geschäftsstil“ (schreibt höflich und formell um); eigene können Sie beliebig anlegen.</li>" +
			"<li>Angehakte Prompts gelten für jedes Diktat, der Reihe nach von oben nach unten (als Kette: die Ausgabe des einen ist die Eingabe des nächsten); ist nichts angehakt, wird der Text so eingefügt, wie er erkannt wurde.</li>" +
			"<li>Der Stift ✎ öffnet den Prompt-Editor: Name, Text und ein Testfeld, das eine Probe direkt aus den Einstellungen durch das laufende Modell schickt. Die Reihenfolge ändert man, indem man den Prompt am Griff links zieht.</li>" +
			"<li>Tipp: kleine Modelle arbeiten deutlich besser, wenn im Prompt Beispiele „Eingabe → Ausgabe“ stehen — alle mitgelieferten sind so geschrieben.</li>" +
			"<li>Scheitert ein Profil (das Modell hat nicht geantwortet), wird der Text ohne es eingefügt: die Leiste zeigt „Eingefügt ohne das Profil …“, und Enter wird in diesem Fall nicht gedrückt.</li>" +
			"</ul>" +
			"<p class=\"wh\">Abhängigkeiten</p>" +
			"<ul>" +
			"<li>Prompts brauchen ein installiertes Bearbeitungsmodell; die Übersetzung braucht es nicht — die macht Whisper allein.</li>" +
			"<li>Das Bearbeitungsmodell wird beim ersten Einsatz geladen und bleibt warm; große Modelle sind auf der CPU spürbar langsamer.</li>" +
			"<li>Schauen Sie vor dem Laden auf die Speicheranzeige: ein „knappes“ Modell bremst das ganze System.</li>" +
			"<li>Ausgegraute Bedienelemente sind Einstellungen, die im aktuellen Modus nichts tun.</li>" +
			"</ul>" +
			"<p class=\"wh\">Installation und Portabilität</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — das Installationsprogramm: ohne Administratorrechte, Verknüpfung im Startmenü, Autostart auf Wunsch, sauberes Entfernen über die Windows-Einstellungen.</li>" +
			"<li>Das Installationsprogramm lädt standardmäßig nichts: der Assistent wählt das Modell beim ersten Start und lädt es. Wird doch eines gewählt — GigaAM v3 für Russisch, Whisper für alle anderen Sprachen —, lässt sich der Download mit einer Taste anhalten, und die Installation läuft trotzdem zu Ende. Dort steht auch der Schalter „Nach Updates suchen“, und die Antwort landet in den Einstellungen der App.</li>" +
			"<li><b>Portabel</b> — kopieren Sie einfach den ganzen Ordner mit der exe (auf einen USB-Stick, an einen anderen Rechner): Einstellungen, Modelle und Protokoll liegen daneben und reisen mit. In die Registry wird nichts geschrieben.</li>" +
			"<li>Ist beim ersten Start kein Erkennungsmodell da, öffnet das Programm den Katalog selbst und wartet auf den Download.</li>" +
			"<li>Voraussetzungen: Windows 10/11 x64, eine CPU mit AVX2 (etwa ab 2013), WebView2 Runtime für das Einstellungsfenster (in Windows 11 enthalten).</li>" +
			"</ul>" +
			"<p class=\"wh\">Menü im Infobereich und Dateien</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Bereit…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Einstellungen…</div><div class=\"mock-mi\">Ausschalten</div><div class=\"mock-mi\">Letztes Ergebnis kopieren</div><hr class=\"mock-sep\"><div class=\"mock-mi\">config.json öffnen</div><div class=\"mock-mi\">Protokoll öffnen</div><div class=\"mock-mi\">Über</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Beenden</div></div>" +
			"<ul>" +
			"<li>Linksklick auf das Symbol im Infobereich — die Einstellungen; Rechtsklick — das Menü. Farben des Symbols: grün — bereit, rot — Aufnahme, orange — Erkennung, grau — ausgeschaltet oder Fehler.</li>" +
			"<li><b>config.json</b> — alle Einstellungen; von Hand geänderte Werte gelten nach <b>Neu einlesen</b> im Bereich „System“. Dort stehen auch „Protokoll öffnen“ und „Einstellungen zurücksetzen“: das Zurücksetzen stellt den Werkszustand her und lässt Modelle, Verlauf und Prompts unberührt.</li>" +
			"<li><b>{log}</b> — das Protokoll, automatisch auf etwa 2 MB begrenzt.</li>" +
			"<li><b>models/</b> — die geladenen Erkennungs- und Bearbeitungsmodelle.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Zwei Minuten Einrichtung",
		"S_WIZ_HELLO_TEXT": "{app} macht aus Ihrer Stimme Text direkt an der Eingabemarke: Tastenkürzel halten, Satz sprechen, loslassen — der Text steht da. Alles läuft auf Ihrem Rechner, der Ton verlässt ihn nie.",
		"S_WIZ_LATER":      "Alles, was wir jetzt wählen, lässt sich später in den Einstellungen ändern.",
		"S_WIZ_T_MODEL":    "Sprache und Modell",
		"S_WIZ_MODEL_TEXT": "Sagen Sie mir, in welcher Sprache Sie diktieren, das Modell suche ich aus. Russisch übernimmt GigaAM, alle anderen Sprachen Whisper.",
		"S_WIZ_T_INPUT":    "Tastenkürzel und Mikrofon",
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
		"S_TITLE": "{app} — Réglages", "S_DICT_HINT": "Termes, noms et abréviations séparés par des virgules — un indice pour l'oreille, pas des commandes. Cela vaut pour Whisper ; le russe via GigaAM l'ignore. Le jeu par défaut suit la langue de reconnaissance jusqu'à ce que vous écriviez le vôtre.",
		"S_TR_DEFAULT": "Changer la langue de sortie du texte", "S_TR_TARGET": "Langue de sortie du texte par défaut", "S_TR_ASK": "Demander la langue de sortie du texte", "S_TR_ASK_NEVER": "Ne pas demander — traduire aussitôt",
		"S_SRCLANG_SUB": "vous la parlez ; elle détermine le modèle de reconnaissance",
		"S_TR_LANGS_SUB": "ces langues deviennent des boutons sur la plaque à l’insertion",
		"S_TR_UNAVAIL": "indisponible — %s ne sait pas traduire",
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
		"S_SEC_CMD": "Commandes vocales", "S_CMD_HINT": "Ce que vous dites devient un saut de ligne, un signe ou une annulation au lieu d'atterrir dans le texte. Reconnues en mots entiers, appliquées de haut en bas, après les remplacements.",
		"S_CMD_ADD": "Ajouter une commande", "S_CMD_PRESET": "Ajouter les habituelles", "S_CMD_PH": "la phrase que vous direz",
		"S_CMD_NEWLINE": "saut de ligne", "S_CMD_PARAGRAPH": "nouveau paragraphe", "S_CMD_TEXT": "insérer du texte", "S_CMD_CANCEL": "annuler la dictée",
		"S_CMD_TEXT_PH": "quoi insérer", "S_CMD_EMPTY": "Aucune commande pour l'instant", "S_CMD_DEL": "Supprimer la commande",
		"S_CMD_P_NEWLINE": "nouvelle ligne", "S_CMD_P_PARAGRAPH": "nouveau paragraphe", "S_CMD_P_CANCEL": "annuler",
		"S_SEC_REPLACE": "Remplacements après la reconnaissance", "S_REPLACE_HINT": "Ce qui a été mal entendu devient ce que vous vouliez dire — juste après la reconnaissance, avant les prompts. Appliqués de haut en bas.",
		"S_REPL_WHOLE_FULL": "Mots entiers uniquement", "S_REPL_CASE_FULL": "Respecter la casse", "S_CMD_ACTION": "Action",
		"S_FM_ADD": "Ajouter",
		"S_TIP_REPL_LANG": "La règle ne s’applique que lorsque vous dictez dans la langue choisie. « toutes les langues » — elle s’applique toujours.",
		"S_TIP_REPL_CASE": "Majuscules et minuscules comptent : « git » et « Git » sont des mots différents. Désactivé — la casse est ignorée.",
		"S_TIP_REPL_WHOLE": "Le remplacement ne s’applique que si le texte forme un mot à part. Désactivé — il agit aussi à l’intérieur d’autres mots.",
		"S_TIP_CMD_ACTION": "Ce qui se passe quand vous prononcez la phrase : saut de ligne, nouveau paragraphe, votre texte ou annulation de la dictée.",
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
		"S_OV_FREE": "Place à vous", "S_OV_FREE_SUB": "la plaque se déplace où vous voulez",
		"S_OVPOS_DRAG_SUB": "faites glisser la plaque à la souris — elle va où vous voulez",
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
		"S_DL": "Télécharger", "S_DEL": "Supprimer",
		"S_M_BASE": "rapide, PC modestes", "S_M_SMALL": "équilibré", "S_M_MED": "plus précis, recommandé", "S_M_TURBO": "précision max sur CPU", "S_M_PARAKEET": "25 langues européennes, ponctue d'elle-même",
		"S_THREADS":  "Threads CPU", "S_MINMS": "Enreg. min, ms", "S_MAXSEC": "Enreg. max, s",
		"S_AUTOSTART": "Démarrer whisper-server automatiquement", "S_PORT": "Port", "S_SERVEREXE": "Chemin de whisper-server", "S_SERVEREXE_SUB": "rempli automatiquement ; à changer seulement si le serveur est ailleurs", "S_EXE_RESET": "Réinitialiser", "S_EXE_WARN": "L'application trouve whisper-server à côté d'elle. Avec un chemin écrit à la main, déplacer le dossier empêchera la reconnaissance de démarrer. Modifier ?", "S_RESET_ALL": "Réinitialiser les réglages", "S_RESET_ALL_SUB": "tout revient d'usine, sauf les modèles et l'historique", "S_RESET_ALL_BTN": "Réinitialiser", "S_RESET_ALL_ASK": "Remettre tous les réglages d'usine ? Les modèles, l'historique et les prompts restent en place.", "S_RELOAD_CFG": "Relire config.json", "S_RELOAD_CFG_SUB": "si vous avez modifié le fichier à la main", "S_RELOAD_CFG_BTN": "Relire", "S_UPD_FOUND": "La version %s est sortie", "S_THEME": "Couleur", "S_THEME_SUB": "la couleur de la fenêtre, du bandeau et de l'icône", "S_THEME_GREEN": "Vert", "S_THEME_AMBER": "Ambre", "S_THEME_BLUE": "Bleu", "S_THEME_PINK": "Rose", "S_THEME_EDITOR": "Éditeur", "S_THEME_NEON": "Néon", "S_WND_MAX": "Occuper l'écran", "S_WND_RESTORE": "Taille précédente", "S_WND_MIN": "Réduire dans la zone de notification", "S_WND_CLOSE": "Fermer la fenêtre", "S_SKIN": "Design", "S_SKIN_SUB": "police, formes, effets et animations", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Doux", "S_SKIN_PAPER": "Document",
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
		"S_DEP_LLAMA": "post-traitement, modèles GGUF",
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
		"S_PRESETS_HINT": "Cliquez sur une langue — le choix des modèles s’ouvre dessous. Les langues sans modèle propre utilisent celui de la détection automatique.",
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
		"S_TR_EN": "traduit vers l’anglais",
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
		"S_OVPOS_SCHEME_SUB": "cliquez sur l’écran — la plaque s’y place",
		"S_OVDRAG": "Glissez-la où vous voulez",
		"S_OVMON": "Écran",
		"S_OVMON_SUB": "sur quel moniteur afficher la vignette",
		"S_OVMON_CURSOR": "L’écran avec le curseur",
		"S_M_NEMOTRON": "écrit pendant que vous parlez : le texte apparaît sur la vignette en direct ; 40 langues, ponctue d’elle-même",
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
		"S_NOT_INSTALLED": "non installÃ©",
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
		"S_SUB_TRTARGET": "le texte y est traduit ; le dialogue sur la plaque la propose en premier",
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
		"S_WIZ_NEED_MODEL": "Téléchargez d'abord un modèle — sans lui il n'y a rien pour reconnaître",
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
		"S_HELP_HTML": "<p class=\"wh\">Comment ça marche</p>" +
			"<p>Maintenez le raccourci — l'enregistrement commence (le bandeau en bas de l'écran montre votre niveau). Relâchez — l'audio est reconnu, traduit si besoin, passé dans les prompts, et le texte final arrive à l'endroit du curseur. Le ✕ du bandeau annule à n'importe quelle étape.</p>" +
			"<p>Le chemin complet : <b>enregistrement → reconnaissance (GigaAM pour le russe, Whisper sinon) → traduction (si activée) → prompts (LLM) → collage</b>. Chaque étape est visible sur le bandeau.</p>" +
			"<p class=\"wh\">Premier lancement</p>" +
			"<p>Le tout premier lancement ouvre un assistant en cinq étapes : la langue de l'interface, la langue de dictée (il choisit et télécharge le modèle), le raccourci et le microphone avec une barre de niveau, un champ pour essayer une dictée et, enfin, le démarrage avec Windows. Vous pouvez le passer, tout fonctionne quand même ; <b>{exe} -wizard</b> le rappelle. Une mise à jour ne l'affiche pas.</p>" +
			"<p class=\"wh\">Le bandeau</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Parlez…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Parlez…</b> — enregistrement : un point rouge et les barres de niveau.</li>" +
			"<li><b>Reconnaissance…</b> — Whisper travaille ; pendant la traduction — « Traduction », pendant les prompts — « Édition : nom (1/2) ».</li>" +
			"<li><b>Inséré : N caractères</b> — terminé ; en cas d'erreur ou de silence, la raison s'affiche brièvement.</li>" +
			"<li>Le ✕ à droite annule à n'importe quelle étape ; le bandeau ne prend jamais le focus. Le bandeau et son animation se désactivent dans « Dictée ».</li>" +
			"<li>Où le bandeau apparaît — en bas, en haut ou au curseur — et s'il affiche le texte reconnu au lieu d'un nombre de caractères, se règle dans « Dictée ».</li>" +
			"<li>Tant que le bandeau pose une question, sa ligne du haut le dit — « En attente de votre réponse » — et le point cesse de clignoter. Chaque réponse porte son numéro : 1…9 en choisissent une, Entrée prend celle qui est mise en avant, Échap annule tout ; les touches sont rappelées à droite, sur la même ligne. Dix secondes avant la limite d'enregistrement, un compte à rebours ambre s'affiche sur le bandeau.</li>" +
			"<li>La barre de titre porte trois boutons : réduire dans la zone de notification, occuper l'écran et fermer. Le même bouton ramène la fenêtre à sa taille précédente, et la taille réglée à la souris est conservée. La fenêtre ne descend pas sous 760×500.</li>" +
			"<li>Les noms longs — appareil, modèle, fichier — sont coupés par des points de suspension sur les cartes de l'État, pour que les cartes s'alignent ; le nom complet apparaît en info-bulle si le pointeur reste sur la carte. Les info-bulles sont dessinées aux couleurs de l'apparence en cours, pas à celles du système.</li>" +
			"<li>L'apparence vient de deux listes dans la section « Système ». « Design » donne la police, les formes, l'épaisseur des bordures, le halo et le caractère des animations ; il y en a trois : « Terminal » (vert, par défaut), « Éditeur » (gris plat, sans halo) et « Néon » (violet, arrondi). « Couleur » n'est proposée qu'au Terminal et ne change que la couleur de la fenêtre, du bandeau et de l'icône : vert, ambre, bleu, rose. Les autres designs ont leurs propres couleurs. Le choix s'applique aussitôt, sans redémarrage.</li>" +
			"</ul>" +
			"<p class=\"wh\">La question de la traduction</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Reconnaissance…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Traduire vers :</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Sans traduction</span></div></div>" +
			"<p>C'est le bandeau lui-même qui demande, sur une deuxième ligne, dès que vous relâchez le raccourci — dans les modes « demander à chaque fois » et « demander avec délai ». Les boutons viennent de « Langues du dialogue » ; la langue cible est mise en avant. Avec un délai, un trait se raccourcit sous ce bouton : à la fin, la langue mise en avant s'applique. <b>Sans traduction</b> insère le texte tel qu'il a été entendu ; le ✕ du bandeau annule tout. Le clavier fonctionne aussi : Entrée prend la réponse mise en avant, 1…9 choisissent un bouton, Échap annule.</p>" +
			"<p class=\"wh\">Insertion sûre</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Reconnaissance…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Le focus a changé — insérer ?</span><span class=\"mock-btn on mock-cd\">Insérer ici</span><span class=\"mock-btn\">Copier</span></div></div>" +
			"<ul>" +
			"<li>La fenêtre cible est retenue au moment où vous appuyez sur le raccourci. Si le focus a changé pendant le traitement, rien n'est collé — le bandeau demande sur sa deuxième ligne : <b>Insérer ici</b> (dans la fenêtre courante), <b>Copier</b> (dans le presse-papiers) ou le ✕. À la fin du décompte l'insertion est annulée et le texte reste dans le dernier résultat.</li>" +
			"<li>Entrée après le collage n'est envoyée que si la fenêtre cible n'a pas changé.</li>" +
			"<li><b>Dernier résultat</b> — le texte final de chaque dictée reste en mémoire jusqu'à la suivante ; le menu de la zone de notification propose « Copier le dernier résultat ». Un collage manqué ou un changement de focus ne fait jamais perdre une dictée.</li>" +
			"</ul>" +
			"<p class=\"wh\">Vérifier le microphone</p>" +
			"<p>Le bouton « Test » de la section Microphone enregistre trois secondes et les décortique : niveau de crête en décibels, part de l'enregistrement qui contient vraiment de la parole, et part d'échantillons écrêtés. La réponse est en mots : bon signal, trop faible — montez le niveau dans Windows, saturation — baissez-le, aucune parole entendue — est-ce le bon micro. Les mêmes mesures sont faites après chaque dictée et écrites dans le journal ; si la reconnaissance revient vide, le bandeau nomme la raison — trop faible, saturation ou silence — au lieu de dire simplement qu'il n'a rien entendu.</p>" +
			"<p class=\"wh\">Coller depuis l'historique</p>" +
			"<p>Chaque entrée de l'historique a un bouton « Coller » : il ramène la fenêtre depuis laquelle vous avez ouvert les réglages et y colle le texte, comme une dictée ordinaire. S'il n'y a nulle part où revenir, le texte est simplement placé dans le presse-papiers et le programme le dit.</p>" +
			"<p class=\"wh\">Les listes dans un seul fichier</p>" +
			"<p>Les remplacements et les commandes vocales peuvent être enregistrés dans un fichier .json et chargés sur un autre ordinateur — les boutons sous la liste des commandes, section « Texte ». Le chargement n'écrase rien : seules les lignes absentes sont ajoutées, et le programme indique combien ont été ajoutées et combien ignorées.</p>" +
			"<p class=\"wh\">Intégrité des fichiers</p>" +
			"<p>Chaque modèle du catalogue a une empreinte SHA-256 de référence. Après le téléchargement, le fichier lui est comparé : s'il ne correspond pas, il est supprimé et le téléchargement peut être repris. Le bouton « Vérifier » de la section « Modèles » compare de la même façon les modèles déjà installés, et lors d'une mise à jour l'installateur téléchargé est vérifié aussi — un fichier étranger ne sera pas lancé.</p>" +
			"<p class=\"wh\">Historique des dictées</p>" +
			"<p>La section « Historique » dans la colonne de gauche conserve ce que vous avez dicté : le texte seul, sur cet ordinateur seulement, l'audio n'est jamais conservé. Désactivée par défaut, elle s'active d'un interrupteur au même endroit. Les entrées restent un nombre de jours et jusqu'à un nombre d'entrées réglables, les plus anciennes disparaissent d'elles-mêmes ; « Ne jamais enregistrer depuis ces programmes » liste, séparés par des virgules, ceux dont rien ne doit être conservé — gestionnaires de mots de passe, applications bancaires. La recherche porte sur le texte et sur le nom du programme, le bouton à côté d'une entrée la met dans le presse-papiers, et « Vider » supprime tout d'un coup avec le fichier <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Commandes vocales</p>" +
			"<p>Sous les remplacements, dans la section Texte, il y a une liste de commandes : ce que vous dites devient une action au lieu de mots. « Nouvelle ligne » et « nouveau paragraphe » posent un saut — les modèles ne le font jamais ; « annuler » jette toute la dictée sans rien insérer ; « insérer du texte » place ce que vous voulez, un émoticône compris. Le bouton à côté de la liste la remplit des formules habituelles dans la langue de l'interface. Les commandes sont reconnues en mots entiers et s'appliquent après les remplacements, si bien que les prompts et la traduction reçoivent déjà le texte fini. Les espaces en trop autour des sauts disparaissent d'eux-mêmes. Le champ du dessous essaie remplacements et commandes sur n'importe quelle phrase : un saut s'affiche comme ⏎.</p>" +
			"<p class=\"wh\">Remplacements après la reconnaissance</p>" +
			"<p>Dans « Texte », vous pouvez lister ce que le modèle entend mal et ce que cela doit devenir : « git hub » → GitHub, des noms propres, des termes maison. Les remplacements s'appliquent juste après la reconnaissance, avant les prompts, pour que l'éditeur reçoive déjà les bons mots. La traduction vers l'anglais se fait pendant la reconnaissance : les remplacements voient donc le texte déjà traduit. Par défaut ils visent les mots entiers et ignorent la casse ; les deux interrupteurs à côté changent cela. Les règles s'appliquent de haut en bas. Le champ du dessous les essaie sur n'importe quelle phrase, sans dicter.</p>" +
			"<p class=\"wh\">Règles par application</p>" +
			"<p>Dans « Dictée », vous pouvez définir des règles pour certains programmes : avec quoi insérer (presse-papiers ou caractère par caractère), s'il faut appuyer sur Entrée, combien de temps attendre avant d'insérer et quels prompts appliquer. Un programme se désigne par son fichier — <b>chrome.exe</b> ; une règle peut en lister plusieurs séparés par des virgules, et une astérisque finale attrape tous les noms qui commencent ainsi. La première règle qui correspond l'emporte ; sans règle, ou sans correspondance, tout fonctionne comme dans les réglages généraux. Le bouton à côté de la liste inscrit le programme où vous avez inséré en dernier.</p>" +
			"<p class=\"wh\">Dictée</p>" +
			"<ul>" +
			"<li><b>Raccourci clavier</b> — le raccourci principal. N'importe quelle combinaison peut être capturée ; les modificateurs gauche et droit sont distingués. Les raccourcis de dictée, de traduction et de profils doivent être uniques — un doublon empêche l'enregistrement.</li>" +
			"<li><b>Mode</b> — maintenir les touches, ou appuyer une fois pour lancer et une fois pour arrêter.</li>" +
			"<li><b>Langue de l'interface</b> — change aussitôt ; « Comme le système » suit Windows.</li>" +
			"<li><b>Langue de reconnaissance</b> — une indication pour Whisper ; « auto » la devine à l'oreille.</li>" +
			"<li><b>Son</b> — signaux de début et de fin : plusieurs jeux plus les sons système de Windows, ▶ les fait écouter.</li>" +
			"<li><b>Entrée après le collage</b> — envoie aussitôt le texte dicté (pratique dans les messageries).</li>" +
			"<li><b>Restaurer le presse-papiers</b> — remet l'ancien contenu entièrement, images, fichiers et texte enrichi compris. Quand le contenu ne peut pas être sauvegardé, le presse-papiers reste intact et le texte est tapé caractère par caractère.</li>" +
			"<li><b>Bandeau et animation</b> — l'indicateur en bas de l'écran ; l'animation peut être coupée.</li>" +
			"<li><b>Insertion caractère par caractère</b> — au lieu de Ctrl+V, les frappes sont simulées, pour les champs qui refusent le collage.</li>" +
			"</ul>" +
			"<p class=\"wh\">Reconnaissance</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Détection automatique</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Russe</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · comme la détection automatique — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Quel modèle pour quelle langue</b> — l’onglet « Langues et modèles » est une liste de langues. Cliquez-en une — les modèles qui la maîtrisent se déploient dessous : l’attribué et le recommandé d’abord, les absents avec leur taille et une flèche de téléchargement. Un clic sur la carte est le choix ; un modèle absent se télécharge tout seul et prend le relais une fois prêt. Les langues sans modèle propre héritent de celui de la détection automatique et s’affichent estompées.</li>" +
			"<li><b>Le catalogue</b> — Whisper : Base (rapide, pour PC modestes), Small (l’équilibre), Medium et Turbo (plus précis et plus lents ; « q5 » est la version quantifiée : un peu plus petite et rapide, presque sans perte), qui traduisent aussi vers l’anglais ; GigaAM v3 est plus précis en russe et ponctue tout seul ; Parakeet v3 couvre 25 langues européennes ; Nemotron 3.5 tape pendant que vous parlez. Les téléchargements viennent des dépôts officiels Hugging Face, chaque fichier vérifié contre son hachage de référence.</li>" +
			"<li><b>Votre propre modèle</b> — un fichier Whisper unique ggml-*.bin ou un dossier de modèle sherpa-onnx (encoder.onnx, decoder.onnx, tokens.txt) conviennent. Placez-le dans le dossier models à côté de l’application et redémarrez-la — le modèle apparaît dans le choix des langues compatibles ; ses capacités étant inconnues, il s’affiche honnêtement, sans barres.</li>" +
			"<li>Le serveur de reconnaissance garde le modèle en mémoire entre les phrases — la première dictée après le démarrage est plus lente (chargement), ensuite la reconnaissance prend une à trois secondes.</li>" +
			"<li><b>Dictionnaire</b> — termes, noms et abréviations séparés par des virgules. Une indication pour « l'oreille » de Whisper afin que les mots rares passent correctement ; ce ne sont pas des commandes.</li>" +
			"<li><b>Microphone</b> — choix de l'appareil avec un niveau en direct (parlez et la barre bouge : l'appareil est bien entendu). Si l'appareil choisi est débranché, celui du système prend le relais ; un enregistrement sans parole n'est jamais envoyé à la reconnaissance — le bandeau annonce « Silence ».</li>" +
			"<li><b>Service</b> — le serveur de reconnaissance démarre tout seul et tourne en local. Le port, le chemin ou un serveur distant se changent ; la reconnaissance redémarre ensuite d'elle-même.</li>" +
			"<li><b>Traduction</b> — tout passe par Whisper : vers l'anglais dans son mode natif, vers les autres langues <b>à titre expérimental</b>, en forçant la langue de sortie (la qualité dépend du couple de langues ; les grandes langues s'en sortent le mieux). Le modèle Turbo n'est pas entraîné pour cela — les réglages préviennent tant qu'il est actif. « Toujours traduire vers la langue cible » traduit chaque dictée sans rien demander. Sans cette case, le mode de question s'applique : toujours ou avec délai — le dialogue de langue apparaît avant la reconnaissance et, le temps écoulé, la cible est retenue. Le raccourci de traduction séparé traduit une seule fois, sans toucher à la dictée normale.</li>" +
			"</ul>" +
			"<p class=\"wh\">Post-traitement (LLM)</p>" +
			"<p>Une deuxième couche facultative : un modèle de langue local (llama.cpp) retouche le texte reconnu selon vos prompts — enlève les mots parasites, change le style, met en forme. Entièrement hors ligne, sur le processeur seulement.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modèles</b> — les modèles d'édition installés ; le bouton radio choisit l'actif (appliqué aussitôt), le ✕ supprime (l'actif aussi — le post-traitement s'éteint alors). La progression des téléchargements s'affiche ici également.</li>" +
			"<li><b>Recherche</b> — les modèles GGUF sur Hugging Face par nom (par exemple « qwen2.5 instruct »). Chaque dépôt indique sa date de mise à jour, son nombre de téléchargements et un ↗ vers sa page ; un clic sur la ligne déplie ses fichiers de quantification. L'indicateur ● ≈N GB se compare à la mémoire <b>libre</b> (affichée au-dessus de la liste).</li>" +
			"<li><b>Quelle quantification :</b> le chiffre correspond aux bits par poids (Q4 — le bon compromis, Q8 — presque non compressé, Q3 — économise la mémoire au prix de la qualité) ; K_M vaut mieux que K_S ; IQ4 est la génération plus récente, meilleure à taille égale. L'indicateur ● ≈N GB estime la mémoire nécessaire (le fichier plus une réserve pour le contexte) : vert ça tient, ambre c'est juste, rouge ça ne tient pas.</li>" +
			"<li>Un modèle de 1,5 à 3 milliards de paramètres édite vite ; 7 à 9 milliards est nettement plus fin mais prend des secondes par passage sur le processeur. Le serveur LLM démarre au premier usage et garde le modèle au chaud.</li>" +
			"</ul>" +
			"<p class=\"wh\">Prompts</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Nettoyage</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Style professionnel</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Un prompt est une consigne pour le modèle d'édition. Deux sont livrés d'origine : « Nettoyage » (retire les mots parasites, les répétitions et les faux départs, remet la ponctuation) et « Style professionnel » (réécrit poliment et formellement) ; ajoutez les vôtres librement.</li>" +
			"<li>Les prompts cochés s'appliquent à chaque dictée, dans l'ordre, de haut en bas (en chaîne : la sortie de l'un devient l'entrée du suivant) ; rien de coché — le texte est inséré tel qu'il a été reconnu.</li>" +
			"<li>Le crayon ✎ ouvre l'éditeur du prompt : nom, texte et un champ d'essai qui passe un exemple par le modèle en marche, depuis les réglages. L'ordre se change en faisant glisser le prompt par la poignée à gauche.</li>" +
			"<li>Astuce : les petits modèles travaillent bien mieux quand le prompt contient des exemples « entrée → sortie » — tous ceux livrés sont écrits ainsi.</li>" +
			"<li>Si un profil échoue (le modèle n'a pas répondu), le texte est inséré sans lui : le bandeau affiche « Inséré sans le profil … » et Entrée n'est pas envoyée dans ce cas.</li>" +
			"</ul>" +
			"<p class=\"wh\">Dépendances</p>" +
			"<ul>" +
			"<li>Les prompts exigent un modèle d'édition installé ; la traduction n'en dépend pas — Whisper la fait seul.</li>" +
			"<li>Le modèle d'édition se charge au premier usage et reste au chaud ; les gros modèles sont nettement plus lents sur le processeur.</li>" +
			"<li>Regardez l'indicateur de mémoire avant de télécharger : un modèle « juste » ralentit tout le système.</li>" +
			"<li>Les contrôles grisés sont les réglages qui ne servent à rien dans le mode courant.</li>" +
			"</ul>" +
			"<p class=\"wh\">Installation et portabilité</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — l'installateur : sans droits administrateur, raccourci dans le menu Démarrer, démarrage automatique si vous le souhaitez, suppression propre depuis les paramètres de Windows.</li>" +
			"<li>L'installateur ne télécharge rien par défaut : l'assistant choisit et récupère le modèle au premier lancement. Si un modèle est tout de même choisi — GigaAM v3 pour le russe, Whisper pour les autres langues — le téléchargement peut être arrêté d'un bouton et l'installation va quand même à son terme. La case « Vérifier les mises à jour » est là aussi, et la réponse est écrite dans les réglages de l'application.</li>" +
			"<li><b>Portable</b> — copiez simplement tout le dossier contenant l'exe (sur une clé USB, vers un autre PC) : réglages, modèles et journal vivent à côté et voyagent avec. Rien n'est écrit dans le registre.</li>" +
			"<li>Au premier lancement sans modèle de reconnaissance, le programme ouvre lui-même le catalogue et attend le téléchargement.</li>" +
			"<li>Prérequis : Windows 10/11 x64, un processeur avec AVX2 (à partir de 2013 environ), WebView2 Runtime pour la fenêtre des réglages (inclus dans Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Zone de notification et fichiers</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Prêt…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Réglages…</div><div class=\"mock-mi\">Désactiver</div><div class=\"mock-mi\">Copier le dernier résultat</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Ouvrir config.json</div><div class=\"mock-mi\">Ouvrir le journal</div><div class=\"mock-mi\">À propos</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Quitter</div></div>" +
			"<ul>" +
			"<li>Clic gauche sur l'icône — les réglages ; clic droit — le menu. Couleurs de l'icône : vert — prêt, rouge — enregistrement, orange — reconnaissance, gris — désactivé ou erreur.</li>" +
			"<li><b>config.json</b> — tous les réglages ; les modifications à la main s'appliquent via <b>Relire</b> dans la section « Système ». On y trouve aussi « Ouvrir le journal » et « Réinitialiser les réglages » : la réinitialisation revient à l'état d'usine et laisse les modèles, l'historique et les prompts en place.</li>" +
			"<li><b>{log}</b> — le journal, limité automatiquement à environ 2 Mo.</li>" +
			"<li><b>models/</b> — les modèles de reconnaissance et d'édition téléchargés.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Deux minutes de configuration",
		"S_WIZ_HELLO_TEXT": "{app} transforme votre voix en texte à l'endroit du curseur : maintenez le raccourci, dites une phrase, relâchez — le texte est là. Tout tourne sur votre machine, le son n'en sort jamais.",
		"S_WIZ_LATER":      "Tout ce que nous choisissons ici se modifie ensuite dans les paramètres.",
		"S_WIZ_T_MODEL":    "Langue et modèle",
		"S_WIZ_MODEL_TEXT": "Dites-moi dans quelle langue vous dicterez, je choisis le modèle. Le russe passe par GigaAM, toutes les autres langues par Whisper.",
		"S_WIZ_T_INPUT":    "Raccourci et microphone",
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
		"S_TITLE": "{app} — Ajustes", "S_DICT_HINT": "Términos, nombres y abreviaturas separados por comas — una pista para el oído, no comandos. Vale para Whisper; el ruso a través de GigaAM lo ignora. El conjunto por defecto sigue al idioma de reconocimiento hasta que escribas el tuyo.",
		"S_TR_DEFAULT": "Cambiar el idioma de salida del texto", "S_TR_TARGET": "Idioma de salida del texto por defecto", "S_TR_ASK": "Preguntar el idioma de salida del texto", "S_TR_ASK_NEVER": "No preguntar — traducir enseguida",
		"S_SRCLANG_SUB": "usted lo habla; determina el modelo de reconocimiento",
		"S_TR_LANGS_SUB": "estos idiomas serán botones en la placa al insertar",
		"S_TR_UNAVAIL": "no disponible — %s no sabe traducir",
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
		"S_SEC_CMD": "Comandos de voz", "S_CMD_HINT": "Lo que dices se convierte en un salto de línea, un signo o una cancelación en vez de acabar en el texto. Se buscan como palabras completas y se aplican de arriba abajo, después de los reemplazos.",
		"S_CMD_ADD": "Añadir un comando", "S_CMD_PRESET": "Añadir los habituales", "S_CMD_PH": "la frase que dirás",
		"S_CMD_NEWLINE": "salto de línea", "S_CMD_PARAGRAPH": "nuevo párrafo", "S_CMD_TEXT": "insertar texto", "S_CMD_CANCEL": "cancelar el dictado",
		"S_CMD_TEXT_PH": "qué insertar", "S_CMD_EMPTY": "Todavía no hay comandos", "S_CMD_DEL": "Eliminar el comando",
		"S_CMD_P_NEWLINE": "nueva línea", "S_CMD_P_PARAGRAPH": "nuevo párrafo", "S_CMD_P_CANCEL": "cancelar",
		"S_SEC_REPLACE": "Reemplazos tras el reconocimiento", "S_REPLACE_HINT": "Lo que se oyó mal se convierte en lo que querías decir, justo tras el reconocimiento y antes de los prompts. Se aplican de arriba abajo.",
		"S_REPL_WHOLE_FULL": "Solo palabras completas", "S_REPL_CASE_FULL": "Distinguir mayúsculas y minúsculas", "S_CMD_ACTION": "Acción",
		"S_FM_ADD": "Añadir",
		"S_TIP_REPL_LANG": "La regla solo actúa cuando dictas en el idioma elegido. «todos los idiomas» — actúa siempre.",
		"S_TIP_REPL_CASE": "Las mayúsculas y minúsculas importan: «git» y «Git» son palabras distintas. Desactivado — se ignora.",
		"S_TIP_REPL_WHOLE": "El reemplazo solo actúa si el texto forma una palabra aparte. Desactivado — también dentro de otras palabras.",
		"S_TIP_CMD_ACTION": "Qué ocurre al decir la frase: salto de línea, nuevo párrafo, tu propio texto o cancelar el dictado.",
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
		"S_OV_FREE": "Sitio propio", "S_OV_FREE_SUB": "la placa se puede arrastrar a cualquier parte",
		"S_OVPOS_DRAG_SUB": "arrastra la placa con el ratón: se coloca donde quieras",
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
		"S_DL": "Descargar", "S_DEL": "Borrar",
		"S_M_BASE": "rápido, PCs modestos", "S_M_SMALL": "equilibrado", "S_M_MED": "más preciso, recomendado", "S_M_TURBO": "máx. precisión en CPU", "S_M_PARAKEET": "25 idiomas europeos, puntúa por sí solo",
		"S_THREADS":  "Hilos de CPU", "S_MINMS": "Grab. mín, ms", "S_MAXSEC": "Grab. máx, s",
		"S_AUTOSTART": "Iniciar whisper-server automáticamente", "S_PORT": "Puerto", "S_SERVEREXE": "Ruta de whisper-server", "S_SERVEREXE_SUB": "se rellena solo; cámbiela únicamente si el servidor está en otro sitio", "S_EXE_RESET": "Restablecer", "S_EXE_WARN": "La aplicación encuentra whisper-server junto a ella. Con una ruta escrita a mano, al mover la carpeta el reconocimiento dejará de arrancar. ¿Cambiarla?", "S_RESET_ALL": "Restablecer los ajustes", "S_RESET_ALL_SUB": "todo vuelve de fábrica, salvo modelos e historial", "S_RESET_ALL_BTN": "Restablecer", "S_RESET_ALL_ASK": "¿Devolver todos los ajustes de fábrica? Los modelos, el historial y los prompts se quedan.", "S_RELOAD_CFG": "Releer config.json", "S_RELOAD_CFG_SUB": "si editó el archivo a mano", "S_RELOAD_CFG_BTN": "Releer", "S_UPD_FOUND": "Está la versión %s", "S_THEME": "Color", "S_THEME_SUB": "el color de la ventana, la barra y el icono", "S_THEME_GREEN": "Verde", "S_THEME_AMBER": "Ámbar", "S_THEME_BLUE": "Azul", "S_THEME_PINK": "Rosa", "S_THEME_EDITOR": "Editor", "S_THEME_NEON": "Neón", "S_WND_MAX": "Ocupar la pantalla", "S_WND_RESTORE": "Volver al tamaño anterior", "S_WND_MIN": "Ocultar en la bandeja", "S_WND_CLOSE": "Cerrar la ventana", "S_SKIN": "Diseño", "S_SKIN_SUB": "tipografía, formas, efectos y animación", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Suave", "S_SKIN_PAPER": "Documento",
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
		"S_DEP_LLAMA": "posprocesado, modelos GGUF",
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
		"S_PRESETS_HINT": "Haga clic en un idioma — debajo se abre la elección de modelos para él. Los idiomas sin modelo propio usan el de la detección automática.",
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
		"S_TR_EN": "traduce al inglés",
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
		"S_OVPOS_SCHEME_SUB": "haz clic en la pantalla: la placa se coloca ahí",
		"S_OVDRAG": "Arrástrela a donde quiera",
		"S_OVMON": "Pantalla",
		"S_OVMON_SUB": "en qué monitor mostrar la placa",
		"S_OVMON_CURSOR": "La pantalla con el cursor",
		"S_M_NEMOTRON": "escribe mientras habla: el texto aparece en la placa en directo; 40 idiomas, puntúa por sí solo",
		"S_M_TINY": "el más pequeño y rápido, para equipos muy modestos; bastante menos preciso",
		"S_STATE_LOADED": "En memoria ahora mismo",
		"S_STATE_LOADED_SUB": "los modelos se descargan solos tras la inactividad",
		"S_STATE_WEEK": "Esta semana",
		"S_ST_SUMMARY": "Resumen", "S_ST_OVERLAY": "Franja en pantalla", "S_ST_BEEP": "Señal sonora", "S_ST_AUTORUN": "Inicio con Windows", "S_ST_POST": "Posprocesado", "S_ST_LOCAL": "local", "S_ST_CHECKED": "comprobado %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "activado", "S_ST_OFF_M": "desactivado", "S_ST_MIC_OK": "la señal es buena", "S_ST_MIC_BAD": "el micrófono no responde", "S_ST_CHECK": "Comprobar", "S_ST_RECOG": "reconocido por %s", "S_ST_VER": "Versión %s", "S_ST_LATEST": "la última", "S_ST_OUTDATED": "desactualizada", "S_ST_UPD_OK": "tiene la última versión", "S_ST_UPD_DL": "Descargando la actualización…",
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
		"S_SUB_TRTARGET": "el texto se traduce a él; el diálogo en la placa lo ofrece primero",
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
		"S_WIZ_NEED_MODEL": "Descarga primero un modelo: sin él no hay con qué reconocer",
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
		"S_HELP_HTML": "<p class=\"wh\">Cómo funciona</p>" +
			"<p>Mantén el atajo — empieza la grabación (la barra en la parte inferior de la pantalla muestra tu nivel). Suelta — el audio se reconoce, se traduce si hace falta, pasa por los prompts y el texto final aparece donde está el cursor. La ✕ de la barra cancela en cualquier momento.</p>" +
			"<p>El camino completo: <b>grabación → reconocimiento (GigaAM para el ruso, Whisper para el resto) → traducción (si está activa) → prompts (LLM) → pegado</b>. Cada paso se ve en la barra.</p>" +
			"<p class=\"wh\">Primer inicio</p>" +
			"<p>El primer inicio abre un asistente de cinco pasos: el idioma de la interfaz, el idioma en el que vas a dictar (él elige y descarga el modelo), el atajo y el micrófono con barra de nivel, un campo para probar un dictado y, por último, iniciar con Windows. Puedes omitirlo y todo sigue funcionando; <b>{exe} -wizard</b> lo trae de vuelta. Al actualizar no aparece.</p>" +
			"<p class=\"wh\">La barra</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Habla…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Habla…</b> — grabando: un punto rojo y las barras de nivel.</li>" +
			"<li><b>Reconociendo…</b> — Whisper trabaja; al traducir — «Traduciendo», al aplicar prompts — «Editando: nombre (1/2)».</li>" +
			"<li><b>Insertado: N caracteres</b> — listo; si hay error o silencio, se muestra el motivo en corto.</li>" +
			"<li>La ✕ de la derecha cancela en cualquier momento; la barra nunca roba el foco. La barra y su animación se apagan en «Dictado».</li>" +
			"<li>Dónde aparece la barra — abajo, arriba o en el cursor — y si muestra el texto reconocido en vez del número de caracteres, se ajusta en «Dictado».</li>" +
			"<li>Mientras la barra pregunta algo, su línea de arriba lo dice — «Esperando tu respuesta» — y el punto deja de parpadear. Cada respuesta lleva su número: 1…9 eligen una, Enter toma la resaltada, Esc cancela todo; las teclas están escritas a la derecha, en la misma fila. Diez segundos antes del límite de grabación, en la barra corre una cuenta atrás ámbar.</li>" +
			"<li>La barra de título lleva tres botones: ocultar en la bandeja, ocupar la pantalla y cerrar. El mismo botón devuelve la ventana al tamaño anterior, y el tamaño que fijaste con el ratón se conserva. La ventana no baja de 760×500.</li>" +
			"<li>Los nombres largos — dispositivo, modelo, archivo — se cortan con puntos suspensivos en las tarjetas de Estado para que queden alineadas; el nombre completo aparece como pista si el puntero se detiene sobre la tarjeta. Las pistas están dibujadas con los colores del aspecto actual, no con los del sistema.</li>" +
			"<li>El aspecto sale de dos listas en la sección «Sistema». «Diseño» fija la tipografía, las formas, el grosor de los bordes, el halo y el carácter de las animaciones; hay tres: «Terminal» (verde, por defecto), «Editor» (gris plano, sin halo) y «Neón» (violeta, redondeado). «Color» solo se ofrece al Terminal y cambia únicamente el color de la ventana, la barra y el icono: verde, ámbar, azul, rosa. Los demás diseños traen sus propios colores. La elección se aplica al momento, sin reiniciar.</li>" +
			"</ul>" +
			"<p class=\"wh\">La pregunta de la traducción</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Reconociendo…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Traducir a:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Sin traducción</span></div></div>" +
			"<p>Pregunta la propia barra, en una segunda línea, en cuanto sueltas el atajo — en los modos «preguntar siempre» y «preguntar con cuenta atrás». Los botones vienen de «Idiomas del diálogo»; el idioma de destino está resaltado. Con cuenta atrás, bajo ese botón se acorta una línea: cuando se agota, se aplica el idioma resaltado. <b>Sin traducción</b> inserta el texto tal como se oyó; la ✕ de la barra cancela toda la operación. El teclado también sirve: Intro toma la respuesta resaltada, 1…9 eligen un botón, Esc cancela.</p>" +
			"<p class=\"wh\">Inserción segura</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Reconociendo…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Cambió el foco, ¿insertar?</span><span class=\"mock-btn on mock-cd\">Insertar aquí</span><span class=\"mock-btn\">Copiar</span></div></div>" +
			"<ul>" +
			"<li>La ventana de destino se recuerda en el momento de pulsar el atajo. Si el foco cambió mientras se procesaba el habla, no se pega nada — la barra pregunta en su segunda línea: <b>Insertar aquí</b> (en la ventana actual), <b>Copiar</b> (al portapapeles) o la ✕. Si se agota la cuenta atrás, la inserción se cancela y el texto queda en el último resultado.</li>" +
			"<li>Enter tras pegar solo se envía si la ventana de destino no ha cambiado.</li>" +
			"<li><b>Último resultado</b> — el texto final de cada dictado se guarda en memoria hasta el siguiente; el menú del área de notificación tiene «Copiar el último resultado». Un pegado fallido o un cambio de foco nunca hacen perder un dictado.</li>" +
			"</ul>" +
			"<p class=\"wh\">Comprobar el micrófono</p>" +
			"<p>El botón «Test» de la sección Micrófono graba tres segundos y los analiza: pico en decibelios, qué parte de la grabación contiene voz de verdad y qué parte de las muestras se recortó. La respuesta llega en palabras: se oye bien, demasiado bajo — sube el nivel en Windows, saturación — bájalo, no se oye voz — ¿está elegido el micrófono correcto? Lo mismo se mide tras cada dictado y se escribe en el registro; si el reconocimiento vuelve vacío, la barra nombra el motivo — bajo, saturación o silencio — en vez de decir solo que no oyó nada.</p>" +
			"<p class=\"wh\">Pegar desde el historial</p>" +
			"<p>Cada entrada del historial tiene el botón «Pegar»: devuelve el foco a la ventana desde la que abrió los ajustes y pega allí el texto, como un dictado normal. Si no hay adónde volver, el texto se queda en el portapapeles y el programa lo avisa.</p>" +
			"<p class=\"wh\">Las listas en un archivo</p>" +
			"<p>Las sustituciones y las órdenes de voz pueden guardarse en un archivo .json y cargarse en otro ordenador: los botones bajo la lista de órdenes en la sección «Texto». La carga no borra nada: solo se añaden las líneas que faltan, y el programa dice cuántas se añadieron y cuántas se omitieron.</p>" +
			"<p class=\"wh\">Integridad de los archivos</p>" +
			"<p>Cada modelo del catálogo tiene un hash SHA-256 de referencia. Tras la descarga el archivo se compara con él: si no coincide, se borra y la descarga puede repetirse. El botón «Comprobar» de la sección «Modelos» compara igual los modelos ya instalados, y al actualizar el programa también se comprueba el instalador descargado: un archivo ajeno no se ejecutará.</p>" +
			"<p class=\"wh\">Historial de dictados</p>" +
			"<p>La sección «Historial» de la columna izquierda guarda lo que has dictado: solo el texto, solo en este equipo, el audio nunca se guarda. Está desactivada por defecto y se activa con un interruptor en el mismo sitio. Las entradas se conservan durante los días y hasta la cantidad que fijes, las viejas caen solas; «No registrar nunca desde estos programas» enumera, separados por comas, aquellos de los que no debe guardarse nada — gestores de contraseñas, banca. La búsqueda cubre el texto y el nombre del programa, el botón junto a una entrada la copia al portapapeles, y «Vaciar» borra todo de golpe junto con el archivo <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Comandos de voz</p>" +
			"<p>Bajo los reemplazos, en la sección Texto, hay una lista de comandos: lo que dices se convierte en una acción en lugar de en palabras. «Nueva línea» y «nuevo párrafo» ponen un salto — los modelos nunca lo hacen; «cancelar» descarta todo el dictado sin insertar nada; «insertar texto» coloca lo que quieras, incluso un emoticono. El botón junto a la lista la rellena con las frases habituales en el idioma de la interfaz. Los comandos se buscan como palabras completas y se aplican tras los reemplazos, así que los prompts y la traducción reciben ya el texto terminado. Los espacios sobrantes junto a los saltos se limpian solos. El campo de abajo prueba reemplazos y comandos con cualquier frase: un salto se muestra como ⏎.</p>" +
			"<p class=\"wh\">Reemplazos tras el reconocimiento</p>" +
			"<p>En «Texto» puedes enumerar lo que el modelo oye mal y en qué debe convertirse: «git hub» → GitHub, apellidos, términos internos. Los reemplazos se aplican justo después del reconocimiento, antes de los prompts, para que el editor ya reciba las palabras correctas. La traducción al inglés ocurre dentro del reconocimiento, así que los reemplazos ven el texto ya traducido. Por defecto buscan palabras completas y no distinguen mayúsculas; los dos interruptores de al lado lo cambian. Las reglas se aplican de arriba abajo. El campo de abajo las prueba con cualquier frase, sin dictar.</p>" +
			"<p class=\"wh\">Reglas por aplicación</p>" +
			"<p>En «Dictado» puedes fijar reglas para programas concretos: con qué insertar (portapapeles o carácter a carácter), si pulsar Intro, cuánto esperar antes de insertar y qué prompts aplicar. El programa se indica por su archivo — <b>chrome.exe</b>; una regla puede listar varios separados por comas, y un asterisco final atrapa todos los nombres que empiecen igual. Gana la primera regla que coincide; sin reglas, o si ninguna coincide, todo funciona como en los ajustes generales. El botón junto a la lista rellena el programa donde insertaste por última vez.</p>" +
			"<p class=\"wh\">Dictado</p>" +
			"<ul>" +
			"<li><b>Atajo de teclado</b> — el atajo principal. Se puede capturar cualquier combinación; los modificadores izquierdo y derecho se distinguen. Los atajos de dictado, traducción y perfiles deben ser únicos — uno repetido impide guardar.</li>" +
			"<li><b>Modo</b> — mantener las teclas, o pulsar una vez para empezar y otra para parar.</li>" +
			"<li><b>Idioma de la interfaz</b> — cambia al instante; «Como el sistema» sigue a Windows.</li>" +
			"<li><b>Idioma de reconocimiento</b> — una pista para Whisper; «auto» lo deduce del habla.</li>" +
			"<li><b>Sonido</b> — avisos de inicio y fin: varios juegos más los sonidos del sistema, ▶ los reproduce.</li>" +
			"<li><b>Enter tras pegar</b> — envía el texto dictado enseguida (cómodo en mensajería).</li>" +
			"<li><b>Restaurar el portapapeles</b> — devuelve el contenido anterior por completo, incluidas imágenes, archivos y texto con formato. Cuando no se puede guardar ese contenido, el portapapeles se deja intacto y el texto se escribe carácter a carácter.</li>" +
			"<li><b>Barra y animación</b> — el indicador en la parte inferior de la pantalla; la animación se puede apagar.</li>" +
			"<li><b>Insertar carácter a carácter</b> — en lugar de Ctrl+V se simulan pulsaciones, para los campos que no admiten pegar.</li>" +
			"</ul>" +
			"<p class=\"wh\">Reconocimiento</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Detección automática</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Ruso</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · como la detección automática — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Qué modelo para qué idioma</b> — la pestaña «Idiomas y modelos» es una lista de idiomas. Haga clic en uno — debajo se despliegan los modelos que lo dominan: el asignado y el recomendado primero, los ausentes con su tamaño y una flecha de descarga. Un clic en la tarjeta es la elección; un modelo ausente se descarga solo y entra en cuanto está listo. Los idiomas sin modelo propio heredan el de la detección automática y se muestran atenuados.</li>" +
			"<li><b>El catálogo</b> — Whisper: Base (rápida, para PC modestos), Small (el equilibrio), Medium y Turbo (más precisas y lentas; «q5» es la versión cuantizada: algo más pequeña y rápida casi sin pérdida), que además traducen al inglés; GigaAM v3 es más precisa en ruso y puntúa sola; Parakeet v3 cubre 25 idiomas europeos; Nemotron 3.5 escribe mientras usted habla. Las descargas vienen de los repositorios oficiales de Hugging Face, cada archivo cotejado con su hash de referencia.</li>" +
			"<li><b>Modelo propio</b> — sirve un archivo único Whisper ggml-*.bin o una carpeta de modelo sherpa-onnx (encoder.onnx, decoder.onnx, tokens.txt). Póngalo en la carpeta models junto a la aplicación y reiníciela — el modelo aparecerá en la elección de los idiomas compatibles; sus capacidades son desconocidas, así que se muestra honestamente, sin barras.</li>" +
			"<li>El servidor de reconocimiento mantiene el modelo en memoria entre frases — el primer dictado tras arrancar tarda más (carga) y luego el reconocimiento lleva de uno a tres segundos.</li>" +
			"<li><b>Diccionario</b> — términos, nombres y abreviaturas separados por comas. Una pista para el «oído» de Whisper, para que las palabras raras salgan bien; no son órdenes.</li>" +
			"<li><b>Micrófono</b> — elección del dispositivo con medidor de nivel (habla y la barra se mueve: el dispositivo se oye). Si desconectas el elegido, se usa el del sistema; una grabación sin voz no se manda a reconocer — la barra dice «Silencio».</li>" +
			"<li><b>Servicio</b> — el servidor de reconocimiento arranca solo y funciona en local. Puedes cambiar el puerto, la ruta o apuntar a un servidor remoto; el reconocedor se reinicia por su cuenta.</li>" +
			"<li><b>Traducción</b> — todo lo traduce Whisper: al inglés en su modo nativo, a otros idiomas <b>de forma experimental</b>, forzando el idioma de salida (la calidad depende del par de idiomas; los idiomas grandes salen mejor). El modelo Turbo no está entrenado para traducir — los ajustes avisan mientras esté activo. «Traducir siempre al idioma de destino» traduce cada dictado sin preguntar. Sin esa casilla manda el modo de pregunta: siempre o con cuenta atrás — el diálogo de idioma aparece antes del reconocimiento y, al agotarse el tiempo, se usa el destino. El atajo de traducción aparte traduce una sola vez, sin tocar el dictado normal.</li>" +
			"</ul>" +
			"<p class=\"wh\">Posprocesado (LLM)</p>" +
			"<p>Una segunda capa opcional: un modelo de lenguaje local (llama.cpp) retoca el texto reconocido según tus prompts — quita muletillas, cambia el estilo, da formato. Totalmente sin conexión, solo en el procesador.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modelos</b> — los modelos de edición instalados; el botón de opción elige el activo (se aplica al momento) y la ✕ borra (también el activo — entonces el posprocesado se apaga). El progreso de las descargas también se ve aquí.</li>" +
			"<li><b>Búsqueda</b> — modelos GGUF en Hugging Face por nombre (por ejemplo «qwen2.5 instruct»). Cada repositorio muestra su fecha de actualización, el número de descargas y un ↗ hacia su página; al pulsar la fila se despliegan sus archivos cuantizados. El indicador ● ≈N GB se compara con la memoria <b>libre</b> (aparece encima de la lista).</li>" +
			"<li><b>Qué cuantización elegir:</b> el número son bits por peso (Q4 — el término medio, Q8 — casi sin comprimir, Q3 — ahorra memoria a costa de calidad); K_M es mejor que K_S; IQ4 es la generación más nueva y gana a las clásicas a igualdad de tamaño. El indicador ● ≈N GB estima la memoria necesaria (el archivo más un margen para el contexto): verde cabe, ámbar va justo, rojo no cabe.</li>" +
			"<li>Un modelo de 1,5–3B edita rápido; uno de 7–9B es bastante más listo pero tarda segundos por pasada en el procesador. El servidor LLM arranca en el primer uso y mantiene el modelo caliente.</li>" +
			"</ul>" +
			"<p class=\"wh\">Prompts</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Limpieza</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Estilo formal</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Un prompt es una instrucción para el modelo editor. Vienen dos de fábrica: «Limpieza» (quita muletillas, repeticiones y arranques en falso, arregla la puntuación) y «Estilo formal» (reescribe con cortesía y formalidad); añade los tuyos sin límite.</li>" +
			"<li>Los prompts marcados se aplican a cada dictado en orden, de arriba abajo (en cadena: la salida de uno es la entrada del siguiente); si no hay ninguno marcado, el texto se inserta tal como se reconoció.</li>" +
			"<li>El lápiz ✎ abre el editor del prompt: nombre, texto y un campo de prueba que pasa un ejemplo por el modelo en marcha desde los propios ajustes. El orden se cambia arrastrando el prompt por el asa de la izquierda.</li>" +
			"<li>Consejo: los modelos pequeños funcionan mucho mejor si el prompt lleva ejemplos «entrada → salida» — todos los incluidos están escritos así.</li>" +
			"<li>Si un perfil falla (el modelo no respondió), el texto se inserta sin él: la barra muestra «Insertado sin el perfil …» y en ese caso no se pulsa Enter.</li>" +
			"</ul>" +
			"<p class=\"wh\">Dependencias</p>" +
			"<ul>" +
			"<li>Los prompts necesitan un modelo editor instalado; la traducción no depende de él — la hace Whisper solo.</li>" +
			"<li>El modelo editor se carga en el primer uso y se queda caliente; los modelos grandes van bastante más lentos en el procesador.</li>" +
			"<li>Mira el indicador de memoria antes de descargar: un modelo «justo» frena todo el sistema.</li>" +
			"<li>Los controles en gris son ajustes que no intervienen en el modo actual.</li>" +
			"</ul>" +
			"<p class=\"wh\">Instalación y portabilidad</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — el instalador: sin permisos de administrador, acceso directo en el menú Inicio, inicio automático si quieres y desinstalación limpia desde la configuración de Windows.</li>" +
			"<li>El instalador no descarga nada por defecto: el asistente elige y baja el modelo en el primer arranque. Si aun así se elige uno — GigaAM v3 para el ruso, Whisper para los demás idiomas — la descarga se puede detener con un botón y la instalación termina igualmente. Allí está también la casilla «Buscar actualizaciones», y la respuesta se escribe en los ajustes de la aplicación.</li>" +
			"<li><b>Portátil</b> — copia toda la carpeta con el exe (a un USB, a otro equipo): ajustes, modelos y registro viven al lado y viajan con él. En el registro de Windows no se escribe nada.</li>" +
			"<li>En el primer arranque sin modelo de reconocimiento, el programa abre el catálogo por su cuenta y espera la descarga.</li>" +
			"<li>Requisitos: Windows 10/11 x64, un procesador con AVX2 (de 2013 en adelante, más o menos) y WebView2 Runtime para la ventana de ajustes (incluido en Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Área de notificación y archivos</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Listo…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Ajustes…</div><div class=\"mock-mi\">Desactivar</div><div class=\"mock-mi\">Copiar el último resultado</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Abrir config.json</div><div class=\"mock-mi\">Abrir el registro</div><div class=\"mock-mi\">Acerca de</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Salir</div></div>" +
			"<ul>" +
			"<li>Clic izquierdo en el icono — los ajustes; clic derecho — el menú. Colores del icono: verde — listo, rojo — grabando, naranja — reconociendo, gris — desactivado o error.</li>" +
			"<li><b>config.json</b> — todos los ajustes; los cambios a mano se aplican con <b>Releer</b> en la sección «Sistema». Allí están también «Abrir el registro» y «Restablecer los ajustes»: el restablecimiento vuelve al estado de fábrica y deja modelos, historial y prompts como están.</li>" +
			"<li><b>{log}</b> — el registro, limitado automáticamente a unos 2 MB.</li>" +
			"<li><b>models/</b> — los modelos de reconocimiento y de edición descargados.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Dos minutos de configuración",
		"S_WIZ_HELLO_TEXT": "{app} convierte tu voz en texto justo donde está el cursor: mantén la combinación, di una frase, suelta y el texto aparece. Todo se procesa en tu equipo; el audio nunca sale de él.",
		"S_WIZ_LATER":      "Todo lo que elijamos ahora se puede cambiar después en los ajustes.",
		"S_WIZ_T_MODEL":    "Idioma y modelo",
		"S_WIZ_MODEL_TEXT": "Dime en qué idioma vas a dictar y yo elijo el modelo. El ruso lo reconoce GigaAM; el resto de idiomas, Whisper.",
		"S_WIZ_T_INPUT":    "Atajo y micrófono",
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
		"S_TITLE": "{app} — Impostazioni", "S_DICT_HINT": "Termini, nomi e abbreviazioni separati da virgole — un suggerimento per l'ascolto, non comandi. Vale per Whisper; il russo tramite GigaAM lo ignora. Il set predefinito segue la lingua di riconoscimento finché non scrivi il tuo.",
		"S_TR_DEFAULT": "Cambiare la lingua di uscita del testo", "S_TR_TARGET": "Lingua di uscita del testo predefinita", "S_TR_ASK": "Chiedere la lingua di uscita del testo", "S_TR_ASK_NEVER": "Non chiedere — tradurre subito",
		"S_SRCLANG_SUB": "la parlate voi; determina il modello di riconoscimento",
		"S_TR_LANGS_SUB": "queste lingue diventano pulsanti sulla targa all’inserimento",
		"S_TR_UNAVAIL": "non disponibile — %s non sa tradurre",
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
		"S_SEC_CMD": "Comandi vocali", "S_CMD_HINT": "Ciò che dici diventa un a capo, un segno o un annullamento invece di finire nel testo. Riconosciuti come parole intere, applicati dall'alto in basso, dopo le sostituzioni.",
		"S_CMD_ADD": "Aggiungi un comando", "S_CMD_PRESET": "Aggiungi quelli soliti", "S_CMD_PH": "la frase che dirai",
		"S_CMD_NEWLINE": "a capo", "S_CMD_PARAGRAPH": "nuovo paragrafo", "S_CMD_TEXT": "inserire testo", "S_CMD_CANCEL": "annullare la dettatura",
		"S_CMD_TEXT_PH": "cosa inserire", "S_CMD_EMPTY": "Ancora nessun comando", "S_CMD_DEL": "Elimina il comando",
		"S_CMD_P_NEWLINE": "nuova riga", "S_CMD_P_PARAGRAPH": "nuovo paragrafo", "S_CMD_P_CANCEL": "annulla",
		"S_SEC_REPLACE": "Sostituzioni dopo il riconoscimento", "S_REPLACE_HINT": "Ciò che è stato sentito male diventa ciò che intendevi — subito dopo il riconoscimento, prima dei prompt. Applicate dall'alto in basso.",
		"S_REPL_WHOLE_FULL": "Solo parole intere", "S_REPL_CASE_FULL": "Distingui le maiuscole", "S_CMD_ACTION": "Azione",
		"S_FM_ADD": "Aggiungi",
		"S_TIP_REPL_LANG": "La regola agisce solo quando detti nella lingua scelta. «tutte le lingue» — agisce sempre.",
		"S_TIP_REPL_CASE": "Maiuscole e minuscole contano: «git» e «Git» sono parole diverse. Disattivato — il caso è ignorato.",
		"S_TIP_REPL_WHOLE": "La sostituzione agisce solo se il testo è una parola a sé. Disattivato — anche dentro altre parole.",
		"S_TIP_CMD_ACTION": "Cosa succede quando pronunci la frase: a capo, nuovo paragrafo, un tuo testo o annullamento della dettatura.",
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
		"S_OV_FREE": "Posto tuo", "S_OV_FREE_SUB": "la targhetta si può trascinare ovunque",
		"S_OVPOS_DRAG_SUB": "trascina la targhetta con il mouse: va dove vuoi",
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
		"S_DL": "Scarica", "S_DEL": "Elimina",
		"S_M_BASE": "veloce, PC modesti", "S_M_SMALL": "bilanciato", "S_M_MED": "più preciso, consigliato", "S_M_TURBO": "massima precisione su CPU", "S_M_PARAKEET": "25 lingue europee, punteggia da sé",
		"S_THREADS":  "Thread CPU", "S_MINMS": "Registrazione min, ms", "S_MAXSEC": "Registrazione max, s",
		"S_AUTOSTART": "Avvia whisper-server automaticamente", "S_PORT": "Porta", "S_SERVEREXE": "Percorso whisper-server", "S_SERVEREXE_SUB": "si compila da sé; cambialo solo se il server sta altrove", "S_EXE_RESET": "Ripristina", "S_EXE_WARN": "L'app trova whisper-server accanto a sé. Con un percorso scritto a mano, spostando la cartella il riconoscimento non parte più. Cambiarlo?", "S_RESET_ALL": "Ripristina le impostazioni", "S_RESET_ALL_SUB": "tutto torna di fabbrica, tranne modelli e cronologia", "S_RESET_ALL_BTN": "Ripristina", "S_RESET_ALL_ASK": "Riportare tutte le impostazioni di fabbrica? Modelli, cronologia e prompt restano.", "S_RELOAD_CFG": "Rileggi config.json", "S_RELOAD_CFG_SUB": "se hai modificato il file a mano", "S_RELOAD_CFG_BTN": "Rileggi", "S_UPD_FOUND": "È uscita la versione %s", "S_THEME": "Colore", "S_THEME_SUB": "il colore della finestra, della barra e dell'icona", "S_THEME_GREEN": "Verde", "S_THEME_AMBER": "Ambra", "S_THEME_BLUE": "Blu", "S_THEME_PINK": "Rosa", "S_THEME_EDITOR": "Editor", "S_THEME_NEON": "Neon", "S_WND_MAX": "Riempi lo schermo", "S_WND_RESTORE": "Torna alla dimensione precedente", "S_WND_MIN": "Riduci nell'area di notifica", "S_WND_CLOSE": "Chiudi la finestra", "S_SKIN": "Design", "S_SKIN_SUB": "carattere, forme, effetti e movimento", "S_SKIN_TERMINAL": "Terminale", "S_SKIN_SOFT": "Morbido", "S_SKIN_PAPER": "Documento",
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
		"S_DEP_LLAMA": "post-elaborazione, modelli GGUF",
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
		"S_PRESETS_HINT": "Fate clic su una lingua — sotto si apre la scelta dei modelli per essa. Le lingue senza modello proprio usano quello del rilevamento automatico.",
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
		"S_TR_EN": "traduce in inglese",
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
		"S_OVPOS_SCHEME_SUB": "fai clic sullo schermo: la targhetta va lì",
		"S_OVDRAG": "Trascinatela dove volete",
		"S_OVMON": "Schermo",
		"S_OVMON_SUB": "su quale monitor mostrare la targhetta",
		"S_OVMON_CURSOR": "Lo schermo con il cursore",
		"S_M_NEMOTRON": "scrive mentre parlate: il testo appare sulla targhetta in diretta; 40 lingue, punteggia da sé",
		"S_M_TINY": "il più piccolo e veloce, per macchine molto modeste; sensibilmente meno preciso",
		"S_STATE_LOADED": "In memoria in questo momento",
		"S_STATE_LOADED_SUB": "i modelli si scaricano da soli dopo l’inattività",
		"S_STATE_WEEK": "Questa settimana",
		"S_ST_SUMMARY": "Riepilogo", "S_ST_OVERLAY": "Targhetta a schermo", "S_ST_BEEP": "Segnale acustico", "S_ST_AUTORUN": "Avvio con Windows", "S_ST_POST": "Post-elaborazione", "S_ST_LOCAL": "in locale", "S_ST_CHECKED": "verificato %s", "S_ST_GB": "%s GB", "S_ST_ON_M": "attivo", "S_ST_OFF_M": "spento", "S_ST_MIC_OK": "il segnale è buono", "S_ST_MIC_BAD": "il microfono non risponde", "S_ST_CHECK": "Controlla", "S_ST_RECOG": "riconosciuto da %s", "S_ST_VER": "Versione %s", "S_ST_LATEST": "l'ultima", "S_ST_OUTDATED": "non aggiornata", "S_ST_UPD_OK": "hai l'ultima versione", "S_ST_UPD_DL": "Scarico l'aggiornamento…",
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
		"S_SUB_TRTARGET": "il testo viene tradotto in essa; il dialogo sulla targa la propone per prima",
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
		"S_WIZ_NEED_MODEL": "Scarica prima un modello: senza non c'è nulla con cui riconoscere",
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
		"S_HELP_HTML": "<p class=\"wh\">Come funziona</p>" +
			"<p>Tieni premuta la scorciatoia — parte la registrazione (la barra in fondo allo schermo mostra il tuo livello). Rilascia — l'audio viene riconosciuto, tradotto se serve, passato nei prompt, e il testo finale arriva dov'è il cursore. La ✕ sulla barra annulla in qualsiasi momento.</p>" +
			"<p>Il percorso completo: <b>registrazione → riconoscimento (GigaAM per il russo, Whisper per il resto) → traduzione (se attiva) → prompt (LLM) → incollaggio</b>. Ogni passo si vede sulla barra.</p>" +
			"<p class=\"wh\">Primo avvio</p>" +
			"<p>Al primissimo avvio si apre una procedura in cinque passi: la lingua dell'interfaccia, la lingua in cui detterai (il modello lo sceglie e lo scarica lui), la scorciatoia e il microfono con la barra del livello, un campo per provare una dettatura e, per ultimo, l'avvio con Windows. Puoi saltarla e funziona lo stesso; <b>{exe} -wizard</b> la riporta. Con un aggiornamento non compare.</p>" +
			"<p class=\"wh\">La barra</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Parla…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Parla…</b> — registrazione: un punto rosso e le barre del livello.</li>" +
			"<li><b>Riconosco…</b> — Whisper sta lavorando; durante la traduzione — «Traduco», durante i prompt — «Modifico: nome (1/2)».</li>" +
			"<li><b>Inserito: N caratteri</b> — fatto; in caso di errore o silenzio compare brevemente il motivo.</li>" +
			"<li>La ✕ a destra annulla in qualsiasi momento; la barra non ruba mai il fuoco. Barra e animazione si spengono in «Dettatura».</li>" +
			"<li>Dove appare la barra — in basso, in alto o al cursore — e se mostra il testo riconosciuto invece del numero di caratteri, si imposta in «Dettatura».</li>" +
			"<li>Finché la barra fa una domanda, la riga in alto lo dice — «Attendo la tua risposta» — e il punto smette di pulsare. Ogni risposta ha il suo numero: 1…9 ne scelgono una, Invio prende quella evidenziata, Esc annulla tutto; i tasti sono scritti a destra, sulla stessa riga. Dieci secondi prima del limite di registrazione, sulla barra scorre un conto alla rovescia ambra.</li>" +
			"<li>La barra del titolo ha tre pulsanti: riduci nell'area di notifica, riempi lo schermo e chiudi. Lo stesso pulsante riporta la finestra alla dimensione precedente, e la dimensione impostata col mouse resta. La finestra non scende sotto 760×500.</li>" +
			"<li>I nomi lunghi — dispositivo, modello, file — sono tagliati con i puntini sulle schede di Stato, così le schede restano allineate; il nome intero compare come suggerimento se il puntatore si ferma sulla scheda. I suggerimenti sono disegnati nei colori dell'aspetto in uso, non in quelli di sistema.</li>" +
			"<li>L'aspetto viene da due elenchi nella sezione «Sistema». «Design» decide il carattere tipografico, le forme, lo spessore dei bordi, l'alone e il modo in cui tutto si muove; sono tre: «Terminale» (verde, predefinito), «Editor» (grigio piatto, senza alone) e «Neon» (viola, arrotondato). «Colore» viene offerto solo al Terminale e cambia soltanto il colore della finestra, della barra e dell'icona: verde, ambra, blu, rosa. Gli altri design portano i propri colori. La scelta vale subito, senza riavvio.</li>" +
			"</ul>" +
			"<p class=\"wh\">La domanda sulla traduzione</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Riconoscimento…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Tradurre in:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Senza traduzione</span></div></div>" +
			"<p>È la barra stessa a chiedere, su una seconda riga, appena rilasci la scorciatoia — nelle modalità «chiedi sempre» e «chiedi con timer». I pulsanti vengono da «Lingue nella finestra»; la lingua di destinazione è evidenziata. Con il timer, sotto quel pulsante si accorcia una linea: quando finisce vale la lingua evidenziata. <b>Senza traduzione</b> inserisce il testo così come è stato sentito; la ✕ della barra annulla tutta l'operazione. Funziona anche la tastiera: Invio prende la risposta evidenziata, 1…9 scelgono un pulsante, Esc annulla.</p>" +
			"<p class=\"wh\">Inserimento sicuro</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Riconoscimento…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Il fuoco è cambiato — inserire?</span><span class=\"mock-btn on mock-cd\">Inserisci qui</span><span class=\"mock-btn\">Copia</span></div></div>" +
			"<ul>" +
			"<li>La finestra di destinazione viene memorizzata nel momento in cui premi la scorciatoia. Se il fuoco è cambiato durante l'elaborazione, non viene incollato nulla — la barra chiede sulla sua seconda riga: <b>Inserisci qui</b> (nella finestra corrente), <b>Copia</b> (negli appunti) oppure la ✕. Allo scadere del tempo l'inserimento viene annullato e il testo resta nell'ultimo risultato.</li>" +
			"<li>Invio dopo l'incollaggio viene premuto solo se la finestra di destinazione è rimasta la stessa.</li>" +
			"<li><b>Ultimo risultato</b> — il testo finale di ogni dettatura resta in memoria fino alla successiva; nel menu dell'area di notifica c'è «Copia l'ultimo risultato». Un incollaggio fallito o un cambio di fuoco non fanno mai perdere una dettatura.</li>" +
			"</ul>" +
			"<p class=\"wh\">Controllare il microfono</p>" +
			"<p>Il pulsante «Test» nella sezione Microfono registra tre secondi e li scompone: picco in decibel, quanta parte della registrazione contiene davvero voce e quanti campioni sono stati tagliati. La risposta arriva a parole: si sente bene, troppo basso — alza il livello in Windows, distorsione — abbassalo, nessuna voce sentita — è scelto il microfono giusto. Le stesse misure vengono fatte dopo ogni dettatura e finiscono nel log; se il riconoscimento torna vuoto, la barra dice il motivo — basso, distorsione o silenzio — invece di dire soltanto che non ha sentito nulla.</p>" +
			"<p class=\"wh\">Incollare dalla cronologia</p>" +
			"<p>Ogni voce della cronologia ha il pulsante «Incolla»: riporta in primo piano la finestra da cui avete aperto le impostazioni e vi incolla il testo, come una dettatura normale. Se non c'è dove tornare, il testo finisce semplicemente negli appunti e il programma lo dice.</p>" +
			"<p class=\"wh\">Le liste in un solo file</p>" +
			"<p>Sostituzioni e comandi vocali si possono salvare in un file .json e caricare su un altro computer: i pulsanti sotto l'elenco dei comandi nella sezione «Testo». Il caricamento non cancella nulla: vengono aggiunte solo le righe che mancano, e il programma dice quante ne ha aggiunte e quante saltate.</p>" +
			"<p class=\"wh\">Integrità dei file</p>" +
			"<p>Per ogni modello del catalogo è noto l'hash SHA-256 di riferimento. Dopo lo scaricamento il file viene confrontato con esso: se non coincide, il file viene eliminato e lo scaricamento si può ripetere. Il pulsante «Controlla» nella sezione «Modelli» confronta allo stesso modo i modelli già installati, e all'aggiornamento del programma viene controllato anche l'installatore scaricato: un file estraneo non verrà avviato.</p>" +
			"<p class=\"wh\">Cronologia delle dettature</p>" +
			"<p>La sezione «Cronologia» nella colonna di sinistra conserva ciò che hai dettato: solo testo, solo su questo computer, l'audio non viene mai salvato. È disattivata per impostazione predefinita e si accende con un interruttore lì accanto. Le voci restano per i giorni e fino al numero che imposti, le più vecchie escono da sole; «Non registrare mai da questi programmi» elenca, separati da virgole, quelli da cui non salvare nulla — gestori di password, home banking. La ricerca copre il testo e il nome del programma, il pulsante accanto a una voce la mette negli appunti, e «Svuota» cancella tutto insieme al file <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Comandi vocali</p>" +
			"<p>Sotto le sostituzioni, nella sezione Testo, c'è un elenco di comandi: ciò che dici diventa un'azione invece che parole. «Nuova riga» e «nuovo paragrafo» inseriscono un a capo — i modelli non lo fanno mai; «annulla» butta via l'intera dettatura senza inserire nulla; «inserire testo» mette quello che vuoi, faccina compresa. Il pulsante accanto all'elenco lo riempie con le formule solite nella lingua dell'interfaccia. I comandi sono riconosciuti come parole intere e si applicano dopo le sostituzioni, così i prompt e la traduzione ricevono già il testo finito. Gli spazi di troppo attorno agli a capo spariscono da soli. Il campo in fondo prova sostituzioni e comandi su qualsiasi frase: un a capo appare come ⏎.</p>" +
			"<p class=\"wh\">Sostituzioni dopo il riconoscimento</p>" +
			"<p>In «Testo» puoi elencare ciò che il modello sente male e in che cosa deve diventare: «git hub» → GitHub, cognomi, termini interni. Le sostituzioni si applicano subito dopo il riconoscimento, prima dei prompt, così l'editor riceve già le parole giuste. La traduzione in inglese avviene dentro il riconoscimento, quindi le sostituzioni vedono il testo già tradotto. Per impostazione predefinita cercano parole intere e ignorano le maiuscole; i due interruttori accanto lo cambiano. Le regole si applicano dall'alto in basso. Il campo in fondo le prova su qualsiasi frase, senza dettare.</p>" +
			"<p class=\"wh\">Regole per applicazione</p>" +
			"<p>In «Dettatura» puoi impostare regole per programmi specifici: con che cosa inserire (appunti o carattere per carattere), se premere Invio, quanto attendere prima di inserire e quali prompt applicare. Il programma si indica con il nome del file — <b>chrome.exe</b>; in una regola se ne possono elencare più d'uno separati da virgole, e un asterisco finale cattura tutti i nomi che iniziano così. Vince la prima regola che corrisponde; senza regole, o se nessuna corrisponde, tutto funziona come nelle impostazioni generali. Il pulsante accanto all'elenco inserisce il programma in cui hai scritto l'ultima volta.</p>" +
			"<p class=\"wh\">Dettatura</p>" +
			"<ul>" +
			"<li><b>Scorciatoia da tastiera</b> — la scorciatoia principale. Si può catturare qualsiasi combinazione; i modificatori sinistro e destro sono distinti. Le scorciatoie di dettatura, traduzione e profili devono essere uniche — un doppione impedisce il salvataggio.</li>" +
			"<li><b>Modo</b> — tenere premuti i tasti, oppure premere una volta per iniziare e una per fermare.</li>" +
			"<li><b>Lingua dell'interfaccia</b> — cambia subito; «Come il sistema» segue Windows.</li>" +
			"<li><b>Lingua di riconoscimento</b> — un suggerimento per Whisper; «auto» la riconosce dal parlato.</li>" +
			"<li><b>Suono</b> — segnali di inizio e fine: diversi set più i suoni di sistema di Windows, ▶ li fa ascoltare.</li>" +
			"<li><b>Invio dopo l'incollaggio</b> — manda subito il testo dettato (comodo nelle chat).</li>" +
			"<li><b>Ripristina gli appunti</b> — rimette per intero il contenuto precedente, immagini, file e testo formattato compresi. Se il contenuto non si può salvare, gli appunti restano intatti e il testo viene digitato carattere per carattere.</li>" +
			"<li><b>Barra e animazione</b> — l'indicatore in fondo allo schermo; l'animazione si può spegnere.</li>" +
			"<li><b>Inserimento carattere per carattere</b> — invece di Ctrl+V vengono simulati i tasti, per i campi che rifiutano di incollare.</li>" +
			"</ul>" +
			"<p class=\"wh\">Riconoscimento</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Rilevamento automatico</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Russo</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · come il rilevamento automatico — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Quale modello per quale lingua</b> — la scheda «Lingue e modelli» è un elenco di lingue. Fatevi clic — sotto si aprono i modelli che la sanno servire: l’assegnato e il consigliato per primi, gli assenti con dimensione e freccia di download. Un clic sulla scheda è la scelta; un modello assente si scarica da solo e subentra appena pronto. Le lingue senza modello proprio ereditano quello del rilevamento automatico e appaiono attenuate.</li>" +
			"<li><b>Il catalogo</b> — Whisper: Base (veloce, per PC deboli), Small (l’equilibrio), Medium e Turbo (più precisi e lenti; «q5» è la versione quantizzata: un po’ più piccola e veloce quasi senza perdita), che traducono anche in inglese; GigaAM v3 è più preciso in russo e mette da sé la punteggiatura; Parakeet v3 copre 25 lingue europee; Nemotron 3.5 scrive mentre parlate. I download vengono dai repository ufficiali Hugging Face, ogni file verificato contro il suo hash di riferimento.</li>" +
			"<li><b>Modello proprio</b> — va bene Whisper come singolo file ggml-*.bin o una cartella di modello sherpa-onnx (encoder.onnx, decoder.onnx, tokens.txt). Mettetela nella cartella models accanto all’applicazione e riavviatela — il modello apparirà nella scelta delle lingue adatte; le sue capacità sono ignote, quindi è mostrato onestamente, senza barre.</li>" +
			"<li>Il server di riconoscimento tiene il modello in memoria tra una frase e l'altra — la prima dettatura dopo l'avvio è più lenta (caricamento), poi il riconoscimento richiede uno-tre secondi.</li>" +
			"<li><b>Dizionario</b> — termini, nomi e abbreviazioni separati da virgole. Un suggerimento per «l'orecchio» di Whisper, perché le parole rare arrivino giuste; non sono comandi.</li>" +
			"<li><b>Microfono</b> — scelta del dispositivo con il livello dal vivo (parla e la barra si muove: il dispositivo si sente). Se scolleghi quello scelto, subentra quello di sistema; una registrazione senza voce non viene nemmeno mandata al riconoscimento — la barra dice «Silenzio».</li>" +
			"<li><b>Servizio</b> — il server di riconoscimento parte da solo e gira in locale. Si possono cambiare porta, percorso o puntare a un server remoto; il riconoscitore poi si riavvia da sé.</li>" +
			"<li><b>Traduzione</b> — traduce tutto Whisper: in inglese nel suo modo nativo, nelle altre lingue <b>in via sperimentale</b>, forzando la lingua di uscita (la qualità dipende dalla coppia di lingue; le lingue grandi riescono meglio). Il modello Turbo non è addestrato per tradurre — le impostazioni avvisano finché è attivo. «Tradurre sempre nella lingua di destinazione» traduce ogni dettatura senza chiedere. Senza quella casella vale la modalità con domanda: sempre o con timer — la finestra della lingua compare prima del riconoscimento e allo scadere vale la destinazione. La scorciatoia di traduzione separata traduce una volta sola, senza toccare la dettatura normale.</li>" +
			"</ul>" +
			"<p class=\"wh\">Post-elaborazione (LLM)</p>" +
			"<p>Un secondo strato facoltativo: un modello linguistico locale (llama.cpp) sistema il testo riconosciuto secondo i tuoi prompt — toglie gli intercalari, cambia lo stile, formatta. Del tutto offline, solo sulla CPU.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modelli</b> — i modelli editor installati; il pulsante sceglie quello attivo (vale subito), la ✕ elimina (anche quello attivo — allora la post-elaborazione si spegne). Anche l'avanzamento degli scaricamenti si vede qui.</li>" +
			"<li><b>Ricerca</b> — modelli GGUF su Hugging Face per nome (per esempio «qwen2.5 instruct»). Ogni repository mostra la data dell'ultimo aggiornamento, il numero di download e un ↗ verso la pagina; un clic sulla riga apre i file quantizzati. L'indicatore ● ≈N GB si confronta con la memoria <b>libera</b> (indicata sopra l'elenco).</li>" +
			"<li><b>Quale quantizzazione:</b> il numero sono i bit per peso (Q4 — la giusta via di mezzo, Q8 — quasi non compresso, Q3 — risparmia memoria a scapito della qualità); K_M è meglio di K_S; IQ4 è la generazione più nuova e a parità di dimensione batte le classiche. L'indicatore ● ≈N GB stima la memoria necessaria (il file più un margine per il contesto): verde ci sta, ambra è al limite, rosso non ci sta.</li>" +
			"<li>Un modello da 1,5–3B modifica in fretta; uno da 7–9B è molto più fine ma su CPU impiega secondi a passata. Il server LLM parte al primo utilizzo e tiene il modello caldo.</li>" +
			"</ul>" +
			"<p class=\"wh\">Prompt</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Pulizia</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Stile formale</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Un prompt è un'istruzione per il modello editor. Due arrivano già pronti: «Pulizia» (toglie intercalari, ripetizioni e false partenze, sistema la punteggiatura) e «Stile formale» (riscrive in modo cortese e formale); i tuoi puoi aggiungerli liberamente.</li>" +
			"<li>I prompt spuntati valgono per ogni dettatura, in ordine dall'alto in basso (a catena: l'uscita di uno è l'ingresso del successivo); se non ne è spuntato nessuno, il testo viene inserito com'è stato riconosciuto.</li>" +
			"<li>La matita ✎ apre l'editor del prompt: nome, testo e un campo di prova che manda un esempio al modello acceso direttamente dalle impostazioni. L'ordine si cambia trascinando il prompt per la maniglia a sinistra.</li>" +
			"<li>Consiglio: i modelli piccoli lavorano molto meglio se nel prompt ci sono esempi «ingresso → uscita» — tutti quelli inclusi sono scritti così.</li>" +
			"<li>Se un profilo fallisce (il modello non ha risposto), il testo viene inserito senza di lui: la barra mostra «Inserito senza il profilo …» e in quel caso Invio non viene premuto.</li>" +
			"</ul>" +
			"<p class=\"wh\">Dipendenze</p>" +
			"<ul>" +
			"<li>I prompt richiedono un modello editor installato; la traduzione non ne dipende — la fa Whisper da sola.</li>" +
			"<li>Il modello editor si carica al primo utilizzo e resta caldo; i modelli grandi su CPU sono sensibilmente più lenti.</li>" +
			"<li>Guarda l'indicatore della memoria prima di scaricare: un modello «al limite» rallenta tutto il sistema.</li>" +
			"<li>I controlli in grigio sono impostazioni che nella modalità corrente non fanno nulla.</li>" +
			"</ul>" +
			"<p class=\"wh\">Installazione e portabilità</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — il programma di installazione: senza diritti di amministratore, collegamento nel menu Start, avvio automatico a richiesta, rimozione pulita dalle impostazioni di Windows.</li>" +
			"<li>Il programma di installazione non scarica nulla per impostazione predefinita: la procedura guidata sceglie e scarica il modello al primo avvio. Se un modello viene comunque scelto — GigaAM v3 per il russo, Whisper per le altre lingue — lo scaricamento si può fermare con un pulsante e l'installazione arriva lo stesso in fondo. Lì c'è anche la casella «Cerca aggiornamenti», e la risposta finisce nelle impostazioni dell'app.</li>" +
			"<li><b>Portabile</b> — copia semplicemente tutta la cartella con l'exe (su una chiavetta, su un altro computer): impostazioni, modelli e registro stanno accanto e viaggiano con lui. Nel registro di sistema non viene scritto nulla.</li>" +
			"<li>Al primo avvio senza modello di riconoscimento il programma apre il catalogo da solo e aspetta lo scaricamento.</li>" +
			"<li>Requisiti: Windows 10/11 x64, una CPU con AVX2 (dal 2013 circa), WebView2 Runtime per la finestra delle impostazioni (incluso in Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Area di notifica e file</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Pronto…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Impostazioni…</div><div class=\"mock-mi\">Disattiva</div><div class=\"mock-mi\">Copia l'ultimo risultato</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Apri config.json</div><div class=\"mock-mi\">Apri il registro</div><div class=\"mock-mi\">Informazioni</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Esci</div></div>" +
			"<ul>" +
			"<li>Clic sinistro sull'icona — le impostazioni; clic destro — il menu. Colori dell'icona: verde — pronto, rosso — registrazione, arancione — riconoscimento, grigio — disattivato o errore.</li>" +
			"<li><b>config.json</b> — tutte le impostazioni; le modifiche a mano valgono dopo <b>Rileggi</b> nella sezione «Sistema». Lì ci sono anche «Apri il registro» e «Ripristina le impostazioni»: il ripristino torna allo stato di fabbrica e lascia modelli, cronologia e prompt dove sono.</li>" +
			"<li><b>{log}</b> — il registro, limitato automaticamente a circa 2 MB.</li>" +
			"<li><b>models/</b> — i modelli di riconoscimento e di editing scaricati.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Due minuti di configurazione",
		"S_WIZ_HELLO_TEXT": "{app} trasforma la voce in testo proprio dove si trova il cursore: tieni premuta la combinazione, pronuncia una frase, rilascia — il testo è lì. Tutto gira sul tuo computer, l'audio non esce mai.",
		"S_WIZ_LATER":      "Tutto quello che scegliamo ora si può cambiare dopo nelle impostazioni.",
		"S_WIZ_T_MODEL":    "Lingua e modello",
		"S_WIZ_MODEL_TEXT": "Dimmi in che lingua detterai e il modello lo scelgo io. Il russo lo riconosce GigaAM, tutte le altre lingue Whisper.",
		"S_WIZ_T_INPUT":    "Scorciatoia e microfono",
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
		"S_TITLE": "{app} — Ustawienia", "S_DICT_HINT": "Terminy, nazwy i skróty oddzielone przecinkami — podpowiedź dla słuchu, nie polecenia. Działa dla Whispera; rosyjska mowa przez GigaAM tego nie używa. Zestaw domyślny podąża za językiem rozpoznawania, dopóki nie wpiszesz własnego.",
		"S_TR_DEFAULT": "Zmień język wyjściowego tekstu", "S_TR_TARGET": "Domyślny język wyjściowego tekstu", "S_TR_ASK": "Pytaj o język wyjściowego tekstu", "S_TR_ASK_NEVER": "Nie pytaj — tłumacz od razu",
		"S_SRCLANG_SUB": "mówisz w nim; określa model rozpoznawania",
		"S_TR_LANGS_SUB": "te języki będą przyciskami na plakietce przy wstawianiu",
		"S_TR_UNAVAIL": "niedostępne — %s nie umie tłumaczyć",
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
		"S_SEC_CMD": "Komendy głosowe", "S_CMD_HINT": "To, co powiesz, zamienia się w złamanie wiersza, znak albo anulowanie, zamiast trafić do tekstu. Rozpoznawane jako całe słowa, stosowane od góry do dołu, po zamianach.",
		"S_CMD_ADD": "Dodaj komendę", "S_CMD_PRESET": "Dodaj typowe", "S_CMD_PH": "fraza, którą wypowiesz",
		"S_CMD_NEWLINE": "złamanie wiersza", "S_CMD_PARAGRAPH": "nowy akapit", "S_CMD_TEXT": "wstawić tekst", "S_CMD_CANCEL": "anulować dyktowanie",
		"S_CMD_TEXT_PH": "co wstawić", "S_CMD_EMPTY": "Na razie brak komend", "S_CMD_DEL": "Usuń komendę",
		"S_CMD_P_NEWLINE": "nowy wiersz", "S_CMD_P_PARAGRAPH": "nowy akapit", "S_CMD_P_CANCEL": "anuluj",
		"S_SEC_REPLACE": "Zamiany po rozpoznaniu", "S_REPLACE_HINT": "To, co zostało źle usłyszane, staje się tym, co miałeś na myśli — zaraz po rozpoznaniu, przed promptami. Stosowane od góry do dołu.",
		"S_REPL_WHOLE_FULL": "Tylko całe słowa", "S_REPL_CASE_FULL": "Uwzględniaj wielkość liter", "S_CMD_ACTION": "Działanie",
		"S_FM_ADD": "Dodaj",
		"S_TIP_REPL_LANG": "Reguła działa tylko, gdy dyktujesz w wybranym języku. „wszystkie języki” — działa zawsze.",
		"S_TIP_REPL_CASE": "Wielkość liter ma znaczenie: „git” i „Git” to różne słowa. Wyłączone — wielkość liter jest pomijana.",
		"S_TIP_REPL_WHOLE": "Zamiana działa tylko, gdy tekst stoi jako osobne słowo. Wyłączone — trafia też wewnątrz innych słów.",
		"S_TIP_CMD_ACTION": "Co się stanie, gdy wypowiesz frazę: nowa linia, nowy akapit, własny tekst albo anulowanie dyktowania.",
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
		"S_DL": "Pobierz", "S_DEL": "Usuń",
		"S_M_BASE": "szybki, słabe komputery", "S_M_SMALL": "zrównoważony", "S_M_MED": "dokładniejszy, polecany", "S_M_TURBO": "najlepsza dokładność na CPU", "S_M_PARAKEET": "25 języków europejskich, sama stawia interpunkcję",
		"S_THREADS":  "Wątki CPU", "S_MINMS": "Min. nagranie, ms", "S_MAXSEC": "Maks. nagranie, s",
		"S_AUTOSTART": "Uruchamiaj whisper-server automatycznie", "S_PORT": "Port", "S_SERVEREXE": "Ścieżka whisper-server", "S_SERVEREXE_SUB": "wypełnia się sama; zmieniaj tylko, jeśli serwer leży gdzie indziej", "S_EXE_RESET": "Przywróć", "S_EXE_WARN": "Aplikacja znajduje whisper-server obok siebie. Przy ręcznie wpisanej ścieżce po przeniesieniu folderu rozpoznawanie przestanie się uruchamiać. Zmienić?", "S_RESET_ALL": "Przywróć ustawienia", "S_RESET_ALL_SUB": "wszystko wraca do fabrycznych, poza modelami i historią", "S_RESET_ALL_BTN": "Przywróć", "S_RESET_ALL_ASK": "Przywrócić wszystkie ustawienia fabryczne? Modele, historia i prompty zostają.", "S_RELOAD_CFG": "Wczytaj config.json ponownie", "S_RELOAD_CFG_SUB": "jeśli plik był zmieniany ręcznie", "S_RELOAD_CFG_BTN": "Wczytaj", "S_UPD_FOUND": "Jest wersja %s", "S_THEME": "Kolor", "S_THEME_SUB": "kolor okna, paska i ikony w zasobniku", "S_THEME_GREEN": "Zielony", "S_THEME_AMBER": "Bursztynowy", "S_THEME_BLUE": "Niebieski", "S_THEME_PINK": "Różowy", "S_THEME_EDITOR": "Edytor", "S_THEME_NEON": "Neon", "S_WND_MAX": "Na cały ekran", "S_WND_RESTORE": "Poprzedni rozmiar", "S_WND_MIN": "Ukryj w zasobniku", "S_WND_CLOSE": "Zamknij okno", "S_SKIN": "Wygląd", "S_SKIN_SUB": "krój pisma, kształty, efekty i animacja", "S_SKIN_TERMINAL": "Terminal", "S_SKIN_SOFT": "Miękki", "S_SKIN_PAPER": "Dokument",
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
		"S_DEP_LLAMA": "obróbka tekstu, modele GGUF",
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
		"S_PRESETS_HINT": "Kliknij język — pod nim otworzy się wybór modeli dla niego. Języki bez własnego modelu używają modelu automatycznego wykrywania.",
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
		"S_TR_EN": "tłumaczy na angielski",
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
		"S_OVMON_SUB": "na którym monitorze pokazywać plakietkę",
		"S_OVMON_CURSOR": "Ekran z kursorem",
		"S_M_NEMOTRON": "pisze w trakcie mowy: tekst pojawia się na plakietce na żywo; 40 języków, sama stawia interpunkcję",
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
		"S_ST_OVERLAY_SUB": "widoczna podczas nagrywania",
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
		"S_SUB_TRTARGET": "na niego tłumaczony jest tekst; dialog na plakietce proponuje go pierwszy",
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
		"S_WIZ_NEED_MODEL": "Najpierw pobierz model — bez niego nie ma czym rozpoznawać",
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
		"S_HELP_HTML": "<p class=\"wh\">Jak to działa</p>" +
			"<p>Przytrzymaj skrót — zaczyna się nagrywanie (pasek na dole ekranu pokazuje twój poziom). Puść — dźwięk zostaje rozpoznany, w razie potrzeby przetłumaczony, przepuszczony przez prompty, a gotowy tekst trafia tam, gdzie stoi kursor. ✕ na pasku przerywa na każdym etapie.</p>" +
			"<p>Cała droga: <b>nagranie → rozpoznawanie (rosyjski — GigaAM, pozostałe języki — Whisper) → tłumaczenie (jeśli włączone) → prompty (LLM) → wklejenie</b>. Każdy etap widać na pasku.</p>" +
			"<p class=\"wh\">Pierwsze uruchomienie</p>" +
			"<p>Przy pierwszym uruchomieniu otwiera się kreator z pięciu kroków: język interfejsu, język dyktowania (model dobierze i pobierze sam), skrót i mikrofon z paskiem poziomu, pole do próbnego dyktowania i na końcu uruchamianie z Windows. Można go pominąć — wszystko i tak działa; <b>{exe} -wizard</b> przywraca go. Przy aktualizacji się nie pojawia.</p>" +
			"<p class=\"wh\">Pasek</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Mów…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Mów…</b> — nagrywanie: czerwona kropka i słupki poziomu.</li>" +
			"<li><b>Rozpoznaję…</b> — Whisper pracuje; przy tłumaczeniu — „Tłumaczę”, przy promptach — „Redaguję: nazwa (1/2)”.</li>" +
			"<li><b>Wstawiono: N znaków</b> — gotowe; przy błędzie albo ciszy pojawia się krótki powód.</li>" +
			"<li>✕ po prawej przerywa na każdym etapie; pasek nigdy nie zabiera fokusu. Pasek i jego animację można wyłączyć w „Dyktowaniu”.</li>" +
			"<li>Gdzie pojawia się pasek — na dole, na górze albo przy kursorze — i czy pokazuje rozpoznany tekst zamiast liczby znaków, ustawia się w „Dyktowaniu”.</li>" +
			"<li>Dopóki pasek o coś pyta, górny wiersz tak właśnie mówi — „Czekam na odpowiedź” — a kropka przestaje pulsować. Każda odpowiedź ma swój numer: 1…9 wybierają jedną, Enter bierze wyróżnioną, Esc anuluje wszystko; klawisze są wypisane po prawej, w tym samym wierszu. Na dziesięć sekund przed limitem nagrania na pasku biegnie bursztynowe odliczanie.</li>" +
			"<li>Pasek tytułu ma trzy przyciski: ukryj w zasobniku, na cały ekran i zamknij. Ten sam przycisk przywraca poprzedni rozmiar, a rozmiar ustawiony myszą zostaje zapamiętany. Okno nie schodzi poniżej 760×500.</li>" +
			"<li>Długie nazwy — urządzenia, modelu, pliku — są ucinane wielokropkiem na kartach „Stanu”, żeby karty stały równo; pełna nazwa pojawia się jako podpowiedź, gdy wskaźnik zatrzyma się na karcie. Podpowiedzi są narysowane w kolorach bieżącego wyglądu, nie systemowych.</li>" +
			"<li>Wygląd bierze się z dwóch list w sekcji „System”. „Wygląd” ustawia krój pisma, kształty, grubość ramek, poświatę i charakter animacji; są trzy: „Terminal” (zielony, domyślny), „Edytor” (płaska szarość, bez poświaty) i „Neon” (fiolet, zaokrąglony). „Kolor” pojawia się tylko przy Terminalu i zmienia jedynie kolor okna, paska i ikony w zasobniku: zielony, bursztynowy, niebieski, różowy. Pozostałe wyglądy mają własne kolory. Wybór działa od razu, bez ponownego uruchamiania.</li>" +
			"</ul>" +
			"<p class=\"wh\">Pytanie o tłumaczenie</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Rozpoznawanie…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Tłumaczyć na:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Bez tłumaczenia</span></div></div>" +
			"<p>Pyta sam pasek, w drugim wierszu, gdy tylko puścisz skrót — w trybach „pytaj za każdym razem” i „pytaj z odliczaniem”. Przyciski biorą się z „Języków w oknie”; język docelowy jest wyróżniony. Przy odliczaniu pod tym przyciskiem skraca się kreska: gdy zniknie, obowiązuje wyróżniony język. <b>Bez tłumaczenia</b> wstawia tekst tak, jak został usłyszany; ✕ na pasku przerywa całą operację. Klawiatura też działa: Enter wybiera wyróżnioną odpowiedź, 1…9 wskazują przycisk, Esc przerywa.</p>" +
			"<p class=\"wh\">Bezpieczne wstawianie</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Rozpoznawanie…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Zmieniło się okno — wstawić?</span><span class=\"mock-btn on mock-cd\">Wstaw tutaj</span><span class=\"mock-btn\">Kopiuj</span></div></div>" +
			"<ul>" +
			"<li>Okno docelowe zostaje zapamiętane w chwili naciśnięcia skrótu. Jeśli fokus zmienił się w trakcie przetwarzania, nic nie zostanie wklejone — pasek pyta w drugim wierszu: <b>Wstaw tutaj</b> (do bieżącego okna), <b>Kopiuj</b> (do schowka) albo ✕. Po upływie czasu wstawianie zostaje odwołane, a tekst zostaje w ostatnim wyniku.</li>" +
			"<li>Enter po wklejeniu jest wysyłany tylko wtedy, gdy okno docelowe się nie zmieniło.</li>" +
			"<li><b>Ostatni wynik</b> — gotowy tekst każdego dyktowania zostaje w pamięci do następnego; w menu w zasobniku jest „Kopiuj ostatni wynik”. Nieudane wklejenie albo zmiana okna nigdy nie gubią dyktowania.</li>" +
			"</ul>" +
			"<p class=\"wh\">Sprawdzenie mikrofonu</p>" +
			"<p>Przycisk „Test” w sekcji Mikrofon nagrywa trzy sekundy i je rozbiera: szczyt w decybelach, jaka część nagrania naprawdę zawiera mowę i ile próbek zostało obciętych. Odpowiedź przychodzi słowami: brzmi dobrze, za cicho — podnieś poziom w Windows, przesterowanie — zmniejsz go, nie słychać mowy — czy wybrany jest właściwy mikrofon. To samo mierzy się po każdym dyktowaniu i trafia do dziennika; gdy rozpoznanie wraca puste, pasek nazywa powód — cicho, przesterowanie albo cisza — zamiast mówić tylko, że nic nie usłyszał.</p>" +
			"<p class=\"wh\">Wstrzymanie nagrania</p>" +
			"<p>W trybie przełącznika (jedno naciśnięcie zaczyna, kolejne kończy) można ustawić osobny skrót do pauzy: karta „Dyktowanie”, wiersz „Wstrzymaj nagranie”. Naciśnięcie zatrzymuje nagranie — plakietka pokazuje „Pauza” i nic nie jest zapisywane; kolejne wznawia je, a wszystko powiedziane wcześniej zostaje. Ograniczenie długości nie zadziała w czasie pauzy.</p>" +
			"<p class=\"wh\">Wklejanie z historii</p>" +
			"<p>Każdy wpis w historii ma przycisk „Wklej”: przywraca okno, z którego otworzyliście ustawienia, i wkleja tam tekst jak zwykłe dyktowanie. Gdy nie ma dokąd wracać, tekst po prostu trafia do schowka, a program o tym mówi.</p>" +
			"<p class=\"wh\">Listy w jednym pliku</p>" +
			"<p>Zamiany i polecenia głosowe można zapisać do jednego pliku .json i wczytać na innym komputerze — przyciski pod listą poleceń w sekcji „Tekst”. Wczytanie niczego nie kasuje: dodawane są tylko wiersze, których jeszcze nie ma, a program powie, ile dodano i ile pominięto.</p>" +
			"<p class=\"wh\">Nienaruszalność plików</p>" +
			"<p>Dla każdego modelu z katalogu znany jest wzorcowy skrót SHA-256. Po pobraniu plik jest z nim porównywany: gdy się nie zgadza, plik zostaje usunięty i pobieranie można powtórzyć. Przycisk „Sprawdź” w sekcji „Modele” tak samo porównuje modele już zainstalowane, a przy aktualizacji programu sprawdzany jest też pobrany instalator — obcy plik się nie uruchomi.</p>" +
			"<p class=\"wh\">Historia dyktowań</p>" +
			"<p>Sekcja „Historia” w lewej kolumnie przechowuje to, co podyktowałeś: tylko tekst, tylko na tym komputerze, dźwięk nigdy nie jest zapisywany. Domyślnie wyłączona, włącza się jednym przełącznikiem w tym samym miejscu. Wpisy trzymają się przez ustaloną liczbę dni i do ustalonej liczby, starsze wypadają same; „Nigdy nie zapisuj z tych programów” wymienia po przecinku te, z których nic nie ma być zapisywane — menedżery haseł, bankowość. Wyszukiwanie obejmuje tekst i nazwę programu, przycisk obok wpisu wkłada go do schowka, a „Wyczyść” usuwa wszystko razem z plikiem <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Komendy głosowe</p>" +
			"<p>Pod zamianami w sekcji „Tekst” jest lista komend: to, co powiesz, zamienia się w działanie zamiast w słowa. „Nowy wiersz” i „nowy akapit” wstawiają złamanie — modele nigdy tego nie robią; „anuluj” wyrzuca całe dyktowanie i nic nie wstawia; „wstawić tekst” wkłada cokolwiek, choćby uśmieszek. Przycisk obok listy wypełnia ją typowymi zwrotami w języku interfejsu. Komendy rozpoznawane są jako całe słowa i działają po zamianach, więc prompty i tłumaczenie dostają już gotowy tekst. Zbędne spacje wokół złamań znikają same. Pole poniżej sprawdza zamiany i komendy na dowolnym zdaniu: złamanie pokazuje się jako ⏎.</p>" +
			"<p class=\"wh\">Zamiany po rozpoznaniu</p>" +
			"<p>W „Tekście” można wypisać, co model słyszy źle i w co ma to zamienić: „git hub” → GitHub, nazwiska, wewnętrzne terminy. Zamiany działają zaraz po rozpoznaniu, przed promptami, więc edytor dostaje już właściwe słowa. Tłumaczenie na angielski dzieje się wewnątrz rozpoznawania, więc zamiany widzą już przetłumaczony tekst. Domyślnie szukają całych słów i nie zważają na wielkość liter; dwa przełączniki obok to zmieniają. Reguły działają od góry do dołu. Pole poniżej sprawdza je na dowolnym zdaniu, bez dyktowania.</p>" +
			"<p class=\"wh\">Reguły dla programów</p>" +
			"<p>W „Dyktowaniu” można ustawić reguły dla wybranych programów: czym wstawiać (schowkiem czy znak po znaku), czy naciskać Enter, ile czekać przed wstawieniem i jakie prompty stosować. Program wskazuje się nazwą pliku — <b>chrome.exe</b>; w jednej regule można wymienić kilka po przecinku, a gwiazdka na końcu łapie wszystkie nazwy o takim początku. Wygrywa pierwsza pasująca reguła; gdy reguł nie ma albo żadna nie pasuje, wszystko działa jak w ustawieniach ogólnych. Przycisk obok listy wpisuje program, do którego wstawiano ostatnio.</p>" +
			"<p class=\"wh\">Dyktowanie</p>" +
			"<ul>" +
			"<li><b>Skrót klawiszowy</b> — główny skrót do dyktowania. Można przechwycić dowolną kombinację; lewe i prawe modyfikatory są rozróżniane. Skróty dyktowania, tłumaczenia i profili muszą być unikalne — powtórzenie blokuje zapis.</li>" +
			"<li><b>Tryb</b> — trzymać klawisze albo nacisnąć raz, by zacząć, i raz, by skończyć.</li>" +
			"<li><b>Język interfejsu</b> — zmienia się od razu; „Jak w systemie” idzie za Windows.</li>" +
			"<li><b>Język rozpoznawania</b> — podpowiedź dla Whispera; „auto” rozpoznaje go z mowy.</li>" +
			"<li><b>Dźwięk</b> — sygnały startu i końca: kilka zestawów oraz dźwięki systemowe Windows, ▶ je odtwarza.</li>" +
			"<li><b>Enter po wklejeniu</b> — od razu wysyła podyktowany tekst (wygodne w komunikatorach).</li>" +
			"<li><b>Przywracanie schowka</b> — oddaje poprzednią zawartość w całości, razem z obrazami, plikami i tekstem z formatowaniem. Gdy zawartości nie da się zapisać, schowek zostaje nietknięty, a tekst jest wpisywany znak po znaku.</li>" +
			"<li><b>Pasek i animacja</b> — wskaźnik na dole ekranu; animację można wyłączyć.</li>" +
			"<li><b>Wstawianie znak po znaku</b> — zamiast Ctrl+V symulowane są naciśnięcia klawiszy, dla pól, które nie przyjmują wklejania.</li>" +
			"</ul>" +
			"<p class=\"wh\">Rozpoznawanie</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Automatyczne wykrywanie</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Rosyjski</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · jak wykrywanie automatyczne — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Który model dla którego języka</b> — karta „Języki i modele” to lista języków. Kliknij język — pod nim rozwiną się modele, które go potrafią: przypisany i polecany najpierw, brakujące z rozmiarem i strzałką pobierania. Kliknięcie karty to wybór; brakujący model pobierze się sam i przejmie pracę, gdy będzie gotowy. Języki bez własnego modelu dziedziczą model automatycznego wykrywania i są przygaszone.</li>" +
			"<li><b>Katalog</b> — Whisper: Base (szybki, dla słabszych PC), Small (równowaga), Medium i Turbo (dokładniejsze i wolniejsze; „q5” to wersja skwantyzowana: nieco mniejsza i szybsza niemal bez straty), tłumaczą też na angielski; GigaAM v3 jest dokładniejszy po rosyjsku i sam stawia interpunkcję; Parakeet v3 obejmuje 25 języków europejskich; Nemotron 3.5 pisze, gdy mówisz. Pobierane z oficjalnych repozytoriów Hugging Face, każdy plik sprawdzany z hashem wzorcowym.</li>" +
			"<li><b>Własny model</b> — pasuje Whisper jako pojedynczy plik ggml-*.bin lub folder modelu sherpa-onnx (encoder.onnx, decoder.onnx, tokens.txt). Umieść go w folderze models obok aplikacji i uruchom ją ponownie — model pojawi się w wyborze pasujących języków; jego możliwości są nieznane, więc pokazywany jest uczciwie, bez pasków.</li>" +
			"<li>Serwer rozpoznawania trzyma model w pamięci między zdaniami — pierwsze dyktowanie po starcie trwa dłużej (ładowanie), potem rozpoznanie zajmuje jedną–trzy sekundy.</li>" +
			"<li><b>Słownik</b> — terminy, nazwy i skróty po przecinku. Podpowiedź dla „ucha” Whispera, żeby rzadkie słowa trafiały poprawnie; to nie są polecenia.</li>" +
			"<li><b>Mikrofon</b> — wybór urządzenia z miernikiem poziomu (mów, a słupek się rusza — czyli urządzenie jest słyszane). Gdy wybrane urządzenie zniknie, wraca systemowe; nagranie bez mowy w ogóle nie trafia do rozpoznawania — pasek mówi „Cisza”.</li>" +
			"<li><b>Usługa</b> — serwer rozpoznawania startuje sam i działa lokalnie. Można zmienić port, ścieżkę albo wskazać serwer zdalny; rozpoznawanie samo się przeładuje.</li>" +
			"<li><b>Tłumaczenie</b> — wszystko tłumaczy Whisper: na angielski w trybie natywnym, na inne języki <b>eksperymentalnie</b>, przez wymuszenie języka wyjściowego (jakość zależy od pary językowej; duże języki wychodzą najlepiej). Model Turbo nie jest uczony tłumaczenia — ustawienia ostrzegają, dopóki jest aktywny. „Zawsze tłumacz na język docelowy” tłumaczy każde dyktowanie bez pytania. Bez tego pola obowiązuje tryb pytania: zawsze albo z odliczaniem — okno wyboru języka pojawia się przed rozpoznaniem, a po upływie czasu obowiązuje język docelowy. Osobny skrót tłumaczenia tłumaczy jednorazowo, nie zmieniając zwykłego dyktowania.</li>" +
			"</ul>" +
			"<p class=\"wh\">Obróbka tekstu (LLM)</p>" +
			"<p>Nieobowiązkowa druga warstwa: lokalny model językowy (llama.cpp) poprawia rozpoznany tekst według twoich promptów — usuwa przerywniki, zmienia styl, formatuje. Całkowicie bez sieci, tylko na procesorze.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modele</b> — zainstalowane modele redaktora; przycisk wyboru wskazuje aktywny (działa od razu), ✕ usuwa (także aktywny — wtedy obróbka się wyłącza). Postęp pobierania też widać tutaj.</li>" +
			"<li><b>Szukaj</b> — modele GGUF na Hugging Face po nazwie (na przykład „qwen2.5 instruct”). Każde repozytorium pokazuje datę ostatniej zmiany, liczbę pobrań i ↗ do strony modelu; kliknięcie wiersza rozwija pliki kwantyzacji. Wskaźnik ● ≈N GB porównuje się z <b>wolną</b> pamięcią (podaną nad listą).</li>" +
			"<li><b>Którą kwantyzację:</b> liczba to bity na wagę (Q4 — złoty środek, Q8 — prawie bez kompresji, Q3 — oszczędza pamięć kosztem jakości); K_M jest lepsze niż K_S; IQ4 to nowsza generacja i przy tym samym rozmiarze bije klasyczne. Wskaźnik ● ≈N GB szacuje potrzebną pamięć (plik plus zapas na kontekst): zielony mieści się, bursztynowy jest na styk, czerwony nie zmieści się.</li>" +
			"<li>Model 1,5–3B redaguje szybko; 7–9B jest wyraźnie mądrzejszy, ale na procesorze każde przejście trwa sekundy. Serwer LLM startuje przy pierwszym użyciu i trzyma model ciepły.</li>" +
			"</ul>" +
			"<p class=\"wh\">Prompty</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Porządki</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Styl formalny</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Prompt to polecenie dla modelu redaktora. Dwa są od razu w zestawie: „Porządki” (usuwa przerywniki, powtórzenia i falstarty, poprawia interpunkcję) i „Styl formalny” (przepisuje uprzejmie i oficjalnie); własne można dodawać bez ograniczeń.</li>" +
			"<li>Zaznaczone prompty działają na każde dyktowanie po kolei, z góry na dół (łańcuchem: wyjście jednego jest wejściem następnego); jeśli nic nie jest zaznaczone, tekst trafia taki, jaki został rozpoznany.</li>" +
			"<li>Ołówek ✎ otwiera edytor promptu: nazwa, treść i pole próby, które przepuszcza przykład przez działający model prosto z ustawień. Kolejność zmienia się przeciąganiem promptu za uchwyt po lewej.</li>" +
			"<li>Wskazówka: małe modele pracują znacznie lepiej, gdy w prompcie są przykłady „wejście → wyjście” — wszystkie dołączone są tak napisane.</li>" +
			"<li>Gdy profil zawiedzie (model nie odpowiedział), tekst zostanie wstawiony bez niego: pasek pokaże „Wstawiono bez profilu …”, a Enter w takim razie nie zostanie naciśnięty.</li>" +
			"</ul>" +
			"<p class=\"wh\">Zależności</p>" +
			"<ul>" +
			"<li>Prompty wymagają zainstalowanego modelu redaktora; tłumaczenie od niego nie zależy — robi je sam Whisper.</li>" +
			"<li>Model redaktora ładuje się przy pierwszym użyciu i zostaje ciepły; duże modele na procesorze są wyraźnie wolniejsze.</li>" +
			"<li>Zerknij na wskaźnik pamięci przed pobraniem: model „na styk” spowalnia cały system.</li>" +
			"<li>Wyszarzone elementy to ustawienia, które w bieżącym trybie nic nie robią.</li>" +
			"</ul>" +
			"<p class=\"wh\">Instalacja i przenośność</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — instalator: bez praw administratora, skrót w menu Start, autostart do wyboru, czyste usunięcie z ustawień Windows.</li>" +
			"<li>Instalator domyślnie nic nie pobiera: model wybierze i pobierze kreator przy pierwszym uruchomieniu. Jeśli model jednak zostanie wybrany — GigaAM v3 dla rosyjskiego, Whisper dla pozostałych języków — pobieranie można zatrzymać przyciskiem, a instalacja i tak dobiegnie końca. Jest tam też pole „Sprawdzaj aktualizacje”, a odpowiedź trafia do ustawień aplikacji.</li>" +
			"<li><b>Przenośność</b> — po prostu skopiuj cały folder z plikiem exe (na pendrive, na inny komputer): ustawienia, modele i dziennik leżą obok i jadą razem z nim. Do rejestru nic nie jest zapisywane.</li>" +
			"<li>Przy pierwszym uruchomieniu bez modelu rozpoznawania program sam otwiera katalog i czeka na pobranie.</li>" +
			"<li>Wymagania: Windows 10/11 x64, procesor z AVX2 (mniej więcej od 2013), WebView2 Runtime dla okna ustawień (jest w Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Zasobnik i pliki</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Gotowe…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Ustawienia…</div><div class=\"mock-mi\">Wyłącz</div><div class=\"mock-mi\">Kopiuj ostatni wynik</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Otwórz config.json</div><div class=\"mock-mi\">Otwórz dziennik</div><div class=\"mock-mi\">O programie</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Zakończ</div></div>" +
			"<ul>" +
			"<li>Lewy klik w ikonę — ustawienia; prawy — menu. Kolory ikony: zielony — gotowe, czerwony — nagrywanie, pomarańczowy — rozpoznawanie, szary — wyłączone albo błąd.</li>" +
			"<li><b>config.json</b> — wszystkie ustawienia; ręczne zmiany działają po <b>Wczytaj</b> w sekcji „System”. Tam są też „Otwórz dziennik” i „Przywróć ustawienia”: przywrócenie wraca do stanu fabrycznego i nie rusza modeli, historii ani promptów.</li>" +
			"<li><b>{log}</b> — dziennik, automatycznie ograniczany do około 2 MB.</li>" +
			"<li><b>models/</b> — pobrane modele rozpoznawania i redaktora.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Dwie minuty na konfigurację",
		"S_WIZ_HELLO_TEXT": "{app} zamienia głos w tekst dokładnie w miejscu kursora: przytrzymaj skrót, powiedz zdanie, puść — tekst jest na miejscu. Wszystko liczy się na twoim komputerze, dźwięk nigdzie nie wychodzi.",
		"S_WIZ_LATER":      "Wszystko, co teraz wybierzemy, można później zmienić w ustawieniach.",
		"S_WIZ_T_MODEL":    "Język i model",
		"S_WIZ_MODEL_TEXT": "Powiedz, w jakim języku będziesz dyktować, a model dobiorę sam. Rosyjski rozpoznaje GigaAM, pozostałe języki — Whisper.",
		"S_WIZ_T_INPUT":    "Skrót i mikrofon",
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
		"ov.silence":       "Тиша — нічого не розпізнано", "ov.server.loading": "Сервер ще завантажується",
		"ov.tooshort":  "Занадто коротко — тримайте клавіші довше",
		"ov.cancelled": "Скасовано", "ov.editing": "Редагую: %s", "ov.translating": "Перекладаю",
		"ov.llm.needed": "Ця мова потребує LLM-модуль", "td.title": "Перекласти на:", "td.plain": "Без перекладу",
		"cap.title": "{app} — сполучення клавіш", "cap.prompt": "Натисніть нове сполучення клавіш\n\nзараз: %s   ·   Esc — скасувати",
		"cap.selected": "Обрано: %s", "cap.cancelled": "Скасовано", "hk.taken": "Сполучення %s зайняте Windows: %s. Диктування може не початися", "hk.lock": "блокування комп’ютера", "hk.desktop": "показати робочий стіл", "hk.explorer": "провідник", "hk.run": "вікно «Виконати»", "hk.settings": "параметри Windows", "hk.search": "пошук", "hk.center": "центр сповіщень", "hk.menu": "меню досвідченого користувача", "hk.clipboard": "журнал буфера обміну", "hk.gamebar": "ігрова панель", "hk.voice": "голосове введення Windows", "hk.project": "проєціювання на екран", "hk.tasks": "подання завдань", "hk.layout": "зміна розкладки", "hk.newdesktop": "новий робочий стіл", "hk.closedesktop": "закрити робочий стіл", "hk.snip": "знімок екрана", "hk.switch": "перемикання вікон", "hk.close": "закрити вікно", "hk.cycle": "перебір вікон", "hk.start": "меню «Пуск»", "hk.taskmgr": "диспетчер завдань", "hk.secure": "екран безпеки",
		"err.hotkey.dup":    "Сполучення «%s» призначено двічі — хоткеї не повинні збігатися",
		"cfg.err.recovered": "config.json пошкоджено (%s).\nФайл збережено як %s, налаштування скинуто до типових.",
		"err.disk.space":    "мало місця на диску: вільно %d МБ, потрібно ~%d МБ",
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
		"un.data": "Видалити також налаштування та завантажені моделі?", "un.done": "{app} видалено.",
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
		"S_TITLE": "{app} — налаштування", "S_DICT_HINT": "Терміни, імена та абревіатури через кому — підказка слуху, не команди. Працює для Whisper; російська мова через GigaAM це ігнорує. Типовий набір іде за мовою розпізнавання, доки ви не впишете своє.",
		"S_TR_DEFAULT": "Змінити мову текстового виводу", "S_TR_TARGET": "Мова текстового виводу за замовчуванням", "S_TR_ASK": "Питати мову текстового виводу", "S_TR_ASK_NEVER": "Не питати — перекладати одразу",
		"S_SRCLANG_SUB": "нею ви говорите; вона визначає модель розпізнавання",
		"S_TR_LANGS_SUB": "ці мови будуть кнопками на плашці під час вставлення",
		"S_TR_UNAVAIL": "недоступно — %s не вміє перекладати",
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
		"S_WIZ_NEED_MODEL": "Спочатку завантажте модель — без неї немає чим розпізнавати",
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
		"S_SEC_CMD": "Голосові команди", "S_CMD_HINT": "Сказане перетворюється на перенесення рядка, знак або скасування замість того, щоб потрапити в текст. Шукаються цілими словами, застосовуються згори вниз, уже після замін.",
		"S_CMD_ADD": "Додати команду", "S_CMD_PRESET": "Додати звичайні", "S_CMD_PH": "фраза, яку ви промовите",
		"S_CMD_NEWLINE": "перенесення рядка", "S_CMD_PARAGRAPH": "новий абзац", "S_CMD_TEXT": "підставити текст", "S_CMD_CANCEL": "скасувати диктування",
		"S_CMD_TEXT_PH": "що підставити", "S_CMD_EMPTY": "Команд поки немає", "S_CMD_DEL": "Видалити команду",
		"S_CMD_P_NEWLINE": "новий рядок", "S_CMD_P_PARAGRAPH": "новий абзац", "S_CMD_P_CANCEL": "скасувати",
		"S_SEC_REPLACE": "Заміни після розпізнавання", "S_REPLACE_HINT": "Те, що почулося неправильно, стає тим, що ви мали на увазі — одразу після розпізнавання, до промптів. Застосовуються згори вниз.",
		"S_REPL_WHOLE_FULL": "Лише цілі слова", "S_REPL_CASE_FULL": "Враховувати регістр", "S_CMD_ACTION": "Дія",
		"S_FM_ADD": "Додати",
		"S_TIP_REPL_LANG": "Правило діє, лише коли диктування йде вибраною мовою. «всі мови» — діє завжди.",
		"S_TIP_REPL_CASE": "Великі й малі літери різняться: «гіт» і «Гіт» — різні слова. Вимкнено — регістр не важливий.",
		"S_TIP_REPL_WHOLE": "Заміна спрацьовує, лише якщо текст стоїть окремим словом. Вимкнено — шукається й усередині інших слів.",
		"S_TIP_CMD_ACTION": "Що станеться, коли ви промовите фразу: перенесення рядка, новий абзац, підстановка свого тексту або скасування диктування.",
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
		"S_OV_FREE": "Своє місце", "S_OV_FREE_SUB": "плашку можна перетягнути будь-куди",
		"S_OVPOS_DRAG_SUB": "тягніть плашку мишею — вона стане будь-де",
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
		"S_DL": "Завантажити", "S_DEL": "Видалити",
		"S_M_BASE": "швидка, для слабких ПК", "S_M_SMALL": "баланс швидкості й точності", "S_M_MED": "точніша, рекомендуємо", "S_M_TURBO": "максимум точності на CPU", "S_M_PARAKEET": "25 європейських мов, сама ставить розділові знаки",
		"S_MIC_CHECK": "Перевірити мікрофон", "S_MIC_CHECK_SUB": "три секунди запису та розбір: гучність, перевантаження, чи є мовлення", "S_MIC_CHECKING": "Перевіряю…",
		"S_MCHECK": "Перевірити встановлені моделі", "S_MCHECK_SUB": "звіряє файли моделей з еталонними хешами", "S_MCHECK_GO": "Перевірити", "S_MCHECK_RUN": "Перевіряю…",
		"S_HIST_INSERT": "Вставити",
		"S_MIC": "Мікрофон", "S_MIC_DEFAULT": "Системний за замовчуванням", "S_MIC_REFRESH": "Оновити список",
		"S_MIC_LEVEL": "Рівень сигналу", "S_MIC_QUIET": "тихо",
		"S_THREADS": "Потоки CPU", "S_MINMS": "Мін. запис, мс", "S_MAXSEC": "Макс. запис, с",
		"S_AUTOSTART": "Запускати whisper-server автоматично", "S_PORT": "Порт", "S_SERVEREXE": "Шлях до whisper-server", "S_SERVEREXE_SUB": "заповнюється самотужки; змінюйте, лише якщо сервер лежить в іншому місці", "S_EXE_RESET": "Скинути", "S_EXE_WARN": "Застосунок знаходить whisper-server поруч із собою. Зі шляхом, вписаним вручну, після перенесення теки розпізнавання перестане запускатися. Змінити?", "S_RESET_ALL": "Скинути налаштування", "S_RESET_ALL_SUB": "усе, крім моделей та історії, повертається до заводського", "S_RESET_ALL_BTN": "Скинути", "S_RESET_ALL_ASK": "Повернути всі налаштування до заводських? Моделі, історія та промпти залишаться.", "S_RELOAD_CFG": "Перечитати config.json", "S_RELOAD_CFG_SUB": "якщо файл правили руками", "S_RELOAD_CFG_BTN": "Перечитати", "S_UPD_FOUND": "Є версія %s", "S_THEME": "Колір", "S_THEME_SUB": "колір вікна, смужки та значка в лотку", "S_THEME_GREEN": "Зелений", "S_THEME_AMBER": "Бурштиновий", "S_THEME_BLUE": "Синій", "S_THEME_PINK": "Рожевий", "S_THEME_EDITOR": "Редактор", "S_THEME_NEON": "Неон", "S_WND_MAX": "Розгорнути на весь екран", "S_WND_RESTORE": "Повернути попередній розмір", "S_WND_MIN": "Згорнути в лоток", "S_WND_CLOSE": "Закрити вікно", "S_SKIN": "Дизайн", "S_SKIN_SUB": "шрифт, форма, ефекти та анімація", "S_SKIN_TERMINAL": "Термінал", "S_SKIN_SOFT": "М'який", "S_SKIN_PAPER": "Документ",
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
		"S_DEP_LLAMA": "постобробка тексту, моделі GGUF",
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
		"S_PRESETS_HINT": "Клацніть мову — під нею розгорнеться вибір моделей для неї. Мови без власної моделі використовують модель автовизначення.",
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
		"S_TR_EN": "перекладає англійською",
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
		"S_OVPOS_SCHEME_SUB": "клацніть по екрану — плашка стане туди",
		"S_OVDRAG": "Перетягніть, куди потрібно",
		"S_OVMON": "Екран",
		"S_OVMON_SUB": "на якому моніторі показувати плашку",
		"S_OVMON_CURSOR": "Екран із курсором",
		"S_M_NEMOTRON": "друкує під час мовлення: текст з’являється на плашці наживо; 40 мов, знаки сама",
		"S_M_TINY": "найменша і найшвидша, для дуже слабких машин; помітно менш точна",
		"S_STATE_LOADED": "Зараз у пам’яті",
		"S_STATE_LOADED_SUB": "моделі вивантажуються самі після простою",
		"S_STATE_WEEK": "За тиждень",
		"S_ST_SUMMARY": "Зведення", "S_ST_OVERLAY": "Плашка на екрані", "S_ST_BEEP": "Звуковий сигнал", "S_ST_AUTORUN": "Запуск із Windows", "S_ST_POST": "Постобробка", "S_ST_LOCAL": "локально", "S_ST_CHECKED": "перевірено %s", "S_ST_GB": "%s ГБ", "S_ST_ON_M": "увімкнений", "S_ST_OFF_M": "вимкнений", "S_ST_MIC_OK": "сигнал у нормі", "S_ST_MIC_BAD": "мікрофон не відповідає", "S_ST_CHECK": "Перевірити", "S_ST_RECOG": "розпізнає %s", "S_ST_VER": "Версія %s", "S_ST_LATEST": "остання", "S_ST_OUTDATED": "не остання", "S_ST_UPD_OK": "у вас остання версія", "S_ST_UPD_DL": "Завантажую оновлення…",
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
		"S_SUB_MINMS":   "відсікає випадкові натискання",
		"S_SUB_ENTER":   "надсилає повідомлення одразу",
		"S_SUB_CLIP":    "зображення й файли повертаються як були",
		"S_SUB_TYPE":    "допомагає там, де поле не приймає вставку",
		"S_SEC_OVERLAY": "Смужка на екрані",
		"S_SUB_THREADS": "більше потоків не завжди швидше — виміряйте на своїй машині",
		"S_SEC_LLM":     "Модель-редактор",
		"S_PUNCT":       "Розділові знаки й великі літери", "S_SUB_PUNCT": "звідки беруться розділові знаки й великі літери",
		"S_PUNCT_MODEL": "від моделі", "S_PUNCT_LLM": "від моделі-редактора", "S_PUNCT_OFF": "прибирати",
		"S_SUB_TRTARGET": "у неї перекладається текст; у діалозі на плашці вона запропонована першою",
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
		"S_NOT_FOUND": "нічого",   "S_HELP_HTML": "<p class=\"wh\">Як це працює</p>" +
			"<p>Тримайте сполучення — починається запис (смужка внизу екрана показує ваш рівень). Відпустіть — звук розпізнається, за потреби перекладається, проходить через промпти, і готовий текст з'являється там, де стоїть курсор. ✕ на смужці скасовує на будь-якому кроці.</p>" +
			"<p>Увесь шлях: <b>запис → розпізнавання (російська — GigaAM, інші мови — Whisper) → переклад (якщо ввімкнено) → промпти (LLM) → вставка</b>. Кожен крок видно на смужці.</p>" +
			"<p class=\"wh\">Перший запуск</p>" +
			"<p>При найпершому запуску відкривається майстер із п'яти кроків: мова інтерфейсу, мова диктування (модель він добере й завантажить сам), сполучення клавіш і мікрофон із живою смужкою рівня, поле для пробного диктування і, останнім, запуск разом із Windows. Майстер можна пропустити — усе працює й без нього; повернути — запуском <b>{exe} -wizard</b>. Під час оновлення він не з'являється.</p>" +
			"<p class=\"wh\">Смужка</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\"></span><span>Говоріть…</span><span class=\"mock-bars\"><i style=\"height:6px\"></i><i style=\"height:13px\"></i><i style=\"height:9px\"></i><i style=\"height:16px\"></i><i style=\"height:7px\"></i><i style=\"height:12px\"></i><i style=\"height:5px\"></i></span><span class=\"mock-x\">✕</span></div></div>" +
			"<ul>" +
			"<li><b>Говоріть…</b> — запис: червона крапка й смужки рівня.</li>" +
			"<li><b>Розпізнаю…</b> — Whisper працює; під час перекладу — «Перекладаю», під час промптів — «Редагую: назва (1/2)».</li>" +
			"<li><b>Вставлено: N символів</b> — готово; при помилці чи тиші коротко пишеться причина.</li>" +
			"<li>✕ праворуч скасовує на будь-якому кроці; смужка ніколи не забирає фокус вводу. Смужку та її анімацію можна вимкнути в «Диктуванні».</li>" +
			"<li>Де з'являється смужка — унизу, угорі чи біля курсора — і чи показує вона сам розпізнаний текст замість кількості символів, налаштовується в «Диктуванні».</li>" +
			"<li>Поки смужка про щось питає, верхній рядок так і каже — «Чекаю на відповідь», і крапка перестає блимати. Кожна відповідь має свій номер: 1…9 обирають одну, Enter бере підсвічену, Esc скасовує все; клавіші підписані праворуч у тому ж рядку. За десять секунд до межі запису на смужці йде бурштиновий зворотний відлік.</li>" +
			"<li>У заголовку вікна три кнопки: згорнути в лоток, розгорнути на весь екран і закрити. Та сама кнопка повертає вікно до попереднього розміру, а розмір, заданий мишею, зберігається. Менше 760×500 вікно не стає.</li>" +
			"<li>Довгі імена — пристрою, моделі, файла — на картках «Стану» обрізаються трикрапкою, щоб картки стояли рівно; повне ім'я показується підказкою, якщо затримати на картці вказівник. Підказки намальовані в кольорах поточного вигляду, а не системні.</li>" +
			"<li>Зовнішній вигляд задається двома списками в розділі «Система». «Дизайн» — це шрифт, форма, товщина рамок, сяйво та характер анімації; дизайнів три: «Термінал» (зелений, за замовчуванням), «Редактор» (плаский сірий, без сяйва) і «Неон» (фіолетовий, зі скругленнями). «Колір» пропонується лише «Терміналу» і змінює тільки колір вікна, смужки та значка в лотку: зелений, бурштиновий, синій, рожевий. Решта дизайнів мають власні кольори. Вибір діє одразу, без перезапуску.</li>" +
			"</ul>" +
			"<p class=\"wh\">Питання про переклад</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Розпізнаю…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span>Перекласти на:</span><span class=\"mock-btn on mock-cd\">EN</span><span class=\"mock-btn\">DE</span><span class=\"mock-btn\">Без перекладу</span></div></div>" +
			"<p>Питає сама смужка, другим рядком, щойно ви відпустили сполучення, — у режимах «питати щоразу» та «питати з відліком». Кнопки беруться з «Мов у вікні»; цільова мова виділена. У режимі з відліком під цією кнопкою коротшає риска: коли вона скінчиться, застосується виділена мова. <b>Без перекладу</b> вставляє текст так, як почуто; хрестик смужки скасовує операцію цілком. Клавіатура теж працює: Enter — виділена відповідь, 1…9 — кнопка за номером, Esc — скасувати.</p>" +
			"<p class=\"wh\">Безпечна вставка</p>" +
			"<div class=\"mock\"><div class=\"mock-pill\"><span class=\"mock-dot\" style=\"background:#ffb347;box-shadow:0 0 8px rgba(255,179,71,.8)\"></span><span>Розпізнаю…</span><span class=\"mock-x\">✕</span></div><div class=\"mock-ask\"><span style=\"color:var(--amber)\">Вікно змінилося — вставити?</span><span class=\"mock-btn on mock-cd\">Вставити тут</span><span class=\"mock-btn\">Копіювати</span></div></div>" +
			"<ul>" +
			"<li>Цільове вікно запам'ятовується в мить натискання сполучення. Якщо фокус змінився, поки оброблялася мова, нічого не вставляється — смужка питає другим рядком: <b>Вставити тут</b> (у поточне вікно), <b>Копіювати</b> (у буфер) або хрестик. Коли відлік вичерпано, вставка скасовується, а текст лишається в останньому результаті.</li>" +
			"<li>Enter після вставки натискається лише тоді, коли цільове вікно не змінилося.</li>" +
			"<li><b>Останній результат</b> — готовий текст кожного диктування лишається в пам'яті до наступного; у меню в треї є «Копіювати останній результат». Невдала вставка чи зміна вікна ніколи не втрачають диктування.</li>" +
			"</ul>" +
			"<p class=\"wh\">Перевірка мікрофона</p>" +
			"<p>Кнопка «Перевірка» на «Мікрофоні» записує три секунди й розбирає їх: пікова гучність у децибелах, частка запису, де справді є мовлення, і частка обрізаних відліків. Відповідь приходить словами: чути добре, надто тихо — додайте гучності у Windows, перевантаження — зменште її, мовлення не чути — чи той мікрофон вибрано. Те саме рахується після кожного диктування й пишеться в журнал; якщо розпізнати не вдалося, смужка назве причину — тихо, перевантаження чи тиша, — а не просто «нічого не почув».</p>" +
			"<p class=\"wh\">Повторна вставка з історії</p>" +
			"<p>У кожного запису в історії є кнопка «Вставити»: вона повертає вікно, з якого ви відкрили налаштування, і вставляє текст туди — як звичайне диктування. Якщо повертатися нікуди, текст просто лягає в буфер обміну, і програма про це скаже.</p>" +
			"<p class=\"wh\">Списки одним файлом</p>" +
			"<p>Заміни й голосові команди можна зберегти в один файл .json і завантажити на іншому комп'ютері — кнопки під списком команд на «Тексті». Завантаження нічого не затирає: додаються лише ті рядки, яких ще немає, а скільки додано й скільки пропущено, програма скаже.</p>" +
			"<p class=\"wh\">Цілісність файлів</p>" +
			"<p>Для кожної моделі з каталогу відомий еталонний хеш SHA-256. Після завантаження файл звіряється з ним: не зійшлося — файл видаляється, і завантаження можна повторити. Кнопка «Перевірити» на «Моделях» так само звіряє вже встановлені моделі, а під час оновлення програми звіряється й завантажений установник — чужий файл не запуститься.</p>" +
			"<p class=\"wh\">Історія диктувань</p>" +
			"<p>Розділ «Історія» в лівому стовпці зберігає те, що ви надиктували: лише текст, лише на цьому комп'ютері, звук не зберігається ніколи. Типово вимкнено — вмикається одним перемикачем там само. Записи тримаються задану кількість днів і до заданої кількості, старі зникають самі; поле «Не записувати з цих програм» перелічує через кому ті, з яких не треба зберігати нічого — менеджери паролів, банк-клієнт. Пошук шукає і за текстом, і за назвою програми, кнопка поруч із записом кладе його в буфер обміну, а «Очистити» видаляє все разом із файлом <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Голосові команди</p>" +
			"<p>Під замінами на «Тексті» — список команд: сказане перетворюється не на слова, а на дію. «Новий рядок» і «новий абзац» ставлять перенесення — моделі його не ставлять ніколи; «скасувати» викидає все диктування, нічого не вставляючи; «підставити текст» кладе будь-що, хоч смайлик. Кнопка поруч зі списком заповнює його звичними фразами мовою інтерфейсу. Команди шукаються цілими словами й застосовуються після замін, тому в промпти та переклад іде вже готовий текст. Зайві пробіли навколо перенесень прибираються самі. Поле внизу перевіряє і заміни, і команди на будь-якій фразі: перенесення показується значком ⏎.</p>" +
			"<p class=\"wh\">Заміни після розпізнавання</p>" +
			"<p>У «Тексті» можна перелічити, що модель чує неправильно і на що це змінювати: «гіт хаб» → GitHub, прізвища, внутрішні терміни. Заміни спрацьовують одразу після розпізнавання, до промптів, тому редактор отримує вже правильні слова. Переклад англійською відбувається всередині розпізнавання, тож заміни бачать уже перекладений текст. Типово шукаються цілі слова й без урахування регістру, два перемикачі поруч це змінюють. Правила застосовуються згори вниз. Поле внизу перевіряє їх на будь-якій фразі, нічого не диктуючи.</p>" +
			"<p class=\"wh\">Правила для програм</p>" +
			"<p>У «Диктуванні» можна задати правила для окремих програм: чим вставляти (буфером чи посимвольно), чи натискати Enter, скільки чекати перед вставкою та які промпти застосовувати. Програма вказується іменем файлу — <b>chrome.exe</b>; в одному правилі їх можна перелічити через кому, а зірочка в кінці ловить усі імена з таким початком. Виграє перше відповідне правило; якщо правил немає або жодне не підійшло, усе працює як у загальних налаштуваннях. Кнопка поруч зі списком підставляє програму, куди вставляли востаннє.</p>" +
			"<p class=\"wh\">Диктування</p>" +
			"<ul>" +
			"<li><b>Сполучення клавіш</b> — головне сполучення для диктування. Можна перехопити будь-яку комбінацію; ліві й праві модифікатори розрізняються. Сполучення диктування, перекладу та профілів мають бути унікальними — повтор блокує збереження.</li>" +
			"<li><b>Режим</b> — тримати клавіші або натиснути раз, щоб почати, і раз, щоб зупинити.</li>" +
			"<li><b>Мова інтерфейсу</b> — перемикається одразу; «Як у системі» йде за Windows.</li>" +
			"<li><b>Мова розпізнавання</b> — підказка для Whisper; «auto» визначає її з мовлення.</li>" +
			"<li><b>Звук</b> — сигнали початку й кінця: кілька наборів та системні звуки Windows, ▶ програє їх.</li>" +
			"<li><b>Enter після вставки</b> — одразу надсилає продиктований текст (зручно в месенджерах).</li>" +
			"<li><b>Відновлення буфера</b> — повертає попередній вміст повністю, разом із зображеннями, файлами й форматованим текстом. Коли вміст не вдається зберегти, буфер лишається недоторканим, а текст набирається символ за символом.</li>" +
			"<li><b>Смужка й анімація</b> — індикатор унизу екрана; анімацію можна вимкнути.</li>" +
			"<li><b>Вставка посимвольно</b> — замість Ctrl+V імітуються натискання клавіш, для полів, які не приймають вставку.</li>" +
			"</ul>" +
			"<p class=\"wh\">Розпізнавання</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><b>Автовизначення</b><span style=\"margin-left:auto\">Turbo (q5)</span></div>" +
			"<div class=\"mock-row\"><b>Російська</b><span style=\"margin-left:auto\">GigaAM v3</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-note\">English · як автовизначення — Turbo (q5)</span><span style=\"margin-left:auto\">▼</span></div></div>" +
			"<ul>" +
			"<li><b>Яка модель якій мові</b> — вкладка «Мови і моделі» — це список мов. Клацніть мову — під нею розгорнуться моделі, що її вміють: призначена й рекомендована першими, відсутні — з розміром і стрілкою завантаження. Клац по картці — це вибір; відсутня модель завантажиться сама й підхопить роботу, щойно буде готова. Мови без власної моделі успадковують модель автовизначення й показані тьмяно.</li>" +
			"<li><b>Каталог</b> — Whisper: Base (швидка, для слабких ПК), Small (баланс), Medium і Turbo (точніші й повільніші; «q5» — квантована версія: трохи менша й швидша майже без втрат), вони ж перекладають англійською; GigaAM v3 точніша російською й сама ставить розділові знаки; Parakeet v3 — 25 європейських мов; Nemotron 3.5 друкує, поки ви говорите. Завантаження — з офіційних репозиторіїв Hugging Face, кожен файл звіряється з еталонним хешем.</li>" +
			"<li><b>Своя модель</b> — підійде Whisper одним файлом ggml-*.bin або тека моделі sherpa-onnx (encoder.onnx, decoder.onnx, tokens.txt). Покладіть її в теку models поруч із застосунком і перезапустіть його — модель з’явиться у виборі відповідних мов; характеристики невідомі, тому показується чесно, без смужок.</li>" +
			"<li>Сервер розпізнавання тримає модель у пам'яті між фразами — перше диктування після запуску повільніше (завантаження), далі розпізнавання займає одну-три секунди.</li>" +
			"<li><b>Словник</b> — терміни, назви та скорочення через кому. Підказка для «вуха» Whisper, щоб рідкісні слова виходили правильно; це не команди.</li>" +
			"<li><b>Мікрофон</b> — вибір пристрою зі шкалою рівня (говоріть — і смужка рухається, отже пристрій чути). Якщо обраний пристрій зник, береться системний; запис без мовлення взагалі не йде на розпізнавання — смужка каже «Тиша».</li>" +
			"<li><b>Службове</b> — сервер розпізнавання запускається сам і працює локально. Порт, шлях або віддалений сервер можна змінити; розпізнавач після цього перезапуститься сам.</li>" +
			"<li><b>Переклад</b> — усе перекладає Whisper: англійською в рідному режимі, іншими мовами <b>експериментально</b>, через примусову мову виводу (якість залежить від пари мов; великі мови виходять найкраще). Модель Turbo цього не вміє — налаштування попереджають, поки вона активна. «Завжди перекладати цільовою мовою» перекладає кожне диктування без запитань. Без цієї позначки діє режим запитання: завжди або з відліком — вікно вибору мови з'являється перед розпізнаванням, а коли час вийде, береться цільова. Окреме сполучення перекладу перекладає один раз, не змінюючи звичайного диктування.</li>" +
			"</ul>" +
			"<p class=\"wh\">Постобробка (LLM)</p>" +
			"<p>Необов'язковий другий шар: локальна мовна модель (llama.cpp) править розпізнаний текст за вашими промптами — прибирає слова-паразити, змінює стиль, форматує. Повністю офлайн, лише на процесорі.</p>" +
			"<div class=\"mock\"><div class=\"mock-row\" style=\"color:var(--dim)\">▸ Qwen/Qwen2.5-3B-Instruct-GGUF<span style=\"margin-left:auto;color:var(--faint)\">2024-09-20</span><span>↓303k</span><span>↗</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q4_k_m.gguf<span style=\"margin-left:auto\">1.9 GB</span><span style=\"color:var(--green)\">● ≈3.7 GB</span><span style=\"color:var(--dim)\">⭳</span></div>" +
			"<div class=\"mock-row\" style=\"padding-left:18px\">q8_0.gguf<span style=\"margin-left:auto\">3.6 GB</span><span style=\"color:var(--amber)\">● ≈8.2 GB</span><span style=\"color:var(--dim)\">⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Моделі</b> — встановлені моделі-редактори; радіокнопка обирає активну (діє одразу), ✕ видаляє (зокрема активну — тоді постобробка вимикається). Поступ завантажень видно теж тут.</li>" +
			"<li><b>Пошук</b> — моделі GGUF на Hugging Face за назвою (наприклад «qwen2.5 instruct»). Кожен репозиторій показує дату оновлення, кількість завантажень і ↗ на сторінку моделі; клац по рядку розкриває файли квантування. Показник ● ≈N GB порівнюється з <b>вільною</b> пам'яттю (вона вказана над списком).</li>" +
			"<li><b>Яке квантування:</b> число — це біти на вагу (Q4 — золота середина, Q8 — майже без стиснення, Q3 — економить пам'ять ціною якості); K_M кращий за K_S; IQ4 — новіше покоління, за однакового розміру кращий за класичні. Показник ● ≈N GB оцінює потрібну пам'ять (файл плюс запас на контекст): зелений — вміститься, бурштиновий — впритул, червоний — не вміститься.</li>" +
			"<li>Модель на 1,5–3B редагує швидко; 7–9B помітно розумніша, але на процесорі кожен прохід триває секунди. Сервер LLM стартує при першому використанні й тримає модель теплою.</li>" +
			"</ul>" +
			"<p class=\"wh\">Промпти</p>" +
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-cb on\">✓</span><b>Чистка</b><span class=\"mock-note\">—</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-cb\"></span><b>Діловий стиль</b><span class=\"mock-note\">ctrl+alt+f</span><span style=\"margin-left:auto;color:var(--dim)\">✎ ✕</span></div></div>" +
			"<ul>" +
			"<li>Промпт — це інструкція для моделі-редактора. Два є з коробки: «Чистка» (прибирає слова-паразити, повтори й фальстарти, лагодить пунктуацію) та «Діловий стиль» (переписує ввічливо й формально); власні додавайте скільки завгодно.</li>" +
			"<li>Позначені промпти діють на кожне диктування по черзі, згори вниз (ланцюжком: вихід одного стає входом наступного); якщо не позначено нічого — текст вставляється таким, як розпізнано.</li>" +
			"<li>Промпт може мати власне сполучення: диктування з ним застосовує лише його, один раз. Олівець ✎ відкриває редактор: назва, текст промпта, сполучення та поле проби ▶, яке пропускає приклад крізь запущену модель прямо з налаштувань.</li>" +
			"<li>Порада: малі моделі працюють значно краще, коли в промпті є приклади «вхід → вихід» — усі вбудовані написані саме так.</li>" +
			"<li>Якщо профіль не спрацював (модель не відповіла), текст вставляється без нього: смужка показує «Вставлено без профілю …», і Enter у такому разі не натискається.</li>" +
			"</ul>" +
			"<p class=\"wh\">Залежності</p>" +
			"<ul>" +
			"<li>Промпти потребують встановленої моделі-редактора; переклад від неї не залежить — його робить сам Whisper.</li>" +
			"<li>Модель-редактор завантажується в пам'ять при першому використанні й лишається теплою; великі моделі на процесорі помітно повільніші.</li>" +
			"<li>Погляньте на показник пам'яті перед завантаженням: модель «впритул» гальмує всю систему.</li>" +
			"<li>Приглушені елементи — це налаштування, які в поточному режимі нічого не роблять.</li>" +
			"</ul>" +
			"<p class=\"wh\">Встановлення й портативність</p>" +
			"<ul>" +
			"<li><b>{setup}</b> — інсталятор: без прав адміністратора, ярлик у меню «Пуск», автозапуск за бажанням, чисте видалення з налаштувань Windows.</li>" +
			"<li>Інсталятор нічого не завантажує за замовчуванням: модель добере й завантажить майстер під час першого запуску. Якщо модель усе ж обрано — GigaAM v3 для російської, Whisper для решти мов — завантаження можна зупинити кнопкою, і встановлення все одно дійде до кінця. Там же галочка «Перевіряти оновлення», і відповідь записується в налаштування застосунку.</li>" +
			"<li><b>Портативність</b> — просто скопіюйте всю теку з exe (на флешку, на інший комп'ютер): налаштування, моделі та журнал лежать поруч і їдуть разом. У реєстр нічого не пишеться.</li>" +
			"<li>При першому запуску без моделі розпізнавання програма сама відкриває каталог і чекає на завантаження.</li>" +
			"<li>Вимоги: Windows 10/11 x64, процесор з AVX2 (приблизно від 2013 року), WebView2 Runtime для вікна налаштувань (входить до Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Трей і файли</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Готово…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Налаштування…</div><div class=\"mock-mi\">Вимкнути</div><div class=\"mock-mi\">Копіювати останній результат</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Відкрити config.json</div><div class=\"mock-mi\">Відкрити журнал</div><div class=\"mock-mi\">Про програму</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Вийти</div></div>" +
			"<ul>" +
			"<li>Лівий клац по значку — налаштування; правий — меню. Кольори значка: зелений — готово, червоний — запис, помаранчевий — розпізнавання, сірий — вимкнено або помилка.</li>" +
			"<li><b>config.json</b> — усі налаштування; правки вручну діють після кнопки <b>Перечитати</b> в розділі «Система». Там же — «Відкрити лог» і «Скинути налаштування»: скидання повертає все до заводського вигляду й не чіпає моделі, історію та промпти.</li>" +
			"<li><b>{log}</b> — журнал, автоматично обмежений приблизно 2 МБ.</li>" +
			"<li><b>models/</b> — завантажені моделі розпізнавання та редактора.</li>" +
			"</ul>",
		"S_WIZ_HELLO":      "Дві хвилини на налаштування",
		"S_WIZ_HELLO_TEXT": "{app} перетворює голос на текст просто в позиції курсора: затиснули сполучення клавіш, сказали фразу, відпустили — текст на місці. Усе рахується на вашому комп'ютері, звук нікуди не йде.",
		"S_WIZ_LATER":      "Усе, що ми зараз оберемо, потім змінюється в налаштуваннях.",
		"S_WIZ_T_MODEL":    "Мова і модель",
		"S_WIZ_MODEL_TEXT": "Скажіть, якою мовою диктуватимете, — модель підберу сам. Російську розпізнає GigaAM, решту мов — Whisper.",
		"S_WIZ_T_INPUT":    "Клавіші та мікрофон",
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
