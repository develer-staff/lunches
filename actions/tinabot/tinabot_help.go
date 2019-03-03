package tinabot

// HelpStr is Tinabot help string
const HelpStr = `Elenco domandi supportati da tinabot9000:

*PER ORDINARE UN PIATTO:*
‘@Tinabot 9000 per <utente> <ordine>‘
*<utente>* può essere ‘me‘ per ordinare per se stessi, oppure il nome di un altro utente slack (che verrà avvisato!). *E' possibile ordinare per ospiti esterni senza utente slack* chiamandoli ‘guest_<nome>‘.

*<ordine>* può essere una serie di stringhe separate da spazi, tinabot9000 cercherà di fare un il meglio che può per capire il piatto tra le voci presenti nel menù.

*E' possibile usare i seguenti caratteri speciali per personalizzare l'ordine:*
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
‘@Tinabot 9000 menu <stringa menu>‘
*<stringa menu>* può essere multilinea. E' sufficiente copiare le celle dal file excel inviato per mail dal tuttobene. Chiunque può impostare il menù.
`
