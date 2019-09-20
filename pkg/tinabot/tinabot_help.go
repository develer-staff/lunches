package tinabot

// HelpStr is Tinabot help string
const HelpStr = `Elenco comandi supportati da tinabot9000:

*PER ORDINARE UN PIATTO:*
‘@Tinabot 9000 per <utente> <ordine>‘
*<utente>* può essere ‘me‘ per ordinare per se stessi, oppure il nome di un altro utente slack (che verrà avvisato!). *E' possibile ordinare per ospiti esterni senza utente slack* chiamandoli ‘guest_<nome>‘.

*<ordine>* può essere una serie di stringhe separate da spazi, tinabot9000 cercherà di fare un il meglio che può per capire il piatto tra le voci presenti nel menù.

*E' possibile usare le seguenti funzionalità speciali per personalizzare l'ordine:*
*&* - E commerciale
Se usato per unire due o più stringhe, permette di personalizzare il secondo, aggiungendo uno o più contorni. Il piatto creato sarà indicato come *un* solo secondo.
es.
‘‘‘
@Tinabot 9000 per me scorfano & piselli

Tinabot 9000:
Trovato: Scorfano con ginger lime (secondi piatti)
Trovato: Piselli (contorni)
Piatto personalizzato: Scorfano con ginger lime con Piselli
Ok, aggiunto 1 piatto per batt
‘‘‘

*+* - Più
Se usato per unire due o più stringhe, permette, con un’unico ordine, di richiedere *più di una portata alla volta*: un primo e un secondo, un secondo e una frutta, etc… I piatti ordinati saranno assegnati allo stesso utente
es. 
‘‘‘
@Tinabot 9000 per me fusilli + peposo
Tinabot 9000:
Trovato: Fusilli con salsiccia pomodoro e olive (primi piatti)
Trovato: Peposo con patate in umido (secondi piatti)
Ok, aggiunti 2 piatti per batt
‘‘‘

*"* - Virgolette
Inserendo una stringa tra virgolette tinabot9000 aggiungerà *testualmente* quello che avete scritto all’ordine. E’ una funzionalità utile da usare in casi particolari (tipo per le insalate o per il “senza glutine”). Mi raccomando, è pensato per essere usato in situazioni particolari, quindi non abusatene!
‘‘‘
@Tinabot 9000 per me "pasta senza glutine al ragù"

Tinabot 9000:
Aggiungo testualmente: 'pasta senza glutine al ragù'
Ok, aggiunto 1 piatto per batt
‘‘‘

*come* - Copia ordine
Indicando "per me come <utente>", tinabot9000 copierà l'ordine dell'utente indicato
‘‘‘
@Tinabot 9000 per me come djeasy

Tinabot 9000:
Ok, copio l'ordine di djeasy:
Baccalà alla livornese con ceci
Ok, aggiunto 1 piatto per batt
‘‘‘

Le funzionalità speciali possono anche essere combinate tra loro

*PER CANCELLARE UN ORDINE:*
‘@Tinabot 9000 per <utente> niente‘
*<utente>* può essere ‘me‘ o il nome di un altro utente slack (che verrà avvisato). 

*PER VEDERE I PIATTI ORDINATI:*
‘@Tinabot 9000 ordine‘

*PER INVIARE LA MAIL AL TUTTOBENE:*
‘@Tinabot 9000 email‘
Verrà fornito un link che autocompone una mail nel client di posta locale. Chiunque può inviare la mail al tuttobene.

*PER VEDERE IL MENÙ DEI PIATTI:*
‘@Tinabot 9000 menu‘

*PER IMPOSTARE IL MENÙ DEI PIATTI:*
‘@Tinabot 9000 setmenu <stringa menu>‘
*<stringa menu>* può essere multilinea. E' sufficiente copiare le celle dal file excel inviato per mail dal tuttobene. Chiunque può impostare il menù.

*PER IMPOSTARE IL REMINDER:*
Nel caso tu abbia attivato la funzionalità reminder, se è impostato un menù valido per il giorno e non hai ancora ordinato, alle 11:50 ti verrà inviato un messaggio privato contenente il menù del giorno.
Ecco come fare:
‘@Tinabot 9000 remind <giorni>‘
*<giorni>* può essere ‘on‘ per indicare tutti i giorni:
‘‘‘
@Tinabot 9000 remind on
Tinabot 9000:
Reminder attivo tutti i giorni
‘‘‘

E' anche possibile specificare i singoli giorni separati da virgola:
‘‘‘
@Tinabot 9000 remind lun, mar
Tinabot 9000:
Reminder attivo lunedì, martedì
‘‘‘

Per disattivare il reminder usare ‘off‘:
‘‘‘
@Tinabot 9000 remind off
Tinabot 9000:
Reminder disattivato
‘‘‘

*PER VEDERE LO STATO DEL REMINDER:*
‘@Tinabot 9000 remind‘

*PER SEGNARE IL PRANZO:*
Tinabot 9000 è in grado di segnare *in automatico* il pranzo sul foglio google di riepilogo, usato dall'amministrazione per tenere traccia dei pasti e dei buoni.
Se hai ordinato il pranzo con Tinabot, *verrà registrato in automatico alle 14:00*.
E' comunque possibile modificare quanto segnato sul foglio manualmente:
‘@Tinabot 9000 segna <cibo>‘
*<cibo>* può essere:
‘P‘, ‘PS‘, ‘PD‘, ‘S‘, ‘SD‘, ‘D‘, ‘PSD‘ oppure ‘Niente‘

es. 
‘‘‘
@Tinabot 9000 segna p
Tinabot 9000:
Ok, segnato 'P' per batt sul foglio dei pranzi
‘‘‘


`
