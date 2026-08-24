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
		"err.webview":        "Das Einstellungsfenster benötigt die Microsoft WebView2 Runtime (in Windows 11 enthalten).\nInstallieren Sie sie oder bearbeiten Sie config.json manuell.",
		"status.loading":     "Modell wird geladen…", "status.ready": "Bereit — %s halten und sprechen",
		"status.recording": "Aufnahme läuft…", "status.transcribing": "Erkenne…", "status.disabled": "Deaktiviert",
		"status.server.restart": "Erkennungsserver abgestürzt, Neustart…", "status.cfg.err": "Fehler in config.json (siehe Log)",
		"menu.settings":         "Einstellungen…", "menu.enable": "Aktivieren", "menu.disable": "Deaktivieren",
		"menu.reload": "config.json neu laden", "menu.open.config": "config.json öffnen", "menu.open.log": "Log öffnen",
		"menu.about": "Über", "menu.quit": "Beenden",
		"ov.speak": "Sprechen…", "ov.transcribing": "Erkenne", "ov.inserted": "Eingefügt: %d Zeichen",
		"ov.err.recognize": "Erkennungsfehler (siehe Log)", "ov.err.paste": "Nicht eingefügt — der Text steht im letzten Ergebnis",
		"ov.moved": "Das Fenster hat gewechselt — der Text liegt in der Zwischenablage",
		"copy.ok": "Kopiert",
		"copy.none": "Nichts zu kopieren",
		"copy.fail": "Kopieren fehlgeschlagen: %s",
		"mic.busy": "Ein Diktat läuft, jetzt geht das nicht", "mic.check.ok": "Klingt gut: Spitze %.0f dB, Sprache in %.0f%% der Aufnahme",
		"mic.check.quiet": "Zu leise: Spitze %.0f dB — Mikrofonpegel in Windows anheben oder näher sitzen", "mic.check.clipped": "Übersteuert: %.1f%% der Abtastwerte abgeschnitten — Mikrofonpegel senken", "mic.check.silent": "Keine Sprache gehört — prüfen Sie, ob das richtige Mikrofon gewählt und nicht stumm ist",
		"ov.quiet": "Zu leise, es war fast nichts zu hören", "ov.clipped": "Übersteuert — der Ton wurde abgeschnitten",
		"ov.cmd.cancelled": "Per Sprache abgebrochen",
		"ov.silence": "Stille — nichts erkannt", "ov.server.loading": "Server lädt noch",
		"ov.tooshort": "Zu kurz — Tasten länger halten",
		"ov.cancelled": "Abgebrochen", "ov.editing": "Bearbeite: %s", "ov.translating": "Übersetze",
		"ov.llm.needed": "Diese Sprache benötigt das LLM-Modul", "td.title": "Übersetzen nach:", "td.plain": "Ohne Übersetzung",
		"cap.title": "{app} — Tastenkürzel", "cap.prompt": "Neue Tastenkombination drücken…\n\n(aktuell: %s)\n\nEsc — Abbrechen",
		"cap.selected": "Gewählt: %s", "cap.cancelled": "Abgebrochen",
		"model.switching": "Modellwechsel — der Erkenner startet neu…", "model.del.active": "Aktives Modell kann nicht gelöscht werden",
		"model.del.ok": "Modell gelöscht",
		"about.text":   "{app} %s\n\nSprache → Text an der Cursorposition.\nCursor in ein Eingabefeld setzen, %s halten, sprechen, loslassen — der Text wird eingefügt.\n\nErkennung: whisper.cpp, vollständig lokal und offline.\nModell: %s (Sprache: %s)\n\nEinstellungen: Klick auf das Tray-Symbol oder config.json.\nLogs: {log} (max. ~2 MB).",
		"status.nomodel": "Kein Erkennungsmodell geladen — wählen Sie eines in den Einstellungen",
		"menu.lastcopy": "Letztes Ergebnis kopieren",
		"ov.copied": "In die Zwischenablage kopiert", "ov.kept": "Abgebrochen — Text bleibt im letzten Ergebnis",
		"ov.llm.skipped": "Eingefügt ohne das Profil „%s“",
		"fd.title": "Fokus gewechselt — einfügen?", "fd.here": "Hier einfügen", "fd.copy": "Kopieren",
		"ov.err.mic": "Mikrofon nicht verfügbar — prüfen Sie das Gerät in den Einstellungen",
		"ov.notranslate": "Das aktive Modell kann nicht übersetzen — wie erkannt eingefügt",
		"ov.engine.fallback": "Die andere Engine startete nicht — es bleibt bei der aktuellen",
		"route.speech": "Sprache auf %s", "route.other": "Andere Sprachen", "route.translate": "Übersetzung",
		"route.lang.auto": "beliebige Sprache",
		"route.why.language": "hier genauer, mit Satzzeichen", "route.why.otherlang": "99 Sprachen",
		"route.why.translate": "nur Whisper übersetzt", "route.why.notinstalled": "das russische Modell ist nicht installiert",
		"route.why.unknownlang": "keine Sprache gesetzt — nur Whisper erkennt sie", "route.why.forced": "in config.json erzwungen",
		"adv.pick": "Ich empfehle %s.", "adv.companion": "%s passt gut dazu — für andere Sprachen und Übersetzung.",
		"adv.ram": "%d MB frei", "status.line": "Bereit · %s · %.1f GB frei",
		"ago.now": "gerade eben", "ago.min": "vor %d Min.", "ago.hour": "vor %d Std.",
		"chars": "%d Zeichen", "inserted.into": "eingefügt in %s",
		"punct.prompt": "Setze Satzzeichen und Großschreibung. Ändere die Wörter nicht, übersetze nicht, füge nichts hinzu. Gib nur den korrigierten Text zurück.",
		"err.sherpa.notfound": "sherpa-Erkennung nicht gefunden: %s",
		"err.sherpa.start": "sherpa-server hat sich beim Start beendet (siehe Protokoll)",
		"err.sherpa.translate": "dieses Modell kann nicht übersetzen",
		"err.sherpa.model": "Modelldatei nicht gefunden: %s — laden Sie sie in den Einstellungen oder korrigieren Sie sherpa_model in config.json",
		"err.hotkey.dup": "Das Kürzel „%s“ ist doppelt vergeben — Kürzel müssen eindeutig sein",
		"cfg.err.recovered": "config.json ist beschädigt (%s).\nDie Datei wurde als %s gesichert, die Einstellungen wurden zurückgesetzt.",
		"err.disk.space": "wenig Speicherplatz: %d MB frei, ~%d MB nötig",
		"err.save": "Einstellungen nicht gespeichert: %s — die alten bleiben",
		"err.port": "Port %d passt nicht: eine Nummer zwischen 1024 und 65535 wählen",
		"err.nolangs": "lassen Sie mindestens eine Sprache für die Übersetzungsfrage stehen",
		"ov.mic.lost": "Mikrofon weg — Aufnahme abgebrochen",
		"err.hash": "die heruntergeladene Datei ist beschädigt — bitte erneut versuchen",
		"models.check.ok": "Geprüfte Modelle: %d — alle Dateien in Ordnung",
		"models.check.none": "Nichts zu prüfen — kein installiertes Modell hat einen Referenz-Hash",
		"models.check.bad": "Beschädigte Dateien: %s — Modell erneut herunterladen",
		"ov.paused": "Pause",
		"status.paused": "Pause — die Aufnahme wartet",
		"hist.insert.gone": "Eintrag nicht gefunden",
		"hist.insert.nowin": "kein Ziel zum Einfügen — der Text liegt in der Zwischenablage",
		"hist.insert.ok": "in „%s“ eingefügt",
		"lists.bad": "diese Datei passt nicht",
		"lists.saved": "gespeichert in %s",
		"lists.added": "hinzugefügt: %d, übersprungen: %d",
		"lists.save.title": "Wohin die Listen speichern",
		"lists.open.title": "Welche Datei laden",
		"un.title": "{app} — Deinstallation", "un.confirm": "{app} von diesem Rechner entfernen?",
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
		"err.webview":        "La fenêtre des réglages nécessite Microsoft WebView2 Runtime (inclus dans Windows 11).\nInstallez-le ou modifiez config.json manuellement.",
		"status.loading":     "Chargement du modèle…", "status.ready": "Prêt — maintenez %s et parlez",
		"status.recording": "Enregistrement…", "status.transcribing": "Reconnaissance…", "status.disabled": "Désactivé",
		"status.server.restart": "Serveur planté, redémarrage…", "status.cfg.err": "Erreur dans config.json (voir log)",
		"menu.settings":         "Réglages…", "menu.enable": "Activer", "menu.disable": "Désactiver",
		"menu.reload": "Recharger config.json", "menu.open.config": "Ouvrir config.json", "menu.open.log": "Ouvrir le log",
		"menu.about": "À propos", "menu.quit": "Quitter",
		"ov.speak": "Parlez…", "ov.transcribing": "Reconnaissance", "ov.inserted": "Inséré : %d caractères",
		"ov.err.recognize": "Erreur de reconnaissance (voir log)", "ov.err.paste": "Non collé — le texte est dans le dernier résultat",
		"ov.moved": "La fenêtre a changé — le texte est dans le presse-papiers",
		"copy.ok": "Copié",
		"copy.none": "Rien à copier",
		"copy.fail": "Copie impossible : %s",
		"mic.busy": "Une dictée est en cours, impossible de vérifier", "mic.check.ok": "Bon signal : crête %.0f dB, parole sur %.0f%% de l'enregistrement",
		"mic.check.quiet": "Trop faible : crête %.0f dB — montez le niveau du micro dans Windows ou rapprochez-vous", "mic.check.clipped": "Saturation : %.1f%% des échantillons écrêtés — baissez le niveau du micro", "mic.check.silent": "Aucune parole entendue — vérifiez que le bon micro est choisi et qu'il n'est pas coupé",
		"ov.quiet": "Trop faible, presque rien n'a été entendu", "ov.clipped": "Saturation — le son a été écrêté",
		"ov.cmd.cancelled": "Annulé à la voix",
		"ov.silence": "Silence — rien reconnu", "ov.server.loading": "Le serveur charge encore",
		"ov.tooshort": "Trop court — maintenez les touches plus longtemps",
		"ov.cancelled": "Annulé", "ov.editing": "Édition : %s", "ov.translating": "Traduction",
		"ov.llm.needed": "Cette langue nécessite le module LLM", "td.title": "Traduire vers :", "td.plain": "Sans traduction",
		"cap.title": "{app} — raccourci", "cap.prompt": "Appuyez une nouvelle combinaison…\n\n(actuel : %s)\n\nÉchap — annuler",
		"cap.selected": "Choisi : %s", "cap.cancelled": "Annulé",
		"model.switching": "Changement de modèle — redémarrage…", "model.del.active": "Impossible de supprimer le modèle actif",
		"model.del.ok": "Modèle supprimé",
		"about.text":   "{app} %s\n\nVoix → texte à la position du curseur.\nPlacez le curseur, maintenez %s, parlez, relâchez — le texte s'insère.\n\nReconnaissance : whisper.cpp, entièrement locale et hors ligne.\nModèle : %s (langue : %s)\n\nRéglages : clic sur l'icône ou config.json.\nLogs : {log} (max ~2 Mo).",
		"status.nomodel": "Aucun modèle de reconnaissance téléchargé — choisissez-en un dans les réglages",
		"menu.lastcopy": "Copier le dernier résultat",
		"ov.copied": "Copié dans le presse-papiers", "ov.kept": "Annulé — le texte reste dans le dernier résultat",
		"ov.llm.skipped": "Inséré sans le profil « %s »",
		"fd.title": "Le focus a changé — insérer ?", "fd.here": "Insérer ici", "fd.copy": "Copier",
		"ov.err.mic": "Microphone indisponible — vérifiez le périphérique dans les réglages",
		"ov.notranslate": "Le modèle actif ne sait pas traduire — inséré tel que reconnu",
		"ov.engine.fallback": "L'autre moteur n'a pas démarré — on garde l'actuel",
		"route.speech": "Parole en %s", "route.other": "Autres langues", "route.translate": "Traduction",
		"route.lang.auto": "n'importe quelle langue",
		"route.why.language": "plus précis ici, avec la ponctuation", "route.why.otherlang": "99 langues",
		"route.why.translate": "seul Whisper traduit", "route.why.notinstalled": "le modèle russe n'est pas installé",
		"route.why.unknownlang": "aucune langue définie — seul Whisper la détecte", "route.why.forced": "forcé dans config.json",
		"adv.pick": "Je recommande %s.", "adv.companion": "%s fait un bon complément — pour les autres langues et la traduction.",
		"adv.ram": "%d Mo libres", "status.line": "Prêt · %s · %.1f Go libres",
		"ago.now": "à l'instant", "ago.min": "il y a %d min", "ago.hour": "il y a %d h",
		"chars": "%d caractères", "inserted.into": "inséré dans %s",
		"punct.prompt": "Ajoute la ponctuation et les majuscules. Ne change pas les mots, ne traduis pas, n'ajoute rien. Renvoie uniquement le texte corrigé.",
		"err.sherpa.notfound": "reconnaissance sherpa introuvable : %s",
		"err.sherpa.start": "sherpa-server s'est arrêté au démarrage (voir le journal)",
		"err.sherpa.translate": "ce modèle ne sait pas traduire",
		"err.sherpa.model": "Fichier de modèle introuvable : %s — téléchargez-le dans les réglages ou corrigez sherpa_model dans config.json",
		"err.hotkey.dup": "Le raccourci « %s » est attribué deux fois — les raccourcis doivent être uniques",
		"cfg.err.recovered": "config.json est corrompu (%s).\nLe fichier a été enregistré sous %s et les réglages ont été réinitialisés.",
		"err.disk.space": "espace disque faible : %d Mo libres, ~%d Mo nécessaires",
		"err.save": "réglages non enregistrés : %s — les anciens sont conservés",
		"err.port": "le port %d ne convient pas : choisissez un nombre entre 1024 et 65535",
		"err.nolangs": "laissez au moins une langue pour la question de traduction",
		"ov.mic.lost": "Micro débranché — enregistrement interrompu",
		"err.hash": "le fichier téléchargé est endommagé — réessayez",
		"models.check.ok": "Modèles vérifiés : %d — tous les fichiers sont intacts",
		"models.check.none": "Rien à vérifier — aucun modèle installé n'a d'empreinte de référence",
		"models.check.bad": "Fichiers endommagés : %s — téléchargez le modèle à nouveau",
		"ov.paused": "Pause",
		"status.paused": "Pause — l'enregistrement attend",
		"hist.insert.gone": "entrée introuvable",
		"hist.insert.nowin": "rien où coller — le texte est dans le presse-papiers",
		"hist.insert.ok": "collé dans « %s »",
		"lists.bad": "ce fichier ne convient pas",
		"lists.saved": "enregistré dans %s",
		"lists.added": "ajoutés : %d, ignorés : %d",
		"lists.save.title": "Où enregistrer les listes",
		"lists.open.title": "Quel fichier charger",
		"un.title": "{app} — Désinstallation", "un.confirm": "Supprimer {app} de cet ordinateur ?",
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
		"err.webview":        "La ventana de ajustes requiere Microsoft WebView2 Runtime (incluido en Windows 11).\nInstálelo o edite config.json manualmente.",
		"status.loading":     "Cargando modelo…", "status.ready": "Listo — mantenga %s y hable",
		"status.recording": "Grabando…", "status.transcribing": "Reconociendo…", "status.disabled": "Desactivado",
		"status.server.restart": "Servidor caído, reiniciando…", "status.cfg.err": "Error en config.json (vea el log)",
		"menu.settings":         "Ajustes…", "menu.enable": "Activar", "menu.disable": "Desactivar",
		"menu.reload": "Recargar config.json", "menu.open.config": "Abrir config.json", "menu.open.log": "Abrir log",
		"menu.about": "Acerca de", "menu.quit": "Salir",
		"ov.speak": "Hable…", "ov.transcribing": "Reconociendo", "ov.inserted": "Insertado: %d caracteres",
		"ov.err.recognize": "Error de reconocimiento (vea el log)", "ov.err.paste": "No se pegó: el texto está en el último resultado",
		"ov.moved": "La ventana cambió: el texto está en el portapapeles",
		"copy.ok": "Copiado",
		"copy.none": "Nada que copiar",
		"copy.fail": "No se pudo copiar: %s",
		"mic.busy": "Hay un dictado en curso, ahora no se puede comprobar", "mic.check.ok": "Se oye bien: pico %.0f dB, voz en el %.0f%% de la grabación",
		"mic.check.quiet": "Demasiado bajo: pico %.0f dB — sube el nivel del micrófono en Windows o acércate", "mic.check.clipped": "Saturación: %.1f%% de muestras recortadas — baja el nivel del micrófono", "mic.check.silent": "No se oye voz — comprueba que el micrófono elegido es el correcto y no está silenciado",
		"ov.quiet": "Demasiado bajo, casi no se oyó nada", "ov.clipped": "Saturación: el sonido se recortó",
		"ov.cmd.cancelled": "Cancelado por voz",
		"ov.silence": "Silencio — nada reconocido", "ov.server.loading": "El servidor aún carga",
		"ov.tooshort": "Demasiado corto: mantén las teclas más tiempo",
		"ov.cancelled": "Cancelado", "ov.editing": "Editando: %s", "ov.translating": "Traduciendo",
		"ov.llm.needed": "Este idioma requiere el módulo LLM", "td.title": "Traducir a:", "td.plain": "Sin traducción",
		"cap.title": "{app} — atajo", "cap.prompt": "Pulse una nueva combinación…\n\n(actual: %s)\n\nEsc — cancelar",
		"cap.selected": "Elegido: %s", "cap.cancelled": "Cancelado",
		"model.switching": "Cambiando modelo — reiniciando…", "model.del.active": "No se puede borrar el modelo activo",
		"model.del.ok": "Modelo borrado",
		"about.text":   "{app} %s\n\nVoz → texto en la posición del cursor.\nColoque el cursor, mantenga %s, hable, suelte — el texto se inserta.\n\nReconocimiento: whisper.cpp, totalmente local y sin conexión.\nModelo: %s (idioma: %s)\n\nAjustes: clic en el icono o config.json.\nLogs: {log} (máx ~2 MB).",
		"status.nomodel": "No hay ningún modelo descargado — elige uno en los ajustes",
		"menu.lastcopy": "Copiar el último resultado",
		"ov.copied": "Copiado al portapapeles", "ov.kept": "Cancelado — el texto queda en el último resultado",
		"ov.llm.skipped": "Insertado sin el perfil «%s»",
		"fd.title": "Cambió el foco, ¿insertar?", "fd.here": "Insertar aquí", "fd.copy": "Copiar",
		"ov.err.mic": "Micrófono no disponible — revisa el dispositivo en los ajustes",
		"ov.notranslate": "El modelo activo no traduce — insertado tal como se reconoció",
		"ov.engine.fallback": "El otro motor no arrancó — se sigue con el actual",
		"route.speech": "Habla en %s", "route.other": "Otros idiomas", "route.translate": "Traducción",
		"route.lang.auto": "cualquier idioma",
		"route.why.language": "más preciso aquí, con puntuación", "route.why.otherlang": "99 idiomas",
		"route.why.translate": "solo Whisper traduce", "route.why.notinstalled": "el modelo ruso no está instalado",
		"route.why.unknownlang": "sin idioma definido — solo Whisper lo detecta", "route.why.forced": "forzado en config.json",
		"adv.pick": "Recomiendo %s.", "adv.companion": "%s es un buen acompañante — para otros idiomas y la traducción.",
		"adv.ram": "%d MB libres", "status.line": "Listo · %s · %.1f GB libres",
		"ago.now": "ahora mismo", "ago.min": "hace %d min", "ago.hour": "hace %d h",
		"chars": "%d caracteres", "inserted.into": "insertado en %s",
		"punct.prompt": "Añade puntuación y mayúsculas. No cambies las palabras, no traduzcas, no añadas nada. Devuelve solo el texto corregido.",
		"err.sherpa.notfound": "no se encuentra el reconocedor sherpa: %s",
		"err.sherpa.start": "sherpa-server se cerró al arrancar (mira el registro)",
		"err.sherpa.translate": "este modelo no puede traducir",
		"err.sherpa.model": "No se encuentra el archivo del modelo: %s — descárgalo en los ajustes o corrige sherpa_model en config.json",
		"err.hotkey.dup": "El atajo «%s» está asignado dos veces — los atajos deben ser únicos",
		"cfg.err.recovered": "config.json está dañado (%s).\nEl archivo se guardó como %s y los ajustes volvieron a los valores por defecto.",
		"err.disk.space": "poco espacio en disco: %d MB libres, ~%d MB necesarios",
		"err.save": "ajustes no guardados: %s — se mantienen los anteriores",
		"err.port": "el puerto %d no sirve: elige un número entre 1024 y 65535",
		"err.nolangs": "deja al menos un idioma para la pregunta de traducción",
		"ov.mic.lost": "Micrófono desconectado: grabación interrumpida",
		"err.hash": "el archivo descargado está dañado: inténtalo de nuevo",
		"models.check.ok": "Modelos comprobados: %d, todos los archivos están intactos",
		"models.check.none": "Nada que comprobar: ningún modelo instalado tiene hash de referencia",
		"models.check.bad": "Archivos dañados: %s — descarga el modelo de nuevo",
		"ov.paused": "Pausa",
		"status.paused": "Pausa: la grabación espera",
		"hist.insert.gone": "no se encuentra la entrada",
		"hist.insert.nowin": "no hay dónde pegar: el texto está en el portapapeles",
		"hist.insert.ok": "pegado en «%s»",
		"lists.bad": "este archivo no sirve",
		"lists.saved": "guardado en %s",
		"lists.added": "añadidos: %d, omitidos: %d",
		"lists.save.title": "Dónde guardar las listas",
		"lists.open.title": "Qué archivo cargar",
		"un.title": "{app} — Desinstalación", "un.confirm": "¿Quitar {app} de este equipo?",
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
		"err.webview":        "La finestra impostazioni richiede Microsoft WebView2 Runtime (incluso in Windows 11).\nInstallalo o modifica config.json manualmente.",
		"status.loading":     "Caricamento modello…", "status.ready": "Pronto — tieni %s e parla",
		"status.recording": "Registrazione…", "status.transcribing": "Riconoscimento…", "status.disabled": "Disattivato",
		"status.server.restart": "Server bloccato, riavvio…", "status.cfg.err": "Errore in config.json (vedi log)",
		"menu.settings":         "Impostazioni…", "menu.enable": "Attiva", "menu.disable": "Disattiva",
		"menu.reload": "Ricarica config.json", "menu.open.config": "Apri config.json", "menu.open.log": "Apri log",
		"menu.about": "Informazioni", "menu.quit": "Esci",
		"ov.speak": "Parla…", "ov.transcribing": "Riconoscimento", "ov.inserted": "Inserito: %d caratteri",
		"ov.err.recognize": "Errore di riconoscimento (vedi log)", "ov.err.paste": "Non incollato: il testo è nell'ultimo risultato",
		"ov.moved": "La finestra è cambiata: il testo è negli appunti",
		"copy.ok": "Copiato",
		"copy.none": "Niente da copiare",
		"copy.fail": "Impossibile copiare: %s",
		"mic.busy": "C'è una dettatura in corso, ora non si può controllare", "mic.check.ok": "Si sente bene: picco %.0f dB, voce nel %.0f%% della registrazione",
		"mic.check.quiet": "Troppo basso: picco %.0f dB — alza il livello del microfono in Windows o avvicinati", "mic.check.clipped": "Distorsione: %.1f%% dei campioni tagliati — abbassa il livello del microfono", "mic.check.silent": "Nessuna voce sentita — controlla che sia scelto il microfono giusto e non sia muto",
		"ov.quiet": "Troppo basso, non si è sentito quasi nulla", "ov.clipped": "Distorsione: il suono è stato tagliato",
		"ov.cmd.cancelled": "Annullato a voce",
		"ov.silence": "Silenzio — nulla riconosciuto", "ov.server.loading": "Il server sta ancora caricando",
		"ov.tooshort": "Troppo breve — tieni premuti i tasti più a lungo",
		"ov.cancelled": "Annullato", "ov.editing": "Modifica: %s", "ov.translating": "Traduzione",
		"ov.llm.needed": "Questa lingua richiede il modulo LLM", "td.title": "Traduci in:", "td.plain": "Senza traduzione",
		"cap.title": "{app} — scorciatoia", "cap.prompt": "Premi una nuova combinazione…\n\n(attuale: %s)\n\nEsc — annulla",
		"cap.selected": "Scelto: %s", "cap.cancelled": "Annullato",
		"model.switching": "Cambio modello — riavvio…", "model.del.active": "Impossibile eliminare il modello attivo",
		"model.del.ok": "Modello eliminato",
		"about.text":   "{app} %s\n\nVoce → testo alla posizione del cursore.\nPosiziona il cursore, tieni %s, parla, rilascia — il testo viene inserito.\n\nRiconoscimento: whisper.cpp, completamente locale e offline.\nModello: %s (lingua: %s)\n\nImpostazioni: clic sull'icona o config.json.\nLog: {log} (max ~2 MB).",
		"status.nomodel": "Nessun modello di riconoscimento scaricato — scegline uno nelle impostazioni",
		"menu.lastcopy": "Copia l'ultimo risultato",
		"ov.copied": "Copiato negli appunti", "ov.kept": "Annullato — il testo resta nell'ultimo risultato",
		"ov.llm.skipped": "Inserito senza il profilo «%s»",
		"fd.title": "Il focus è cambiato — inserire?", "fd.here": "Inserisci qui", "fd.copy": "Copia",
		"ov.err.mic": "Microfono non disponibile — controlla il dispositivo nelle impostazioni",
		"ov.notranslate": "Il modello attivo non traduce — inserito così com'è stato riconosciuto",
		"ov.engine.fallback": "L'altro motore non è partito — si resta su quello attuale",
		"route.speech": "Parlato in %s", "route.other": "Altre lingue", "route.translate": "Traduzione",
		"route.lang.auto": "qualsiasi lingua",
		"route.why.language": "qui è più preciso, con la punteggiatura", "route.why.otherlang": "99 lingue",
		"route.why.translate": "solo Whisper traduce", "route.why.notinstalled": "il modello russo non è installato",
		"route.why.unknownlang": "nessuna lingua impostata — solo Whisper la riconosce", "route.why.forced": "forzato in config.json",
		"adv.pick": "Consiglio %s.", "adv.companion": "%s è un buon compagno — per le altre lingue e la traduzione.",
		"adv.ram": "%d MB liberi", "status.line": "Pronto · %s · %.1f GB liberi",
		"ago.now": "proprio ora", "ago.min": "%d min fa", "ago.hour": "%d h fa",
		"chars": "%d caratteri", "inserted.into": "inserito in %s",
		"punct.prompt": "Aggiungi punteggiatura e maiuscole. Non cambiare le parole, non tradurre, non aggiungere nulla. Restituisci solo il testo corretto.",
		"err.sherpa.notfound": "riconoscitore sherpa non trovato: %s",
		"err.sherpa.start": "sherpa-server si è chiuso durante l'avvio (vedi il registro)",
		"err.sherpa.translate": "questo modello non sa tradurre",
		"err.sherpa.model": "File del modello non trovato: %s — scaricalo nelle impostazioni o correggi sherpa_model in config.json",
		"err.hotkey.dup": "La scorciatoia «%s» è assegnata due volte — le scorciatoie devono essere uniche",
		"cfg.err.recovered": "config.json è danneggiato (%s).\nIl file è stato salvato come %s e le impostazioni sono tornate ai valori predefiniti.",
		"err.disk.space": "poco spazio su disco: %d MB liberi, ~%d MB necessari",
		"err.save": "impostazioni non salvate: %s — restano quelle di prima",
		"err.port": "la porta %d non va bene: scegli un numero tra 1024 e 65535",
		"err.nolangs": "lascia almeno una lingua per la domanda sulla traduzione",
		"ov.mic.lost": "Microfono scollegato: registrazione interrotta",
		"err.hash": "il file scaricato è danneggiato — riprova",
		"models.check.ok": "Modelli controllati: %d — tutti i file sono integri",
		"models.check.none": "Niente da controllare: nessun modello installato ha un hash di riferimento",
		"models.check.bad": "File danneggiati: %s — scarica di nuovo il modello",
		"ov.paused": "Pausa",
		"status.paused": "Pausa — la registrazione aspetta",
		"hist.insert.gone": "voce non trovata",
		"hist.insert.nowin": "non c'è dove incollare: il testo è negli appunti",
		"hist.insert.ok": "incollato in «%s»",
		"lists.bad": "questo file non va bene",
		"lists.saved": "salvato in %s",
		"lists.added": "aggiunti: %d, saltati: %d",
		"lists.save.title": "Dove salvare le liste",
		"lists.open.title": "Quale file caricare",
		"un.title": "{app} — Disinstallazione", "un.confirm": "Rimuovere {app} da questo computer?",
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
		"err.webview":        "Okno ustawień wymaga Microsoft WebView2 Runtime (zawarty w Windows 11).\nZainstaluj go lub edytuj config.json ręcznie.",
		"status.loading":     "Ładowanie modelu…", "status.ready": "Gotowy — przytrzymaj %s i mów",
		"status.recording": "Nagrywanie…", "status.transcribing": "Rozpoznawanie…", "status.disabled": "Wyłączone",
		"status.server.restart": "Serwer padł, restart…", "status.cfg.err": "Błąd w config.json (zobacz log)",
		"menu.settings":         "Ustawienia…", "menu.enable": "Włącz", "menu.disable": "Wyłącz",
		"menu.reload": "Przeładuj config.json", "menu.open.config": "Otwórz config.json", "menu.open.log": "Otwórz log",
		"menu.about": "O programie", "menu.quit": "Zakończ",
		"ov.speak": "Mów…", "ov.transcribing": "Rozpoznawanie", "ov.inserted": "Wstawiono: %d znaków",
		"ov.err.recognize": "Błąd rozpoznawania (zobacz log)", "ov.err.paste": "Nie wklejono — tekst jest w ostatnim wyniku",
		"ov.moved": "Okno się zmieniło — tekst jest w schowku",
		"copy.ok": "Skopiowano",
		"copy.none": "Nie ma czego kopiować",
		"copy.fail": "Nie udało się skopiować: %s",
		"mic.busy": "Trwa dyktowanie, teraz nie można sprawdzić", "mic.check.ok": "Brzmi dobrze: szczyt %.0f dB, mowa w %.0f%% nagrania",
		"mic.check.quiet": "Za cicho: szczyt %.0f dB — podnieś poziom mikrofonu w Windows albo usiądź bliżej", "mic.check.clipped": "Przesterowanie: obcięto %.1f%% próbek — zmniejsz poziom mikrofonu", "mic.check.silent": "Nie słychać mowy — sprawdź, czy wybrany jest właściwy mikrofon i czy nie jest wyciszony",
		"ov.quiet": "Za cicho, prawie nic nie było słychać", "ov.clipped": "Przesterowanie — dźwięk został obcięty",
		"ov.cmd.cancelled": "Anulowano głosem",
		"ov.silence": "Cisza — nic nie rozpoznano", "ov.server.loading": "Serwer wciąż się ładuje",
		"ov.tooshort": "Za krótko — przytrzymaj klawisze dłużej",
		"ov.cancelled": "Anulowano", "ov.editing": "Edycja: %s", "ov.translating": "Tłumaczenie",
		"ov.llm.needed": "Ten język wymaga modułu LLM", "td.title": "Tłumacz na:", "td.plain": "Bez tłumaczenia",
		"cap.title": "{app} — skrót", "cap.prompt": "Naciśnij nową kombinację…\n\n(obecna: %s)\n\nEsc — anuluj",
		"cap.selected": "Wybrano: %s", "cap.cancelled": "Anulowano",
		"model.switching": "Zmiana modelu — restart…", "model.del.active": "Nie można usunąć aktywnego modelu",
		"model.del.ok": "Model usunięty",
		"about.text":   "{app} %s\n\nGłos → tekst w pozycji kursora.\nUstaw kursor, przytrzymaj %s, mów, puść — tekst zostanie wstawiony.\n\nRozpoznawanie: whisper.cpp, w pełni lokalnie i offline.\nModel: %s (język: %s)\n\nUstawienia: kliknij ikonę lub config.json.\nLogi: {log} (maks. ~2 MB).",
		"status.nomodel": "Nie pobrano żadnego modelu — wybierz go w ustawieniach",
		"menu.lastcopy": "Kopiuj ostatni wynik",
		"ov.copied": "Skopiowano do schowka", "ov.kept": "Anulowano — tekst został w ostatnim wyniku",
		"ov.llm.skipped": "Wstawiono bez profilu „%s”",
		"fd.title": "Zmieniło się okno — wstawić?", "fd.here": "Wstaw tutaj", "fd.copy": "Kopiuj",
		"ov.err.mic": "Mikrofon niedostępny — sprawdź urządzenie w ustawieniach",
		"ov.notranslate": "Aktywny model nie tłumaczy — wstawiono tak, jak rozpoznano",
		"ov.engine.fallback": "Drugi silnik nie wystartował — zostaje bieżący",
		"route.speech": "Mowa w %s", "route.other": "Inne języki", "route.translate": "Tłumaczenie",
		"route.lang.auto": "dowolny język",
		"route.why.language": "tu dokładniej, ze znakami", "route.why.otherlang": "99 języków",
		"route.why.translate": "tylko Whisper tłumaczy", "route.why.notinstalled": "model rosyjski nie jest zainstalowany",
		"route.why.unknownlang": "nie ustawiono języka — rozpozna go tylko Whisper", "route.why.forced": "wymuszone w config.json",
		"adv.pick": "Polecam %s.", "adv.companion": "%s dobrze uzupełni — dla innych języków i tłumaczenia.",
		"adv.ram": "%d MB wolnych", "status.line": "Gotowe · %s · %.1f GB wolnych",
		"ago.now": "przed chwilą", "ago.min": "%d min temu", "ago.hour": "%d godz. temu",
		"chars": "%d znaków", "inserted.into": "wstawiono do %s",
		"punct.prompt": "Dodaj znaki interpunkcyjne i wielkie litery. Nie zmieniaj słów, nie tłumacz, nic nie dopisuj. Zwróć tylko poprawiony tekst.",
		"err.sherpa.notfound": "nie znaleziono rozpoznawania sherpa: %s",
		"err.sherpa.start": "sherpa-server zakończył się przy starcie (zobacz dziennik)",
		"err.sherpa.translate": "ten model nie umie tłumaczyć",
		"err.sherpa.model": "Nie znaleziono pliku modelu: %s — pobierz go w ustawieniach albo popraw sherpa_model w config.json",
		"err.hotkey.dup": "Skrót „%s” jest przypisany dwa razy — skróty muszą być unikalne",
		"cfg.err.recovered": "config.json jest uszkodzony (%s).\nPlik zapisano jako %s, a ustawienia wróciły do domyślnych.",
		"err.disk.space": "mało miejsca na dysku: %d MB wolnych, potrzeba ~%d MB",
		"err.save": "ustawienia nie zapisane: %s — zostają poprzednie",
		"err.port": "port %d się nie nadaje: wybierz numer od 1024 do 65535",
		"err.nolangs": "zostaw przynajmniej jeden język do pytania o tłumaczenie",
		"ov.mic.lost": "Mikrofon odłączony — nagranie przerwane",
		"err.hash": "pobrany plik jest uszkodzony — spróbuj ponownie",
		"models.check.ok": "Sprawdzone modele: %d — wszystkie pliki są całe",
		"models.check.none": "Nie ma czego sprawdzać — żaden zainstalowany model nie ma wzorcowego skrótu",
		"models.check.bad": "Uszkodzone pliki: %s — pobierz model ponownie",
		"ov.paused": "Pauza",
		"status.paused": "Pauza — nagranie czeka",
		"hist.insert.gone": "nie znaleziono wpisu",
		"hist.insert.nowin": "nie ma gdzie wkleić — tekst jest w schowku",
		"hist.insert.ok": "wklejono w „%s”",
		"lists.bad": "ten plik nie pasuje",
		"lists.saved": "zapisano w %s",
		"lists.added": "dodano: %d, pominięto: %d",
		"lists.save.title": "Gdzie zapisać listy",
		"lists.open.title": "Który plik wczytać",
		"un.title": "{app} — Odinstalowanie", "un.confirm": "Usunąć {app} z tego komputera?",
		"un.data": "Usunąć także ustawienia i pobrane modele?", "un.done": "{app} został usunięty.",
		"srv.restarting": "Ponowne uruchamianie rozpoznawania z nowymi ustawieniami…",
	}

	settingsStrings["de"] = map[string]string{
		"S_TITLE": "{app} — Einstellungen", "S_TAB_GENERAL": "Allgemein", "S_TAB_REC": "Erkennung",
		"S_TAB_PROC": "Nachbearbeitung", "S_TAB_SERVER": "Server", "S_TAB_ABOUT": "Über",
		"S_PIPE": "Stimme ▸ Erkennung ▸ Bearbeitung ▸ Einfügen",
		"S_DICT": "Erkennungswörterbuch", "S_DICT_HINT": "Begriffe, Namen und Abkürzungen, durch Kommas getrennt — ein Hinweis fürs Gehör, keine Befehle.",
		"S_TR": "Übersetzung", "S_TR_HINT": "Übersetzung durch Whisper: nach Englisch — nativ, andere Sprachen — durch erzwungene Ausgabesprache (Qualität variiert).",
		"S_TR_DEFAULT": "Immer in die Zielsprache übersetzen", "S_TR_TARGET": "Zielsprache", "S_TR_ASK": "Sprachauswahl", "S_TR_ASK_NEVER": "Nicht fragen (Standard verwenden)",
		"S_TR_ASK_ALWAYS": "Jedes Mal fragen", "S_TR_ASK_TIMEOUT": "Mit Timeout fragen", "S_TR_SECONDS": "Timeout, Sek.",
		"S_TR_LANGS": "Sprachen im Dialog",
		"S_LLM":      "Verarbeitungsmodi", "S_LLM_HINT": "Der Standardmodus läuft auf dem Haupt-Shortcut; ein Modus kann eigenen Hotkey haben. Prompt-Profile bearbeitet ein zweites neuronales Netz (offline).",
		"S_PROF_ASIS": "Unverändert", "S_PROF_WT": "Übersetzen → English (schnell)",
		"S_PROF_ADD": "Profil hinzufügen", "S_PROF_NAME": "Name", "S_PROF_PROMPT": "Prompt", "S_PROF_HOTKEY": "Hotkey",
		"S_PROF_SET": "Festlegen…", "S_PROF_CLEAR": "Löschen", "S_PROF_TEST": "Test",
		"S_HOTKEY": "Tastenkürzel", "S_CHANGE": "Ändern…", "S_UILANG": "Sprache der Oberfläche", "S_AUTO": "Wie im System",
		"S_SEC_SOUND": "Ton", "S_SEC_BEHAVIOR": "Verhalten", "S_BEEP": "Tonsignale der Aufnahme", "S_SOUND": "Signalton",
		"S_SND_SPEECH": "System (Sprache)", "S_SND_CHIME": "Glöckchen", "S_SND_SOFT": "Sanft", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Enter nach dem Einfügen drücken (Auto-Senden)", "S_RESTORE": "Zwischenablage nach Einfügen wiederherstellen",
		"S_NAV_HISTORY": "Verlauf", "S_HIST_ON": "Verlauf der Diktate führen", "S_HIST_ON_SUB": "nur Text, auf diesem Rechner; Ton wird nie gespeichert",
		"S_HIST_DAYS": "Wie viele Tage aufbewahren", "S_HIST_MAX": "Wie viele Einträge aufbewahren",
		"S_HIST_SKIP": "Aus diesen Programmen nie aufzeichnen", "S_HIST_SKIP_SUB": "durch Komma getrennt: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Einträge", "S_HIST_CLEAR": "Leeren", "S_HIST_COPY": "Kopieren",
		"S_HIST_FIND": "Im Verlauf suchen…", "S_HIST_EMPTY": "Noch kein Verlauf", "S_HIST_ASK": "Den gesamten Diktatverlauf löschen?",
		"S_SEC_CMD": "Sprachbefehle", "S_CMD_HINT": "Gesagtes wird zu einem Zeilenumbruch, einem Zeichen oder einem Abbruch, statt im Text zu landen. Als ganze Wörter erkannt, von oben nach unten angewendet, nach den Ersetzungen.",
		"S_CMD_ADD": "Befehl hinzufügen", "S_CMD_PRESET": "Übliche hinzufügen", "S_CMD_PH": "neue Zeile",
		"S_CMD_NEWLINE": "Zeilenumbruch", "S_CMD_PARAGRAPH": "neuer Absatz", "S_CMD_TEXT": "Text einfügen", "S_CMD_CANCEL": "Diktat abbrechen",
		"S_CMD_TEXT_PH": "was einfügen", "S_CMD_EMPTY": "Noch keine Befehle", "S_CMD_DEL": "Befehl löschen",
		"S_CMD_P_NEWLINE": "neue Zeile", "S_CMD_P_PARAGRAPH": "neuer Absatz", "S_CMD_P_CANCEL": "abbrechen",
		"S_SEC_REPLACE": "Ersetzungen nach der Erkennung", "S_REPLACE_HINT": "Falsch Gehörtes wird zu dem, was gemeint war — direkt nach der Erkennung, vor den Prompts. Von oben nach unten angewendet.",
		"S_REPL_ADD": "Ersetzung hinzufügen", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "ganze Wörter", "S_REPL_CASE": "Groß-/Kleinschreibung", "S_REPL_EMPTY": "Noch keine Ersetzungen",
		"S_REPL_DEL": "Ersetzung löschen", "S_REPL_TEST_PH": "Satz eingeben, um Ersetzungen und Befehle zu prüfen",
		"S_SEC_RULES": "Regeln pro Programm", "S_RULES_HINT": "Für einzelne Programme kann das Einfügen anders laufen. Die erste passende Regel gewinnt.",
		"S_RULE_ADD": "Regel hinzufügen", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "Einfügen: wie sonst", "S_RULE_ENTER_INH": "Enter: wie sonst", "S_RULE_DELAY_NONE": "ohne Verzögerung", "S_RULE_PROMPT_INH": "Prompts: wie sonst",
		"S_RULE_CLIP": "Zwischenablage", "S_RULE_TYPE": "zeichenweise", "S_RULE_ENTER_ON": "mit Enter", "S_RULE_ENTER_OFF": "ohne Enter",
		"S_RULE_NOPROMPT": "ohne Prompts", "S_RULE_LAST": "zuletzt eingefügt in: %s", "S_RULE_EMPTY": "Noch keine Regeln",
		"S_RULE_DEL": "Regel löschen", "S_RULE_PROMPTS": "Prompts",
		"S_PASTE_DELAY": "Verzögerung vor dem Einfügen", "S_PASTE_DELAY_SUB": "wenn das Programm den Text noch nicht annimmt",
		"S_OVPOS": "Wo die Leiste erscheint", "S_OVPOS_SUB": "am Cursor — neben der Eingabestelle; zeigt die App sie nicht, dann neben dem Mauszeiger",
		"S_OVPOS_BOTTOM": "Unten am Bildschirm", "S_OVPOS_TOP": "Oben am Bildschirm", "S_OVPOS_CARET": "Am Cursor",
		"S_OVTEXT": "Erkannten Text anzeigen", "S_OVTEXT_SUB": "auf der Leiste nach dem Einfügen, statt der Zeichenzahl",
		"S_OVERLAY": "Bildschirmanzeige", "S_ANIM": "Aufnahme-/Erkennungsanimation", "S_TYPEMODE": "Zeichenweise Eingabe (für Felder ohne Einfügen)",
		"S_RECLANG": "Erkennungssprache", "S_RECAUTO": "Automatisch",
		"S_MODELS": "Erkennungsmodelle", "S_DL": "Laden", "S_DEL": "Löschen",
		"S_M_BASE": "schnell, für schwache PCs", "S_M_SMALL": "ausgewogen", "S_M_MED": "genauer, empfohlen", "S_M_TURBO": "beste Genauigkeit auf CPU",
		"S_M_CUSTOM": "benutzerdefiniert (aus config.json)",
		"S_THREADS":  "CPU-Threads", "S_MINMS": "Min. Aufnahme, ms", "S_MAXSEC": "Max. Aufnahme, s",
		"S_AUTOSTART": "whisper-server automatisch starten", "S_PORT": "Port", "S_SERVEREXE": "Pfad zu whisper-server",
		"S_SERVERURL": "Externer Server (URL)", "S_URLHINT": "Falls gesetzt, wird kein lokaler Server gestartet",
		"S_SAVED": "Gespeichert",
		"S_ABOUT_HTML": "<p><b>Stimme → Text an der Cursorposition.</b></p><p>Cursor in ein Eingabefeld setzen, Shortcut halten, sprechen, loslassen — der Text wird eingefügt.</p><p>Vollständig lokal und offline. Technik: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; Modelle von Hugging Face.</p><p>Logs überschreiten nie ~2 MB.</p>",
		"S_VERSION":    "Version",
		"S_LVL_SIMPLE": "einfach", "S_LVL_ALL": "alles", "S_SEARCH": "Einstellung finden…",
		"S_GRP_WORK": "Arbeit", "S_GRP_REC": "Erkennung", "S_GRP_OTHER": "Sonstiges",
		"S_NAV_STATE": "Status", "S_NAV_DICT": "Diktat", "S_NAV_MIC": "Mikrofon", "S_NAV_MODELS": "Modelle",
		"S_NAV_TEXT": "Text", "S_NAV_TR": "Übersetzung", "S_NAV_SYSTEM": "System", "S_NAV_ABOUT": "Über",
		"S_STATE_HINT": "halten und sprechen — der Text landet dort, wo der Cursor steht",
		"S_STATE_RU": "Russische Sprache", "S_STATE_OTHER": "Andere Sprachen", "S_STATE_PROC": "Nachbearbeitung",
		"S_CHANGE_MODEL": "Wechseln", "S_PICK_MODEL": "Auswählen", "S_STATE_GET": "Laden",
		"S_RETRY": "Erneut versuchen", "S_BERR_OPEN": "Servereinstellungen öffnen",
		"S_STATE_LAST": "Letztes Diktat", "S_STATE_COPY": "Kopieren", "S_STATE_MEM": "Speicher",
		"S_STATE_MEM_SUB": "Modelle bleiben geladen, der erste Satz kommt ohne Verzögerung",
		"S_HOTMODE": "Modus", "S_HOTMODE_HOLD": "halten", "S_HOTMODE_TOGGLE": "umschalten",
		"S_SUB_HOTMODE": "Tasten halten oder einmal drücken zum Starten und noch einmal zum Beenden",
		"S_SUB_MINMS": "ignoriert versehentliche Tastendrücke",
		"S_SUB_ENTER": "schickt die Nachricht sofort ab",
		"S_SUB_CLIP": "Bilder und Dateien kommen unverändert zurück",
		"S_SUB_TYPE": "hilft dort, wo ein Feld das Einfügen verweigert",
		"S_SEC_OVERLAY": "Bildschirmanzeige",
		"S_MIC_CHECK": "Mikrofon prüfen", "S_MIC_CHECK_SUB": "drei Sekunden Aufnahme und ein Urteil: Pegel, Übersteuerung, ob Sprache da ist", "S_MIC_CHECKING": "Prüfe…",
		"S_PAUSE": "Aufnahme pausieren", "S_PAUSE_SUB": "im Umschaltmodus: einmal drücken und die Aufnahme hält an, noch einmal und sie läuft weiter",
		"S_MCHECK": "Installierte Modelle prüfen", "S_MCHECK_SUB": "vergleicht die Modelldateien mit den Referenz-Hashes", "S_MCHECK_GO": "Prüfen", "S_MCHECK_RUN": "Prüfe…",
		"S_HIST_INSERT": "Einfügen",
		"S_LISTS_HINT": "Ersetzungen und Befehle in einer Datei — zum Mitnehmen auf einen anderen Rechner", "S_LISTS_EXPORT": "In Datei speichern", "S_LISTS_IMPORT": "Aus Datei laden",
		"S_MIC": "Mikrofon", "S_MIC_DEFAULT": "Systemstandard", "S_MIC_REFRESH": "Liste aktualisieren",
		"S_MIC_LEVEL": "Eingangspegel", "S_MIC_QUIET": "still",
		"S_ADV_TITLE": "Modell auswählen", "S_F_ALL": "alle", "S_F_RU": "Russisch",
		"S_F_MULTI": "viele Sprachen", "S_F_PUNCT": "setzt Satzzeichen", "S_F_FIT": "passt in den Speicher",
		"S_ADV_LANGQ": "In welcher Sprache diktieren Sie", "S_ADV_PRIOQ": "Was ist wichtiger",
		"S_ADV_ACC": "Wegen der Genauigkeit gewählt.", "S_ADV_SPEED": "Wegen der Geschwindigkeit gewählt.",
		"S_ADV_TRQ": "Übersetzung nötig", "S_ADV_GO": "Empfehlen",
		"S_ADV_PRIMARY": "Haupt", "S_ADV_COMPANION": "zweites", "S_ADV_HAVE": "schon da", "S_ADV_APPLY": "Übernehmen",
		"S_ADV_ASK": "Geladen werden: %s — insgesamt %s. Starten?",
		"S_SUB_THREADS": "mehr Threads sind nicht immer schneller — messen Sie es auf Ihrem Rechner",
		"S_SEC_LLM": "Bearbeitungsmodell",
		"S_PUNCT": "Satzzeichen und Großschreibung", "S_SUB_PUNCT": "woher Satzzeichen und Großschreibung kommen",
		"S_PUNCT_MODEL": "vom Modell", "S_PUNCT_LLM": "vom Bearbeitungsmodell", "S_PUNCT_OFF": "entfernen",
		"S_SUB_DICT": "Wörterbuch", "S_SUB_PROMPTS": "Prompts",
		"S_TR_TURBO": "⚠ Das aktive Turbo-Modell ist nicht für die Übersetzung ins Englische trainiert — wählen Sie zum Übersetzen ein anderes Modell im Reiter „Modelle“.",
		"S_SUB_TRTARGET": "Englisch ist für Whisper nativ, andere Zielsprachen sind experimentell",
		"S_TR_EXP": "außer Englisch erzwingt die App die Ausgabesprache, statt zu übersetzen — der Text kann in der gesprochenen Sprache bleiben",
		"S_REMOTE_ABOUT": "Ein entfernter Server ist eingestellt: Audio wird dorthin gesendet, und das Versprechen oben gilt so lange nicht.",
		"S_UPD": "Aktualisierungen", "S_UPD_CHECK": "Nach Updates suchen", "S_UPD_AUTO": "Beim Start prüfen",
		"S_SUB_UPD": "die einzige Netzwerkanfrage außer Modell-Downloads",
		"S_UPD_NONE": "Sie haben die neueste Version", "S_UPD_AVAIL": "Version %s ist verfügbar.",
		"S_UPD_GO": "Aktualisieren", "S_UPD_ERR": "Update-Prüfung fehlgeschlagen", "S_UPD_DL": "Update wird geladen…",
		"S_SEC_SERVICE": "Dienst", "S_SUB_AUTOSTART": "ausschalten, wenn Sie den Server selbst starten",
		"S_SUB_PORT": "die Erkennung startet sich selbst neu",
		"S_MODEL_READY": "Modell geladen — wählen Sie es aus, um zu wechseln",
		"S_FIT_OK": "passt", "S_FIT_WARN": "knapp", "S_FIT_BAD": "zu wenig RAM", "S_RAM": "Arbeitsspeicher:",
		"S_HF_PH": "Modellname — z. B. qwen2.5 instruct",
		"S_NO_LLM": "Noch keine Modelle installiert — im Reiter „Suche“ eines finden und laden.",
		"S_NO_LLM_PROF": "Prompts stehen zur Verfügung, sobald ein Modell installiert ist (siehe Reiter „Modelle“ und „Suche“).",
		"S_UPDATED": "Modell zuletzt aktualisiert", "S_PROF_EDIT": "Bearbeiten", "S_PROF_CLOSE": "Einklappen",
		"S_CONFIRM_DEL": "Modell „%s“ löschen? Es kann erneut geladen werden.", "S_FREE": "frei",
		"S_DEL_ACTIVE": "Das aktive Modell „%s“ löschen? Die Erkennung hält an, bis Sie ein anderes wählen — herunterladen können Sie es gleich hier.",
		"S_WIZ_NEED_MODEL": "Laden Sie zuerst ein Modell — ohne Modell gibt es nichts zu erkennen",
		"S_REMOTE_WARN": "Audio wird an diesen Server gesendet. Der lokale Modus ist aus.",
		"S_REMOTE_ASK": "Audio wird nicht mehr auf diesem Rechner verarbeitet, sondern an %s gesendet. Fernmodus einschalten?",
		"S_REMOTE_BADGE": "EXTERN",
		"S_OK": "Ja", "S_CANCEL": "Abbrechen", "S_DL_START": "Laden", "S_DL_CANCEL": "Download abbrechen",
		"S_DL_ASK": "Das Modell „%s“ ist nicht geladen (%s). Jetzt laden?",
		"S_NOT_FOUND": "nichts", "S_MORE": "%d weitere Einstellungen", "S_LESS": "%d Einstellungen einklappen",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autor und Entwickler von {app} — einem lokalen Diktierwerkzeug für Windows: Stimme wird direkt am Cursor zu Text, ohne Cloud, ohne Abo.</p>" +
			"<p>Das Projekt ist offen: Quellcode, Build-Pipeline und aktuelle Releases liegen auf GitHub.</p>" +
			"<ul>" +
			"<li>Repository: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Autorenprofil: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Fehler gefunden oder eine Idee — eröffnen Sie ein Issue im Repository.</p>",
		"S_HELP_HTML": "<p class=\"wh\">Wie es funktioniert</p>" +
			"<p>Kürzel halten — die Aufnahme läuft (die Leiste am unteren Bildschirmrand zeigt Ihren Pegel). Loslassen — der Ton wird erkannt, bei Bedarf übersetzt und durch die Prompts geschickt, und der fertige Text landet an der Cursorposition. Das ✕ auf der Leiste bricht in jedem Schritt ab.</p>" +
			"<p>Der ganze Weg: <b>Aufnahme → Erkennung (Whisper) → Übersetzung (falls aktiv) → Prompts (LLM) → Einfügen</b>. Jeder Schritt ist auf der Leiste zu sehen.</p>" +
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
			"<p class=\"wh\">Aufnahme pausieren</p>" +
			"<p>Im Umschaltmodus (einmal drücken startet, noch einmal stoppt) lässt sich ein eigenes Kürzel für die Pause festlegen: Reiter „Diktat“, Zeile „Aufnahme pausieren“. Ein Druck hält die Aufnahme an — die Leiste zeigt „Pause“ und es wird nichts aufgezeichnet; noch ein Druck und es geht weiter, alles vor der Pause bleibt erhalten. Die Längenbegrenzung greift während der Pause nicht.</p>" +
			"<p class=\"wh\">Aus dem Verlauf einfügen</p>" +
			"<p>Jeder Eintrag im Verlauf hat die Schaltfläche „Einfügen“: sie holt das Fenster zurück, aus dem Sie die Einstellungen geöffnet haben, und fügt den Text dort ein wie ein gewöhnliches Diktat. Gibt es kein solches Fenster, landet der Text einfach in der Zwischenablage, und das Programm sagt es.</p>" +
			"<p class=\"wh\">Die Listen in einer Datei</p>" +
			"<p>Ersetzungen und Sprachbefehle lassen sich in eine .json-Datei speichern und auf einem anderen Rechner laden — die Schaltflächen unter der Befehlsliste im Reiter „Text“. Beim Laden wird nichts überschrieben: nur die Zeilen, die noch fehlen, kommen hinzu, und das Programm nennt die Zahl der hinzugefügten und der übersprungenen.</p>" +
			"<p class=\"wh\">Unversehrtheit der Dateien</p>" +
			"<p>Zu jedem Modell aus dem Katalog gehört ein bekannter SHA-256-Referenz-Hash. Nach dem Herunterladen wird die Datei damit verglichen: passt sie nicht, wird sie gelöscht und der Download kann wiederholt werden. Die Schaltfläche „Prüfen“ im Reiter „Modelle“ vergleicht die bereits installierten Modelle genauso, und beim Update wird auch das heruntergeladene Installationsprogramm geprüft — eine fremde Datei wird nicht gestartet.</p>" +
			"<p class=\"wh\">Verlauf der Diktate</p>" +
			"<p>Der Abschnitt „Verlauf“ in der linken Spalte bewahrt auf, was Sie diktiert haben: nur Text, nur auf diesem Rechner, Ton wird nie gespeichert. Standardmäßig aus, eingeschaltet wird er mit einem Schalter an derselben Stelle. Einträge bleiben eine eingestellte Zahl von Tagen und bis zu einer eingestellten Anzahl, Älteres fällt von selbst heraus; „Aus diesen Programmen nie aufzeichnen“ listet durch Komma getrennt jene, aus denen nichts gespeichert werden soll — Passwortmanager, Banking. Die Suche greift auf Text und Programmnamen, die Schaltfläche neben einem Eintrag legt ihn in die Zwischenablage, und „Leeren“ entfernt alles samt der Datei <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Sprachbefehle</p>" +
			"<p>Unter den Ersetzungen auf der Registerkarte „Text“ steht eine Liste von Befehlen: Gesagtes wird zur Handlung statt zu Wörtern. „Neue Zeile“ und „neuer Absatz“ setzen einen Umbruch — Modelle tun das nie; „abbrechen“ verwirft das ganze Diktat und fügt nichts ein; „Text einfügen“ setzt beliebiges ein, auch ein Smiley. Die Schaltfläche neben der Liste füllt sie mit den üblichen Wendungen in der Sprache der Oberfläche. Befehle werden als ganze Wörter erkannt und laufen nach den Ersetzungen, damit Prompts und Übersetzung schon den fertigen Text bekommen. Überflüssige Leerzeichen um die Umbrüche verschwinden von selbst. Das Feld darunter probiert Ersetzungen und Befehle an jedem Satz: ein Umbruch erscheint als ⏎.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">Gleichgewicht aus Tempo und Genauigkeit</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">genauer, empfohlen</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">beste Genauigkeit auf der CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modelle</b> — der Katalog: Base (schnell, für schwache Rechner), Small (ausgewogen), Medium und Turbo (genauer, langsamer; „q5“ ist eine quantisierte Fassung — etwas kleiner und schneller, fast ohne Qualitätsverlust) sowie GigaAM v3 für Russisch. Der Radioknopf wählt das aktive Modell (gilt sofort, die Erkennung startet neu); bei einem fehlenden Modell fragt das Programm, ob es geladen werden soll.</li>" +
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
			"<li>Ein Prompt kann ein eigenes Kürzel haben: ein Diktat damit wendet nur diesen einen an. Der Stift ✎ öffnet den Editor: Name, Prompttext, Kürzel und ein Testfeld ▶, das eine Probe direkt aus den Einstellungen durch das laufende Modell schickt.</li>" +
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
			"<li><b>Portabel</b> — kopieren Sie einfach den ganzen Ordner mit der exe (auf einen USB-Stick, an einen anderen Rechner): Einstellungen, Modelle und Protokoll liegen daneben und reisen mit. In die Registry wird nichts geschrieben.</li>" +
			"<li>Ist beim ersten Start kein Erkennungsmodell da, öffnet das Programm den Katalog selbst und wartet auf den Download.</li>" +
			"<li>Voraussetzungen: Windows 10/11 x64, eine CPU mit AVX2 (etwa ab 2013), WebView2 Runtime für das Einstellungsfenster (in Windows 11 enthalten).</li>" +
			"</ul>" +
			"<p class=\"wh\">Menü im Infobereich und Dateien</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Bereit…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Einstellungen…</div><div class=\"mock-mi\">Ausschalten</div><div class=\"mock-mi\">Letztes Ergebnis kopieren</div><hr class=\"mock-sep\"><div class=\"mock-mi\">config.json neu lesen</div><div class=\"mock-mi\">config.json öffnen</div><div class=\"mock-mi\">Protokoll öffnen</div><div class=\"mock-mi\">Über</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Beenden</div></div>" +
			"<ul>" +
			"<li>Linksklick auf das Symbol im Infobereich — die Einstellungen; Rechtsklick — das Menü. Farben des Symbols: grün — bereit, rot — Aufnahme, orange — Erkennung, grau — ausgeschaltet oder Fehler.</li>" +
			"<li><b>config.json</b> — alle Einstellungen; von Hand geänderte Werte gelten nach „config.json neu lesen“ im Menü.</li>" +
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
		"S_AUTORUN":    "Mit Windows starten",
		"S_AUTORUN_SUB": "Ein Eintrag im Autostart des aktuellen Benutzers",
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
		"S_TITLE": "{app} — Réglages", "S_TAB_GENERAL": "Général", "S_TAB_REC": "Reconnaissance",
		"S_TAB_PROC": "Post-traitement", "S_TAB_SERVER": "Serveur", "S_TAB_ABOUT": "À propos",
		"S_PIPE": "voix ▸ reconnaissance ▸ édition ▸ insertion",
		"S_DICT": "Dictionnaire de reconnaissance", "S_DICT_HINT": "Termes, noms et abréviations séparés par des virgules — un indice pour l'oreille, pas des commandes.",
		"S_TR": "Traduction", "S_TR_HINT": "Traduction par Whisper : vers l'anglais — mode natif, autres langues — langue de sortie forcée (qualité variable).",
		"S_TR_DEFAULT": "Toujours traduire vers la langue cible", "S_TR_TARGET": "Langue cible", "S_TR_ASK": "Choix de langue", "S_TR_ASK_NEVER": "Ne pas demander (défaut)",
		"S_TR_ASK_ALWAYS": "Demander à chaque fois", "S_TR_ASK_TIMEOUT": "Demander avec délai", "S_TR_SECONDS": "Délai, s",
		"S_TR_LANGS": "Langues du dialogue",
		"S_LLM":      "Modes de traitement", "S_LLM_HINT": "Le mode par défaut s'applique au raccourci principal ; un mode peut avoir son propre raccourci. Les profils sont réécrits par un second réseau neuronal (hors ligne).",
		"S_PROF_ASIS": "Tel quel", "S_PROF_WT": "Traduire → English (rapide)",
		"S_PROF_ADD": "Ajouter un profil", "S_PROF_NAME": "Nom", "S_PROF_PROMPT": "Prompt", "S_PROF_HOTKEY": "Raccourci",
		"S_PROF_SET": "Définir…", "S_PROF_CLEAR": "Effacer", "S_PROF_TEST": "Test",
		"S_HOTKEY": "Raccourci clavier", "S_CHANGE": "Changer…", "S_UILANG": "Langue de l'interface", "S_AUTO": "Comme le système",
		"S_SEC_SOUND": "Son", "S_SEC_BEHAVIOR": "Comportement", "S_BEEP": "Signaux sonores", "S_SOUND": "Son du signal",
		"S_SND_SPEECH": "Système (voix)", "S_SND_CHIME": "Clochette", "S_SND_SOFT": "Doux", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Appuyer Entrée après insertion (envoi auto)", "S_RESTORE": "Restaurer le presse-papiers après insertion",
		"S_NAV_HISTORY": "Historique", "S_HIST_ON": "Conserver l'historique des dictées", "S_HIST_ON_SUB": "texte seul, sur cet ordinateur ; l'audio n'est jamais conservé",
		"S_HIST_DAYS": "Combien de jours conserver", "S_HIST_MAX": "Combien d'entrées conserver",
		"S_HIST_SKIP": "Ne jamais enregistrer depuis ces programmes", "S_HIST_SKIP_SUB": "séparés par des virgules : keepass.exe, 1password.exe",
		"S_HIST_LIST": "Entrées", "S_HIST_CLEAR": "Vider", "S_HIST_COPY": "Copier",
		"S_HIST_FIND": "Chercher dans l'historique…", "S_HIST_EMPTY": "Pas encore d'historique", "S_HIST_ASK": "Supprimer tout l'historique des dictées ?",
		"S_SEC_CMD": "Commandes vocales", "S_CMD_HINT": "Ce que vous dites devient un saut de ligne, un signe ou une annulation au lieu d'atterrir dans le texte. Reconnues en mots entiers, appliquées de haut en bas, après les remplacements.",
		"S_CMD_ADD": "Ajouter une commande", "S_CMD_PRESET": "Ajouter les habituelles", "S_CMD_PH": "nouvelle ligne",
		"S_CMD_NEWLINE": "saut de ligne", "S_CMD_PARAGRAPH": "nouveau paragraphe", "S_CMD_TEXT": "insérer du texte", "S_CMD_CANCEL": "annuler la dictée",
		"S_CMD_TEXT_PH": "quoi insérer", "S_CMD_EMPTY": "Aucune commande pour l'instant", "S_CMD_DEL": "Supprimer la commande",
		"S_CMD_P_NEWLINE": "nouvelle ligne", "S_CMD_P_PARAGRAPH": "nouveau paragraphe", "S_CMD_P_CANCEL": "annuler",
		"S_SEC_REPLACE": "Remplacements après la reconnaissance", "S_REPLACE_HINT": "Ce qui a été mal entendu devient ce que vous vouliez dire — juste après la reconnaissance, avant les prompts. Appliqués de haut en bas.",
		"S_REPL_ADD": "Ajouter un remplacement", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "mots entiers", "S_REPL_CASE": "casse", "S_REPL_EMPTY": "Aucun remplacement pour l'instant",
		"S_REPL_DEL": "Supprimer le remplacement", "S_REPL_TEST_PH": "tapez une phrase pour tester remplacements et commandes",
		"S_SEC_RULES": "Règles par application", "S_RULES_HINT": "L'insertion peut fonctionner autrement pour certains programmes. La première règle qui correspond l'emporte.",
		"S_RULE_ADD": "Ajouter une règle", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "insertion : comme partout", "S_RULE_ENTER_INH": "Entrée : comme partout", "S_RULE_DELAY_NONE": "sans délai", "S_RULE_PROMPT_INH": "prompts : comme partout",
		"S_RULE_CLIP": "presse-papiers", "S_RULE_TYPE": "caractère par caractère", "S_RULE_ENTER_ON": "avec Entrée", "S_RULE_ENTER_OFF": "sans Entrée",
		"S_RULE_NOPROMPT": "sans prompts", "S_RULE_LAST": "dernière insertion : %s", "S_RULE_EMPTY": "Aucune règle pour l'instant",
		"S_RULE_DEL": "Supprimer la règle", "S_RULE_PROMPTS": "Prompts",
		"S_PASTE_DELAY": "Délai avant l'insertion", "S_PASTE_DELAY_SUB": "quand le programme n'est pas encore prêt",
		"S_OVPOS": "Où afficher le bandeau", "S_OVPOS_SUB": "au curseur — près de l'endroit où vous tapez ; si l'application ne le montre pas, près du pointeur",
		"S_OVPOS_BOTTOM": "En bas de l'écran", "S_OVPOS_TOP": "En haut de l'écran", "S_OVPOS_CARET": "Au curseur",
		"S_OVTEXT": "Afficher le texte reconnu", "S_OVTEXT_SUB": "sur le bandeau après l'insertion, au lieu du nombre de caractères",
		"S_OVERLAY": "Indicateur à l'écran", "S_ANIM": "Animation d'enregistrement", "S_TYPEMODE": "Saisie caractère par caractère",
		"S_RECLANG": "Langue de reconnaissance", "S_RECAUTO": "Auto",
		"S_MODELS": "Modèles de reconnaissance", "S_DL": "Télécharger", "S_DEL": "Supprimer",
		"S_M_BASE": "rapide, PC modestes", "S_M_SMALL": "équilibré", "S_M_MED": "plus précis, recommandé", "S_M_TURBO": "précision max sur CPU",
		"S_M_CUSTOM": "personnalisé (config.json)",
		"S_THREADS":  "Threads CPU", "S_MINMS": "Enreg. min, ms", "S_MAXSEC": "Enreg. max, s",
		"S_AUTOSTART": "Démarrer whisper-server automatiquement", "S_PORT": "Port", "S_SERVEREXE": "Chemin de whisper-server",
		"S_SERVERURL": "Serveur externe (URL)", "S_URLHINT": "Si défini, le serveur local ne démarre pas",
		"S_SAVED": "Enregistré",
		"S_ABOUT_HTML": "<p><b>Voix → texte à la position du curseur.</b></p><p>Placez le curseur, maintenez le raccourci, parlez, relâchez — le texte s'insère.</p><p>Entièrement local et hors ligne. Technologies : <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b> ; modèles depuis Hugging Face.</p><p>Les logs ne dépassent jamais ~2 Mo.</p>",
		"S_VERSION":    "Version",
		"S_LVL_SIMPLE": "simple", "S_LVL_ALL": "tout", "S_SEARCH": "Trouver un réglage…",
		"S_GRP_WORK": "Travail", "S_GRP_REC": "Reconnaissance", "S_GRP_OTHER": "Autre",
		"S_NAV_STATE": "État", "S_NAV_DICT": "Dictée", "S_NAV_MIC": "Microphone", "S_NAV_MODELS": "Modèles",
		"S_NAV_TEXT": "Texte", "S_NAV_TR": "Traduction", "S_NAV_SYSTEM": "Système", "S_NAV_ABOUT": "À propos",
		"S_STATE_HINT": "maintenez et parlez — le texte arrive là où se trouve le curseur",
		"S_STATE_RU": "Parole russe", "S_STATE_OTHER": "Autres langues", "S_STATE_PROC": "Post-traitement",
		"S_CHANGE_MODEL": "Changer", "S_PICK_MODEL": "Choisir", "S_STATE_GET": "Télécharger",
		"S_RETRY": "Réessayer", "S_BERR_OPEN": "Ouvrir les réglages du serveur",
		"S_STATE_LAST": "Dernière dictée", "S_STATE_COPY": "Copier", "S_STATE_MEM": "Mémoire",
		"S_STATE_MEM_SUB": "les modèles restent chargés, la première phrase part sans délai",
		"S_HOTMODE": "Mode", "S_HOTMODE_HOLD": "maintenir", "S_HOTMODE_TOGGLE": "bascule",
		"S_SUB_HOTMODE": "maintenez les touches, ou appuyez une fois pour lancer et une fois pour arrêter",
		"S_SUB_MINMS": "ignore les appuis accidentels",
		"S_SUB_ENTER": "envoie le message aussitôt",
		"S_SUB_CLIP": "les images et les fichiers reviennent tels quels",
		"S_SUB_TYPE": "utile là où un champ refuse le collage",
		"S_SEC_OVERLAY": "Bandeau à l'écran",
		"S_MIC_CHECK": "Vérifier le microphone", "S_MIC_CHECK_SUB": "trois secondes d'enregistrement, puis un verdict : niveau, saturation, présence de parole", "S_MIC_CHECKING": "Vérification…",
		"S_PAUSE": "Mettre l'enregistrement en pause", "S_PAUSE_SUB": "en mode bascule : une pression fige l'enregistrement, une autre le relance",
		"S_MCHECK": "Vérifier les modèles installés", "S_MCHECK_SUB": "compare les fichiers des modèles aux empreintes de référence", "S_MCHECK_GO": "Vérifier", "S_MCHECK_RUN": "Vérification…",
		"S_HIST_INSERT": "Coller",
		"S_LISTS_HINT": "Remplacements et commandes dans un seul fichier — pour passer à un autre ordinateur", "S_LISTS_EXPORT": "Enregistrer dans un fichier", "S_LISTS_IMPORT": "Charger depuis un fichier",
		"S_MIC": "Microphone", "S_MIC_DEFAULT": "Par défaut du système", "S_MIC_REFRESH": "Actualiser la liste",
		"S_MIC_LEVEL": "Niveau d'entrée", "S_MIC_QUIET": "silence",
		"S_ADV_TITLE": "Choisir un modèle", "S_F_ALL": "tous", "S_F_RU": "russe",
		"S_F_MULTI": "plusieurs langues", "S_F_PUNCT": "ponctue", "S_F_FIT": "tient en mémoire",
		"S_ADV_LANGQ": "Dans quelle langue dictez-vous", "S_ADV_PRIOQ": "Qu'est-ce qui compte le plus",
		"S_ADV_ACC": "Choisi pour la précision.", "S_ADV_SPEED": "Choisi pour la vitesse.",
		"S_ADV_TRQ": "Traduction nécessaire", "S_ADV_GO": "Recommander",
		"S_ADV_PRIMARY": "principal", "S_ADV_COMPANION": "second", "S_ADV_HAVE": "déjà là", "S_ADV_APPLY": "Appliquer",
		"S_ADV_ASK": "Seront téléchargés : %s — %s au total. Commencer ?",
		"S_SUB_THREADS": "plus de threads n'est pas toujours plus rapide — mesurez sur votre machine",
		"S_SEC_LLM": "Modèle d'édition",
		"S_PUNCT": "Ponctuation et majuscules", "S_SUB_PUNCT": "d'où viennent la ponctuation et les majuscules",
		"S_PUNCT_MODEL": "du modèle", "S_PUNCT_LLM": "par le modèle d'édition", "S_PUNCT_OFF": "retirer",
		"S_SUB_DICT": "Dictionnaire", "S_SUB_PROMPTS": "Prompts",
		"S_TR_TURBO": "⚠ Le modèle Turbo actif n'est pas entraîné pour la traduction vers l'anglais — choisissez un autre modèle dans l'onglet « Modèles » pour traduire.",
		"S_SUB_TRTARGET": "l'anglais est natif pour Whisper, les autres cibles sont expérimentales",
		"S_TR_EXP": "hors anglais, l'application impose la langue de sortie au lieu de traduire — le texte peut rester dans la langue parlée",
		"S_REMOTE_ABOUT": "Un serveur distant est configuré : l'audio y est envoyé, et la promesse ci-dessus ne tient pas tant qu'il est actif.",
		"S_UPD": "Mises à jour", "S_UPD_CHECK": "Rechercher des mises à jour", "S_UPD_AUTO": "Vérifier au démarrage",
		"S_SUB_UPD": "la seule requête réseau en dehors du téléchargement des modèles",
		"S_UPD_NONE": "Vous avez la dernière version", "S_UPD_AVAIL": "La version %s est disponible.",
		"S_UPD_GO": "Mettre à jour", "S_UPD_ERR": "Échec de la vérification", "S_UPD_DL": "Téléchargement de la mise à jour…",
		"S_SEC_SERVICE": "Service", "S_SUB_AUTOSTART": "désactivez si vous lancez le serveur vous-même",
		"S_SUB_PORT": "la reconnaissance redémarre toute seule",
		"S_MODEL_READY": "Modèle téléchargé — choisissez-le pour basculer",
		"S_FIT_OK": "tient", "S_FIT_WARN": "juste", "S_FIT_BAD": "mémoire insuffisante", "S_RAM": "Mémoire de l'ordinateur :",
		"S_HF_PH": "Nom du modèle — par ex. qwen2.5 instruct",
		"S_NO_LLM": "Aucun modèle installé pour l'instant — trouvez-en un dans l'onglet « Recherche ».",
		"S_NO_LLM_PROF": "Les prompts deviennent disponibles dès qu'un modèle est installé (onglets « Modèles » et « Recherche »).",
		"S_UPDATED": "Dernière mise à jour du modèle", "S_PROF_EDIT": "Modifier", "S_PROF_CLOSE": "Replier",
		"S_CONFIRM_DEL": "Supprimer le modèle « %s » ? Il pourra être téléchargé à nouveau.", "S_FREE": "libre",
		"S_DEL_ACTIVE": "Supprimer le modèle actif « %s » ? La reconnaissance s'arrête jusqu'à ce que vous en choisissiez un autre — vous pouvez le télécharger ici même.",
		"S_WIZ_NEED_MODEL": "Téléchargez d'abord un modèle — sans lui il n'y a rien pour reconnaître",
		"S_REMOTE_WARN": "L'audio sera envoyé à ce serveur. Le mode local est désactivé.",
		"S_REMOTE_ASK": "L'audio ne sera plus traité sur cet ordinateur : il sera envoyé à %s. Activer le mode distant ?",
		"S_REMOTE_BADGE": "DISTANT",
		"S_OK": "Oui", "S_CANCEL": "Annuler", "S_DL_START": "Télécharger", "S_DL_CANCEL": "Annuler le téléchargement",
		"S_DL_ASK": "Le modèle « %s » n'est pas téléchargé (%s). Commencer le téléchargement ?",
		"S_NOT_FOUND": "rien", "S_MORE": "%d réglages de plus", "S_LESS": "Replier %d réglages",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Auteur et développeur de {app} — un outil de dictée local pour Windows : la voix devient du texte à l'endroit du curseur, sans nuage ni abonnement.</p>" +
			"<p>Le projet est ouvert : code source, chaîne de compilation et versions récentes sont sur GitHub.</p>" +
			"<ul>" +
			"<li>Dépôt : <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Profil de l'auteur : <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Un bug ou une idée — ouvrez une issue dans le dépôt.</p>",
		"S_HELP_HTML": "<p class=\"wh\">Comment ça marche</p>" +
			"<p>Maintenez le raccourci — l'enregistrement commence (le bandeau en bas de l'écran montre votre niveau). Relâchez — l'audio est reconnu, traduit si besoin, passé dans les prompts, et le texte final arrive à l'endroit du curseur. Le ✕ du bandeau annule à n'importe quelle étape.</p>" +
			"<p>Le chemin complet : <b>enregistrement → reconnaissance (Whisper) → traduction (si activée) → prompts (LLM) → collage</b>. Chaque étape est visible sur le bandeau.</p>" +
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
			"<p>Le bouton « Test » de l'onglet Microphone enregistre trois secondes et les décortique : niveau de crête en décibels, part de l'enregistrement qui contient vraiment de la parole, et part d'échantillons écrêtés. La réponse est en mots : bon signal, trop faible — montez le niveau dans Windows, saturation — baissez-le, aucune parole entendue — est-ce le bon micro. Les mêmes mesures sont faites après chaque dictée et écrites dans le journal ; si la reconnaissance revient vide, le bandeau nomme la raison — trop faible, saturation ou silence — au lieu de dire simplement qu'il n'a rien entendu.</p>" +
			"<p class=\"wh\">Mettre l'enregistrement en pause</p>" +
			"<p>En mode bascule (une pression démarre, une autre arrête), un raccourci à part peut mettre en pause : onglet « Dictée », ligne « Mettre l'enregistrement en pause ». Une pression fige l'enregistrement — le bandeau affiche « Pause » et plus rien n'est enregistré ; une autre le relance, et tout ce qui a été dit avant reste. La limite de durée ne se déclenche pas pendant la pause.</p>" +
			"<p class=\"wh\">Coller depuis l'historique</p>" +
			"<p>Chaque entrée de l'historique a un bouton « Coller » : il ramène la fenêtre depuis laquelle vous avez ouvert les réglages et y colle le texte, comme une dictée ordinaire. S'il n'y a nulle part où revenir, le texte est simplement placé dans le presse-papiers et le programme le dit.</p>" +
			"<p class=\"wh\">Les listes dans un seul fichier</p>" +
			"<p>Les remplacements et les commandes vocales peuvent être enregistrés dans un fichier .json et chargés sur un autre ordinateur — les boutons sous la liste des commandes, onglet « Texte ». Le chargement n'écrase rien : seules les lignes absentes sont ajoutées, et le programme indique combien ont été ajoutées et combien ignorées.</p>" +
			"<p class=\"wh\">Intégrité des fichiers</p>" +
			"<p>Chaque modèle du catalogue a une empreinte SHA-256 de référence. Après le téléchargement, le fichier lui est comparé : s'il ne correspond pas, il est supprimé et le téléchargement peut être repris. Le bouton « Vérifier » de l'onglet « Modèles » compare de la même façon les modèles déjà installés, et lors d'une mise à jour l'installateur téléchargé est vérifié aussi — un fichier étranger ne sera pas lancé.</p>" +
			"<p class=\"wh\">Historique des dictées</p>" +
			"<p>La section « Historique » dans la colonne de gauche conserve ce que vous avez dicté : le texte seul, sur cet ordinateur seulement, l'audio n'est jamais conservé. Désactivée par défaut, elle s'active d'un interrupteur au même endroit. Les entrées restent un nombre de jours et jusqu'à un nombre d'entrées réglables, les plus anciennes disparaissent d'elles-mêmes ; « Ne jamais enregistrer depuis ces programmes » liste, séparés par des virgules, ceux dont rien ne doit être conservé — gestionnaires de mots de passe, applications bancaires. La recherche porte sur le texte et sur le nom du programme, le bouton à côté d'une entrée la met dans le presse-papiers, et « Vider » supprime tout d'un coup avec le fichier <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Commandes vocales</p>" +
			"<p>Sous les remplacements, dans l'onglet Texte, il y a une liste de commandes : ce que vous dites devient une action au lieu de mots. « Nouvelle ligne » et « nouveau paragraphe » posent un saut — les modèles ne le font jamais ; « annuler » jette toute la dictée sans rien insérer ; « insérer du texte » place ce que vous voulez, un émoticône compris. Le bouton à côté de la liste la remplit des formules habituelles dans la langue de l'interface. Les commandes sont reconnues en mots entiers et s'appliquent après les remplacements, si bien que les prompts et la traduction reçoivent déjà le texte fini. Les espaces en trop autour des sauts disparaissent d'eux-mêmes. Le champ du dessous essaie remplacements et commandes sur n'importe quelle phrase : un saut s'affiche comme ⏎.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">équilibre vitesse / précision</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">plus précis, recommandé</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">meilleure précision sur CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modèles</b> — le catalogue : Base (rapide, pour machines modestes), Small (équilibré), Medium et Turbo (plus précis, plus lents ; « q5 » est une version quantifiée — un peu plus petite et plus rapide, presque sans perte) ainsi que GigaAM v3 pour le russe. Le bouton radio choisit le modèle actif (appliqué aussitôt, la reconnaissance redémarre) ; pour un modèle absent, le programme demande s'il faut le télécharger.</li>" +
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
			"<li>Un prompt peut avoir son propre raccourci : dicter avec lui n'applique que celui-là, une fois. Le crayon ✎ ouvre l'éditeur : nom, texte du prompt, raccourci et un champ d'essai ▶ qui passe un exemple par le modèle en marche, depuis les réglages.</li>" +
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
			"<li><b>Portable</b> — copiez simplement tout le dossier contenant l'exe (sur une clé USB, vers un autre PC) : réglages, modèles et journal vivent à côté et voyagent avec. Rien n'est écrit dans le registre.</li>" +
			"<li>Au premier lancement sans modèle de reconnaissance, le programme ouvre lui-même le catalogue et attend le téléchargement.</li>" +
			"<li>Prérequis : Windows 10/11 x64, un processeur avec AVX2 (à partir de 2013 environ), WebView2 Runtime pour la fenêtre des réglages (inclus dans Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Zone de notification et fichiers</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Prêt…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Réglages…</div><div class=\"mock-mi\">Désactiver</div><div class=\"mock-mi\">Copier le dernier résultat</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Relire config.json</div><div class=\"mock-mi\">Ouvrir config.json</div><div class=\"mock-mi\">Ouvrir le journal</div><div class=\"mock-mi\">À propos</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Quitter</div></div>" +
			"<ul>" +
			"<li>Clic gauche sur l'icône — les réglages ; clic droit — le menu. Couleurs de l'icône : vert — prêt, rouge — enregistrement, orange — reconnaissance, gris — désactivé ou erreur.</li>" +
			"<li><b>config.json</b> — tous les réglages ; les modifications à la main s'appliquent via « Relire config.json » dans le menu.</li>" +
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
		"S_AUTORUN":    "Démarrer avec Windows",
		"S_AUTORUN_SUB": "Une entrée au démarrage de l'utilisateur courant",
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
		"S_TITLE": "{app} — Ajustes", "S_TAB_GENERAL": "General", "S_TAB_REC": "Reconocimiento",
		"S_TAB_PROC": "Posprocesado", "S_TAB_SERVER": "Servidor", "S_TAB_ABOUT": "Acerca de",
		"S_PIPE": "voz ▸ reconocimiento ▸ edición ▸ inserción",
		"S_DICT": "Diccionario de reconocimiento", "S_DICT_HINT": "Términos, nombres y abreviaturas separados por comas — una pista para el oído, no comandos.",
		"S_TR": "Traducción", "S_TR_HINT": "La traducción la hace Whisper: al inglés — modo nativo, otros idiomas — forzando el idioma de salida (calidad variable).",
		"S_TR_DEFAULT": "Traducir siempre al idioma de destino", "S_TR_TARGET": "Idioma destino", "S_TR_ASK": "Elección de idioma", "S_TR_ASK_NEVER": "No preguntar (predeterminado)",
		"S_TR_ASK_ALWAYS": "Preguntar siempre", "S_TR_ASK_TIMEOUT": "Preguntar con tiempo límite", "S_TR_SECONDS": "Tiempo, s",
		"S_TR_LANGS": "Idiomas del diálogo",
		"S_LLM":      "Modos de procesado", "S_LLM_HINT": "El modo predeterminado usa el atajo principal; un modo puede tener su propio atajo. Los perfiles los reescribe una segunda red neuronal (sin conexión).",
		"S_PROF_ASIS": "Tal cual", "S_PROF_WT": "Traducir → English (rápido)",
		"S_PROF_ADD": "Añadir perfil", "S_PROF_NAME": "Nombre", "S_PROF_PROMPT": "Prompt", "S_PROF_HOTKEY": "Atajo",
		"S_PROF_SET": "Definir…", "S_PROF_CLEAR": "Borrar", "S_PROF_TEST": "Prueba",
		"S_HOTKEY": "Atajo de teclado", "S_CHANGE": "Cambiar…", "S_UILANG": "Idioma de la interfaz", "S_AUTO": "Como el sistema",
		"S_SEC_SOUND": "Sonido", "S_SEC_BEHAVIOR": "Comportamiento", "S_BEEP": "Señales sonoras", "S_SOUND": "Sonido de señal",
		"S_SND_SPEECH": "Sistema (voz)", "S_SND_CHIME": "Campanilla", "S_SND_SOFT": "Suave", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Pulsar Enter tras insertar (envío auto)", "S_RESTORE": "Restaurar portapapeles tras insertar",
		"S_NAV_HISTORY": "Historial", "S_HIST_ON": "Guardar el historial de dictados", "S_HIST_ON_SUB": "solo texto, en este equipo; el audio nunca se guarda",
		"S_HIST_DAYS": "Cuántos días guardar", "S_HIST_MAX": "Cuántas entradas guardar",
		"S_HIST_SKIP": "No registrar nunca desde estos programas", "S_HIST_SKIP_SUB": "separados por comas: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Entradas", "S_HIST_CLEAR": "Vaciar", "S_HIST_COPY": "Copiar",
		"S_HIST_FIND": "Buscar en el historial…", "S_HIST_EMPTY": "Todavía no hay historial", "S_HIST_ASK": "¿Eliminar todo el historial de dictados?",
		"S_SEC_CMD": "Comandos de voz", "S_CMD_HINT": "Lo que dices se convierte en un salto de línea, un signo o una cancelación en vez de acabar en el texto. Se buscan como palabras completas y se aplican de arriba abajo, después de los reemplazos.",
		"S_CMD_ADD": "Añadir un comando", "S_CMD_PRESET": "Añadir los habituales", "S_CMD_PH": "nueva línea",
		"S_CMD_NEWLINE": "salto de línea", "S_CMD_PARAGRAPH": "nuevo párrafo", "S_CMD_TEXT": "insertar texto", "S_CMD_CANCEL": "cancelar el dictado",
		"S_CMD_TEXT_PH": "qué insertar", "S_CMD_EMPTY": "Todavía no hay comandos", "S_CMD_DEL": "Eliminar el comando",
		"S_CMD_P_NEWLINE": "nueva línea", "S_CMD_P_PARAGRAPH": "nuevo párrafo", "S_CMD_P_CANCEL": "cancelar",
		"S_SEC_REPLACE": "Reemplazos tras el reconocimiento", "S_REPLACE_HINT": "Lo que se oyó mal se convierte en lo que querías decir, justo tras el reconocimiento y antes de los prompts. Se aplican de arriba abajo.",
		"S_REPL_ADD": "Añadir un reemplazo", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "palabras completas", "S_REPL_CASE": "mayúsculas", "S_REPL_EMPTY": "Todavía no hay reemplazos",
		"S_REPL_DEL": "Eliminar el reemplazo", "S_REPL_TEST_PH": "escribe una frase para probar reemplazos y comandos",
		"S_SEC_RULES": "Reglas por aplicación", "S_RULES_HINT": "La inserción puede funcionar de otra forma en programas concretos. Gana la primera regla que coincide.",
		"S_RULE_ADD": "Añadir una regla", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "inserción: como siempre", "S_RULE_ENTER_INH": "Intro: como siempre", "S_RULE_DELAY_NONE": "sin retraso", "S_RULE_PROMPT_INH": "prompts: como siempre",
		"S_RULE_CLIP": "portapapeles", "S_RULE_TYPE": "carácter a carácter", "S_RULE_ENTER_ON": "con Intro", "S_RULE_ENTER_OFF": "sin Intro",
		"S_RULE_NOPROMPT": "sin prompts", "S_RULE_LAST": "última inserción: %s", "S_RULE_EMPTY": "Todavía no hay reglas",
		"S_RULE_DEL": "Eliminar la regla", "S_RULE_PROMPTS": "Prompts",
		"S_PASTE_DELAY": "Retraso antes de insertar", "S_PASTE_DELAY_SUB": "cuando el programa aún no acepta el texto",
		"S_OVPOS": "Dónde mostrar la barra", "S_OVPOS_SUB": "en el cursor: junto al punto de escritura; si la aplicación no lo expone, junto al puntero",
		"S_OVPOS_BOTTOM": "Abajo en la pantalla", "S_OVPOS_TOP": "Arriba en la pantalla", "S_OVPOS_CARET": "En el cursor",
		"S_OVTEXT": "Mostrar el texto reconocido", "S_OVTEXT_SUB": "en la barra tras insertar, en vez del número de caracteres",
		"S_OVERLAY": "Indicador en pantalla", "S_ANIM": "Animación de grabación", "S_TYPEMODE": "Escritura carácter a carácter",
		"S_RECLANG": "Idioma de reconocimiento", "S_RECAUTO": "Auto",
		"S_MODELS": "Modelos de reconocimiento", "S_DL": "Descargar", "S_DEL": "Borrar",
		"S_M_BASE": "rápido, PCs modestos", "S_M_SMALL": "equilibrado", "S_M_MED": "más preciso, recomendado", "S_M_TURBO": "máx. precisión en CPU",
		"S_M_CUSTOM": "personalizado (config.json)",
		"S_THREADS":  "Hilos de CPU", "S_MINMS": "Grab. mín, ms", "S_MAXSEC": "Grab. máx, s",
		"S_AUTOSTART": "Iniciar whisper-server automáticamente", "S_PORT": "Puerto", "S_SERVEREXE": "Ruta de whisper-server",
		"S_SERVERURL": "Servidor externo (URL)", "S_URLHINT": "Si se define, no se inicia el servidor local",
		"S_SAVED": "Guardado",
		"S_ABOUT_HTML": "<p><b>Voz → texto en la posición del cursor.</b></p><p>Coloque el cursor, mantenga el atajo, hable, suelte — el texto se inserta.</p><p>Totalmente local y sin conexión. Tecnologías: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; modelos de Hugging Face.</p><p>Los logs nunca superan ~2 MB.</p>",
		"S_VERSION":    "Versión",
		"S_LVL_SIMPLE": "simple", "S_LVL_ALL": "todo", "S_SEARCH": "Buscar un ajuste…",
		"S_GRP_WORK": "Trabajo", "S_GRP_REC": "Reconocimiento", "S_GRP_OTHER": "Otros",
		"S_NAV_STATE": "Estado", "S_NAV_DICT": "Dictado", "S_NAV_MIC": "Micrófono", "S_NAV_MODELS": "Modelos",
		"S_NAV_TEXT": "Texto", "S_NAV_TR": "Traducción", "S_NAV_SYSTEM": "Sistema", "S_NAV_ABOUT": "Acerca de",
		"S_STATE_HINT": "mantén y habla — el texto aparece donde está el cursor",
		"S_STATE_RU": "Habla en ruso", "S_STATE_OTHER": "Otros idiomas", "S_STATE_PROC": "Posprocesado",
		"S_CHANGE_MODEL": "Cambiar", "S_PICK_MODEL": "Elegir", "S_STATE_GET": "Descargar",
		"S_RETRY": "Reintentar", "S_BERR_OPEN": "Abrir los ajustes del servidor",
		"S_STATE_LAST": "Último dictado", "S_STATE_COPY": "Copiar", "S_STATE_MEM": "Memoria",
		"S_STATE_MEM_SUB": "los modelos siguen cargados, la primera frase no espera",
		"S_HOTMODE": "Modo", "S_HOTMODE_HOLD": "mantener", "S_HOTMODE_TOGGLE": "alternar",
		"S_SUB_HOTMODE": "mantén las teclas, o pulsa una vez para empezar y otra para parar",
		"S_SUB_MINMS": "ignora pulsaciones accidentales",
		"S_SUB_ENTER": "envía el mensaje enseguida",
		"S_SUB_CLIP": "las imágenes y los archivos vuelven tal como estaban",
		"S_SUB_TYPE": "ayuda donde un campo no admite pegar",
		"S_SEC_OVERLAY": "Aviso en pantalla",
		"S_MIC_CHECK": "Comprobar el micrófono", "S_MIC_CHECK_SUB": "tres segundos de grabación y un veredicto: nivel, saturación, si hay voz", "S_MIC_CHECKING": "Comprobando…",
		"S_PAUSE": "Pausar la grabación", "S_PAUSE_SUB": "en modo conmutador: una pulsación congela la grabación, otra la reanuda",
		"S_MCHECK": "Comprobar los modelos instalados", "S_MCHECK_SUB": "compara los archivos de los modelos con los hashes de referencia", "S_MCHECK_GO": "Comprobar", "S_MCHECK_RUN": "Comprobando…",
		"S_HIST_INSERT": "Pegar",
		"S_LISTS_HINT": "Sustituciones y órdenes en un archivo, para llevarlas a otro ordenador", "S_LISTS_EXPORT": "Guardar en un archivo", "S_LISTS_IMPORT": "Cargar desde un archivo",
		"S_MIC": "Micrófono", "S_MIC_DEFAULT": "Predeterminado del sistema", "S_MIC_REFRESH": "Actualizar la lista",
		"S_MIC_LEVEL": "Nivel de entrada", "S_MIC_QUIET": "silencio",
		"S_ADV_TITLE": "Elegir un modelo", "S_F_ALL": "todos", "S_F_RU": "ruso",
		"S_F_MULTI": "varios idiomas", "S_F_PUNCT": "puntúa", "S_F_FIT": "cabe en memoria",
		"S_ADV_LANGQ": "¿En qué idioma dictas?", "S_ADV_PRIOQ": "¿Qué importa más?",
		"S_ADV_ACC": "Elegido por precisión.", "S_ADV_SPEED": "Elegido por velocidad.",
		"S_ADV_TRQ": "Hace falta traducción", "S_ADV_GO": "Recomendar",
		"S_ADV_PRIMARY": "principal", "S_ADV_COMPANION": "segundo", "S_ADV_HAVE": "ya está", "S_ADV_APPLY": "Aplicar",
		"S_ADV_ASK": "Se descargarán: %s — %s en total. ¿Empezar?",
		"S_SUB_THREADS": "más hilos no siempre es más rápido — mídelo en tu equipo",
		"S_SEC_LLM": "Modelo editor",
		"S_PUNCT": "Puntuación y mayúsculas", "S_SUB_PUNCT": "de dónde salen la puntuación y las mayúsculas",
		"S_PUNCT_MODEL": "del modelo", "S_PUNCT_LLM": "del modelo editor", "S_PUNCT_OFF": "quitar",
		"S_SUB_DICT": "Diccionario", "S_SUB_PROMPTS": "Prompts",
		"S_TR_TURBO": "⚠ El modelo Turbo activo no está entrenado para traducir al inglés — elige otro modelo en la pestaña «Modelos» para traducir.",
		"S_SUB_TRTARGET": "el inglés es nativo para Whisper, los demás destinos son experimentales",
		"S_TR_EXP": "salvo el inglés, la aplicación fuerza el idioma de salida en vez de traducir: el texto puede quedarse en el idioma hablado",
		"S_REMOTE_ABOUT": "Hay un servidor remoto configurado: el audio se envía allí y la promesa anterior no se cumple mientras esté activo.",
		"S_UPD": "Actualizaciones", "S_UPD_CHECK": "Buscar actualizaciones", "S_UPD_AUTO": "Comprobar al iniciar",
		"S_SUB_UPD": "la única petición de red aparte de la descarga de modelos",
		"S_UPD_NONE": "Tienes la última versión", "S_UPD_AVAIL": "La versión %s está disponible.",
		"S_UPD_GO": "Actualizar", "S_UPD_ERR": "Fallo al comprobar", "S_UPD_DL": "Descargando la actualización…",
		"S_SEC_SERVICE": "Servicio", "S_SUB_AUTOSTART": "desactívalo si arrancas el servidor tú mismo",
		"S_SUB_PORT": "el reconocedor se reinicia solo",
		"S_MODEL_READY": "Modelo descargado — elígelo para cambiar",
		"S_FIT_OK": "cabe", "S_FIT_WARN": "justo", "S_FIT_BAD": "falta memoria", "S_RAM": "Memoria del equipo:",
		"S_HF_PH": "Nombre del modelo — p. ej. qwen2.5 instruct",
		"S_NO_LLM": "Todavía no hay modelos instalados — busca uno en la pestaña «Búsqueda».",
		"S_NO_LLM_PROF": "Los prompts estarán disponibles en cuanto haya un modelo instalado (pestañas «Modelos» y «Búsqueda»).",
		"S_UPDATED": "Última actualización del modelo", "S_PROF_EDIT": "Editar", "S_PROF_CLOSE": "Plegar",
		"S_CONFIRM_DEL": "¿Eliminar el modelo «%s»? Se puede descargar de nuevo.", "S_FREE": "libre",
		"S_DEL_ACTIVE": "¿Eliminar el modelo activo «%s»? El reconocimiento se detiene hasta que elijas otro, y puedes descargarlo aquí mismo.",
		"S_WIZ_NEED_MODEL": "Descarga primero un modelo: sin él no hay con qué reconocer",
		"S_REMOTE_WARN": "El audio se enviará a este servidor. El modo local está desactivado.",
		"S_REMOTE_ASK": "El audio dejará de procesarse en este equipo y se enviará a %s. ¿Activar el modo remoto?",
		"S_REMOTE_BADGE": "REMOTO",
		"S_OK": "Sí", "S_CANCEL": "Cancelar", "S_DL_START": "Descargar", "S_DL_CANCEL": "Cancelar la descarga",
		"S_DL_ASK": "El modelo «%s» no está descargado (%s). ¿Empezar la descarga?",
		"S_NOT_FOUND": "nada", "S_MORE": "%d ajustes más", "S_LESS": "Plegar %d ajustes",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autor y desarrollador de {app} — una herramienta de dictado local para Windows: la voz se convierte en texto justo en el cursor, sin nubes ni suscripciones.</p>" +
			"<p>El proyecto es abierto: el código, la compilación y las versiones están en GitHub.</p>" +
			"<ul>" +
			"<li>Repositorio: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Perfil del autor: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>¿Has encontrado un fallo o tienes una idea? Abre un issue en el repositorio.</p>",
		"S_HELP_HTML": "<p class=\"wh\">Cómo funciona</p>" +
			"<p>Mantén el atajo — empieza la grabación (la barra en la parte inferior de la pantalla muestra tu nivel). Suelta — el audio se reconoce, se traduce si hace falta, pasa por los prompts y el texto final aparece donde está el cursor. La ✕ de la barra cancela en cualquier momento.</p>" +
			"<p>El camino completo: <b>grabación → reconocimiento (Whisper) → traducción (si está activa) → prompts (LLM) → pegado</b>. Cada paso se ve en la barra.</p>" +
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
			"<p>El botón «Test» de la pestaña Micrófono graba tres segundos y los analiza: pico en decibelios, qué parte de la grabación contiene voz de verdad y qué parte de las muestras se recortó. La respuesta llega en palabras: se oye bien, demasiado bajo — sube el nivel en Windows, saturación — bájalo, no se oye voz — ¿está elegido el micrófono correcto? Lo mismo se mide tras cada dictado y se escribe en el registro; si el reconocimiento vuelve vacío, la barra nombra el motivo — bajo, saturación o silencio — en vez de decir solo que no oyó nada.</p>" +
			"<p class=\"wh\">Pausar la grabación</p>" +
			"<p>En modo conmutador (una pulsación empieza, otra termina) se puede asignar un atajo aparte para la pausa: pestaña «Dictado», fila «Pausar la grabación». Una pulsación congela la grabación —el rótulo muestra «Pausa» y no se graba nada—, otra la reanuda y conserva todo lo dicho antes. El límite de duración no salta durante la pausa.</p>" +
			"<p class=\"wh\">Pegar desde el historial</p>" +
			"<p>Cada entrada del historial tiene el botón «Pegar»: devuelve el foco a la ventana desde la que abrió los ajustes y pega allí el texto, como un dictado normal. Si no hay adónde volver, el texto se queda en el portapapeles y el programa lo avisa.</p>" +
			"<p class=\"wh\">Las listas en un archivo</p>" +
			"<p>Las sustituciones y las órdenes de voz pueden guardarse en un archivo .json y cargarse en otro ordenador: los botones bajo la lista de órdenes en la pestaña «Texto». La carga no borra nada: solo se añaden las líneas que faltan, y el programa dice cuántas se añadieron y cuántas se omitieron.</p>" +
			"<p class=\"wh\">Integridad de los archivos</p>" +
			"<p>Cada modelo del catálogo tiene un hash SHA-256 de referencia. Tras la descarga el archivo se compara con él: si no coincide, se borra y la descarga puede repetirse. El botón «Comprobar» de la pestaña «Modelos» compara igual los modelos ya instalados, y al actualizar el programa también se comprueba el instalador descargado: un archivo ajeno no se ejecutará.</p>" +
			"<p class=\"wh\">Historial de dictados</p>" +
			"<p>La sección «Historial» de la columna izquierda guarda lo que has dictado: solo el texto, solo en este equipo, el audio nunca se guarda. Está desactivada por defecto y se activa con un interruptor en el mismo sitio. Las entradas se conservan durante los días y hasta la cantidad que fijes, las viejas caen solas; «No registrar nunca desde estos programas» enumera, separados por comas, aquellos de los que no debe guardarse nada — gestores de contraseñas, banca. La búsqueda cubre el texto y el nombre del programa, el botón junto a una entrada la copia al portapapeles, y «Vaciar» borra todo de golpe junto con el archivo <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Comandos de voz</p>" +
			"<p>Bajo los reemplazos, en la pestaña Texto, hay una lista de comandos: lo que dices se convierte en una acción en lugar de en palabras. «Nueva línea» y «nuevo párrafo» ponen un salto — los modelos nunca lo hacen; «cancelar» descarta todo el dictado sin insertar nada; «insertar texto» coloca lo que quieras, incluso un emoticono. El botón junto a la lista la rellena con las frases habituales en el idioma de la interfaz. Los comandos se buscan como palabras completas y se aplican tras los reemplazos, así que los prompts y la traducción reciben ya el texto terminado. Los espacios sobrantes junto a los saltos se limpian solos. El campo de abajo prueba reemplazos y comandos con cualquier frase: un salto se muestra como ⏎.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">equilibrio entre velocidad y precisión</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">más preciso, recomendado</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">la mejor precisión en CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modelos</b> — el catálogo: Base (rápido, para equipos modestos), Small (equilibrado), Medium y Turbo (más precisos y más lentos; «q5» es una versión cuantizada — algo más pequeña y rápida, casi sin pérdida) y GigaAM v3 para el ruso. El botón de opción elige el modelo activo (se aplica al momento y el reconocedor se reinicia); si falta el modelo, el programa pregunta si lo descarga.</li>" +
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
			"<li>Un prompt puede tener su propio atajo: dictar con él aplica solo ese, una vez. El lápiz ✎ abre el editor: nombre, texto del prompt, atajo y un campo de prueba ▶ que pasa un ejemplo por el modelo en marcha desde los propios ajustes.</li>" +
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
			"<li><b>Portátil</b> — copia toda la carpeta con el exe (a un USB, a otro equipo): ajustes, modelos y registro viven al lado y viajan con él. En el registro de Windows no se escribe nada.</li>" +
			"<li>En el primer arranque sin modelo de reconocimiento, el programa abre el catálogo por su cuenta y espera la descarga.</li>" +
			"<li>Requisitos: Windows 10/11 x64, un procesador con AVX2 (de 2013 en adelante, más o menos) y WebView2 Runtime para la ventana de ajustes (incluido en Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Área de notificación y archivos</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Listo…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Ajustes…</div><div class=\"mock-mi\">Desactivar</div><div class=\"mock-mi\">Copiar el último resultado</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Releer config.json</div><div class=\"mock-mi\">Abrir config.json</div><div class=\"mock-mi\">Abrir el registro</div><div class=\"mock-mi\">Acerca de</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Salir</div></div>" +
			"<ul>" +
			"<li>Clic izquierdo en el icono — los ajustes; clic derecho — el menú. Colores del icono: verde — listo, rojo — grabando, naranja — reconociendo, gris — desactivado o error.</li>" +
			"<li><b>config.json</b> — todos los ajustes; los cambios a mano se aplican con «Releer config.json» en el menú.</li>" +
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
		"S_AUTORUN":    "Iniciar con Windows",
		"S_AUTORUN_SUB": "Una entrada en el inicio del usuario actual",
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
		"S_TITLE": "{app} — Impostazioni", "S_TAB_GENERAL": "Generale", "S_TAB_REC": "Riconoscimento",
		"S_TAB_PROC": "Post-elaborazione", "S_TAB_SERVER": "Server", "S_TAB_ABOUT": "Info",
		"S_PIPE": "voce ▸ riconoscimento ▸ modifica ▸ inserimento",
		"S_DICT": "Dizionario di riconoscimento", "S_DICT_HINT": "Termini, nomi e abbreviazioni separati da virgole — un suggerimento per l'ascolto, non comandi.",
		"S_TR": "Traduzione", "S_TR_HINT": "La traduzione è fatta da Whisper: verso l'inglese — modalità nativa, altre lingue — forzando la lingua di output (qualità variabile).",
		"S_TR_DEFAULT": "Traduci sempre nella lingua di destinazione", "S_TR_TARGET": "Lingua di destinazione", "S_TR_ASK": "Scelta della lingua", "S_TR_ASK_NEVER": "Non chiedere (predefinita)",
		"S_TR_ASK_ALWAYS": "Chiedi ogni volta", "S_TR_ASK_TIMEOUT": "Chiedi con timeout", "S_TR_SECONDS": "Timeout, s",
		"S_TR_LANGS": "Lingue nel dialogo",
		"S_LLM":      "Modalità di elaborazione", "S_LLM_HINT": "La modalità predefinita usa la scorciatoia principale; una modalità può avere il proprio hotkey. I profili sono riscritti da una seconda rete neurale (offline).",
		"S_PROF_ASIS": "Così com'è", "S_PROF_WT": "Traduci → English (veloce)",
		"S_PROF_ADD": "Aggiungi profilo", "S_PROF_NAME": "Nome", "S_PROF_PROMPT": "Prompt", "S_PROF_HOTKEY": "Hotkey",
		"S_PROF_SET": "Imposta…", "S_PROF_CLEAR": "Cancella", "S_PROF_TEST": "Prova",
		"S_HOTKEY": "Scorciatoia da tastiera", "S_CHANGE": "Cambia…", "S_UILANG": "Lingua dell'interfaccia", "S_AUTO": "Come il sistema",
		"S_SEC_SOUND": "Suono", "S_SEC_BEHAVIOR": "Comportamento", "S_BEEP": "Segnali acustici", "S_SOUND": "Suono del segnale",
		"S_SND_SPEECH": "Sistema (voce)", "S_SND_CHIME": "Campanello", "S_SND_SOFT": "Morbido", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Premi Invio dopo l'inserimento (invio auto)", "S_RESTORE": "Ripristina appunti dopo l'inserimento",
		"S_NAV_HISTORY": "Cronologia", "S_HIST_ON": "Conservare la cronologia delle dettature", "S_HIST_ON_SUB": "solo testo, su questo computer; l'audio non viene mai salvato",
		"S_HIST_DAYS": "Quanti giorni conservare", "S_HIST_MAX": "Quante voci conservare",
		"S_HIST_SKIP": "Non registrare mai da questi programmi", "S_HIST_SKIP_SUB": "separati da virgole: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Voci", "S_HIST_CLEAR": "Svuota", "S_HIST_COPY": "Copia",
		"S_HIST_FIND": "Cerca nella cronologia…", "S_HIST_EMPTY": "Ancora nessuna cronologia", "S_HIST_ASK": "Eliminare tutta la cronologia delle dettature?",
		"S_SEC_CMD": "Comandi vocali", "S_CMD_HINT": "Ciò che dici diventa un a capo, un segno o un annullamento invece di finire nel testo. Riconosciuti come parole intere, applicati dall'alto in basso, dopo le sostituzioni.",
		"S_CMD_ADD": "Aggiungi un comando", "S_CMD_PRESET": "Aggiungi quelli soliti", "S_CMD_PH": "nuova riga",
		"S_CMD_NEWLINE": "a capo", "S_CMD_PARAGRAPH": "nuovo paragrafo", "S_CMD_TEXT": "inserire testo", "S_CMD_CANCEL": "annullare la dettatura",
		"S_CMD_TEXT_PH": "cosa inserire", "S_CMD_EMPTY": "Ancora nessun comando", "S_CMD_DEL": "Elimina il comando",
		"S_CMD_P_NEWLINE": "nuova riga", "S_CMD_P_PARAGRAPH": "nuovo paragrafo", "S_CMD_P_CANCEL": "annulla",
		"S_SEC_REPLACE": "Sostituzioni dopo il riconoscimento", "S_REPLACE_HINT": "Ciò che è stato sentito male diventa ciò che intendevi — subito dopo il riconoscimento, prima dei prompt. Applicate dall'alto in basso.",
		"S_REPL_ADD": "Aggiungi una sostituzione", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "parole intere", "S_REPL_CASE": "maiuscole", "S_REPL_EMPTY": "Ancora nessuna sostituzione",
		"S_REPL_DEL": "Elimina la sostituzione", "S_REPL_TEST_PH": "scrivi una frase per provare sostituzioni e comandi",
		"S_SEC_RULES": "Regole per applicazione", "S_RULES_HINT": "L'inserimento può funzionare diversamente in certi programmi. Vince la prima regola che corrisponde.",
		"S_RULE_ADD": "Aggiungi una regola", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "inserimento: come sempre", "S_RULE_ENTER_INH": "Invio: come sempre", "S_RULE_DELAY_NONE": "senza ritardo", "S_RULE_PROMPT_INH": "prompt: come sempre",
		"S_RULE_CLIP": "appunti", "S_RULE_TYPE": "carattere per carattere", "S_RULE_ENTER_ON": "con Invio", "S_RULE_ENTER_OFF": "senza Invio",
		"S_RULE_NOPROMPT": "senza prompt", "S_RULE_LAST": "ultimo inserimento: %s", "S_RULE_EMPTY": "Ancora nessuna regola",
		"S_RULE_DEL": "Elimina la regola", "S_RULE_PROMPTS": "Prompt",
		"S_PASTE_DELAY": "Ritardo prima di inserire", "S_PASTE_DELAY_SUB": "quando il programma non è ancora pronto",
		"S_OVPOS": "Dove mostrare la barra", "S_OVPOS_SUB": "al cursore — accanto al punto in cui scrivi; se l'app non lo espone, accanto al puntatore",
		"S_OVPOS_BOTTOM": "In basso sullo schermo", "S_OVPOS_TOP": "In alto sullo schermo", "S_OVPOS_CARET": "Al cursore",
		"S_OVTEXT": "Mostrare il testo riconosciuto", "S_OVTEXT_SUB": "sulla barra dopo l'inserimento, invece del numero di caratteri",
		"S_OVERLAY": "Indicatore a schermo", "S_ANIM": "Animazione di registrazione", "S_TYPEMODE": "Digitazione carattere per carattere",
		"S_RECLANG": "Lingua di riconoscimento", "S_RECAUTO": "Auto",
		"S_MODELS": "Modelli di riconoscimento", "S_DL": "Scarica", "S_DEL": "Elimina",
		"S_M_BASE": "veloce, PC modesti", "S_M_SMALL": "bilanciato", "S_M_MED": "più preciso, consigliato", "S_M_TURBO": "massima precisione su CPU",
		"S_M_CUSTOM": "personalizzato (config.json)",
		"S_THREADS":  "Thread CPU", "S_MINMS": "Registrazione min, ms", "S_MAXSEC": "Registrazione max, s",
		"S_AUTOSTART": "Avvia whisper-server automaticamente", "S_PORT": "Porta", "S_SERVEREXE": "Percorso whisper-server",
		"S_SERVERURL": "Server esterno (URL)", "S_URLHINT": "Se impostato, il server locale non parte",
		"S_SAVED": "Salvato",
		"S_ABOUT_HTML": "<p><b>Voce → testo alla posizione del cursore.</b></p><p>Posiziona il cursore, tieni la scorciatoia, parla, rilascia — il testo viene inserito.</p><p>Completamente locale e offline. Tecnologie: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; modelli da Hugging Face.</p><p>I log non superano mai ~2 MB.</p>",
		"S_VERSION":    "Versione",
		"S_LVL_SIMPLE": "semplice", "S_LVL_ALL": "tutto", "S_SEARCH": "Trova un'impostazione…",
		"S_GRP_WORK": "Lavoro", "S_GRP_REC": "Riconoscimento", "S_GRP_OTHER": "Altro",
		"S_NAV_STATE": "Stato", "S_NAV_DICT": "Dettatura", "S_NAV_MIC": "Microfono", "S_NAV_MODELS": "Modelli",
		"S_NAV_TEXT": "Testo", "S_NAV_TR": "Traduzione", "S_NAV_SYSTEM": "Sistema", "S_NAV_ABOUT": "Informazioni",
		"S_STATE_HINT": "tieni premuto e parla — il testo arriva dov'è il cursore",
		"S_STATE_RU": "Parlato russo", "S_STATE_OTHER": "Altre lingue", "S_STATE_PROC": "Post-elaborazione",
		"S_CHANGE_MODEL": "Cambia", "S_PICK_MODEL": "Scegli", "S_STATE_GET": "Scarica",
		"S_RETRY": "Riprova", "S_BERR_OPEN": "Apri le impostazioni del server",
		"S_STATE_LAST": "Ultima dettatura", "S_STATE_COPY": "Copia", "S_STATE_MEM": "Memoria",
		"S_STATE_MEM_SUB": "i modelli restano caricati, la prima frase non aspetta",
		"S_HOTMODE": "Modo", "S_HOTMODE_HOLD": "tenere premuto", "S_HOTMODE_TOGGLE": "interruttore",
		"S_SUB_HOTMODE": "tieni premuti i tasti, oppure premi una volta per iniziare e una per fermare",
		"S_SUB_MINMS": "ignora le pressioni accidentali",
		"S_SUB_ENTER": "invia subito il messaggio",
		"S_SUB_CLIP": "immagini e file tornano come erano",
		"S_SUB_TYPE": "serve dove un campo rifiuta di incollare",
		"S_SEC_OVERLAY": "Avviso a schermo",
		"S_MIC_CHECK": "Controlla il microfono", "S_MIC_CHECK_SUB": "tre secondi di registrazione e un verdetto: livello, distorsione, se c'è voce", "S_MIC_CHECKING": "Controllo…",
		"S_PAUSE": "Metti in pausa la registrazione", "S_PAUSE_SUB": "in modalità interruttore: una pressione ferma la registrazione, un'altra la riprende",
		"S_MCHECK": "Controlla i modelli installati", "S_MCHECK_SUB": "confronta i file dei modelli con gli hash di riferimento", "S_MCHECK_GO": "Controlla", "S_MCHECK_RUN": "Controllo…",
		"S_HIST_INSERT": "Incolla",
		"S_LISTS_HINT": "Sostituzioni e comandi in un solo file, da portare su un altro computer", "S_LISTS_EXPORT": "Salva su file", "S_LISTS_IMPORT": "Carica da file",
		"S_MIC": "Microfono", "S_MIC_DEFAULT": "Predefinito di sistema", "S_MIC_REFRESH": "Aggiorna l'elenco",
		"S_MIC_LEVEL": "Livello d'ingresso", "S_MIC_QUIET": "silenzio",
		"S_ADV_TITLE": "Scegli un modello", "S_F_ALL": "tutti", "S_F_RU": "russo",
		"S_F_MULTI": "più lingue", "S_F_PUNCT": "mette la punteggiatura", "S_F_FIT": "sta in memoria",
		"S_ADV_LANGQ": "In che lingua detti", "S_ADV_PRIOQ": "Cosa conta di più",
		"S_ADV_ACC": "Scelto per la precisione.", "S_ADV_SPEED": "Scelto per la velocità.",
		"S_ADV_TRQ": "Serve la traduzione", "S_ADV_GO": "Consiglia",
		"S_ADV_PRIMARY": "principale", "S_ADV_COMPANION": "secondo", "S_ADV_HAVE": "già presente", "S_ADV_APPLY": "Applica",
		"S_ADV_ASK": "Verranno scaricati: %s — %s in tutto. Iniziare?",
		"S_SUB_THREADS": "più thread non è sempre più veloce — misura sulla tua macchina",
		"S_SEC_LLM": "Modello editor",
		"S_PUNCT": "Punteggiatura e maiuscole", "S_SUB_PUNCT": "da dove arrivano punteggiatura e maiuscole",
		"S_PUNCT_MODEL": "dal modello", "S_PUNCT_LLM": "dal modello editor", "S_PUNCT_OFF": "togliere",
		"S_SUB_DICT": "Dizionario", "S_SUB_PROMPTS": "Prompt",
		"S_TR_TURBO": "⚠ Il modello Turbo attivo non è addestrato per la traduzione in inglese — per tradurre scegli un altro modello nella scheda «Modelli».",
		"S_SUB_TRTARGET": "l'inglese è nativo per Whisper, le altre lingue di destinazione sono sperimentali",
		"S_TR_EXP": "tranne l'inglese, l'app forza la lingua di uscita invece di tradurre: il testo può restare nella lingua parlata",
		"S_REMOTE_ABOUT": "È impostato un server remoto: l'audio viene inviato lì e la promessa qui sopra non vale finché è attivo.",
		"S_UPD": "Aggiornamenti", "S_UPD_CHECK": "Cerca aggiornamenti", "S_UPD_AUTO": "Controlla all'avvio",
		"S_SUB_UPD": "l'unica richiesta di rete oltre allo scaricamento dei modelli",
		"S_UPD_NONE": "Hai l'ultima versione", "S_UPD_AVAIL": "È disponibile la versione %s.",
		"S_UPD_GO": "Aggiorna", "S_UPD_ERR": "Controllo non riuscito", "S_UPD_DL": "Scaricamento dell'aggiornamento…",
		"S_SEC_SERVICE": "Servizio", "S_SUB_AUTOSTART": "disattiva se avvii il server da solo",
		"S_SUB_PORT": "il riconoscitore si riavvia da sé",
		"S_MODEL_READY": "Modello scaricato — scegli per passare a lui",
		"S_FIT_OK": "ci sta", "S_FIT_WARN": "al limite", "S_FIT_BAD": "memoria insufficiente", "S_RAM": "Memoria del computer:",
		"S_HF_PH": "Nome del modello — per es. qwen2.5 instruct",
		"S_NO_LLM": "Ancora nessun modello installato — cercane uno nella scheda «Ricerca».",
		"S_NO_LLM_PROF": "I prompt diventano disponibili appena è installato un modello (schede «Modelli» e «Ricerca»).",
		"S_UPDATED": "Ultimo aggiornamento del modello", "S_PROF_EDIT": "Modifica", "S_PROF_CLOSE": "Riduci",
		"S_CONFIRM_DEL": "Eliminare il modello «%s»? Si potrà scaricare di nuovo.", "S_FREE": "liberi",
		"S_DEL_ACTIVE": "Eliminare il modello attivo «%s»? Il riconoscimento si ferma finché non ne scegli un altro, che puoi scaricare qui stesso.",
		"S_WIZ_NEED_MODEL": "Scarica prima un modello: senza non c'è nulla con cui riconoscere",
		"S_REMOTE_WARN": "L'audio verrà inviato a questo server. La modalità locale è spenta.",
		"S_REMOTE_ASK": "L'audio non sarà più elaborato su questo computer e verrà inviato a %s. Attivare la modalità remota?",
		"S_REMOTE_BADGE": "REMOTO",
		"S_OK": "Sì", "S_CANCEL": "Annulla", "S_DL_START": "Scarica", "S_DL_CANCEL": "Annulla lo scaricamento",
		"S_DL_ASK": "Il modello «%s» non è scaricato (%s). Iniziare lo scaricamento?",
		"S_NOT_FOUND": "niente", "S_MORE": "altre %d impostazioni", "S_LESS": "Riduci %d impostazioni",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autore e sviluppatore di {app} — uno strumento di dettatura locale per Windows: la voce diventa testo proprio dov'è il cursore, senza cloud e senza abbonamenti.</p>" +
			"<p>Il progetto è aperto: codice sorgente, build e ultime release stanno su GitHub.</p>" +
			"<ul>" +
			"<li>Repository: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Profilo dell'autore: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Hai trovato un bug o hai un'idea — apri una issue nel repository.</p>",
		"S_HELP_HTML": "<p class=\"wh\">Come funziona</p>" +
			"<p>Tieni premuta la scorciatoia — parte la registrazione (la barra in fondo allo schermo mostra il tuo livello). Rilascia — l'audio viene riconosciuto, tradotto se serve, passato nei prompt, e il testo finale arriva dov'è il cursore. La ✕ sulla barra annulla in qualsiasi momento.</p>" +
			"<p>Il percorso completo: <b>registrazione → riconoscimento (Whisper) → traduzione (se attiva) → prompt (LLM) → incollaggio</b>. Ogni passo si vede sulla barra.</p>" +
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
			"<p>Il pulsante «Test» nella scheda Microfono registra tre secondi e li scompone: picco in decibel, quanta parte della registrazione contiene davvero voce e quanti campioni sono stati tagliati. La risposta arriva a parole: si sente bene, troppo basso — alza il livello in Windows, distorsione — abbassalo, nessuna voce sentita — è scelto il microfono giusto. Le stesse misure vengono fatte dopo ogni dettatura e finiscono nel log; se il riconoscimento torna vuoto, la barra dice il motivo — basso, distorsione o silenzio — invece di dire soltanto che non ha sentito nulla.</p>" +
			"<p class=\"wh\">Mettere in pausa la registrazione</p>" +
			"<p>In modalità interruttore (una pressione avvia, un'altra ferma) si può assegnare una scorciatoia a parte per la pausa: scheda «Dettatura», riga «Metti in pausa la registrazione». Una pressione ferma la registrazione — la targhetta mostra «Pausa» e non viene registrato nulla — un'altra la riprende, e tutto ciò che è stato detto prima resta. Il limite di durata non scatta durante la pausa.</p>" +
			"<p class=\"wh\">Incollare dalla cronologia</p>" +
			"<p>Ogni voce della cronologia ha il pulsante «Incolla»: riporta in primo piano la finestra da cui avete aperto le impostazioni e vi incolla il testo, come una dettatura normale. Se non c'è dove tornare, il testo finisce semplicemente negli appunti e il programma lo dice.</p>" +
			"<p class=\"wh\">Le liste in un solo file</p>" +
			"<p>Sostituzioni e comandi vocali si possono salvare in un file .json e caricare su un altro computer: i pulsanti sotto l'elenco dei comandi nella scheda «Testo». Il caricamento non cancella nulla: vengono aggiunte solo le righe che mancano, e il programma dice quante ne ha aggiunte e quante saltate.</p>" +
			"<p class=\"wh\">Integrità dei file</p>" +
			"<p>Per ogni modello del catalogo è noto l'hash SHA-256 di riferimento. Dopo lo scaricamento il file viene confrontato con esso: se non coincide, il file viene eliminato e lo scaricamento si può ripetere. Il pulsante «Controlla» nella scheda «Modelli» confronta allo stesso modo i modelli già installati, e all'aggiornamento del programma viene controllato anche l'installatore scaricato: un file estraneo non verrà avviato.</p>" +
			"<p class=\"wh\">Cronologia delle dettature</p>" +
			"<p>La sezione «Cronologia» nella colonna di sinistra conserva ciò che hai dettato: solo testo, solo su questo computer, l'audio non viene mai salvato. È disattivata per impostazione predefinita e si accende con un interruttore lì accanto. Le voci restano per i giorni e fino al numero che imposti, le più vecchie escono da sole; «Non registrare mai da questi programmi» elenca, separati da virgole, quelli da cui non salvare nulla — gestori di password, home banking. La ricerca copre il testo e il nome del programma, il pulsante accanto a una voce la mette negli appunti, e «Svuota» cancella tutto insieme al file <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Comandi vocali</p>" +
			"<p>Sotto le sostituzioni, nella scheda Testo, c'è un elenco di comandi: ciò che dici diventa un'azione invece che parole. «Nuova riga» e «nuovo paragrafo» inseriscono un a capo — i modelli non lo fanno mai; «annulla» butta via l'intera dettatura senza inserire nulla; «inserire testo» mette quello che vuoi, faccina compresa. Il pulsante accanto all'elenco lo riempie con le formule solite nella lingua dell'interfaccia. I comandi sono riconosciuti come parole intere e si applicano dopo le sostituzioni, così i prompt e la traduzione ricevono già il testo finito. Gli spazi di troppo attorno agli a capo spariscono da soli. Il campo in fondo prova sostituzioni e comandi su qualsiasi frase: un a capo appare come ⏎.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">equilibrio tra velocità e precisione</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">più preciso, consigliato</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">la migliore precisione su CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modelli</b> — il catalogo: Base (veloce, per macchine modeste), Small (equilibrato), Medium e Turbo (più precisi e più lenti; «q5» è una versione quantizzata — un po' più piccola e rapida, quasi senza perdite) e GigaAM v3 per il russo. Il pulsante scelto indica il modello attivo (vale subito, il riconoscitore si riavvia); se il modello manca, il programma chiede se scaricarlo.</li>" +
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
			"<li>Un prompt può avere la sua scorciatoia: dettare con quella applica solo lui, una volta. La matita ✎ apre l'editor: nome, testo del prompt, scorciatoia e un campo di prova ▶ che manda un esempio al modello acceso direttamente dalle impostazioni.</li>" +
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
			"<li><b>Portabile</b> — copia semplicemente tutta la cartella con l'exe (su una chiavetta, su un altro computer): impostazioni, modelli e registro stanno accanto e viaggiano con lui. Nel registro di sistema non viene scritto nulla.</li>" +
			"<li>Al primo avvio senza modello di riconoscimento il programma apre il catalogo da solo e aspetta lo scaricamento.</li>" +
			"<li>Requisiti: Windows 10/11 x64, una CPU con AVX2 (dal 2013 circa), WebView2 Runtime per la finestra delle impostazioni (incluso in Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Area di notifica e file</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Pronto…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Impostazioni…</div><div class=\"mock-mi\">Disattiva</div><div class=\"mock-mi\">Copia l'ultimo risultato</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Rileggi config.json</div><div class=\"mock-mi\">Apri config.json</div><div class=\"mock-mi\">Apri il registro</div><div class=\"mock-mi\">Informazioni</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Esci</div></div>" +
			"<ul>" +
			"<li>Clic sinistro sull'icona — le impostazioni; clic destro — il menu. Colori dell'icona: verde — pronto, rosso — registrazione, arancione — riconoscimento, grigio — disattivato o errore.</li>" +
			"<li><b>config.json</b> — tutte le impostazioni; le modifiche a mano valgono dopo «Rileggi config.json» nel menu.</li>" +
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
		"S_AUTORUN":    "Avvia con Windows",
		"S_AUTORUN_SUB": "Una voce nell'avvio automatico dell'utente corrente",
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
		"S_TITLE": "{app} — Ustawienia", "S_TAB_GENERAL": "Ogólne", "S_TAB_REC": "Rozpoznawanie",
		"S_TAB_PROC": "Postprodukcja", "S_TAB_SERVER": "Serwer", "S_TAB_ABOUT": "O programie",
		"S_PIPE": "głos ▸ rozpoznawanie ▸ edycja ▸ wstawianie",
		"S_DICT": "Słownik rozpoznawania", "S_DICT_HINT": "Terminy, nazwy i skróty oddzielone przecinkami — podpowiedź dla słuchu, nie polecenia.",
		"S_TR": "Tłumaczenie", "S_TR_HINT": "Tłumaczy Whisper: na angielski — tryb natywny, inne języki — przez wymuszenie języka wyjściowego (jakość bywa różna).",
		"S_TR_DEFAULT": "Zawsze tłumacz na język docelowy", "S_TR_TARGET": "Język docelowy", "S_TR_ASK": "Wybór języka", "S_TR_ASK_NEVER": "Nie pytaj (domyślny)",
		"S_TR_ASK_ALWAYS": "Pytaj za każdym razem", "S_TR_ASK_TIMEOUT": "Pytaj z limitem czasu", "S_TR_SECONDS": "Limit, s",
		"S_TR_LANGS": "Języki w oknie dialogowym",
		"S_LLM":      "Tryby przetwarzania", "S_LLM_HINT": "Tryb domyślny działa na głównym skrócie; tryb może mieć własny skrót. Profile przepisuje druga sieć neuronowa (offline).",
		"S_PROF_ASIS": "Bez zmian", "S_PROF_WT": "Tłumacz → English (szybki)",
		"S_PROF_ADD": "Dodaj profil", "S_PROF_NAME": "Nazwa", "S_PROF_PROMPT": "Prompt", "S_PROF_HOTKEY": "Skrót",
		"S_PROF_SET": "Ustaw…", "S_PROF_CLEAR": "Wyczyść", "S_PROF_TEST": "Test",
		"S_HOTKEY": "Skrót klawiszowy", "S_CHANGE": "Zmień…", "S_UILANG": "Język interfejsu", "S_AUTO": "Jak w systemie",
		"S_SEC_SOUND": "Dźwięk", "S_SEC_BEHAVIOR": "Zachowanie", "S_BEEP": "Sygnały dźwiękowe", "S_SOUND": "Dźwięk sygnału",
		"S_SND_SPEECH": "Systemowy (mowa)", "S_SND_CHIME": "Dzwonek", "S_SND_SOFT": "Miękki", "S_SND_MARIMBA": "Marimba",
		"S_SND_BLIP": "Blip", "S_SND_POP": "Pop",
		"S_AUTOENTER": "Naciśnij Enter po wstawieniu (auto-wysyłka)", "S_RESTORE": "Przywróć schowek po wstawieniu",
		"S_NAV_HISTORY": "Historia", "S_HIST_ON": "Przechowuj historię dyktowań", "S_HIST_ON_SUB": "tylko tekst, na tym komputerze; dźwięk nigdy nie jest zapisywany",
		"S_HIST_DAYS": "Ile dni przechowywać", "S_HIST_MAX": "Ile wpisów przechowywać",
		"S_HIST_SKIP": "Nigdy nie zapisuj z tych programów", "S_HIST_SKIP_SUB": "po przecinku: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Wpisy", "S_HIST_CLEAR": "Wyczyść", "S_HIST_COPY": "Kopiuj",
		"S_HIST_FIND": "Szukaj w historii…", "S_HIST_EMPTY": "Na razie brak historii", "S_HIST_ASK": "Usunąć całą historię dyktowań?",
		"S_SEC_CMD": "Komendy głosowe", "S_CMD_HINT": "To, co powiesz, zamienia się w złamanie wiersza, znak albo anulowanie, zamiast trafić do tekstu. Rozpoznawane jako całe słowa, stosowane od góry do dołu, po zamianach.",
		"S_CMD_ADD": "Dodaj komendę", "S_CMD_PRESET": "Dodaj typowe", "S_CMD_PH": "nowy wiersz",
		"S_CMD_NEWLINE": "złamanie wiersza", "S_CMD_PARAGRAPH": "nowy akapit", "S_CMD_TEXT": "wstawić tekst", "S_CMD_CANCEL": "anulować dyktowanie",
		"S_CMD_TEXT_PH": "co wstawić", "S_CMD_EMPTY": "Na razie brak komend", "S_CMD_DEL": "Usuń komendę",
		"S_CMD_P_NEWLINE": "nowy wiersz", "S_CMD_P_PARAGRAPH": "nowy akapit", "S_CMD_P_CANCEL": "anuluj",
		"S_SEC_REPLACE": "Zamiany po rozpoznaniu", "S_REPLACE_HINT": "To, co zostało źle usłyszane, staje się tym, co miałeś na myśli — zaraz po rozpoznaniu, przed promptami. Stosowane od góry do dołu.",
		"S_REPL_ADD": "Dodaj zamianę", "S_REPL_FROM_PH": "git hub", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "całe słowa", "S_REPL_CASE": "wielkość liter", "S_REPL_EMPTY": "Na razie brak zamian",
		"S_REPL_DEL": "Usuń zamianę", "S_REPL_TEST_PH": "wpisz zdanie, aby sprawdzić zamiany i komendy",
		"S_SEC_RULES": "Reguły dla programów", "S_RULES_HINT": "W wybranych programach wstawianie może działać inaczej. Wygrywa pierwsza pasująca reguła.",
		"S_RULE_ADD": "Dodaj regułę", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "wstawianie: jak zwykle", "S_RULE_ENTER_INH": "Enter: jak zwykle", "S_RULE_DELAY_NONE": "bez opóźnienia", "S_RULE_PROMPT_INH": "prompty: jak zwykle",
		"S_RULE_CLIP": "schowek", "S_RULE_TYPE": "znak po znaku", "S_RULE_ENTER_ON": "z Enterem", "S_RULE_ENTER_OFF": "bez Entera",
		"S_RULE_NOPROMPT": "bez promptów", "S_RULE_LAST": "ostatnie wstawienie: %s", "S_RULE_EMPTY": "Na razie brak reguł",
		"S_RULE_DEL": "Usuń regułę", "S_RULE_PROMPTS": "Prompty",
		"S_PASTE_DELAY": "Opóźnienie przed wstawieniem", "S_PASTE_DELAY_SUB": "gdy program nie zdąża przyjąć tekstu",
		"S_OVPOS": "Gdzie pokazywać pasek", "S_OVPOS_SUB": "przy kursorze — obok miejsca pisania; jeśli aplikacja go nie pokazuje, obok wskaźnika myszy",
		"S_OVPOS_BOTTOM": "Na dole ekranu", "S_OVPOS_TOP": "Na górze ekranu", "S_OVPOS_CARET": "Przy kursorze",
		"S_OVTEXT": "Pokazywać rozpoznany tekst", "S_OVTEXT_SUB": "na pasku po wstawieniu, zamiast liczby znaków",
		"S_OVERLAY": "Wskaźnik na ekranie", "S_ANIM": "Animacja nagrywania", "S_TYPEMODE": "Wpisywanie znak po znaku",
		"S_RECLANG": "Język rozpoznawania", "S_RECAUTO": "Auto",
		"S_MODELS": "Modele rozpoznawania", "S_DL": "Pobierz", "S_DEL": "Usuń",
		"S_M_BASE": "szybki, słabe komputery", "S_M_SMALL": "zrównoważony", "S_M_MED": "dokładniejszy, polecany", "S_M_TURBO": "najlepsza dokładność na CPU",
		"S_M_CUSTOM": "własny (config.json)",
		"S_THREADS":  "Wątki CPU", "S_MINMS": "Min. nagranie, ms", "S_MAXSEC": "Maks. nagranie, s",
		"S_AUTOSTART": "Uruchamiaj whisper-server automatycznie", "S_PORT": "Port", "S_SERVEREXE": "Ścieżka whisper-server",
		"S_SERVERURL": "Serwer zewnętrzny (URL)", "S_URLHINT": "Jeśli ustawiony, lokalny serwer nie startuje",
		"S_SAVED": "Zapisano",
		"S_ABOUT_HTML": "<p><b>Głos → tekst w pozycji kursora.</b></p><p>Ustaw kursor, przytrzymaj skrót, mów, puść — tekst zostanie wstawiony.</p><p>W pełni lokalnie i offline. Technologie: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; modele z Hugging Face.</p><p>Logi nigdy nie przekraczają ~2 MB.</p>",
		"S_VERSION":    "Wersja",
		"S_LVL_SIMPLE": "proste", "S_LVL_ALL": "wszystko", "S_SEARCH": "Znajdź ustawienie…",
		"S_GRP_WORK": "Praca", "S_GRP_REC": "Rozpoznawanie", "S_GRP_OTHER": "Inne",
		"S_NAV_STATE": "Stan", "S_NAV_DICT": "Dyktowanie", "S_NAV_MIC": "Mikrofon", "S_NAV_MODELS": "Modele",
		"S_NAV_TEXT": "Tekst", "S_NAV_TR": "Tłumaczenie", "S_NAV_SYSTEM": "System", "S_NAV_ABOUT": "O programie",
		"S_STATE_HINT": "przytrzymaj i mów — tekst trafia tam, gdzie stoi kursor",
		"S_STATE_RU": "Mowa rosyjska", "S_STATE_OTHER": "Inne języki", "S_STATE_PROC": "Obróbka tekstu",
		"S_CHANGE_MODEL": "Zmień", "S_PICK_MODEL": "Dobierz", "S_STATE_GET": "Pobierz",
		"S_RETRY": "Spróbuj ponownie", "S_BERR_OPEN": "Otwórz ustawienia serwera",
		"S_STATE_LAST": "Ostatnie dyktowanie", "S_STATE_COPY": "Kopiuj", "S_STATE_MEM": "Pamięć",
		"S_STATE_MEM_SUB": "modele zostają w pamięci, pierwsze zdanie nie czeka",
		"S_HOTMODE": "Tryb", "S_HOTMODE_HOLD": "przytrzymanie", "S_HOTMODE_TOGGLE": "przełącznik",
		"S_SUB_HOTMODE": "trzymaj klawisze albo naciśnij raz, by zacząć, i raz, by skończyć",
		"S_SUB_MINMS": "pomija przypadkowe naciśnięcia",
		"S_SUB_ENTER": "od razu wysyła wiadomość",
		"S_SUB_CLIP": "obrazy i pliki wracają bez zmian",
		"S_SUB_TYPE": "pomaga tam, gdzie pole nie przyjmuje wklejania",
		"S_SEC_OVERLAY": "Pasek na ekranie",
		"S_MIC_CHECK": "Sprawdź mikrofon", "S_MIC_CHECK_SUB": "trzy sekundy nagrania i werdykt: poziom, przesterowanie, czy jest mowa", "S_MIC_CHECKING": "Sprawdzam…",
		"S_PAUSE": "Wstrzymaj nagranie", "S_PAUSE_SUB": "w trybie przełącznika: jedno naciśnięcie zatrzymuje nagranie, kolejne je wznawia",
		"S_MCHECK": "Sprawdź zainstalowane modele", "S_MCHECK_SUB": "porównuje pliki modeli z wzorcowymi skrótami", "S_MCHECK_GO": "Sprawdź", "S_MCHECK_RUN": "Sprawdzam…",
		"S_HIST_INSERT": "Wklej",
		"S_LISTS_HINT": "Zamiany i polecenia w jednym pliku — do przeniesienia na inny komputer", "S_LISTS_EXPORT": "Zapisz do pliku", "S_LISTS_IMPORT": "Wczytaj z pliku",
		"S_MIC": "Mikrofon", "S_MIC_DEFAULT": "Domyślny systemowy", "S_MIC_REFRESH": "Odśwież listę",
		"S_MIC_LEVEL": "Poziom wejścia", "S_MIC_QUIET": "cisza",
		"S_ADV_TITLE": "Dobierz model", "S_F_ALL": "wszystkie", "S_F_RU": "rosyjski",
		"S_F_MULTI": "wiele języków", "S_F_PUNCT": "stawia znaki", "S_F_FIT": "mieści się w pamięci",
		"S_ADV_LANGQ": "W jakim języku dyktujesz", "S_ADV_PRIOQ": "Co jest ważniejsze",
		"S_ADV_ACC": "Wybrany dla dokładności.", "S_ADV_SPEED": "Wybrany dla szybkości.",
		"S_ADV_TRQ": "Potrzebne tłumaczenie", "S_ADV_GO": "Poleć",
		"S_ADV_PRIMARY": "główny", "S_ADV_COMPANION": "drugi", "S_ADV_HAVE": "już jest", "S_ADV_APPLY": "Zastosuj",
		"S_ADV_ASK": "Zostaną pobrane: %s — razem %s. Zacząć?",
		"S_SUB_THREADS": "więcej wątków nie zawsze znaczy szybciej — zmierz na swoim komputerze",
		"S_SEC_LLM": "Model redaktora",
		"S_PUNCT": "Znaki i wielkie litery", "S_SUB_PUNCT": "skąd biorą się znaki interpunkcyjne i wielkie litery",
		"S_PUNCT_MODEL": "z modelu", "S_PUNCT_LLM": "od modelu redaktora", "S_PUNCT_OFF": "usuwać",
		"S_SUB_DICT": "Słownik", "S_SUB_PROMPTS": "Prompty",
		"S_TR_TURBO": "⚠ Aktywny model Turbo nie jest uczony tłumaczenia na angielski — do tłumaczenia wybierz inny model w zakładce „Modele”.",
		"S_SUB_TRTARGET": "angielski jest dla Whispera natywny, pozostałe języki docelowe są eksperymentalne",
		"S_TR_EXP": "poza angielskim aplikacja wymusza język wyjścia zamiast tłumaczyć — tekst może zostać w języku mowy",
		"S_REMOTE_ABOUT": "Ustawiony jest serwer zdalny: dźwięk trafia do niego, a obietnica powyżej wtedy nie obowiązuje.",
		"S_UPD": "Aktualizacje", "S_UPD_CHECK": "Sprawdź aktualizacje", "S_UPD_AUTO": "Sprawdzaj przy starcie",
		"S_SUB_UPD": "jedyne zapytanie sieciowe poza pobieraniem modeli",
		"S_UPD_NONE": "Masz najnowszą wersję", "S_UPD_AVAIL": "Dostępna jest wersja %s.",
		"S_UPD_GO": "Aktualizuj", "S_UPD_ERR": "Sprawdzenie nie powiodło się", "S_UPD_DL": "Pobieranie aktualizacji…",
		"S_SEC_SERVICE": "Usługa", "S_SUB_AUTOSTART": "wyłącz, jeśli sam uruchamiasz serwer",
		"S_SUB_PORT": "rozpoznawanie samo się przeładuje",
		"S_MODEL_READY": "Model pobrany — wybierz go, żeby przełączyć",
		"S_FIT_OK": "mieści się", "S_FIT_WARN": "na styk", "S_FIT_BAD": "za mało pamięci", "S_RAM": "Pamięć komputera:",
		"S_HF_PH": "Nazwa modelu — np. qwen2.5 instruct",
		"S_NO_LLM": "Nie ma jeszcze żadnego modelu — znajdź go w zakładce „Szukaj”.",
		"S_NO_LLM_PROF": "Prompty staną się dostępne, gdy będzie zainstalowany model (zakładki „Modele” i „Szukaj”).",
		"S_UPDATED": "Ostatnia aktualizacja modelu", "S_PROF_EDIT": "Edytuj", "S_PROF_CLOSE": "Zwiń",
		"S_CONFIRM_DEL": "Usunąć model „%s”? Będzie można pobrać go ponownie.", "S_FREE": "wolne",
		"S_DEL_ACTIVE": "Usunąć aktywny model „%s”? Rozpoznawanie zatrzyma się, dopóki nie wybierzesz innego — pobrać go można tutaj.",
		"S_WIZ_NEED_MODEL": "Najpierw pobierz model — bez niego nie ma czym rozpoznawać",
		"S_REMOTE_WARN": "Dźwięk będzie wysyłany na ten serwer. Tryb lokalny jest wyłączony.",
		"S_REMOTE_ASK": "Dźwięk przestanie być przetwarzany na tym komputerze i będzie wysyłany na %s. Włączyć tryb zdalny?",
		"S_REMOTE_BADGE": "ZDALNY",
		"S_OK": "Tak", "S_CANCEL": "Anuluj", "S_DL_START": "Pobierz", "S_DL_CANCEL": "Przerwij pobieranie",
		"S_DL_ASK": "Model „%s” nie jest pobrany (%s). Zacząć pobieranie?",
		"S_NOT_FOUND": "nic", "S_MORE": "jeszcze %d ustawień", "S_LESS": "Zwiń %d ustawień",
		"S_AUTHOR_HTML": "<p style=\"font-size:15px;letter-spacing:2px\"><b>Vitalii Yemets</b></p>" +
			"<p>Autor i deweloper {app} — lokalnego narzędzia do dyktowania dla Windows: głos zamienia się w tekst dokładnie tam, gdzie stoi kursor, bez chmury i bez abonamentu.</p>" +
			"<p>Projekt jest otwarty: kod, budowanie i najnowsze wydania są na GitHubie.</p>" +
			"<ul>" +
			"<li>Repozytorium: <span class=\"lnk\" onclick=\"appRepoLink()\">github.com/Vitalii-Yemets/holdtotype</span></li>" +
			"<li>Profil autora: <span class=\"lnk\" onclick=\"appAuthorLink()\">github.com/Vitalii-Yemets</span></li>" +
			"</ul>" +
			"<p>Znalazłeś błąd albo masz pomysł — załóż zgłoszenie w repozytorium.</p>",
		"S_HELP_HTML": "<p class=\"wh\">Jak to działa</p>" +
			"<p>Przytrzymaj skrót — zaczyna się nagrywanie (pasek na dole ekranu pokazuje twój poziom). Puść — dźwięk zostaje rozpoznany, w razie potrzeby przetłumaczony, przepuszczony przez prompty, a gotowy tekst trafia tam, gdzie stoi kursor. ✕ na pasku przerywa na każdym etapie.</p>" +
			"<p>Cała droga: <b>nagranie → rozpoznawanie (Whisper) → tłumaczenie (jeśli włączone) → prompty (LLM) → wklejenie</b>. Każdy etap widać na pasku.</p>" +
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
			"<p>Przycisk „Test” na karcie Mikrofon nagrywa trzy sekundy i je rozbiera: szczyt w decybelach, jaka część nagrania naprawdę zawiera mowę i ile próbek zostało obciętych. Odpowiedź przychodzi słowami: brzmi dobrze, za cicho — podnieś poziom w Windows, przesterowanie — zmniejsz go, nie słychać mowy — czy wybrany jest właściwy mikrofon. To samo mierzy się po każdym dyktowaniu i trafia do dziennika; gdy rozpoznanie wraca puste, pasek nazywa powód — cicho, przesterowanie albo cisza — zamiast mówić tylko, że nic nie usłyszał.</p>" +
			"<p class=\"wh\">Wstrzymanie nagrania</p>" +
			"<p>W trybie przełącznika (jedno naciśnięcie zaczyna, kolejne kończy) można ustawić osobny skrót do pauzy: karta „Dyktowanie”, wiersz „Wstrzymaj nagranie”. Naciśnięcie zatrzymuje nagranie — plakietka pokazuje „Pauza” i nic nie jest zapisywane; kolejne wznawia je, a wszystko powiedziane wcześniej zostaje. Ograniczenie długości nie zadziała w czasie pauzy.</p>" +
			"<p class=\"wh\">Wklejanie z historii</p>" +
			"<p>Każdy wpis w historii ma przycisk „Wklej”: przywraca okno, z którego otworzyliście ustawienia, i wkleja tam tekst jak zwykłe dyktowanie. Gdy nie ma dokąd wracać, tekst po prostu trafia do schowka, a program o tym mówi.</p>" +
			"<p class=\"wh\">Listy w jednym pliku</p>" +
			"<p>Zamiany i polecenia głosowe można zapisać do jednego pliku .json i wczytać na innym komputerze — przyciski pod listą poleceń na karcie „Tekst”. Wczytanie niczego nie kasuje: dodawane są tylko wiersze, których jeszcze nie ma, a program powie, ile dodano i ile pominięto.</p>" +
			"<p class=\"wh\">Nienaruszalność plików</p>" +
			"<p>Dla każdego modelu z katalogu znany jest wzorcowy skrót SHA-256. Po pobraniu plik jest z nim porównywany: gdy się nie zgadza, plik zostaje usunięty i pobieranie można powtórzyć. Przycisk „Sprawdź” na karcie „Modele” tak samo porównuje modele już zainstalowane, a przy aktualizacji programu sprawdzany jest też pobrany instalator — obcy plik się nie uruchomi.</p>" +
			"<p class=\"wh\">Historia dyktowań</p>" +
			"<p>Sekcja „Historia” w lewej kolumnie przechowuje to, co podyktowałeś: tylko tekst, tylko na tym komputerze, dźwięk nigdy nie jest zapisywany. Domyślnie wyłączona, włącza się jednym przełącznikiem w tym samym miejscu. Wpisy trzymają się przez ustaloną liczbę dni i do ustalonej liczby, starsze wypadają same; „Nigdy nie zapisuj z tych programów” wymienia po przecinku te, z których nic nie ma być zapisywane — menedżery haseł, bankowość. Wyszukiwanie obejmuje tekst i nazwę programu, przycisk obok wpisu wkłada go do schowka, a „Wyczyść” usuwa wszystko razem z plikiem <b>{app}-history.json</b>.</p>" +
			"<p class=\"wh\">Komendy głosowe</p>" +
			"<p>Pod zamianami na karcie „Tekst” jest lista komend: to, co powiesz, zamienia się w działanie zamiast w słowa. „Nowy wiersz” i „nowy akapit” wstawiają złamanie — modele nigdy tego nie robią; „anuluj” wyrzuca całe dyktowanie i nic nie wstawia; „wstawić tekst” wkłada cokolwiek, choćby uśmieszek. Przycisk obok listy wypełnia ją typowymi zwrotami w języku interfejsu. Komendy rozpoznawane są jako całe słowa i działają po zamianach, więc prompty i tłumaczenie dostają już gotowy tekst. Zbędne spacje wokół złamań znikają same. Pole poniżej sprawdza zamiany i komendy na dowolnym zdaniu: złamanie pokazuje się jako ⏎.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">równowaga szybkości i dokładności</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">dokładniejszy, polecany</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">najlepsza dokładność na CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Modele</b> — katalog: Base (szybki, na słabe komputery), Small (zrównoważony), Medium i Turbo (dokładniejsze, wolniejsze; „q5” to wersja skwantyzowana — nieco mniejsza i szybsza, prawie bez straty jakości) oraz GigaAM v3 dla rosyjskiego. Przycisk wyboru wskazuje aktywny model (działa od razu, rozpoznawanie startuje na nowo); przy brakującym modelu program pyta, czy go pobrać.</li>" +
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
			"<li>Prompt może mieć własny skrót: dyktowanie nim stosuje tylko ten jeden, raz. Ołówek ✎ otwiera edytor: nazwa, treść promptu, skrót i pole próby ▶, które przepuszcza przykład przez działający model prosto z ustawień.</li>" +
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
			"<li><b>Przenośność</b> — po prostu skopiuj cały folder z plikiem exe (na pendrive, na inny komputer): ustawienia, modele i dziennik leżą obok i jadą razem z nim. Do rejestru nic nie jest zapisywane.</li>" +
			"<li>Przy pierwszym uruchomieniu bez modelu rozpoznawania program sam otwiera katalog i czeka na pobranie.</li>" +
			"<li>Wymagania: Windows 10/11 x64, procesor z AVX2 (mniej więcej od 2013), WebView2 Runtime dla okna ustawień (jest w Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Zasobnik i pliki</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Gotowe…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Ustawienia…</div><div class=\"mock-mi\">Wyłącz</div><div class=\"mock-mi\">Kopiuj ostatni wynik</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Przeładuj config.json</div><div class=\"mock-mi\">Otwórz config.json</div><div class=\"mock-mi\">Otwórz dziennik</div><div class=\"mock-mi\">O programie</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Zakończ</div></div>" +
			"<ul>" +
			"<li>Lewy klik w ikonę — ustawienia; prawy — menu. Kolory ikony: zielony — gotowe, czerwony — nagrywanie, pomarańczowy — rozpoznawanie, szary — wyłączone albo błąd.</li>" +
			"<li><b>config.json</b> — wszystkie ustawienia; ręczne zmiany działają po „Przeładuj config.json” w menu.</li>" +
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
		"S_AUTORUN":    "Uruchamiaj z Windows",
		"S_AUTORUN_SUB": "Wpis w autostarcie bieżącego użytkownika",
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
		"err.webview":        "Для вікна налаштувань потрібен Microsoft WebView2 Runtime (входить до Windows 11).\nЗараз відкриється сторінка завантаження — встановіть його та відкрийте налаштування знову.",
		"status.loading":     "Завантаження моделі…", "status.ready": "Готово — утримуйте %s і говоріть",
		"status.recording": "Йде запис…", "status.transcribing": "Розпізнаю…", "status.disabled": "Вимкнено",
		"status.server.restart": "Сервер розпізнавання впав, перезапускаю…", "status.cfg.err": "Помилка в config.json (див. лог)",
		"status.nomodel":        "Модель розпізнавання не завантажена — оберіть її в налаштуваннях",
		"menu.settings":         "Налаштування…", "menu.enable": "Увімкнути", "menu.disable": "Вимкнути",
		"menu.reload": "Перечитати config.json", "menu.open.config": "Відкрити config.json", "menu.open.log": "Відкрити лог",
		"menu.about": "Про застосунок", "menu.quit": "Вихід",
		"menu.lastcopy": "Копіювати останній результат",
		"ov.copied":     "Скопійовано в буфер обміну", "ov.kept": "Скасовано — текст в «Останньому результаті»",
		"ov.llm.skipped": "Вставлено без профілю «%s»",
		"fd.title":       "Фокус змінився — вставити?", "fd.here": "Вставити сюди", "fd.copy": "Копіювати",
		"ov.speak": "Говоріть…", "ov.transcribing": "Розпізнаю", "ov.inserted": "Вставлено: %d символів",
		"ov.err.mic":       "Мікрофон недоступний — перевірте пристрій у налаштуваннях",
		"ov.err.recognize": "Помилка розпізнавання (див. лог)", "ov.err.paste": "Не вставилося — текст у «Останньому результаті»",
		"ov.moved": "Вікно змінилося — текст у буфері обміну",
		"copy.ok": "Скопійовано",
		"copy.none": "Нема чого копіювати",
		"copy.fail": "Не вдалося скопіювати: %s",
		"mic.busy": "Триває диктування, зараз перевірити не можна", "mic.check.ok": "Чути добре: пік %.0f дБ, мовлення на %.0f%% запису",
		"mic.check.quiet": "Надто тихо: пік %.0f дБ — додайте гучності мікрофона у Windows або сядьте ближче", "mic.check.clipped": "Перевантаження: обрізано %.1f%% відліків — зменште гучність мікрофона", "mic.check.silent": "Мовлення не чути — перевірте, чи вибрано той мікрофон і чи не вимкнений він",
		"ov.quiet": "Надто тихо, майже нічого не чути", "ov.clipped": "Перевантаження — звук обрізано",
		"ov.cmd.cancelled": "Скасовано голосом",
		"ov.silence": "Тиша — нічого не розпізнано", "ov.server.loading": "Сервер ще завантажується",
		"ov.tooshort": "Занадто коротко — тримайте клавіші довше",
		"ov.cancelled": "Скасовано", "ov.editing": "Редагую: %s", "ov.translating": "Перекладаю",
		"ov.llm.needed": "Ця мова потребує LLM-модуль", "td.title": "Перекласти на:", "td.plain": "Без перекладу",
		"cap.title": "{app} — сполучення клавіш", "cap.prompt": "Натисніть нове сполучення клавіш…\n\n(зараз: %s)\n\nEsc — скасувати",
		"cap.selected": "Обрано: %s", "cap.cancelled": "Скасовано",
		"err.hotkey.dup":    "Сполучення «%s» призначено двічі — хоткеї не повинні збігатися",
		"cfg.err.recovered": "config.json пошкоджено (%s).\nФайл збережено як %s, налаштування скинуто до типових.",
		"err.disk.space":    "мало місця на диску: вільно %d МБ, потрібно ~%d МБ",
		"err.save": "не вдалося зберегти налаштування: %s — лишив попередні",
		"err.port": "порт %d не підходить: потрібен номер від 1024 до 65535",
		"err.nolangs": "залиште хоча б одну мову в списку для питання про переклад",
		"ov.mic.lost": "Мікрофон відключено — запис перервано",
		"err.hash": "завантажений файл пошкоджено — спробуйте ще раз",
		"models.check.ok": "Перевірено моделей: %d — усі файли цілі",
		"models.check.none": "Немає що перевіряти — жодна встановлена модель не має еталонного хешу",
		"models.check.bad": "Пошкоджені файли: %s — завантажте модель ще раз",
		"ov.paused": "Пауза",
		"status.paused": "Пауза — запис чекає",
		"hist.insert.gone": "запис не знайдено",
		"hist.insert.nowin": "нікуди вставляти — текст скопійовано в буфер",
		"hist.insert.ok": "вставлено в «%s»",
		"lists.bad": "файл не підходить",
		"lists.saved": "збережено в %s",
		"lists.added": "додано: %d, пропущено: %d",
		"lists.save.title": "Куди зберегти списки",
		"lists.open.title": "Звідки завантажити списки",
		"un.title":          "{app} — видалення", "un.confirm": "Видалити {app} з цього комп'ютера?",
		"un.data": "Видалити також налаштування та завантажені моделі?", "un.done": "{app} видалено.",
		"model.switching": "Перемикаю модель — розпізнавач перезапускається…", "model.del.active": "Не можна видалити активну модель",
		"model.del.ok": "Модель видалено",
		"about.text":   "{app} %s\n\nГолос → текст у позицію курсора.\nПоставте курсор у поле введення, утримуйте %s, скажіть фразу, відпустіть — текст вставиться сам.\n\nРозпізнавання: whisper.cpp, повністю локально й офлайн.\nМодель: %s (мова: %s)\n\nНалаштування: клік по іконці в треї або config.json.\nЛоги: {log} (макс. ~2 МБ).",
		"ov.notranslate": "Активна модель не перекладає — вставлено як розпізнано",
		"ov.engine.fallback": "Другий рушій не запустився — лишаємось на поточному",
		"route.speech": "Мовлення %s", "route.other": "Інші мови", "route.translate": "Переклад",
		"route.lang.auto": "будь-яка мова",
		"route.why.language": "тут точніше, з розділовими", "route.why.otherlang": "99 мов",
		"route.why.translate": "перекладає лише Whisper", "route.why.notinstalled": "російську модель не встановлено",
		"route.why.unknownlang": "мову не задано — розпізнає лише Whisper", "route.why.forced": "примусово в config.json",
		"adv.pick": "Раджу %s.", "adv.companion": "%s добре доповнить — для інших мов і перекладу.",
		"adv.ram": "%d МБ вільно", "status.line": "Готово · %s · %.1f ГБ вільно",
		"ago.now": "щойно", "ago.min": "%d хв тому", "ago.hour": "%d год тому",
		"chars": "%d символів", "inserted.into": "вставлено в %s",
		"punct.prompt": "Додай розділові знаки та великі літери. Не змінюй слова, не перекладай, нічого не додавай. Поверни лише виправлений текст.",
		"err.sherpa.notfound": "розпізнавач sherpa не знайдено: %s",
		"err.sherpa.start": "sherpa-server завершився під час запуску (див. журнал)",
		"err.sherpa.translate": "ця модель не вміє перекладати",
		"err.sherpa.model": "Файл моделі не знайдено: %s — завантажте його в налаштуваннях або виправте sherpa_model у config.json",
		"srv.restarting": "Перезапускаю розпізнавач із новими налаштуваннями…",
	}
	settingsStrings["uk"] = map[string]string{
		"S_TITLE": "{app} — налаштування", "S_TAB_GENERAL": "Основні", "S_TAB_REC": "Розпізнавання",
		"S_TAB_PROC": "Постобробка", "S_TAB_SERVER": "Сервер", "S_TAB_ABOUT": "Про програму",
		"S_PIPE": "голос ▸ розпізнавання ▸ редагування ▸ вставлення",
		"S_DICT": "Словник розпізнавання", "S_DICT_HINT": "Терміни, імена та абревіатури через кому — підказка слуху, не команди.",
		"S_TR": "Переклад", "S_TR_HINT": "Переклад виконує Whisper: англійською — штатним режимом, іншими мовами — експериментально, примусовою мовою виводу (якість залежить від пари мов).",
		"S_TR_TURBO":   "⚠ Активна модель Turbo не навчена перекладу англійською — для перекладу оберіть іншу модель на вкладці «Моделі».",
		"S_TR_DEFAULT": "Завжди перекладати цільовою мовою", "S_TR_TARGET": "Цільова мова", "S_TR_ASK": "Вибір мови", "S_TR_ASK_NEVER": "Не запитувати (цільова мова)",
		"S_TR_ASK_ALWAYS": "Запитувати щоразу", "S_TR_ASK_TIMEOUT": "Запитувати з таймаутом", "S_TR_SECONDS": "Таймаут, с",
		"S_TR_LANGS": "Мови в діалозі",
		"S_LLM":      "Профілі обробки", "S_LLM_HINT": "Позначені профілі застосовуються по черзі, зверху вниз, під час звичайного диктування. Нічого не позначено — текст вставляється як є. Хоткей профілю застосовує разово лише його.",
		"S_PROF_ASIS": "Як є", "S_PROF_WT": "Переклад → English (швидкий)",
		"S_PROF_ADD": "Додати профіль", "S_PROF_NAME": "Ім'я", "S_PROF_PROMPT": "Промпт", "S_PROF_HOTKEY": "Хоткей",
		"S_PROF_SET": "Задати…", "S_PROF_CLEAR": "Скинути", "S_PROF_TEST": "Перевірка",
		"S_PROF_EDIT": "Редагувати", "S_PROF_CLOSE": "Згорнути",
		"S_CONFIRM_DEL": "Видалити модель «%s»? Її можна буде завантажити знову.", "S_FREE": "вільно",
		"S_DEL_ACTIVE": "Видалити активну модель «%s»? Розпізнавання зупиниться, доки ви не виберете іншу — завантажити її можна тут же.",
		"S_WIZ_NEED_MODEL": "Спочатку завантажте модель — без неї немає чим розпізнавати",
		"S_SUB_MODELS": "Моделі", "S_SUB_SEARCH": "Пошук", "S_SUB_PROMPTS": "Промпти",
		"S_SUB_PARAMS": "Параметри", "S_SUB_DICT": "Словник", "S_SUB_TR": "Переклад",
		"S_SUB_INFO": "Інформація", "S_SUB_HELP": "Довідка",
		"S_UPD": "Оновлення", "S_UPD_CHECK": "Перевірити оновлення", "S_UPD_AUTO": "Перевіряти під час запуску",
		"S_UPD_NONE": "Встановлено останню версію", "S_UPD_AVAIL": "Доступна версія %s.",
		"S_UPD_GO": "Оновити", "S_UPD_ERR": "Не вдалося перевірити оновлення", "S_UPD_DL": "Завантажую оновлення…",
		"S_SUB_AUTHOR": "Про автора",
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
		"S_NO_LLM":      "Ще не встановлено жодної моделі — знайдіть і завантажте на вкладці «Пошук».",
		"S_NO_LLM_PROF": "Промпти стануть доступні після встановлення моделі (вкладки «Моделі» та «Пошук»).",
		"S_UPDATED":     "Дата останнього оновлення моделі",
		"S_HOTKEY":      "Сполучення клавіш", "S_CHANGE": "Змінити…", "S_UILANG": "Мова інтерфейсу", "S_AUTO": "Як у системі",
		"S_SEC_SOUND": "Звук", "S_SEC_BEHAVIOR": "Поведінка", "S_BEEP": "Звукові сигнали запису", "S_SOUND": "Сигнал",
		"S_SND_SPEECH": "Системний (мовлення)", "S_SND_CHIME": "Дзвіночок", "S_SND_SOFT": "М'який", "S_SND_MARIMBA": "Марімба",
		"S_SND_BLIP": "Бліп", "S_SND_POP": "Поп",
		"S_AUTOENTER": "Enter після вставлення (автовідправлення)", "S_RESTORE": "Відновлювати буфер обміну після вставлення",
		"S_NAV_HISTORY": "Історія", "S_HIST_ON": "Зберігати історію диктувань", "S_HIST_ON_SUB": "лише текст, на цьому комп'ютері; звук не зберігається ніколи",
		"S_HIST_DAYS": "Скільки днів зберігати", "S_HIST_MAX": "Скільки записів зберігати",
		"S_HIST_SKIP": "Не записувати з цих програм", "S_HIST_SKIP_SUB": "через кому: keepass.exe, 1password.exe",
		"S_HIST_LIST": "Записи", "S_HIST_CLEAR": "Очистити", "S_HIST_COPY": "Копіювати",
		"S_HIST_FIND": "Знайти в історії…", "S_HIST_EMPTY": "Історії поки немає", "S_HIST_ASK": "Видалити всю історію диктувань?",
		"S_SEC_CMD": "Голосові команди", "S_CMD_HINT": "Сказане перетворюється на перенесення рядка, знак або скасування замість того, щоб потрапити в текст. Шукаються цілими словами, застосовуються згори вниз, уже після замін.",
		"S_CMD_ADD": "Додати команду", "S_CMD_PRESET": "Додати звичайні", "S_CMD_PH": "новий рядок",
		"S_CMD_NEWLINE": "перенесення рядка", "S_CMD_PARAGRAPH": "новий абзац", "S_CMD_TEXT": "підставити текст", "S_CMD_CANCEL": "скасувати диктування",
		"S_CMD_TEXT_PH": "що підставити", "S_CMD_EMPTY": "Команд поки немає", "S_CMD_DEL": "Видалити команду",
		"S_CMD_P_NEWLINE": "новий рядок", "S_CMD_P_PARAGRAPH": "новий абзац", "S_CMD_P_CANCEL": "скасувати",
		"S_SEC_REPLACE": "Заміни після розпізнавання", "S_REPLACE_HINT": "Те, що почулося неправильно, стає тим, що ви мали на увазі — одразу після розпізнавання, до промптів. Застосовуються згори вниз.",
		"S_REPL_ADD": "Додати заміну", "S_REPL_FROM_PH": "гіт хаб", "S_REPL_TO_PH": "GitHub",
		"S_REPL_WHOLE": "цілі слова", "S_REPL_CASE": "регістр", "S_REPL_EMPTY": "Замін поки немає",
		"S_REPL_DEL": "Видалити заміну", "S_REPL_TEST_PH": "напишіть фразу, щоб перевірити заміни й команди",
		"S_SEC_RULES": "Правила для програм", "S_RULES_HINT": "Для окремих програм вставка може працювати інакше. Виграє перше відповідне правило.",
		"S_RULE_ADD": "Додати правило", "S_RULE_PH": "chrome.exe, msedge.exe",
		"S_RULE_PASTE_INH": "вставка: як усюди", "S_RULE_ENTER_INH": "Enter: як усюди", "S_RULE_DELAY_NONE": "без затримки", "S_RULE_PROMPT_INH": "промпти: як усюди",
		"S_RULE_CLIP": "буфер обміну", "S_RULE_TYPE": "посимвольно", "S_RULE_ENTER_ON": "з Enter", "S_RULE_ENTER_OFF": "без Enter",
		"S_RULE_NOPROMPT": "без промптів", "S_RULE_LAST": "остання вставка: %s", "S_RULE_EMPTY": "Правил поки немає",
		"S_RULE_DEL": "Видалити правило", "S_RULE_PROMPTS": "Промпти",
		"S_PASTE_DELAY": "Затримка перед вставкою", "S_PASTE_DELAY_SUB": "коли програма не встигає прийняти текст",
		"S_OVPOS": "Де показувати смужку", "S_OVPOS_SUB": "біля курсора — поряд із місцем введення; якщо застосунок його не показує, поряд із вказівником миші",
		"S_OVPOS_BOTTOM": "Унизу екрана", "S_OVPOS_TOP": "Угорі екрана", "S_OVPOS_CARET": "Біля курсора",
		"S_OVTEXT": "Показувати розпізнаний текст", "S_OVTEXT_SUB": "у смужці після вставки, замість кількості символів",
		"S_OVERLAY": "Екранний індикатор", "S_ANIM": "Анімація запису/розпізнавання", "S_TYPEMODE": "Посимвольне введення (для полів без вставлення)",
		"S_RECLANG": "Мова розпізнавання", "S_RECAUTO": "Автовизначення",
		"S_MODELS": "Моделі розпізнавання", "S_DL": "Завантажити", "S_DEL": "Видалити",
		"S_M_BASE": "швидка, для слабких ПК", "S_M_SMALL": "баланс швидкості й точності", "S_M_MED": "точніша, рекомендуємо", "S_M_TURBO": "максимум точності на CPU",
		"S_M_CUSTOM": "власна (з config.json)",
		"S_MIC_CHECK": "Перевірити мікрофон", "S_MIC_CHECK_SUB": "три секунди запису та розбір: гучність, перевантаження, чи є мовлення", "S_MIC_CHECKING": "Перевіряю…",
		"S_PAUSE": "Пауза в записі", "S_PAUSE_SUB": "працює в режимі фіксації: натиснули — запис завмер, натиснули ще раз — пішов далі",
		"S_MCHECK": "Перевірити встановлені моделі", "S_MCHECK_SUB": "звіряє файли моделей з еталонними хешами", "S_MCHECK_GO": "Перевірити", "S_MCHECK_RUN": "Перевіряю…",
		"S_HIST_INSERT": "Вставити",
		"S_LISTS_HINT": "Заміни й команди одним файлом — перенести на інший комп'ютер", "S_LISTS_EXPORT": "Зберегти у файл", "S_LISTS_IMPORT": "Завантажити з файлу",
		"S_MIC":      "Мікрофон", "S_MIC_DEFAULT": "Системний за замовчуванням", "S_MIC_REFRESH": "Оновити список",
		"S_MIC_LEVEL": "Рівень сигналу", "S_MIC_QUIET": "тихо",
		"S_THREADS": "Потоки CPU", "S_MINMS": "Мін. запис, мс", "S_MAXSEC": "Макс. запис, с",
		"S_AUTOSTART": "Запускати whisper-server автоматично", "S_PORT": "Порт", "S_SERVEREXE": "Шлях до whisper-server",
		"S_SERVERURL": "Зовнішній сервер (URL)", "S_URLHINT": "Якщо задано — свій сервер не запускається",
		"S_SAVED": "Збережено",
		"S_ABOUT_HTML": "<p><b>Голос → текст у позицію курсора.</b></p><p>Поставте курсор у поле введення, утримуйте сполучення клавіш, скажіть фразу, відпустіть — текст вставиться сам.</p><p>Повністю локально й офлайн. Технології: <b>Go + WinAPI</b>, <b>WebView2</b>, <b>whisper.cpp</b>, <b>llama.cpp</b>, <b>miniaudio</b>; моделі з Hugging Face.</p><p>Логи не перевищують ~2 МБ.</p>",
		"S_VERSION":    "Версія",
		"S_LVL_SIMPLE": "просто", "S_LVL_ALL": "усе", "S_SEARCH": "Знайти налаштування…",
		"S_GRP_WORK": "Робота", "S_GRP_REC": "Розпізнавання", "S_GRP_OTHER": "Інше",
		"S_NAV_STATE": "Стан", "S_NAV_DICT": "Диктування", "S_NAV_MIC": "Мікрофон", "S_NAV_MODELS": "Моделі",
		"S_NAV_TEXT": "Текст", "S_NAV_TR": "Переклад", "S_NAV_SYSTEM": "Система", "S_NAV_ABOUT": "Про програму",
		"S_STATE_HINT": "тримайте й говоріть — текст з'явиться там, де стоїть курсор",
		"S_STATE_RU": "Російська мова", "S_STATE_OTHER": "Інші мови", "S_STATE_PROC": "Постобробка",
		"S_CHANGE_MODEL": "Змінити", "S_PICK_MODEL": "Підібрати", "S_STATE_GET": "Завантажити",
		"S_RETRY": "Повторити", "S_BERR_OPEN": "Відкрити налаштування сервера",
		"S_STATE_LAST": "Останнє диктування", "S_STATE_COPY": "Копіювати", "S_STATE_MEM": "Пам'ять",
		"S_STATE_MEM_SUB": "моделі лишаються в пам'яті, перша фраза йде без затримки",
		"S_HOTMODE": "Режим", "S_HOTMODE_HOLD": "утримання", "S_HOTMODE_TOGGLE": "перемикач",
		"S_SUB_HOTMODE": "тримайте клавіші або натисніть раз, щоб почати, і раз, щоб зупинити",
		"S_SUB_MINMS": "відсікає випадкові натискання",
		"S_SUB_ENTER": "надсилає повідомлення одразу",
		"S_SUB_CLIP": "зображення й файли повертаються як були",
		"S_SUB_TYPE": "допомагає там, де поле не приймає вставку",
		"S_SEC_OVERLAY": "Смужка на екрані",
		"S_ADV_TITLE": "Підібрати модель", "S_F_ALL": "усі", "S_F_RU": "російська",
		"S_F_MULTI": "багато мов", "S_F_PUNCT": "ставить розділові", "S_F_FIT": "вміщується в пам'ять",
		"S_ADV_LANGQ": "Якою мовою ви диктуєте", "S_ADV_PRIOQ": "Що важливіше",
		"S_ADV_ACC": "Обрано за точністю.", "S_ADV_SPEED": "Обрано за швидкістю.",
		"S_ADV_TRQ": "Потрібен переклад", "S_ADV_GO": "Порадити",
		"S_ADV_PRIMARY": "основна", "S_ADV_COMPANION": "друга", "S_ADV_HAVE": "уже є", "S_ADV_APPLY": "Застосувати",
		"S_ADV_ASK": "Буде завантажено: %s — разом %s. Почати?",
		"S_SUB_THREADS": "більше потоків не завжди швидше — виміряйте на своїй машині",
		"S_SEC_LLM": "Модель-редактор",
		"S_PUNCT": "Розділові знаки й великі літери", "S_SUB_PUNCT": "звідки беруться розділові знаки й великі літери",
		"S_PUNCT_MODEL": "від моделі", "S_PUNCT_LLM": "від моделі-редактора", "S_PUNCT_OFF": "прибирати",
		"S_SUB_TRTARGET": "англійська для Whisper рідна, інші цілі — експериментальні",
		"S_TR_EXP": "крім англійської застосунок примусово задає мову виводу, а не перекладає — текст може лишитися мовою мовлення",
		"S_REMOTE_ABOUT": "Задано віддалений сервер: звук іде на нього, і обіцянка вище не діє, поки він увімкнений.",
		"S_SUB_UPD": "єдиний мережевий запит, окрім завантаження моделей",
		"S_SEC_SERVICE": "Службове", "S_SUB_AUTOSTART": "вимкніть, якщо запускаєте сервер самі",
		"S_SUB_PORT": "розпізнавач перезапуститься сам",
		"S_MODEL_READY": "Модель завантажено — оберіть її, щоб перемкнутися",
		"S_REMOTE_WARN": "Звук ітиме на цей сервер. Локальний режим вимкнено.",
		"S_REMOTE_ASK": "Аудіо перестане оброблятися на цьому комп'ютері й надсилатиметься на %s. Увімкнути віддалений режим?",
		"S_REMOTE_BADGE": "ВІДДАЛЕНО",
		"S_OK": "Так", "S_CANCEL": "Скасувати", "S_DL_START": "Завантажити", "S_DL_CANCEL": "Скасувати завантаження",
		"S_DL_ASK": "Модель «%s» не завантажена (%s). Почати завантаження?",
		"S_NOT_FOUND": "нічого", "S_MORE": "Ще %d налаштувань", "S_LESS": "Згорнути %d налаштувань",
		"S_HELP_HTML": "<p class=\"wh\">Як це працює</p>" +
			"<p>Тримайте сполучення — починається запис (смужка внизу екрана показує ваш рівень). Відпустіть — звук розпізнається, за потреби перекладається, проходить через промпти, і готовий текст з'являється там, де стоїть курсор. ✕ на смужці скасовує на будь-якому кроці.</p>" +
			"<p>Увесь шлях: <b>запис → розпізнавання (Whisper) → переклад (якщо ввімкнено) → промпти (LLM) → вставка</b>. Кожен крок видно на смужці.</p>" +
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
			"<p class=\"wh\">Пауза в записі</p>" +
			"<p>У режимі фіксації (натиснули — пише, натиснули ще раз — зупинилася) можна задати окреме сполучення для паузи: на «Диктуванні», рядок «Пауза в записі». Натиснули — запис завмер, плашка показує «Пауза», і нічого не записується; натиснули ще раз — запис пішов далі, а все сказане до паузи зберігається. Обмеження довжини запису на паузі не спрацьовує.</p>" +
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
			"<div class=\"mock\"><div class=\"mock-row\"><span class=\"mock-radio on\"></span><b>Small</b><span class=\"mock-note\">рівновага швидкості й точності</span><span style=\"margin-left:auto\">466 MB</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Medium (q5)</b><span class=\"mock-note\">точніша, рекомендована</span><span style=\"margin-left:auto\">539 MB ⭳</span></div>" +
			"<div class=\"mock-row\"><span class=\"mock-radio\"></span><b>Turbo (q5)</b><span class=\"mock-note\">найкраща точність на CPU</span><span style=\"margin-left:auto\">574 MB ⭳</span></div></div>" +
			"<ul>" +
			"<li><b>Моделі</b> — каталог: Base (швидка, для слабких ПК), Small (збалансована), Medium і Turbo (точніші, повільніші; «q5» — квантована збірка: трохи менша й швидша майже без утрати якості), а також GigaAM v3 для російської. Радіокнопка обирає активну модель (діє одразу, розпізнавач перезапускається); якщо моделі немає, програма запитує, чи завантажити її.</li>" +
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
			"<li><b>Портативність</b> — просто скопіюйте всю теку з exe (на флешку, на інший комп'ютер): налаштування, моделі та журнал лежать поруч і їдуть разом. У реєстр нічого не пишеться.</li>" +
			"<li>При першому запуску без моделі розпізнавання програма сама відкриває каталог і чекає на завантаження.</li>" +
			"<li>Вимоги: Windows 10/11 x64, процесор з AVX2 (приблизно від 2013 року), WebView2 Runtime для вікна налаштувань (входить до Windows 11).</li>" +
			"</ul>" +
			"<p class=\"wh\">Трей і файли</p>" +
			"<div class=\"mock\" style=\"max-width:290px\"><div class=\"mock-mi dim\">{app} — Готово…</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Налаштування…</div><div class=\"mock-mi\">Вимкнути</div><div class=\"mock-mi\">Копіювати останній результат</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Перечитати config.json</div><div class=\"mock-mi\">Відкрити config.json</div><div class=\"mock-mi\">Відкрити журнал</div><div class=\"mock-mi\">Про програму</div><hr class=\"mock-sep\"><div class=\"mock-mi\">Вийти</div></div>" +
			"<ul>" +
			"<li>Лівий клац по значку — налаштування; правий — меню. Кольори значка: зелений — готово, червоний — запис, помаранчевий — розпізнавання, сірий — вимкнено або помилка.</li>" +
			"<li><b>config.json</b> — усі налаштування; правки вручну діють після «Перечитати config.json» у меню.</li>" +
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
		"S_AUTORUN":    "Запускати разом із Windows",
		"S_AUTORUN_SUB": "Запис в автозапуску поточного користувача",
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
