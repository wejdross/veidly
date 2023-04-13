const defaultLocale = ["en", "US"]
const defaultLang = defaultLocale[0]

export const lsKey = "sweelang"

// allowed languages for api
const apiLang = {
    en: true,
    pl: true,
    de: true
}

const apiLocale = {
    "en": ["en", "EN"],
    "pl": ["pl", "PL"],
    "de": ["de", "DE"]
}

export function getLangLabel(c) {
    let l = getSupportedLanguage()
    return (c.Translations && c.Translations[l]) || c.En
}

export function getLocaleFromNavigator() {

    let l = navigator.language;
    let parts = l.split("-")
    if(parts.length < 2) {
        parts = l.split("_")
    }
    return parts 
}

export function getSupportedLanguage() {
    let l = undefined

    if (!localStorage.getItem(lsKey)) {
        let parts = getLocaleFromNavigator()
        if(parts.length == 2) {
            l = parts[0]
        }
    } else {
        l = localStorage.getItem(lsKey)
    }
    if (!apiLang[l]) {
        l = defaultLang
    }
    return l
}

export function getSupportedLocale() {
    let parts = getLocaleFromNavigator()
    let c = defaultLocale[1]
    if(parts.length == 2) {
        c = parts[1]
    }
    let ret = [getSupportedLanguage(), c]
    return ret
}

function validateLocale(lo) {
    for(let k in lo) { 
        if(k in locale2)
            throw new Error("duplicated locale: " + k)
    }
    return lo
}

/* used in sentences like 'occurring each eachFdayIndex[x] ' */
export const eachFdayIndex = {
    MO: {
        pl: "Poniedziałek",
        en: "Monday",
        de: "Montag"
    },
    TU: {
        pl: "Wtorek",
        en: "Tuesday",
        de: "Dienstag"
    },
    WE: {
        pl: "Środę",
        en: "Wednesday",
        de: "Mittwoch"
    },
    TH: {
        pl: "Czwartek",
        en: "Thursday",
        de: "Donnerstag"
    },
    FR: {
        pl: "Piątek",
        en: "Friday",
        de: "Freitag"
    },
    SA: {
        pl: "Sobotę",
        en: "Saturday",
        de: "Samstag"
    },
    SU: {
        pl: "Niedziele",
        en: "Sunday",
        de: "Sonntag"
    }
}

export const sectionTitles = {
    MY_ACHIEVEMENTS: {
        pl: "Moje osiągnięcia",
        en: "My achievements",
        de: "Meine Leistungen"
    },
}
export const locale2 = {
    "INVALID_PASS_ERR": {
        en: "Invalid password, is should have at least 8 characters, one uppercase, one lowercase, one numeric, and one special character",
        de: "Invalid password, is should have at least 8 characters, one uppercase, one lowercase, one numeric, and one special character",
        pl: "Prawidłowe hasło powinno mieć conajmniej 8 znaków, 1 wielką literę, 1 małą, 1 liczbę, oraz jeden znak specjalny"
    },
    "INVOICE_DISCLAIMER": {
        pl: "Twoje dane powinny zawierać Nazwę, nip, oraz adres twojej firmy",
        en: "Your data should contain Name of your company, ID, and address",
        de: "Your data should contain Name of your company, ID, and address",
    },
    "WITHOUT_INVOICE_DATA": {
        pl: "Bez danych do fakturowania nie będziemy w stanie wystawiać ci faktur",
        en: "Without invoicing data we won't be able to issue you our invoices",
        de: "Without invoicing data we won't be able to issue you our invoices"
    },
    "PROVIDE_INVOICE_DATA": {
        pl: "Podaj dane do faktur",
        en: "Provide data for invoices",
        de: "Provide data for invoices"
    },
    "INVOICE_DATA": {
        pl: "Dane do fakturowania",
        en: "Invoicing data",
        de: "Invoicing data"
    },
    "INVOICES": {
        pl: "Faktury",
        en: "Invoices",
        de: "Invoices"
    },
    "PLACE_SUPPORTS_DISABLED": {
        pl: "Obiekt dostosowany dla niepełnosprawnych",
        en: "Place adapted for disabled",
        de: ""
    },
    "TRAINING_SUPPORTS_DISABLED": {
        pl: "Trening dostosowany dla niepełnosprawnych",
        en: "Training adapted for disabled",
        de: ""
    },
    "TRAINING_SCHEDULE": {
        pl: "Rozkład pojedyńczych zajęć",
        en: "Single training schedule",
        de: ""
    },
    "TIME_OF_STARTING_COURSE": {
        "pl": "Czas rozpoczęcia kursu",
        "en": "Starting time for course"
    },
    "DATE_OF_STARTING_COURSE": {
        "pl": "Data rozpoczęcia kursu",
        "en": "Staring date for course"
    },
    "TIME_OF_ENDING_COURSE": {
        "pl": "Czas zakończnia kursu",
        "en": "Ending time for course"
    },
    "DATE_OF_ENDING_COURSE": {
        "pl": "Data zakończenia kursu",
        "en": "Ending date for course"
    },
    "DELETE_THIS_COURSE": {
        "pl": "usuń ten kurs",
        "en": "remove this course",
        "de": ""
    },
    "CHANGE_SCHEDULE": {
        "pl": "Zmień rozkład zajęć w obrębie kursu",
        "en": "Change schedule for training",
        "de": ""
    },
    "MOVE_THIS_COURSE": {
        "pl": "Przesuń ten kurs",
        "en": "Move this course",
        "de": ""
    },
    "ABOUT_ME": {
        "pl": "O mnie",
        "en": "About me",
        "de": "Über mich"
    },
    "ACCEPT": {
        "pl": "Potwierdź",
        "en": "Confirm",
        "de": "Confirm"
    },
    "ACCEPTED": {
        "pl": "zaakceptowano",
        "en": "accepted",
        "de": "akzeptiert"
    },
    "ACCEPTING_RSV": {
        "pl": "Aktceptowanie rezerwacji",
        "en": "Accepting reservation",
        "de": "Annahme der Reservierung"
    },
    "ACCEPTING_RSV_IS_NOT_NECESSARY": {
        "pl": "Zaakceptowanie tej rezerwacji nie jest konieczne",
        "en": "Accepting this reservation is not necessary",
        "de": "Die Annahme dieser Reservierung ist erforderlich"
    },
    "ACCEPT_USER": {
        "pl": "Zaakceptuj użytkownika",
        "en": "Accept user",
        "de": "Benutzer annehmen"
    },
    "ACCOUNT": {
        "pl": "konto",
        "en": "account",
        "de": "Konto"
    },
    "ACCOUNT_STEP_DESC": {
        "pl": "Twoje osobiste konto, to Twoje centrum dowodzenia całym biznesem.",
        "en": "Your personal account is all in one solution for managing Your business",
        "de": "Dein persönliches Konto dient dazu, dein Gescheft zu verwalten"
    },
    "ACCURACY": {
        "pl": "Trafność",
        "en": "Accuracy",
        "de": "Richtigkeit"
    },
    "ACHIEVEMENTS": {
        "pl": "Osiągnięcia",
        "en": "Achievements",
        "de": "Leistungen"
    },
    "ACTIVATE": {
        "pl": "Aktywuj",
        "en": "Activate",
        "de": "Aktiviere"
    },
    "ACTIVATION": {
        "pl": "Aktywacja",
        "en": "Activation",
        "de": "Aktivierung"
    },
    "ACTIVATION_WARN": {
        "pl": "Poprzez aktywacje swojego konta, twoje treningi ponownie będą widoczne w searchu i będzie się można na nie zapisać",
        "en": "By activating Your account, Your trainings will again be visible in the search and customers can sign up for them",
        "de": "Wenn du dein Konto aktivierst, werden deine Trainings wieder in den Suchergebnissen sichtbar und Kunden werden sich für sie anmelden können"
    },
    "ACTIVE": {
        "pl": "Aktywne",
        "en": "Active",
        "de": "Aktiv"
    },
    "ACTIVITY": {
        "pl": "Zajęcia",
        "en": "Activity",
        "de": "Aktivitäten"
    },
    "ADD": {
        "pl": "Dodaj",
        "en": "Add",
        "de": "Hinzufügen"
    },
    "ADD_FMT": {
        "pl": "Dodaj %s",
        "en": "Add %s",
        "de": "Füge %s hinzu"
    },
    "ADD_HOUR": {
        "pl": "Dodaj godzinę",
        "en": "Add hour",
        "de": "Uhrzeit hinzufügen"
    },
    "ADD_LIMIT": {
        "pl": "Wystarczy dodać limit!",
        "en": "Add a limit!",
        "de": "Setze ein Limit!"
    },
    "ADD_NEW_LINK": {
        "pl": "Dodaj nowy link",
        "en": "Add new link",
        "de": "Einen neuen Link hinzufügen"
    },
    "ADD_OR_SELECT_CAT": {
        "pl": "Dodaj lub wybierz kategorie...",
        "en": "Add or choose category...",
        "de": "Kategorie wählen oder hinzufügen"
    },
    "ADD_PAYOUT_DATA": {
        "pl": "Podaj dane do wypłat",
        "en": "Provide payment details",
        "de": "Zahlungsinformationen angeben"
    },
    "ADD_REVIEW": {
        "pl": "Wystaw opinię",
        "en": "Add review",
        "de": "Bewertung hinzufügen"
    },
    "ADD_SECTION": {
        "pl": "Dodaj sekcję",
        "en": "Add section",
        "de": "Sektion hinzufügen"
    },
    "ADD_TRAINING": {
        "pl": "Dodaj zajęcia",
        "en": "Add training",
        "de": "Training hinzufügen"
    },
    "ADJUST_AGE": {
        "pl": "Dostosuj wiek:",
        "en": "Adjust age:",
        "de": "Alter anpassen"
    },
    "ADJUST_DISTANCE": {
        "pl": "Dostosuj odległość:",
        "en": "Adjust distance:",
        "de": "Entfernung anpassen:"
    },
    "ADJUST_GROUP_SIZE": {
        "pl": "Dobierz rozmiar grupy:",
        "en": "Adjust group size:",
        "de": "Die Größe der Gruppe anpassen:"
    },
    "ADJUST_HOURS": {
        "pl": "Dobierz godziny:",
        "en": "Adjust hours:",
        "de": "Uhrzeit anpassen"
    },
    "ADJUST_LEVEL": {
        "pl": "Dostosuj poziom:",
        "en": "Adjust level:",
        "de": "Niveau anpassen:"
    },
    "ADJUST_PRICE": {
        "pl": "Dostosuj cenę:",
        "en": "Adjust price:",
        "de": "Preis anpassen:"
    },
    "ADVANCED_CALENDAR": {
        "pl": "Zaawansowany kalendarz",
        "en": "Advanced calendar"
    },
    "ADVANCED_CALENDAR_MARKETING": {
        "pl": "Dodaj pojedyncze zajęcia, kursy wielodniowe lub ustal regularne powtarzanie swoich zajęć - wszystko w jednym miejscu. Ta intuicyjna platforma ułatwia dodawanie i zarządzanie Twoją ofertą treningową",
        "en": "Add single classes, multi-day courses or set up regular repeats of your classes - all in one place. This intuitive platform makes it easy to add and manage your training offerings"
    },
    "CHAT_MARKETING": {
        "pl": "Bezpieczny kontakt przez aplikację. Dzięki możliwościom blokowania, mediacji oraz pomocy klientom, masz pełną kontrolę nad swoją pracą. Bezpieczny kontakt z klientami to dla nas priorytet.",
        "en": "Secure contact through the app. With blocking, mediation and customer support capabilities, you are in full control of your work. Secure contact with customers is a priority for us."
    },
    "ADVANCED_LEVEL": {
        "pl": "zaawansowany",
        "en": "advanced",
        "de": "fortgeschrittene"
    },
    "AFTER_ACCEPT_TIME": {
        "pl": "Po zaakceptowaniu będziesz miał czas do ",
        "en": "After accepting you will have time until ",
        "de": "Nach der Annahme der Reservierung hast du Zeit bis "
    },
    "AFTER_PAYMENT_QR": {
        "pl": "Po akceptacji przez instruktora otrzymasz unikalny kod QR,\n        dzięki któremu instruktor zweryfikuje, że wszystko jest w\nporządku.",
        "en": "After payment, You'll receive unique QR code, which You\nshould show instructor so he can verify that everything is ok.",
        "de": " Nach der Bezahlung bekommst du einen einmaligen QR-Code, \n       mit dem der Trainer verifizieren kann, ob deine Reservierung in Ordnung ist"
    },
    "AFTER_REJECT_YOU_WONT_ACCEPT": {
        "pl": "Po odrzuceniu tego użytkownika nie będziesz mógł go już zaakceptować",
        "en": "After rejection of this user, you won't be able to accept him again.",
        "de": "Nach der Ablehnung des Benutzers wirst du ihn nicht wieder akzeptieren können"
    },
    "AGE": {
        "pl": "wiek",
        "en": "age",
        "de": "Alter"
    },
    "ALLOW_LAST_MINUTE": {
        "pl": "Zezwalaj na rezerwacje last minute",
        "en": "Allow last minute reservations",
        "de": "last minute Reservierungen zulassen"
    },
    "ALL_LIMITS": {
        "pl": "Wszystkie limity",
        "en": "All limits",
        "de": "Alle Einschränkungen"
    },
    "ALL_READY": {
        "pl": "Gotowe!",
        "en": "Ready!",
        "de": "Fertig"
    },
    "ALL_READY_DESC": {
        "pl": "Poniżej odkryjesz pełną listę opcji dostepnych od Ciebie od samego początku!",
        "en": "Below You can discover full list of our features waiting for You! Start earning",
        "de": "Unten findest du die ganze Liste mit Optionen, die dir von Anfang an zur Verfügung stehen"
    },
    "AND_I_ACK": {
        "pl": "Oraz wyrażam zgodę na",
        "en": "And i acknowledge",
        "de": "Und ich stimme zu"
    },
    "ANON": {
        "pl": "anon",
        "en": "anon",
        "de": "anon"
    },
    "ANY": {
        "pl": "dowolny",
        "en": "any",
        "de": "jedes"
    },
    "APPLY": {
        "pl": "Zastosuj",
        "en": "Apply",
        "de": "Anwenden"
    },
    "ARE_YOU_INSTRUCTOR": {
        "pl": "Jesteś trenerem?",
        "en": "Are you an instructor?",
        "de": "Bist du ein Trainer?"
    },
    "ARE_YOU_SURE": {
        "pl": "Na pewno?",
        "en": "Are you sure?",
        "de": "Bist du sicher?"
    },
    "ARE_YOU_SURE_YOU_WANT_TO_CANCEL": {
        "pl": "Na pewno anulować płatność?",
        "en": "Are you sure you want to cancel payment?",
        "de": "Bist du sicher, dass du die Zahlung stornieren möchtest?"
    },
    "ASK_FOR_REFUND": {
        "pl": "Jeżeli w tym momencie anulujesz rezerwację to zwrócimy ci\ntylko część środków.\n        Jeżeli wystąpiły nadzwyczajne okoliczności,\n    możesz poprosić o pełny zwrot swoich pieniędzy poprzez zgłoszenie\ntej rezerwacji na adres ",
        "en": "If you cancel reservation at this point, we will only\nrefund you the part of the funds.\n        If extraordinary circumstances occurred,\n    You can ask for a full refund by reporting this reservation at ",
        "de": "Wenn du die Reservierung jetzt stornierst, wird dir nur ein Teil des Geldes zurückerstattet.\n       Wenn außergewöhnliche Umstände eingetreten sind, kannst du eine vollständige Rückerstattung beantragen, \n       indem du diese Reservierung anmeldest bei "
    },
    "ASSIGNED_TRAININGS_FMT": {
        "pl": "Treningi przydzielone do %s",
        "en": "Trainings assigned to %s",
        "de": "Trainings zugewiesen an %s"
    },
    "ASSIGN_TRAINING": {
        "pl": "Przydziel trening",
        "en": "Assign training",
        "de": "Training zuweisen"
    },
    "AS_YOU_LIKE": {
        "pl": "Tak jak lubisz",
        "en": "Fit for you",
        "de": "Angepasst an dich!"
    },
    "AT_ANY_MOMENT_YOU_CAN_DEACTIVATE": {
        "pl": "W każdym momencie możesz deaktywować swoje konto instruktora",
        "en": "You can deactivate Your instructor account at any time",
        "de": "Du kannst dein Trainerkonto jederzeit deaktivieren"
    },
    "AT_ANY_MOMENT_YOU_MAY_DELETE_USERs_ACCOUNT": {
        "pl": "W każdym momencie możesz usunąć swoje konto",
        "en": "You can delete Your account at any time",
        "de": "Du kannst dein Konto jederzeit löschen"
    },
    "AT_ANY_TIME_YOU_CAN_CHECK": {
        "pl": "W dowolnym momencie możesz sprawdzić status, oraz detale swojej rezerwacji:",
        "en": "At any time you can check the status and details of your reservation:",
        "de": "Jederzeit kannst du den Status und die Details deiner Reservierung überprüfen"
    },
    "AUTH_REQUIRED": {
        "pl": "Wymagana autoryzacja",
        "en": "Authorization required",
        "de": "Autorisierung ist erforderlich"
    },
    "AUTOGEN": {
        "pl": "Wygeneruj automatycznie",
        "en": "Generate automatically",
        "de": "Automatisch generieren"
    },
    "AVAILABLE_SPOTS": {
        "pl": "Wolne miejsca",
        "en": "Available spots",
        "de": "Freie Plätze"
    },
    "AWAITING_DECISION": {
        "pl": "Oczekiwanie, instruktor ma jeszcze ",
        "en": "Awaiting, instructor has still",
        "de": "Warten, der Trainer hat noch"
    },
    "AWAITING_PAYOUT": {
        "pl": "opłacono, wypłacimy instruktorowi środki ",
        "en": "Reservation payed, funds will be transferred to the intructor account ",
        "de": "Reservierung bezahlt, das Geld wird auf das Konto des Trainers überwiesen "
    },
    "BASIC_LEVEL": {
        "pl": "podstawowy",
        "en": "novice",
        "de": "basislevel"
    },
    "BECOME_INSTRUCTOR": {
        "pl": "Zostań instruktorem",
        "en": "Become instructor",
        "de": "Trainer werden"
    },
    "BEGIN_TRAINING": {
        "pl": "Trenuj!",
        "en": "Begin training!",
        "de": "Beginne zu trainieren!"
    },
    "BENEFITS_INFO": {
        "pl": "Jeśli szukasz informacji o benefitach i funkcjonalnościach platformy, skorzystaj z poniższych linków",
        "en": "Informations about Trainer and Customer benefits:",
        "de": "Wenn du nach Informationen über Vorteile und Funktionalitäten der Plattform suchst, benutze die folgenden Links"
    },
    "BOOK_TRAINING": {
        "pl": "Zarezerwuj trening",
        "en": "Find training for you",
        "de": "Finde ein Training für dich"
    },
    "BUY_CARNET": {
        "pl": "Kup karnet",
        "en": "Buy carnet",
        "de": "Kaufe das Abonnement"
    },
    "CALENDAR": {
        "pl": "Kalendarz!",
        "en": "Calendar!",
        "de": "Kalender!"
    },
    "CANCEL": {
        "pl": "Anuluj",
        "en": "Cancel",
        "de": "Abbrechen"
    },
    "CANCELLATION": {
        "pl": "anulowanie",
        "en": "cancellation",
        "de": "Stornierung"
    },
    "CANCELLED": {
        "pl": "Anulowane",
        "en": "Cancelled",
        "de": "Storniert"
    },
    "CANCELLED_RSVS": {
        "pl": "Rezerwacje anulowane",
        "en": "Cancelled reservations",
        "de": "Stornierte Reservierungen"
    },
    "CANCEL_PAYMENT": {
        "pl": "Anuluj płatność",
        "en": "Cancel payment",
        "de": "Zahlung stornieren"
    },
    "CANCEL_RSV": {
        "pl": "Anuluj rezerwację",
        "en": "Cancel reservation",
        "de": "Reservierung stornieren"
    },
    "CANCEL_RSVL": {
        "pl": "Anuluj rezerwację",
        "en": "Cancel reservation",
        "de": "Reservierung stornieren"
    },
    "CANNOT_DELETE_ACCOUNT_YET": {
        "pl": "Nie możesz jeszcze usunąć konta, gdyż wciąż masz niezrealizowane (aktywne) rezerwacje.",
        "en": "You cannot delete Your account yet because you still have active reservations.",
        "de": "Du kannst dein Konto noch nicht löschen, weil du noch aktive Reservierungen hast."
    },
    "CANT_ADD_VACATION_IN_THE_PAST": {
        "pl": "Nie można dodać wolnego w przeszłości",
        "en": "Can't add vacation in the past",
        "de": "Man kann den Urlaub nicht in Vergangenheit hinzufügen"
    },
    "CANT_CONFIRM_RSV": {
        "pl": "Nie można potwierdzić rezerwacji",
        "en": "Can't confirm reservation",
        "de": "Die Reservierung kann nicht bestätigt werden"
    },
    "CANT_OBTAIN_CONTACT": {
        "pl": "Nie można pobrać danych kontaktowych",
        "en": "Failed to obtain contact details",
        "de": "Die Kontaktdaten konnten nicht abgerufen werden"
    },
    "CANT_UNDO": {
        "pl": "Tej operacji nie da się odwrócić",
        "en": "This operation can't be reverted",
        "de": "Dieser Vorgang kann nicht rückgängig gemacht werden."
    },
    "CARD_TO_WHICH_WE_TRANSFER": {
        "pl": " kartę płatniczą na którą będziemy przelewać Ci środki za zrealizowane treningi",
        "en": " card to which we will transfer Your salary",
        "de": " Karte, auf die wir dein Gehalt überweisen werden"
    },
    "CARD_TYPE": {
        "pl": "Typ karty (visa / mastercard)",
        "en": "Type of the card ( Visa / Mastercard )",
        "de": "Typ der Karte (Visa / Mastercard)"
    },
    "CARNET": {
        "pl": "Karnet",
        "en": "Carnet",
        "de": "Abonnement"
    },
    "CARNETS": {
        "pl": "Karnety",
        "en": "Carnets",
        "de": "Abonnements"
    },
    "CARNETS_MARKETING": {
        "pl": "Z Veidly, jako trener możesz oferować swoim klientom wygodne karnety w formie kodu QR w telefonie, których można użyć w ograniczonym czasie lub na określoną ilość wejść. To wyjątkowo przyjazna dla użytkownika opcja, która przyciąga i zatrzymuje klientów.",
        "en": "With Veidly, as a trainer you can offer your clients convenient passes in the form of a QR code on their phone, which can be used for a limited time or for a certain number of entries. This is an extremely user-friendly option that attracts and retains clients."
    },
    "CARNETS_ALLOW_FOR_CL": {
        "pl": "Karnety pozwalają na wejścia na twoje treningi bez konieczności rezerwowania",
        "en": "Carnets allow for your clients to enter trainings without having to make reservation every time",
        "de": "Die Abonnements ermöglichen deinen Kunden die Teilnahme an den Trainings, ohne jedes Mal eine Reservierung vornehmen zu müssen"
    },
    "CARNETS_AVAILABLE_FOR_TRAINING": {
        "pl": "Karnety ważne na ten trening",
        "en": "Carnets available for this training",
        "de": "Verfügbare Abonnements für dieses Training"
    },
    "CARNETS_BOUGHT_FOR_TRAININGS": {
        "pl": "Karnety kupione na twoje treningi",
        "en": "Carnets bought for your trainings",
        "de": "Abonnements gekauft für deine Trainings"
    },
    "CARNETS_INSTEAD_OF_SIGNING_IN": {
        "pl": "Zamiast zapisania się na trening twój klient będzie mógł kupić karnet na określone przez ciebie treningi z datą ważności, lub/i ograniczoną ilością wejść",
        "en": "Instead of signing up for training, your client will beable to buy carnet with limited period of validity, or limited number of entries",
        "de": "Anstatt ein Training zu buchen, kann dein Kunde ein Abonnement für bestimmte Trainings mit begrenzter Gültigkeitsdauer oder begrenzter Zahl der Einträge kaufen"
    },
    "CARNET_ISNT_VALID": {
        "pl": "Karnet nie jest aktualny",
        "en": "Carnet isn't valid",
        "de": "Das Abonnement ist ungültig"
    },
    "CARNET_ISSUED_FOR": {
        "pl": "Karnet wydany dla",
        "en": "Carnet issued for",
        "de": "Das Abonnement ausgestellt für"
    },
    "CARNET_IS_VALID_FOR": {
        "pl": "Karnet jest ważny na następujące treningi",
        "en": "Carnet is valid for following trainings",
        "de": "das Abonnement ist für folgende Trainings gültig"
    },
    "CARNET_MUST_BE_PAID_TO_BE_VALIDATED": {
        "pl": "Karnet musi być opłacony by był uznany przez instruktora",
        "en": "Carnet must be paid to be validated by an instructor",
        "de": "Das Abonnement muss bezahlt werden, um vom Trainer validiert zu werden"
    },
    "CARNET_OWNER_INFO": {
        "pl": "Dane właściela karnetu",
        "en": "Carnet owner info",
        "de": "Daten des Inhabers vom Abonnement"
    },
    "CARNET_PAID": {
        "pl": "Karnet opłacony",
        "en": "Carnet already paid",
        "de": "Das Abonnement ist bezahlt"
    },
    "CARNET_PRICE": {
        "pl": "Cena karnetu",
        "en": "Carnet price",
        "de": "Preis des Abonnements"
    },
    "CARNET_VALID_UNTIL": {
        "pl": "Karnet ważny do",
        "en": "Carnet valid until",
        "de": "Abonnement gültig bis"
    },
    "CERTAINTY": {
        "pl": "Pewność!",
        "en": "Certainty!",
        "de": "Sicherheit"
    },
    "CHANGE_PASSWORD": {
        "pl": "Podaj nowe hasło by je zmienić",
        "en": "Provide us with your new password",
        "de": "Ein neues Passwort eingeben"
    },
    "CHANNEL": {
        "pl": "Kanał",
        "en": "Channel",
        "de": "Channel"
    },
    "CHAT": {
        "pl": "Czat",
        "en": "Chat",
        "de": "Chat"
    },
    "CHATROOM_NAME": {
        "pl": "Nazwa kanału",
        "en": "Channel name",
        "de": "Channel name"
    },
    "CHEAPEST_SOLUTION": {
        "pl": "Najtańsza oferta!",
        "en": "Cheapest portal in Internet!",
        "de": "Das günstigste Angebot"
    },
    "CHECK": {
        "pl": "Sprawdź",
        "en": "Check",
        "de": "Überprüfe"
    },
    "CHECK_YOUR_MAIL": {
        "pl": "Sprawdź swoją skrzynkę pocztową @",
        "en": "Check your inbox @",
        "de": "Prüfe deine E-Mail-Box @"
    },
    "CLIENT": {
        "pl": "Klient",
        "en": "Client",
        "de": "Kunde"
    },
    "CLIENT_REMARKS": {
        "pl": "Uwagi widoczne dla klientów",
        "en": "Notes visible for customers",
        "de": "Anmerkungen sichtbar für Kunden"
    },
    "CLOSE": {
        "pl": "Zamknij",
        "en": "Close",
        "de": "Schliessen"
    },
    "CLOSE_AND_TRY_AGAIN": {
        "pl": "Zamknij i spróbuj ponownie",
        "en": "Close and try again",
        "de": "Schliesse und versuch erneut"
    },
    "CLOSE_EDITOR": {
        "pl": "Zamknij edytor",
        "en": "Close editor",
        "de": "Editor schliessen"
    },
    "CODE": {
        "pl": "Kod",
        "en": "Code",
        "de": "Code"
    },
    "CODE_ON_BOOKING": {
        "pl": "Ten kod użytkownik będzie mógł wprowadzić przy dodawaniu rezerwacji",
        "en": "User will be able to enter this code during booking",
        "de": "Der Benutzer kann diesen Code bei der Buchung eingeben"
    },
    "COMPANY": {
        "pl": "Firma",
        "en": "Company",
        "de": "Firma"
    },
    "COMPLETED": {
        "pl": "Zrealizowane",
        "en": "Completed",
        "de": "Abgeschlossen"
    },
    "COMPLETELY_DELETE_ACC": {
        "pl": "Całkowicie usuń swoje konto w systemie",
        "en": "Completely delete Your account in the system",
        "de": "Lösche dein Konto vollständig im System"
    },
    "COMPLIES_WITH_PCI_DSS": {
        "pl": "spełniającej wymagania międzynarodowego standardu PCI DSS",
        "en": "who complies with the international PCI DSS standard",
        "de": "die dem internationalen PCI DSS-Standard entsprechen"
    },
    "CONFIGURATION": {
        "pl": "Konfiguracja",
        "en": "Configuration",
        "de": "Konfiguration"
    },
    "CONFIGURE_INSTR_ACC": {
        "pl": "Konfiguracja konta instruktora",
        "en": "Configuring an instructor account",
        "de": "Konfiguration des Trainerkontos"
    },
    "CONFIRM": {
        "pl": "Na pewno?",
        "en": "Are you sure?",
        "de": "Bist du sicher?"
    },
    "CONFIRMED_RSVS": {
        "pl": "Potwierdzone rezerwacje",
        "en": "Confirmed reservations",
        "de": "Bestätigte Reservierungen"
    },
    "CONFIRMING_RSV": {
        "pl": "Potwierdzanie rezerwacji",
        "en": "Confirming reservation",
        "de": "die Reservierung wird bestätigt"
    },
    "CONFIRM_ENTRANCE": {
        "pl": "Potwierdź wejście",
        "en": "Confirm the entrance",
        "de": "Bestätige den Eintritt"
    },
    "CONFIRM_NEW_PASSWORD": {
        "pl": "Potwierdź nowe hasło",
        "en": "Confirm new password",
        "de": "Bestätige das neue Passwort"
    },
    "CONFIRM_NOT_REQUIRED": {
        "pl": "Ta rezerwacja nie wymaga potwierdzenia od instruktora",
        "en": "This reservation doesn't require confirmation from the instructor",
        "de": "Diese Reservierung erfordert keine Bestätigung des Trainers"
    },
    "CONFIRM_PASSWORD": {
        "en": "Confirm password",
        "pl": "Potwierdź hasło",
        "de": "Passwort bestätigen"
    },
    "CONFIRM_REQUIRED": {
        "pl": "Potwierdzenie od instruktora jest wymagane",
        "en": "Confirmation from the instructor is required",
        "de": "Eine Bestätigung des Trainers ist erforderlich"
    },
    "CONFIRM_RSV_MANUALLY": {
        "pl": "Potwierdzaj ręcznie każdą rezerwację",
        "en": "Confirm every reservation manually",
        "de": "Jede Buchung manuell bestätigen"
    },
    "CONFLICT": {
        "pl": "Konflikt",
        "en": "Conflict",
        "de": "Konflikt"
    },
    "CONTACT_DATA": {
        "pl": "Dane kontaktowe",
        "en": "Contact data",
        "de": "Kontaktdaten"
    },
    "CONTACT_DATA_DISCLAIMER": {
        "pl": "Dane kontaktowe przekażemy instruktorowi po to by się mógł z Tobą skontaktować w sprawie nieścisłości lub pytań.",
        "en": "Contact details will be handed over to the instructor so that they can contact You in case of troubles or questions.",
        "de": "Deine Kontaktdaten werden wir an den Trainer weitergeben, damit er dich im Fall von Fragen oder Unstimmigkeiten kontaktieren kann"
    },
    "CONTACT_INSTRUCTOR": {
        "pl": "Skontaktuj się z instruktorem",
        "en": "Contact instructor",
        "de": "Kontaktiere den Trainer"
    },
    "CONTACT_PASS": {
        "pl": "Podane tu dane kontaktowe przekażemy twoim klientom którzyzrobili rezerwacje",
        "en": "We will only pass your contact data to the customers who have made the reservations",
        "de": "Deine Kontaktdaten werden nur an die Kunden weitergegeben, die das Training reserviert haben"
    },
    "CONTACT_PASS_2": {
        "pl": "Lub jeżeli to ty składasz zamówienie przekażemy je instruktorowi by się z tobą mógł skontaktować",
        "en": "Or if you placed the order, we will pass it on only to the instructor so that they can contact you",
        "de": "Wenn du die Bestellung aufgibst, leiten wir sie an den Trainer weiter, damit er dich kontaktieren kann"
    },
    "CONTACT_PASS_CUST": {
        "pl": "Podane tu dane kontaktowe możesz przekazać instruktorowi by się z tobą mógł skontaktować",
        "en": "You can pass contact data to the instructor so that they can contact you",
        "de": "Du kannst die Kontaktdaten dem Trainer weitergeben, damit er dich kontaktieren kann"
    },
    "CONTACT_SUPPORT": {
        "pl": "Gdyby to nie wystarczyło, skontaktuj się z supportem:",
        "en": "If that doesnt suffice, contact support:",
    },
    "CONTENT": {
        "pl": "Zawartość",
        "en": "Content",
        "de": "Inhalt"
    },
    "CONTINUE": {
        "pl": "kontynuuj",
        "en": "continue",
    },
    "COULDNT_FIND_TRAININGS": {
        "pl": "Nie znaleźliśmy żadnych treningów dla podanych parametrów :(",
        "en": "Failed to find trainings for given parameters :(",
        "de": "Es konnten keine Trainings für die angegebenen Parameter gefunden werden "
    },
    "COULDNT_GET_DATA": {
        "pl": "Nie udało się pobrać danych z api",
        "en": "Failed to download data from api",
        "de": "Das Abrufen der Daten von API fehlgeschlagen"
    },
    "COULDNT_LOGIN": {
        "pl": "Nieudana próba logowania",
        "en": "Couldn't login",
        "de": "Anmeldung fehlgeschlagen"
    },
    "COULDNT_REGISTER": {
        "pl": "Wystąpił błąd podczas tworzenia konta",
        "en": "Couldn't create account",
        "de": "Das Konto konnte nicht erstellt werden"
    },
    "COULDNT_RESEND_EMAIL": {
        "pl": "Nie udało się ponownie wysłać emaila",
        "en": "We Couldn't resend you confirmation email",
        "de": "Wir konnten dir keine neue Bestätigungs-E-Mail senden"
    },
    "COULDNT_RESET_PASSWOD": {
        "pl": "Nie udało się zresetować hasła: ",
        "en": "Couldnt reset password: ",
        "de": "Das Passwort konnte nicht zurückgesetzt werden: "
    },
    "COULDNT_SAVE_OCC": {
        "pl": "Nie udało się zapisać występowania",
        "en": "Failed to save occurrence",
        "de": "Auftreten konnte nicht gespeichert werden"
    },
    "COULDNT_SEND_EMAIL": {
        "pl": "Nie mogliśmy wysłać emaila:",
        "en": "We couldnt send you an email:",
        "de": "Wir konnten dir keine E-Mail schicken:"
    },
    "COULDNT_VERIFY_TID": {
        "pl": "Nie można było zweryfikować id treningu",
        "en": "Couldn't verify training ID",
        "de": "Trainings-ID konnte nicht verifiziert werden"
    },
    "CREATE_ACCOUNT": {
        "pl": "Załóż konto",
        "en": "Create account",
        "de": "Konto erstellen"
    },
    "CREATE_ACCOUNT_TO_CONTACT_INSTR": {
        "pl": "Stwórz konto w systemie jeżeli chcesz się skontaktować z instruktorem",
        "en": "Create account to contact instructor",
        "de": "Erstelle ein Konto, um den Trainer zu kontaktieren"
    },
    "CREATE_QR": {
        "pl": "Stwórz kod QR",
        "en": "Create QR code",
        "de": "QR-Code erstellen"
    },
    "CREATE_QR_EXTENDED": {
        "pl": "Wygeneruj kod QR który instruktor może chcieć zweryfikować przed rozpoczęciem zajęć.",
        "en": "Generate QR code that the instructor may want to verify before starting the training.",
        "de": "Generiere einen QR-Code, den der Trainer vor dem Beginn des Trainings überprüfen kann."
    },
    "CREATE_YOUR_OFFER": {
        "pl": "Stwórz ofertę",
        "en": "Create an offer",
        "de": "Erstelle ein Angebot"
    },
    "CREAT_ACC_AGREEMENT": {
        "pl": "Zgadzam się na",
        "en": "By creating account I accept",
        "de": "Mit Erstellung des Kontos akzeptiere ich"
    },
    "CREDIT_CARD_OR_DEBIT": {
        "pl": "Karta kredytowa lub debetowa",
        "en": "Credit card or debit",
        "de": "Kreditkarte oder EC-Karte"
    },
    "CURRENCY": {
        "pl": "Waluta",
        "en": "Currency",
        "de": "Währung"
    },
    "CURRENTLY_FOR_PAYOUTS": {
        "pl": "Obecnie do wypłacania ci środków używamy karty:",
        "en": "Currently for payouts we are using:",
        "de": "Derzeit verwenden wir für Auszahlungen:"
    },
    "CURRENT_DECISION": {
        "pl": "Aktualna decyzja:",
        "en": "Current Decision:",
        "de": "Aktuelle Entscheidung"
    },
    "CUSTOMER_BENEFITS": {
        "pl": "Korzyści dla klienta",
        "en": "Customer benefits",
        "de": "Vorteile für Kunden"
    },
    "CUSTOMER_HAPPY": {
        "pl": "Klient zadowolony?",
        "en": "Customer happy?",
        "de": "Ist der Kunde zufrieden?"
    },
    "DAILY": {
        "pl": "Codziennie",
        "en": "Daily",
        "de": "Täglich"
    },
    "DATA_THAT_WE_KNOW_AND_STORE": {
        "pl": "Dane które my znamy i zapiszemy u nas w systemie to:",
        "en": "Data that we know and store in our system is:",
        "de": "Daten, die wir kennen und in unserem System speichern, sind:"
    },
    "DAY": {
        "pl": "dzień",
        "en": "day",
        "de": "Tag"
    },
    "DAYS": {
        "pl": "dni",
        "en": "days",
        "de": "Tage"
    },
    "DCS": {
        "pl": "Kody rabatowe",
        "en": "Discount codes",
        "de": "Rabattcodes"
    },
    "DC_NAME": {
        "pl": "Kod rabatowy",
        "en": "Discount code",
        "de": "Rabattcode"
    },
    "DCS_MARKETING": {
        "pl": "Możesz łatwo i szybko stworzyć swoją własną promocję i zwiększyć sprzedaż",
        "en": "You can easily and quickly create your own promotion and increase sales"
    },
    "DEACTIVATE": {
        "pl": "Deaktywuj",
        "en": "Deactivate",
        "de": "Deaktiviere"
    },
    "DEACTIVATE_INSTR_WARN": {
        "pl": "Spowoduje to nieodwracalne usunięcie twoich treningów, rezerwacj i danych płatności",
        "en": "This will permanently delete Your trainings, bookings and payment details",
        "de": "Dadurch werden deine Trainings, Buchungen und Zahlungsdaten dauerhaft gelöscht"
    },
    "DEACTIVATING": {
        "pl": "Deaktywowanie",
        "en": "Deactivating",
        "de": "Deaktivierung"
    },
    "DEACTIVATION": {
        "pl": "Deaktywacja",
        "en": "Deactivation",
        "de": "Deaktivierung"
    },
    "DEACTIVATION_EFFECTS": {
        "pl": "Spowoduje to, że nie będziesz listowany w wynikach searchu, oraz klienci nie będą mogli się do ciebie zapisać na rezerwacje.",
        "en": "This will prevent you from being listed in the search results and customers will not be able to sign up for reservations.",
        "de": "Dadurch wird verhindert, dass du in den Suchergebnissen aufgeführt wirst und die Kunden werden sich nicht für deine Trainings anmelden können."
    },
    "DEACTIVATION_WARN": {
        "pl": "Deaktywacja konta spowoduje, że nie będziesz listowany w wynikach searchu, oraz klienci nie będą mogli się do ciebie zapisać na rezerwacje.",
        "en": "Deactivating Your account will prevent you from being listed in the search results and customers wont be able to sign up for reservations.",
        "de": "Wenn du dein Konto deaktivierst, wirst du nicht mehr in den Suchergebnissen aufgelistet und Kunden werden sich nicht mehr für deine Trainings anmelden können."
    },
    "DECUCT_WARNING": {
        "pl": "Jak odrzucisz tą rezerwację w tym momencie to przy\nnastępnej wypłacie potrącimy ci koszta tego zwrotu.\n        Na pewno nie chcesz dogadać się z użytkownikiem by ustalić np.\ninny termin?",
        "en": "If you reject reservation at this point, we will deduct\nyou the costs of this return during next payout.\n        Are you sure, that You don't want to get along with the user to\ndetermine, for example, another date?",
        "de": "Wenn du die Reservierung jetzt ablehnst, werden die Kosten der Rückerstattung bei der nächsten Auszahlung von deinem Konto abgezogen. \n       Bist du sicher, dass du dich mit dem Kunden nicht auf einen anderen Termin einigen willst?"
    },
    "DELETE": {
        "pl": "Usuń",
        "en": "Delete",
        "de": "Löschen"
    },
    "DELETE_ACCOUNT": {
        "pl": "Usuń konto",
        "en": "Delete account",
        "de": "Konto löschen"
    },
    "DELETE_ACC_WARN": {
        "pl": "W momencie w którym usuniesz u nas konto stracisz dostęp\ndo historii swoich rezerwacji.\n        Do swoich aktywnych rezerwacji wciąż będziesz się mógł odnieść\nprzez np. linki które ci wysłaliśmy w mailu",
        "en": "The moment You delete Your account You lose access to Your\nreservation history.\n        You will still be able to access Your active reservations\nthrough the links we sent You in the mail",
        "de": "In dem Moment, in dem du dein Konto löschst, verlierst du den Zugang zu deinem\nReservierungsverlauf. Du kannst weiterhin auf deine aktiven Reservierungen zugreifen\nüber die Links, die wir dir per Mail geschickt haben."
    },
    "DELETE_CONFIRM": {
        "pl": "Na pewno usunąć %s",
        "en": "Are you sure you want to remove %s",
        "de": "Bist du sicher, dass du %s entfernen willst"
    },
    "DELETE_FMT": {
        "pl": "Usuń %s",
        "en": "Delete %s",
        "de": "Lösche %s"
    },
    "DELETE_TRAINING": {
        "pl": "Usuń trening",
        "en": "Delete training",
        "de": "Training löschen"
    },
    "DELETE_USER_ACC": {
        "pl": "Usuń swoje konto użytkownika",
        "en": "Delete Your user account",
        "de": "Lösche dein Benutzerkonto"
    },
    "DELETING_ACCOUNT": {
        "pl": "Usuwanie konta",
        "en": "Deleting account",
        "de": "Löschung des Kontos"
    },
    "DELETING_CARD": {
        "pl": "Usuwanie karty",
        "en": "Deleting card",
        "de": "Löschung der Karte"
    },
    "DELETING_IMG": {
        "pl": "Usuwanie zdjęcia",
        "en": "Deleting image",
        "de": "Foto löschen"
    },
    "DESCRIBE_YOUR_ISSUE": {
        "pl": "Opisz problem",
        "en": "Describe your problem",
        "de": "Beschreibe das Problem"
    },
    "DESCRIPTION": {
        "pl": "Opis",
        "en": "Description",
        "de": "Beschreibung"
    },
    "DETAILS": {
        "pl": "Detale",
        "en": "Details",
        "de": "Details"
    },
    "DIDNT_FIND_DC": {
        "pl": "Nie znaleźliśmy podanego kodu",
        "en": "Couldn't find given code",
        "de": "Wir konnten den angegebenen Code nicht finden"
    },
    "DISCIPLINE": {
        "pl": "Dyscyplina",
        "en": "Discipline",
        "de": "Disziplin"
    },
    "DISCOUNT": {
        "pl": "Zniżka",
        "en": "Discount",
        "de": "Ermäßigung"
    },
    "DISCOVER_OUR_OFFER": {
        "pl": "Trenerze, trenerko! - odkryj naszą ofertę!",
        "en": "Dear trainer, discover our offer!",
        "de": "Lieber Trainer, liebe Trainerin - entdecke unser Angebot!"
    },
    "DISPLAY_NAME": {
        "pl": "Twoja nazwa",
        "en": "Your name",
        "de": "Your name"
    },
    "DOCS": {
        "pl": "Pomoc & regulaminy",
        "en": "Help & Terms of use",
    },
    "DOCUMENTATION": {
        "pl": "Dokumentacja serwisu",
        "en": "Documentation",
        "de": "Die Dokumentation"
    },
    "DONE": {
        "pl": "Gotowe",
        "en": "Done",
        "de": "Fertig"
    },
    "DONE_PAYOUT": {
        "pl": "Wypłacono środki instruktorowi",
        "en": "Money transferred to the instructor",
        "de": "Geldüberweisung an den Trainer"
    },
    "DONT_REPEAT": {
        "pl": "Nie powtarzaj",
        "en": "Don't repeat",
        "de": "Nicht wiederholen"
    },
    "DO_YOU_HAVE_DC": {
        "pl": "Masz kod rabatowy?",
        "en": "Do you have discount code?",
        "de": "Hast du einen Rabattcode?"
    },
    "DURATION": {
        "pl": "Czas trwania",
        "en": "Duration",
        "de": "Dauer"
    },
    "EACH": {
        "pl": "co ",
        "en": "each ",
        "de": "jede/n "
    },
    "EDIT": {
        "pl": "Edytuj",
        "en": "Edit",
        "de": "Bearbeite"
    },
    "EDITOR_COMPAT_WARN": {
        "pl": "Ten edytor nie jest przystosowany do aktualnego formatu występowania",
        "en": "This editor is not suited for your occurrence format",
        "de": "Dieser Editor ist nicht für das Format des Auftretens geeignet"
    },
    "EDIT_FMT": {
        "pl": "Edytuj %s",
        "en": "Edit %s",
        "de": "Bearbeite %s"
    },
    "EDIT_SECTION": {
        "pl": "Edytuj sekcję",
        "en": "Edit section",
        "de": "Sektion bearbeiten"
    },
    "EMAIL": {
        "pl": "Email",
        "en": "Email",
        "de": "E-Mail"
    },
    "EMPHASIZE_SESS": {
        "pl": "Wyróżnij sesję kolorem",
        "en": "Emphasize session with color",
        "de": "Die Sektion bunt markieren"
    },
    "END": {
        "pl": "Koniec",
        "en": "End",
        "de": "Ende"
    },
    "END_DATE": {
        "pl": "Data zakończenia",
        "en": "End date",
        "de": "Enddatum"
    },
    "END_HOUR": {
        "pl": "Godzina zakończenia",
        "en": "End hour",
        "de": "Die Endzeit"
    },
    "END_PREVIEW": {
        "pl": "Zakończ podgląd",
        "en": "End preview",
        "de": "Vorschau schliessen"
    },
    "ENTER_CONTACT": {
        "pl": "Podaj dane do kontaktu",
        "en": "Enter contact details",
        "de": "Gib deine Kontaktdaten ein"
    },
    "ENTER_DC": {
        "pl": "Wprowadź kod rabatowy otrzymany od instruktora",
        "en": "Enter discount code received from the instructor",
        "de": "Gib den vom Trainer erhaltenen Rabattcode ein"
    },
    "ENTER_MORE_ACCURATE_ADDR": {
        "pl": "Podaj dokładniejszy adres",
        "en": "Enter more accurate address",
        "de": "genauere Adresse eingeben"
    },
    "ENTER_NEW_ADDRESS": {
        "pl": "podaj nowy adres",
        "en": "enter new address",
        "de": "neue Adresse eingeben"
    },
    "ENTER_YOUR_NAME": {
        "pl": "Podaj swoje imię i nazwisko",
        "en": "Enter your name and surname",
        "de": "Namen und Nachnamen eingeben"
    },
    "ERROR": {
        "pl": "Błąd",
        "en": "Error",
        "de": "Fehler"
    },
    "EVERYTHING_IS_READY": {
        "pl": "Wygląda na to, że masz wszystko poustawiane jak należy!",
        "en": "It looks like you have set up Your account perfectly!",
        "de": "Es sieht so aus, als hättest du dein Konto perfekt eingerichtet!"
    },
    "EVERY_WEEK": {
        "pl": ", co tydzień",
        "en": ", every week",
        "de": ", jede Woche"
    },
    "EXISTING_RSV_WARN": {
        "pl": "Istniejące rezerwacje nie zostaną usunięte i pozostaną widoczne w harmonogramie",
        "en": "Existing reservations won't be removed and will remain visible in the schedule",
        "de": "Bestehende Reservierungen werden nicht entfernt und bleiben im Zeitplan sichtbar"
    },
    "EXPLORE": {
        "pl": "Odkrywaj!",
        "en": "Explore!",
        "de": "Entdecke!"
    },
    "FAILED_REASON": {
        "pl": "Nie udało się, powód: ",
        "en": "Failed, reason: ",
        "de": "Fehlgeschlagen, Grund: "
    },
    "FAILED_TO_CONFIRM": {
        "pl": "Nie udało się potwierdzić",
        "en": "Failed to confirm",
        "de": "Bestätigung fehlgeschlagen"
    },
    "FAILED_TO_CONFIRM_CARNET": {
        "pl": "Nie można potwierdzić karnetu",
        "en": "Failed to confirm the carnet",
        "de": "Die Bestätigung des Abonnements ist fehlgeschlagen"
    },
    "FAILED_TO_DOWNLOAD_PRICING": {
        "pl": "Nie udało się pobrać danych płatności",
        "en": "Failed to download pricing info",
        "de": "Herunterladen der Zahlungsdaten ist fehlgeschlagen "
    },
    "FAILED_TO_FETCH_CARNETS": {
        "pl": "Nie udało się pobrać karnetów",
        "en": "Failed to fetch carnets",
        "de": "Das Abrufen der Abonnements ist fehlgeschlagen"
    },
    "FAILED_TO_FETCH_LOCATION_DATA": {
        "pl": "Nie udało się pobrać danych lokalizacji",
        "en": "Failed to fetch localization data",
        "de": "Lokalisierungsdaten konnten nicht abgerufen werden"
    },
    "FAILED_TO_LOGIN_VIA_OAUTH": {
        "pl": "Nie udało się zalogować poprzez oauth: ",
        "en": "Failed to login via oauth: ",
        "de": "Anmeldung über oauth ist fehlgeschlagen: "
    },
    "FETCHING_PAYMENTS": {
        "pl": "pobieranie płatności",
        "en": "fetching payments",
        "de": "Abrufen der Zahlung"
    },
    "FIELD_IS_REQUIRED": {
        "pl": "pole jest wymagane",
        "en": "field is required",
        "de": "dieses Feld ist erforderlich"
    },
    "FILL_YOUR_NAME": {
        "pl": "Twoje imię...",
        "en": "Your name...",
        "de": "Dein Name"
    },
    "FILTERING": {
        "pl": "Filtrowanie",
        "en": "Filtering",
        "de": "Filtern"
    },
    "FIND_ALL_ACTIVITIES": {
        "pl": "Dzięki naszej wyszukiwarce z łatwoscią znajdziesz usługi w\nokolicy.\n        Bez znaczenia czy to boks, taniec czy kurs wspinaczki w Tatrach!\n        Wiele opcji filtrowania i sortowania tylko pomoże i przyspieszy\nproces!",
        "en": "Thanks to our advanced search engine You can easily find\nall activities in Your neighbourhood.\n        It doesn't matter if it's boxing, dance or hiking in Alps!\n        Search, filter and sort trainings right now!",
        "de": "Dank unserer Suchmaschine kannst du problemlos die Aktivitäten in deiner Nähe finden. \n        Egal ob Boxen, Tanz oder Kletterkurs in den Alpen! Suche, filter und sortiere Sie jetzt Trainings!"
    },
    "FINISHED": {
        "pl": "Zakończone",
        "en": "Finished",
        "de": "Abgeschlossen"
    },
    "FIRST_NAME": {
        "pl": "Imię",
        "en": "First name",
        "de": "First name"
    },
    "FIRST_STEP_ACCOUNT": {
        "pl": "Po pierwsze konto!",
        "en": "First step: Account!",
        "de": "Der erste Schritt: dein Konto"
    },
    "FORGOT_PASSWORD": {
        "pl": "Zapomniałeś hasła?",
        "en": "Forgot password?",
        "de": "Passwort vergessen?"
    },
    "FOR_CUSTOMER": {
        "pl": "Dla klienta",
        "en": "For customer",
        "de": "Für den Kunden"
    },
    "FOR_INSTRUCTOR": {
        "pl": "Jestem trenerem",
        "en": "I'm instructor",
        "de": "Für den Trainer"
    },
    "FOR_TRAINEE": {
        "pl": "Chcę trenować",
        "en": "I want to train",
        "de": "Für den Trainierenden"
    },
    "FROM": {
        "pl": "Od",
        "en": "From",
        "de": "Von"
    },
    "GEAR": {
        "pl": "Sprzęt",
        "en": "Gear",
        "de": "Die Ausrüstung"
    },
    "GEAR_WHICH_CUSTOMER_MUST_HAVE": {
        "pl": "Sprzęt który klient musi posiadać",
        "en": "Gear that customer must have",
        "de": "Ausrüstung, die der Kunde haben muss"
    },
    "GEAR_WHICH_YOU_HAVE": {
        "pl": "Sprzęt który ty posiadasz",
        "en": "Gear that you have",
        "de": "Ausrüstung, die du hast"
    },
    "GEAR_WHICH_YOU_RECOMMEND_TO_CUSTOMER": {
        "pl": "Sprzęt który rekomendujesz klientowi",
        "en": "Gear that you recommend to the customer",
        "de": "Ausrüstung, die du deinen Kunden empfiehlst"
    },
    "GIVE_INFO_TO_CLIENTS": {
        "pl": "Przekaż klientom jakieś informacje o sobie",
        "en": "Give some information about yourself to your clients",
        "de": "Gib ein paar Informationen über dich für deine Kunden"
    },
    "GOTO_CARNET": {
        "pl": "Przejdź do karnetu",
        "en": "Goto carnet",
        "de": "Gehe zum Abonnement"
    },
    "GOTO_CONFIG": {
        "pl": "Przejdź do konfiguracji",
        "en": "Go to configuration",
        "de": "Gehe zur Konfiguration"
    },
    "GOTO_PAYMENT": {
        "pl": "Przejdź do zakładki płatności",
        "en": "Go to the payment tab",
        "de": "Gehe zu der Registerkarte - Zahlung"
    },
    "GOTO_RSV": {
        "pl": "Przejdź do rezerwacji",
        "en": "Go to reservation",
        "de": "Gehe zu den Reservierungen"
    },
    "GRID": {
        "pl": "Siatka",
        "en": "Grid",
        "de": "Raster"
    },
    "GROUP_TRAININGS": {
        "pl": "Treningi grupowe",
        "en": "Group trainings",
        "de": "Gruppenschulung"
    },
    "GROUP_TRAININGS_DESC": {
        "pl": "Dzięki Veidly, jako trener możesz cieszyć się korzyściami, takimi jak liczenie użytkowników Twoich zajęć, niższe koszty wynajmu przestrzeni oraz zwiększenie atrakcyjności Twojej oferty",
        "en": "With Veidly, as a trainer you can enjoy benefits such as counting users of your classes, lower space rental costs, and increasing the appeal of your offerings",
    },
    "GRP_NAME": {
        "pl": "Limit",
        "en": "Limit",
        "de": "Limit"
    },
    "HELP": {
        "pl": "Pomoc",
        "en": "Help",
        "de": "Hilfe"
    },
    "HELP_OTHERS": {
        "pl": "Wspieraj innych!",
        "en": "Help others!",
        "de": "Helfe den anderen!"
    },
    "HERE": {
        "pl": "Tutaj",
        "en": "Here",
        "de": "Hier"
    },
    "HIDE_MULTIDAY": {
        "pl": "Ukryj wielodniowe kursy",
        "en": "Hide multi-day courses",
        "de": "Mehrtägige Kurse ausblenden"
    },
    "HOLIDAYS": {
        "pl": "Wakacje",
        "en": "Holidays",
        "de": "Urlaub"
    },
    "HOLIDAYS_DESC": {
        "pl": "Chcesz odpocząć i zregenerować siły, albo potrzebujesz skupić się na przygotowaniach do ważnego wyzwania? Z Veidly to proste! Wystarczy kilka kliknięć, aby zawiesić lub odwołać swoje zajęcia. Bez zbędnych formalności, szybko i wygodnie.",
        "en": "Want to rest and recuperate, or need to focus on preparing for an important challenge? With Veidly it's easy! With just a few clicks, you can suspend or cancel your activities. No paperwork, fast and convenient.",
    },
    "HOLIDAYS_WHENEVER_YOU_LIKE_IT": {
        "pl": "Urlop kiedy chcesz!",
        "en": "Holidays whenever you like it!",
        "de": "Urlaub wann du willst!"
    },
    "HOURS_TO_MAKE_DECISION": {
        "pl": "godzin na decyzję",
        "en": "hours to make decision",
        "de": "Stunden, um Entscheidung zu treffen"
    },
    "HOW_CAN_WE_HELP": {
        "pl": "W czym możemy pomóc?",
        "en": "How can we help?",
        "de": "Was können wir für dich tun?"
    },
    "HOW_DID_YOU_LIKE_TRAINING": {
        "pl": "Jak ci się podobał trening?",
        "en": "How did you like this training?",
        "de": "Wie hat dir das Training gefallen?"
    },
    "HOW_IT_WORKS": {
        "pl": "Jak to działa?",
        "en": "How it works?",
        "de": "Wie funktioniert das?"
    },
    "HOW_VEIDLY_WORKS": {
        "pl": "Jak działa Veidly?",
        "en": "How does Veidly work?",
        "de": "Wie funktioniert Veidly"
    },
    "IF_DIDNT_FIND_EMAIL": {
        "pl": "Jeżeli nie znalazłeś naszego maila to sprawdź folder SPAM",
        "en": "If you didn't find our email, then check SPAM",
        "de": "Wenn du unsere E-Mail nicht gefunden hast, prüfe bitte die Spam-Mails oder"
    },
    "IF_INSTR_AGRESS": {
        "pl": "jak instruktor wyrazi zgodę na rezerwacje",
        "en": "only if instructor agrees on reservation",
        "de": "nur wenn der Trainer mit der Buchung einverstanden ist"
    },
    "IF_U_DONT_CONFIRM_BEFORE_FMT": {
        "pl": "Jak nie potwierdzisz rezerwacji do %s zostanie ona odrzucona automatycznie",
        "en": "If you don't confirm your reservation until %s It will be rejected automatically",
        "de": "Wenn du die Reservierung bis zum %s nicht bestätigst, wird sie automatisch storniert"
    },
    "IF_YOU_DONT_SHARE_CONTACT": {
        "pl": "Jak nie pozwolisz na udostępnianie swoich danych\nkontaktowych,\n        użytkownicy będą mogli się z tobą skontaktować dopiero po\nzrobieniu rezerwacji",
        "en": "If you do not allow Your contact details to be shared,\n        users will only be able to contact you after booking",
        "de": " Wenn du der Weitergabe deiner Kontaktdaten nicht zustimmst, werden die Nutzer dich erst nach der Reservierung kontaktieren können. "
    },
    "IF_YOU_REMOVE_CARD": {
        "pl": "Do czasu dodania innej nie będziesz mógł otrzymać wypłaty",
        "en": "You wont be able to receive a payments until another one is added",
        "de": "Du wirst keine Zahlung erhalten können, bis eine weitere hinzugefügt wird"
    },
    "IF_YOU_TRY_TO_RECREATE": {
        "pl": "Jeżeli ponownie zechcesz otworzyć u nas konto instruktora, to możemy wymagać od ciebie kontaktu mailowego z supportem by uniknąć nadużyć",
        "en": "In case you want to reopen instructor's account again, we may require you to contact the support by e-mail to avoid abuses",
        "de": "Falls du das Trainerkonto wieder öffnen möchtest, kann es erforderlich sein, den Support per E-Mail zu kontaktieren, um Missbräuche zu vermeiden"
    },
    "IF_YOU_WANT_TO_PERFORM_THIS_ACTION": {
        "pl": "By móc wykonać to żądanie",
        "en": "If you want to perform this action",
        "de": "Um diesen Befehl umzusetzen "
    },
    "IMAGE_IS_TOO_BIG": {
        "pl": "Zdjęcie jest zbyt duże (max 5MB)",
        "en": "Photo is too large (max 5MB)",
        "de": "Das Bild ist zu groß (max 5MB)"
    },
    "INCOMING": {
        "pl": "Nachodzące",
        "en": "Incoming",
        "de": "Eingehend"
    },
    "INCOMING_VACATION": {
        "pl": "Nadchodzące wakacje",
        "en": "Incoming vacations",
        "de": "Kommender Urlaub"
    },
    "INSTANT_REFUND": {
        "pl": "Natychmiastowy zwrot pieniędzy",
        "en": "instant refund",
        "de": "sofortige Rückerstattung"
    },
    "INSTRUCTOR": {
        "pl": "Trener",
        "en": "Instructor",
        "de": "Trainer"
    },
    "INSTRUCTORS_PROFILE": {
        "pl": "Profil instruktora",
        "en": "Instructor's profile",
        "de": "Trainerprofil"
    },
    "INSTRUCTOR_CANCEL_RSV": {
        "pl": "Anuluj rezerwację jako instruktor",
        "en": "Cancel reservation as instructor",
        "de": "Storniere die Reservierung als Trainer"
    },
    "INSTRUCTOR_CARNETS": {
        "pl": "Karnety instruktora",
        "en": "Instructor's carnets",
        "de": "Abonnements des Trainers"
    },
    "INSTRUCTOR_CONTACT_INFO": {
        "pl": "Dane kontaktowe instruktora",
        "en": "Instructor's contact info",
        "de": "Die Kontaktdaten des Trainers"
    },
    "INSTRUCTOR_DECISION": {
        "pl": "Decyzja instruktora",
        "en": "Instructor's decision",
        "de": "Entscheidung des Trainers"
    },
    "INSTRUCTOR_GEAR": {
        "pl": "Sprzęt instruktora",
        "en": "Instructor's gear",
        "de": "Ausrüstung des Trainers"
    },
    "INSTRUCTOR_IS_NOT_PROVIDING": {
        "pl": "Instruktor nie udostępnia danych kontaktowych przed wykonaniem rezerwacji",
        "en": "Instructor is not providing contact details before making reservation",
        "de": "Der Trainer gibt keine Kontaktdaten vor der Reservierung des Trainings an"
    },
    "INSTRUCTOR_SCHEDULE": {
        "pl": "Harmonogram instruktora",
        "en": "Instructor schedule",
        "de": "Der Zeitplan des Trainers"
    },
    "INSTRUCTOR_SINCE": {
        "pl": "Trener od",
        "en": "Instructor since",
        "de": "Trainer seit"
    },
    "INSTR_ALSO_OFFERS_CARNETS": {
        "pl": "Trener oferuje także karnety na ten trening:",
        "en": "Instructor also offers carnets for this training:",
        "de": "Auch der Trainer bietet Abonnements für das Training an"
    },
    "INSTR_BENEFITS": {
        "pl": "Korzyści dla trenera",
        "en": "Instructor benefits",
        "de": "Vorteile für Trainer"
    },
    "INSTR_PROFILE": {
        "pl": "Profil",
        "en": "Profile",
        "de": "Profil"
    },
    "INSTR_SPACE": {
        "pl": "Strefa trenera",
        "en": "Instructor space",
        "de": "Trainer Zone"
    },
    "INTERESTED_IN_LEGAL": {
        "pl": "Jeśli interesują Cię aspekty prawne, skorzystaj z poniższych linków:",
        "en": "If You're intrested in legal aspects please visit below link:",
        "de": "Wenn du Interesse an rechtlichen Aspekten hast, besuche die folgenden Links:"
    },
    "INTERMEDIATE_LEVEL": {
        "pl": "średnio-zaawansowany",
        "en": "medium",
        "de": "mittelstufe"
    },
    "INVALID_DATA": {
        "pl": "Wprowadzono nieprawidłowe dane",
        "en": "Invalid input data",
        "de": "Falsche Daten eingegeben"
    },
    "INVALID_EMAIL": {
        "pl": "Nieprawidłowy email",
        "en": "Invalid email",
        "de": "Ungültige E-Mail-Adresse"
    },
    "INVALID_PASSWORD": {
        "pl": "Nieprawidłowe hasło",
        "en": "Invalid password",
        "de": "Ungültiges Passwort"
    },
    "INVALID_SEARCH_QUERY": {
        "pl": "Nieprawidłowe parametry searchu - spróbuj jeszcze raz!",
        "en": "Invalid search parameters - try again!",
        "de": "Ungültige Suchparameter - versuchen Sie es erneut!"
    },
    "IN_CASE_OF_PROBLEMS_DOCS": {
        "pl": "W razie jakichkolwiek problemów z serwisem, proszę zapoznaj się z naszą dokumentacją:",
        "en": "In case of any troubles with our service, please read the documentation:",
        "de": "Im Fall von Problemen mit unserem Service lies bitte unsere Dokumentation:"
    },
    "IN_CASE_OF_QUESTIONS_SUPPORTS": {
        "pl": "W razie wątpliwości skontaktuj się z supportem:",
        "en": "When in doubt, contact Support:",
        "de": "Im Zweifelsfall kontaktiere bitte den Support:"
    },
    "IN_HRS": {
        "pl": " w godzinach ",
        "en": " in hours ",
        "de": "in Stunden"
    },
    "IN_PROGRESS_RSVS": {
        "pl": "Rezerwacje w trakcie realizacji",
        "en": "Reservations in progress",
        "de": "Reservierungen in der Realisation"
    },
    "ISSUED_FOR": {
        "pl": "Wydana dla",
        "en": "Issued for",
        "de": "Ausgestellt für"
    },
    "ISSUE_REPORTED": {
        "pl": "Został zgłoszony problem - rozwiązujemy go",
        "en": "Issue has been reported - we're on it",
        "de": "Ein Problem wurde angemeldet - wir lösen es"
    },
    "IT_SEEMS_YOU_WERE_LOGGED_OUT": {
        "pl": "Wygląda na to, że zostałeś wylogowany",
        "en": "It seems that you have been logged out",
        "de": "Es scheint, dass du abgemeldet worden bist"
    },
    "JOIN": {
        "pl": "Dołącz",
        "en": "Join",
        "de": "Join"
    },
    "JOIN_US": {
        "pl": "Dołącz do nas!",
        "en": "Join us!",
        "de": "Mach mit!"
    },
    "KNOWN_LANGS": {
        "pl": "Znane języki",
        "en": "Known langs",
        "de": "Bekannte Sprachen"
    },
    "LANG": {
        "pl": "Język",
        "en": "Language",
        "de": "Sprache"
    },
    "LAST_4_DIGITS": {
        "pl": "Ostatnie 4 cyfry:",
        "en": "Last 4 digits:",
        "de": "Die letzten 4 Zahlen"
    },
    "LAST_4_DIGITS_AKA": {
        "pl": "Ostatnie 4 liczby na karcie (tzw. podsumowanie karty)",
        "en": "Last 4 digits on card (aka card summary)",
        "de": "Die letzten 4 Ziffern der Karte (auch Kartenzusammenfassung genannt)"
    },
    "LAST_NAME": {
        "pl": "Nazwisko",
        "en": "Last name",
        "de": "Nachname"
    },
    "LEGAL_DOCS": {
        "pl": "Umowy",
        "en": "Legal documents",
        "de": "Verträge"
    },
    "LEVEL": {
        "pl": "Poziom",
        "en": "Level",
        "de": "Das Niveau"
    },
    "LEVEL_OF_DIFF": {
        "pl": "Poziom trudności",
        "en": "Difficulty level",
        "de": "Schwierigkeitsgrad"
    },
    "LIMITS": {
        "pl": "Limity",
        "en": "Limits",
        "de": "Limits"
    },
    "LIMITS_DETAILS": {
        "pl": "Limity (np. \"salka A\", albo imię prowadzącego) pozwolą ci określić ile osób może w jednym momencie przebywać na treningu",
        "en": "Limits (eg \"Room A\", or trainer name) allow you to limit amount of people participating in your trainings at any given moment.",
        "de": "Limits (z. B. \"Raum A\" oder Name des Trainers) ermöglichen es, die Anzahl der Personen, die zu einem bestimmten Zeitpunkt an dem Training teilnehmen, zu bestimmen."
    },
    "LIMIT_PEOPLE": {
        "pl": "Dla dowolnych treningów w obrębie tego Limitu, na zajęcia w jednym momencie może się zapisać max. %d %s",
        "en": "For any trainings belonging to this Limit, max %d %s will be allowed to register at the same time",
        "de": "Für alle Trainings, die zu diesem Limit gehören, können sich maximal %d %s Teilnehmer zur gleichen Zeit anmelden"
    },
    "LIMIT_TRAININGS": {
        "pl": "Jeżeli 2 lub więcej treningów (należących do tego limitu) nachodzi na siebie czasowo, lub ilością osób to maksymalnie na tylko %d z nich będzie można się zapisać",
        "en": "If 2 or more trainings [belonging to this Limit] overlap with each other, then your clients can only sign up for max. %d of them",
        "de": "Wenn sich 2 oder mehr Trainings [die zu diesem Limit gehören] überschneiden, dann können sich deine Kunden nur für max. %d davon anmelden"
    },
    "LINK_EXPIRED": {
        "pl": "Link wygasł",
        "en": "Link expired",
        "de": "Link ist abgelaufen"
    },
    "LINK_TO_PAYMENT": {
        "pl": "Link do płatności",
        "en": "Link to payment",
        "de": "Ling zur Zahlung"
    },
    "LIST": {
        "pl": "Lista",
        "en": "List",
        "de": "Liste"
    },
    "LOCATION": {
        "pl": "Lokalizacja",
        "en": "Location",
        "de": "Standort"
    },
    "LOC_REQUIRED": {
        "pl": "Lokacja jest wymagana",
        "en": "Location is required",
        "de": "Standort ist erforderlich"
    },
    "LOGIN": {
        "pl": "Zaloguj",
        "en": "Login",
        "de": "Anmelden"
    },
    "LOGIN_AGAIN": {
        "pl": "Zaloguj się ponownie",
        "en": "Login again",
        "de": "Melde dich erneut an"
    },
    "LOGIN_EMAIL": {
        "pl": "Email logowania",
        "en": "Login email",
        "de": "Anmeldung per E-Mail"
    },
    "LOGOUT": {
        "pl": "Wyloguj",
        "en": "Logout",
        "de": "Abmelden"
    },
    "MAKE_DECISION_ABOUT_RSV": {
        "pl": "Podejmij decyzję na temat tej rezerwacji",
        "en": "Make a decision about this reservation",
        "de": "Treffe eine Entscheidung über diese Reservierung"
    },
    "MAKE_SURE_INPUT_IS_OK": {
        "pl": "Zwróć uwagę na poprawność danych jakie wprowadzasz",
        "en": "Make sure your input data is correct",
        "de": "Achte auf die Richtigkeit der Daten, die du eingibst"
    },
    "MAP": {
        "pl": "Mapa",
        "en": "Map",
        "de": "die Karte"
    },
    "MARKETING": {
        "pl": "Zachęcaj!",
        "en": "Marketing!",
        "de": "Marketing!"
    },
    "MARKETING_DESC": {
        "pl": "Karnety, kody promocyjne - wszystko w zasięgu paru kliknięć",
        "en": "Vouchers, discount codes - everything to help You run Your business",
        "de": "Abonnements, Rabattcodes - alles in der Reichweite von einigen Klicks"
    },
    "MATCHING_TRAININGS": {
        "pl": "Znalezione treningi",
        "en": "Matching trainings",
        "de": "Gefundede Trainings"
    },
    "MAX_AGE": {
        "pl": "Max. wiek",
        "en": "Max. age",
        "de": "Max. Alter"
    },
    "MAX_ALLOWED_CHARS": {
        "pl": "Maksymalna ilość dozwolonych znaków w nazwie to 128",
        "en": "Maximum number of characters allowed in the name is 128",
        "de": "Die maximale Zahl der Zeichen im Namen beträgt 128"
    },
    "MAX_AMOUNT_OF_PEOPLE": {
        "pl": "Max. ilość osób",
        "en": "Max. amount of people",
        "de": "Max. Zahl der Teilnehmer"
    },
    "MAX_CAPACITY": {
        "pl": "Max. ilość osób na zajęciach",
        "en": "Max. people on training",
        "de": "Max. Teilnehmerzahl"
    },
    "MAX_NO_TRAININGS": {
        "pl": "Max. ilość treningów",
        "en": "Max amount of trainings",
        "de": "Max. Zahl von Trainings"
    },
    "MAX_NUMBER_OF_ENTRIES": {
        "pl": "Max ilość wejść",
        "en": "Max number of entries",
        "de": "Maximale Anzahl der Einträge"
    },
    "MAX_PEOPLE_DESC": {
        "pl": "Max ilość osób malejąco",
        "en": "Max number of people descending",
        "de": "Maximale Zahl der Teilnehmer absteigend"
    },
    "MAYBE_YOU_USE_PASS": {
        "pl": "Być może masz skonfigurowane logowanie/rejestracje przez hasło?",
        "en": "Maybe you have configured login / registration via the password?",
        "de": "Deine Anmeldung kann mit einem Passwort konfiguriert sein"
    },
    "MEET_PLATFORM": {
        "pl": "Poznaj platformę",
        "en": "Meet the plafrom",
        "de": "Erkunde die Platform"
    },
    "MGMT": {
        "pl": "Zarządzaj",
        "en": "Management",
        "de": "Verwaltung"
    },
    "MINUTES": {
        "pl": "minut",
        "en": "minutes",
        "de": "Minuten"
    },
    "MIN_AGE": {
        "pl": "Min. wiek",
        "en": "Min. age",
        "de": "Min. Alter"
    },
    "MISSING_INFO": {
        "pl": "Brakuje nam jeszcze paru ważnych informacji",
        "en": "We are still missing few important informations",
        "de": "Wir brauchen noch ein paar wichige Informationen"
    },
    "MODIFY": {
        "pl": "Zmodyfikuj",
        "en": "Modify",
        "de": "Verändere"
    },
    "MORE": {
        "pl": "Więcej...",
        "en": "More...",
        "de": "Mehr..."
    },
    "MULTIDAY_TRAININGS": {
        "pl": "Treningi wielodniowe!",
        "en": "Multiday trainings!",
        "de": "Mehrtägige Trainings!"
    },
    "MY_CARNETS": {
        "pl": "Moje karnety",
        "en": "My carnets",
        "de": "Meine Abonnements"
    },
    "MY_DATA": {
        "pl": "Moje Dane",
        "en": "My Data",
        "de": "Meine Daten"
    },
    "MY_DATA_INSTR": {
        "pl": "Te dane są związane z twoim kontem klienta.\n        Jako że także jesteś instruktorem, niektóre pola są widoczne i\nedytowalne w twoim",
        "en": "This data is related to your user account,\n        since you are also an instructor, some of this data can be\nviewed and changed in your",
        "de": "Diese Daten sind mit deinem Kundenkonto verbunden. \n       Da du Trainer bist, sind manche Felder sichtbar und editierbar in deinem\n"
    },
    "MY_DATA_INSTR_PROFILE": {
        "pl": "Profilu",
        "en": "Profile",
        "de": "des Profils"
    },
    "MY_LINKS": {
        "pl": "Moje linki",
        "en": "My links",
        "de": "Meine Links"
    },
    "MY_TRAININGS": {
        "pl": "Moje treningi",
        "en": "My trainings",
        "de": "Meine Trainings"
    },
    "NAME": {
        "pl": "Nazwa",
        "en": "Name",
        "de": "Name"
    },
    "NAME_AND_LAST_NAME_ON_CARD": {
        "pl": "Imię i nazwisko na karcie",
        "en": "Name and surname on the card",
        "de": "Name und Nachname auf der Karte"
    },
    "NAME_OR_NICK": {
        "pl": "Imię albo nick",
        "en": "Name or nick",
        "de": "Name oder Nick"
    },
    "NEAREST_DATE": {
        "pl": "Najbliższy wolny termin:",
        "en": "Nearest available date:",
        "de": "Der nächste freie Termin"
    },
    "NEW_CHATROOM": {
        "pl": "Nowy kanał",
        "en": "New channel",
        "de": "New channel"
    },
    "NEW_PASSWORD": {
        "pl": "Nowe hasło",
        "en": "New password",
        "de": "Neues Passwort"
    },
    "NEW_TRAINING": {
        "pl": "Nowy trening",
        "en": "New training",
        "de": "Neues Training"
    },
    "NEXT": {
        "pl": "Dalej",
        "en": "Next",
        "de": "Weiter"
    },
    "NO": {
        "pl": "Nie",
        "en": "No",
        "de": "Nein"
    },
    "NONE": {
        "pl": "brak",
        "en": "none",
        "de": "kein/e"
    },
    "NONSTANDARD": {
        "pl": "Niestandartowe",
        "en": "Non-standard",
        "de": "Nicht-Standard"
    },
    "NOONE_ELSE_WILL_SEE": {
        "pl": "Nikt inny nie będzie miał do nich wglądu",
        "en": "No one else will see them",
        "de": "Kein anderer wird Zugang zu den Daten haben"
    },
    "NOONE_REVIEWED_THIS": {
        "pl": "Jeszcze nikt nie wystawił opinii :(",
        "en": "No one reviewed this training so far :(",
        "de": "Bislang hat noch niemand das Training bewertet :("
    },
    "NOONE_WILL_BE_ABLE_TO_SIGN_UP_DURING_VACATION": {
        "pl": "Jak weźmiesz sobie wolne to w tym czasie nikt nie będzie się mógł zapisać na twoje treningi",
        "en": "No one will be able to sign up on your trainings during your leave.",
        "de": "Während deines Urlaubs wird sich niemand für deine Trainings anmelden können"
    },
    "NOTIFY_EMAIL": {
        "pl": "Twój email do notyfikacji",
        "en": "Your email used for notifications",
        "de": "Your email used for notifications"
    },
    "NOTIFY_EMAIL_DISCLAIMER": {
        "pl": "Twój email nie zostanie nikomu przekazany, użyjemy go po to by ci wysyłać notyfikacje z kanału",
        "en": "Your email won't be shared with other channel members, and will be only used by us to notify you about channel events",
        "de": "Your email won't be shared with other channel members, and will be only used by us to notify you about channel events"
    },
    "NOT_APPLICABLE": {
        "pl": "Nie dotyczy",
        "en": "Not applicable",
        "de": "Nicht zutreffend"
    },
    "NO_AVAILABLE_SLOTS_FOR_THIS_TRAINING": {
        "pl": "Nie ma już wolnych miejsc na ten trening",
        "en": "There are no available slots for this training",
        "de": "Es gibt keine freinen Plätze für dieses Training"
    },
    "NO_CARNETS_AVAILABLE": {
        "pl": "Brak karnetów na te zajęcia",
        "en": "No carnets available",
        "de": "Keine Abonnements verfügbar"
    },
    "NO_INCOMING_RSV": {
        "pl": "Nie masz żadnych nadchodzących rezerwacji",
        "en": "You have no incoming reservations",
        "de": "Du hast keine eingehenden Reservierungen"
    },
    "NO_LIMIT_PEOPLE": {
        "pl": "W obrębie tego Limitu nie ma ograniczenia na ilość osób przebywających na zajęciach w jednym momencie",
        "en": "There is not limit on amount of people present on the trainings in this Limit",
        "de": "Es gibt keine Begrenzung der Anzahl der Teilnehmer an den Trainings innerhalb dieses Limits"
    },
    "NO_LIMIT_TRAININGS": {
        "pl": "W obrębie tego Limitu nie ma ograniczeń co do ilości treningów nachodzących na siebie",
        "en": "There is no limit on amount of trainings happening simultaneously within this Limit",
        "de": "Es gibt keine Begrenzung für die Anzahl der gleichzeitig stattfindenden Trainings innerhalb des Limits"
    },
    "NO_OAUTH_CODE": {
        "pl": "Nie otrzymaliśmy kodu od providera Oauth",
        "en": "We have not received code from the OAuth Provider",
        "de": "Wir haben keinen Code vom Oauth Provider erhalten"
    },
    "NO_RESERVATIONS": {
        "pl": "Brak rezerwacji na te zajęcia",
        "en": "No reservations for this training",
        "de": "Keine Reservierungen für dieses Training"
    },
    "NO_USES": {
        "pl": "Ilość użyć",
        "en": "Number of uses",
        "de": "Anzahl der Verwendungen"
    },
    "NUMBER_OF_ENTRIES": {
        "pl": "Ilość wejść",
        "en": "Number of entries",
        "de": "Anzahl der Einträge"
    },
    "NUMBER_OF_ROWS_PER_PAGE": {
        "pl": "Ilość wierszy na stronie",
        "en": "Number of rows on the page",
        "de": "Die Anzahl der Zeilen auf der Seite"
    },
    "OCCURRENCE": {
        "pl": "Występowanie",
        "en": "Occurrence",
        "de": "Auftreten"
    },
    "OF_ACCOUNT": {
        "pl": " konta",
        "en": " account",
        "de": " des Kontos"
    },
    "OLD_PASSWORD": {
        "pl": "Stare hasło",
        "en": "Old password",
        "de": "Altes Passwort"
    },
    "ONCE_AGAIN": {
        "pl": "Spróbuj jeszcze raz",
        "en": "Try again",
        "de": "Erneut versuchen"
    },
    "ONCE_YOU_ENTER_CARD": {
        "pl": "jak wprowadzisz tu dane swojej karty, przekażemy te dane bezpośrednio",
        "en": "once you enter Your card details here, we will pass it on directly",
        "de": "Wenn du hier deine Kartendaten eingibst, werden wir diese direkt weiterleiten "
    },
    "OR": {
        "pl": "lub",
        "en": "or",
        "de": "oder"
    },
    "ORDERED_PAYED_YOURS": {
        "pl": "Zamówione, zapłacone - Twoje!",
        "en": "Booked, payed - Yours!",
        "de": "Gebucht, bezahlt - Deins!"
    },
    "OR_IF_YOU_HAVE_ACCOUNT_YOU_CAN_CHECK": {
        "pl": "Lub jeżeli masz konto to możesz to też sprawdzić w",
        "en": "Or if you have an account, then you can also see",
        "de": "Oder wenn du ein Konto hast, kannst du es prüfen in"
    },
    "OR_LOGIN_VIA": {
        "pl": "Lub zaloguj poprzez",
        "en": "Or login via",
        "de": "Oder melde dich an via"
    },
    "OR_SIGN_UP_WITH": {
        "pl": "Lub dołącz poprzez",
        "en": "Or sign up with...",
        "de": "Oder melden Sie sich an mit"
    },
    "OR_YOU_CAN_EMBED_IT": {
        "pl": "Lub możesz zintegrować go ze swoją stroną",
        "en": "Or you can embed it into your website",
        "de": "Oder du kannst es mit deiner Homepage integrieren"
    },
    "OUR_RATE": {
        "pl": "5% + VAT/płatność",
        "en": "5% + VAT/payment",
        "de": "5% + VAT/Zahlung"
    },
    "PASSWORD": {
        "pl": "Hasło",
        "en": "Password",
        "de": "Passwort"
    },
    "PASSWORD_MANAGED_BY_GOOGLE": {
        "pl": "Hasło zarządzane przez google",
        "en": "Password managed by google",
        "de": "Von Google verwaltetes Passwort "
    },
    "PASSWORD_RECOVERY": {
        "pl": "Odzyskiwanie hasła",
        "en": "Password recovery",
        "de": "Passwort-Wiederherstellung"
    },
    "PASSWORD_RECOVERY_SUB": {
        "pl": "Aby odzyskać hasło, podaj adres email którego użyłeś podczas rejestracji",
        "en": "We will send reset link to the following email address",
        "de": "Gib und deine E-Mail-Adresse und wir schicken dir einen Rücksetzlink"
    },
    "PASS_RESET": {
        "pl": "Reset hasła",
        "en": "Password reset",
        "de": "Passwort wird zurückgesetzt"
    },
    "PAYMENT": {
        "pl": "Płatność ",
        "en": "Payment ",
        "de": "Zahlung "
    },
    "PAYMENTS": {
        "pl": "Płatności",
        "en": "Payments",
        "de": "Zahlungen"
    },
    "PAYMENT_BROKER": {
        "pl": "bramce płatniczej",
        "en": "to the payment broker",
        "de": "an den Zahlungsmittler"
    },
    "PAYMENT_STATUS": {
        "pl": "Status płatności",
        "en": "Payment status",
        "de": "Zahlungsstatus"
    },
    "PAYOUTS": {
        "pl": "Wynagrodzenia",
        "en": "Payouts",
        "de": "Auszahlungen"
    },
    "PAYOUTS_DESC": {
        "pl": "Wynagrodzenie za trening zostanie przesłane w 24h po\nukończonym treningu, chyba że klient poprosi o reklamację.\n        Z tego właśnie powodu będziemy potrzebowali danych do karty\npłatniczej - nie przechowujemy tych danych,\n        nie mamy do nich dostepu.",
        "en": "Your salary for training will be transferred in 24h after\ntraining finish, if customer won't have any compliants.\n        We don't know Your credit/debt card number details, all is\nprocessed and kept by our payment broker.",
        "de": "Deine Bezahlung für die Trainings wird innerhalb von 24 Stunden nach dem abgeschlossenen Training überwiesen, \n       es sei denn der Kunde hat eine Reklamation erhoben. Aus diesem Grund brauchen die Daten deiner Kredit-/EC-Karte. \n       Diese Daten werden von uns nicht gespeichert, sondern nur zum Zweck der Zahlungen erhoben."
    },
    "PAYOUT_IN_PROGRESS": {
        "pl": "Wypłacanie środków instruktorowi",
        "en": "Instructor payout in progress",
        "de": "Auszahlung an den Trainer erfolgt"
    },
    "PERIOD_OF_VALIDITY": {
        "pl": "Okres ważności",
        "en": "Period of validity",
        "de": "Gültigkeitsdatum"
    },
    "PERSONAL_GROUP_CARNETS": {
        "pl": "Treningi: Personalne, grupowe i karnety",
        "en": "Trainings: Personal, group and carnets",
        "de": "Trainings: persönlich, in der Guppe und Abonnements"
    },
    "PERSONAL_TRAININGS": {
        "pl": "Treningi personalne",
        "en": "Personal trainings",
        "de": "Persönliche Trainings"
    },
    "PHONE": {
        "pl": "Numer telefonu",
        "en": "Phone number",
        "de": "Telefonnummer"
    },
    "PHOTO_DISPLAY": {
        "pl": "Będą się one wyświetlały w siatce przy prezentacji treningu",
        "en": "They will be displayed in the grid during training presentation ",
        "de": "Sie werden während der Trainingansicht im Raster angezeigt "
    },
    "PHOTO_NUM_FMT": {
        "pl": "Możesz dodać max jedno główne zdjęcie + %d dodatkowych",
        "en": "You can add one main picture + %d additional",
        "de": "Du kannst ein Hauptbild + %d zusätzliche Bilder hinzufügen"
    },
    "PHOTO_SECTIONS": {
        "pl": "Twoje zdjęcia (max 6)",
        "en": "Your photos (max 6)",
        "de": "Deine Bilder (max 6)"
    },
    "PICTURE_TOO_LARGE": {
        "pl": "Rozmiar zdjęcia jest zbyt duży (max: 5MB)",
        "en": "Picture is too large (max: 5MB)",
        "de": "Das Bild ist zu groß (max: 5MB)"
    },
    "PREVIEW": {
        "pl": "Podgląd",
        "en": "Summary",
        "de": "Ansicht"
    },
    "PRICE": {
        "pl": "Cena",
        "en": "Price",
        "de": "Preis"
    },
    "PRICE_ASC": {
        "pl": "Cena rosnąco",
        "en": "Price ascending",
        "de": "Preis aufsteigend"
    },
    "PRICE_DESC": {
        "pl": "Cena malejąco",
        "en": "Price descending",
        "de": "Preis absteigend"
    },
    "PRICE_VAL": {
        "pl": "Cena musi być większa niż ",
        "en": "The price must be greater than ",
        "de": "Der Preis muss höher sein als "
    },
    "PRIVACY_POLICY": {
        "pl": "Politykę prywatności",
        "en": "Privacy policy",
        "de": "Datenschutzbestimmungen"
    },
    "PROBLEM_FETCHING_TRAINING": {
        "pl": "Wystąpił problem przy pobieraniu treningu",
        "en": "There was problem during fetching your training data",
        "de": "Es ist ein Problem beim Abrufen deiner Trainingsdaten aufgetreten"
    },
    "PROBLEM_SHOULD_BE_SOLVED_SOON": {
        "pl": "Problem powinien być wkrótce rozwiązany",
        "en": "Problem should be solved soon",
        "de": "Das Problem sollte bald behoben werden"
    },
    "QR_HAS_BEEN_USED": {
        "pl": "Ten kod QR już został użyty",
        "en": "This QR code has already been used",
        "de": "Dieder QR-Code wurde schon verwendet"
    },
    "QUICK_AND_SECURE_PAYMENTS": {
        "pl": "Szybkie i bezpieczne płatności sprawią,\n        że opłacenie wybranego treningu potrwa mniej niż minutę!",
        "en": "Fast and secure payments makes it possible to get desired\ntraining in less than a minute!",
        "de": "Schnelle und sichere Zahlungsmethoden ermöglichen es, das ausgewählte Training in weniger als eine Minute zu bezahlen!"
    },
    "QUICK_LOOK_AT_CALENDAR": {
        "pl": "2 szybkie spojrzenia na harmonogram i już wiesz który\ntermin najbardziej Ci odpowiada.\n        Sam zdecydujesz czy chcesz pojedynczy trening, kurs, czy karnet.",
        "en": "Quick look at calendar and you know which date suits you best.\n        Buy personal training, carnet, course in easy and comfortable way!",
        "de": "Kurzer Blik auf den Zeitplan und du kannst entscheiden, welcher Termin für dich der beste ist.\n        Du entscheidest selbst, ob du ein Training oder das Abonnement kaufen willst"
    },
    "RATE_TRAININGS": {
        "pl": "Oceń zajęcia, pomóż innym zdecydować!",
        "en": "Rate trainings",
        "de": "Bewerte die Trainings!"
    },
    "REACTIVATE_ACCOUNT": {
        "pl": "Reaktywuj konto",
        "en": "Reactivate account",
        "de": "Konto reaktivieren"
    },
    "REACTIVATING": {
        "pl": "Reaktywowanie",
        "en": "Reactivating",
        "de": "Reaktivierung"
    },
    "READ_FURTHER": {
        "pl": "Czytaj dalej...",
        "en": "Read further...",
        "de": "Lies weiter..."
    },
    "RECOMMENDED_GEAR": {
        "pl": "Polecany sprzęt",
        "en": "Recommended gear",
        "de": "Empfohlene Ausrüstung"
    },
    "REFER_CREATE_YOUR_OFFER": {
        "pl": "Możesz przejść teraz do harmonogramu i stworzyć swoją ofertę",
        "en": "You can now go to the schedule and create your offer",
        "de": "Jetzt kannst du zum Zeitplan gehen und dein Angebot erstellen"
    },
    "REFUNDED": {
        "pl": "Zwróciliśmy pieniądze klientowi",
        "en": "reservation has been refunded",
        "de": "Geld für das gebuchte Training wurde dem Kunden zurückerstattet"
    },
    "REFUNDING": {
        "pl": "robimy zwrot pieniędzy klientowi",
        "en": "Making refund customer's account",
        "de": "Wir erstatten das Geld zurück"
    },
    "REGISTER": {
        "pl": "Zarejestruj",
        "en": "Register",
        "de": "Konto erstellen"
    },
    "REGISTER_AS_INSTR": {
        "pl": "Zarejestruj jako trener",
        "en": "Register as instructor",
        "de": "Als Trainer anmelden"
    },
    "REJECT": {
        "pl": "Odrzuć",
        "en": "Reject",
        "de": "Ablehnen"
    },
    "REJECTED": {
        "pl": "odrzucono",
        "en": "rejected",
        "de": "abgelehnt"
    },
    "REJECT_CHANGES": {
        "pl": "odrzuć zmiany",
        "en": "reject changes",
        "de": "Änderungen ablehnen"
    },
    "REJECT_USER": {
        "pl": "Odrzuć użytkownika: ",
        "en": "Reject user: ",
        "de": "Benutzer ablehnen: "
    },
    "REMAINING": {
        "pl": "Pozostało",
        "en": "Remaining",
        "de": "Übrige"
    },
    "REMAINING_ENTRIES": {
        "pl": "Pozostałe wejścia",
        "en": "Remaining entries",
        "de": "Übrige Einträge"
    },
    "REMOVING_TRAINING": {
        "pl": "Usuwanie treningu",
        "en": "Removing training",
        "de": "Entfernen des Trainings"
    },
    "REPEATED_EVERY": {
        "pl": ", powtarzane co",
        "en": ", repeated every",
        "de": ", wiederholt jede/n "
    },
    "REPEATING": {
        "pl": "Powtarzanie",
        "en": "Repeating",
        "de": "Wiederholung"
    },
    "REPORT_ISSUE": {
        "pl": "Zgłoś problem",
        "en": "Report issue",
        "de": "Problem melden"
    },
    "REPORT_RSV": {
        "pl": "Zgłoś rezerwację",
        "en": "Report reservation",
        "de": "Reservierung anmelden"
    },
    "REQUIRED_GEAR": {
        "pl": "Wymagany sprzęt",
        "en": "Required gear",
        "de": "Erforderliche Ausrüstung"
    },
    "RESENT_EMAIL": {
        "pl": "Wysłaliśmy ci jeszcze raz email",
        "en": "We have sent you confirmation email again",
        "de": "Wir haben die Bestätigungs-E-Mail erneut gesendet"
    },
    "RESERVATION": {
        "pl": "Rezerwacja",
        "en": "Reservation",
        "de": "Reservierung"
    },
    "RESERVATIONS": {
        "pl": "Rezerwacje",
        "en": "Reservations",
        "de": "Buchungen"
    },
    "RESET": {
        "pl": "Resetuj",
        "en": "Reset",
        "de": "Zurücksetzen"
    },
    "RESET_APP": {
        "pl": "Zresetuj aplikację",
        "en": "Reset application",
        "de": "App zurücksetzen"
    },
    "RESET_FILTERS": {
        "pl": "Zresetuj filtry",
        "en": "Reset all filters",
        "de": "Alle Filter zurücksetzen"
    },
    "RETURN": {
        "pl": "Powrót",
        "en": "Return",
        "de": "Zurück"
    },
    "REVIEWS": {
        "pl": "Opinie",
        "en": "Reviews",
        "de": "Bewerungen"
    },
    "REVIEW_HAS_BEEN_SAVED": {
        "pl": "Twoja opinia została zapisana",
        "en": "Review has been saved",
        "de": "Deine Bewertung wurde gespeichert"
    },
    "RM_ELEM": {
        "pl": "Usuń element",
        "en": "Remove element",
        "de": "Das Element entfernen"
    },
    "RSV_CONFIRMED": {
        "pl": "Rezerwacja została potwierdzona",
        "en": "Reservation has been confirmed",
        "de": "Die Reservierung wurde bestätigt"
    },
    "RSV_NOT_PAYED": {
        "pl": "Rezerwacja nie została opłacona",
        "en": "Reservation has not been paid for",
        "de": "Die Reservierung wurde nicht bezahlt"
    },
    "RSV_SIGN_UP": {
        "pl": "Zapisz się",
        "en": "Sign up",
        "de": "Melde dich an"
    },
    "SAVE": {
        "pl": "Zapisz",
        "en": "Save",
        "de": "Speicher"
    },
    "SAVE_CHANGES": {
        "pl": "Zapisz zmiany",
        "en": "Save changes",
        "de": "Änderungen spreichern"
    },
    "SCHEDULE": {
        "pl": "Grafik",
        "en": "Schedule",
        "de": "Zeitplan"
    },
    "SEARCH": {
        "pl": "Wyszukaj",
        "en": "Search",
        "de": "Suche"
    },
    "SEARCH_AND_BOOK": {
        "pl": "ZNAJDŹ, ZAREZERWUJ - EWOLUUJ",
        "en": "FIND IT, BOOK IT, EVOLVE",
        "de": "Training suchen und reservieren"
    },
    "SECURE_CERTAIN_EXPRESS": {
        "pl": "Bezpieczne, pewne, ekspresowe!",
        "en": "Secure, certain and express!",
        "de": "Sicher, gewiss und schnell!"
    },
    "SELECT_DIFF": {
        "pl": "Dobierz poziom trudndości",
        "en": "Select difficulty level",
        "de": "Schwierigkeitsgrad festlegen"
    },
    "SEND_EMAIL": {
        "pl": "Wyślij email",
        "en": "Send email",
        "de": "E-Mail senden"
    },
    "SEND_EMAIL_AGAIN": {
        "pl": "Wyślij email jeszcze raz",
        "en": "Resend email",
        "de": "E-Mail erneut senden"
    },
    "SENT_EMAIL": {
        "pl": "Wysłaliśmy ci email z dalszymi instrukcjami",
        "en": "We sent you an email with further instructions",
        "de": "Wir haben dir eine E-Mail mit weiteren Anweisungen geschickt"
    },
    "SENT_EMAIL_TO": {
        "pl": "Wysłaliśmy wiadomość email na adres: ",
        "en": "We sent email to: ",
        "de": "Wir schickten eine Mail an"
    },
    "SERVER_IS_NOT_ABLE_TO_PROCESS": {
        "pl": "Serwer obecnie nie jest w stanie przetworzyć tego żądania",
        "en": "Server is currently unable to process this request",
        "de": "Der Server kann diese Anfrage derzeit nicht bearbeiten"
    },
    "SET_NAME_AND_SURNAME": {
        "pl": "Ustaw imię i nazwisko",
        "en": "Set name and surname",
        "de": "Den Namen und Nachnamen festlegen"
    },
    "SET_NEW_PASSWORD": {
        "pl": "Ustaw nowe hasło",
        "en": "Set new password",
        "de": "Neues Passwort festlegen"
    },
    "SHARE_PRIOR_BOOKING": {
        "pl": "Udostępniaj te dane przed złożeniem rezerwacji",
        "en": "Share this data prior to booking",
        "de": "Teile diese Daten vor der Buchung mit"
    },
    "SIGN_IN_FOR_THIS_TRAINING": {
        "pl": "Zapisz się na ten trening",
        "en": "Sign in for this training",
        "de": "Melde dich für das Training an"
    },
    "SIMILAR_OFFERS": {
        "pl": "Podobne oferty",
        "en": "Similar offers",
        "de": "Ähnliche Angebote"
    },
    "SIMULTANEOUS_TRAININGS": {
        "pl": "Jeden trener na dwie sale?",
        "en": "Simultanous trainings?",
        "de": "Gleichzeitige Trainings?"
    },
    "SINCE_WHEN_INSTR": {
        "pl": "Od kiedy jesteś instruktorem?",
        "en": "Since when are you an instructor?",
        "de": "Seit wann bist du Trainer?"
    },
    "SKIING_DANCING_BOX": {
        "pl": "Znajdź idealny trening dla siebie",
        "en": "Find perfect training for You",
    },
    "SOMETHING_WENT_WRONG": {
        "pl": "Coś poszło nie tak",
        "en": "Something went wrong",
        "de": "Etwas ist schiefgelaufen"
    },
    "SOMETHING_WENT_WRONG_CONTACT": {
        "pl": "Coś poszło nie tak, wkrótce się z tobą skontaktujemy!",
        "en": "Something went wrong, we will contact you soon!",
        "de": "Etwas ist schiefgelaufen, wir kontaktieren dich bald!"
    },
    "SOMETHING_WENT_WRONG_FETCH_VACATION": {
        "pl": "Uups, coś poszło nie tak podczas pobierania twoich urlopów",
        "en": "Whoops, something went wrong while fetching Your vacation data",
        "de": "Uups, beim Abrufen deines Urlaubs ist etwas schiefgelaufen"
    },
    "SORTING": {
        "pl": "Sortowanie",
        "en": "Sorting",
        "de": "Sortierung"
    },
    "START": {
        "pl": "Start",
        "en": "Start",
        "de": "Start"
    },
    "START_DATE": {
        "pl": "Data rozpoczęcia",
        "en": "Start date",
        "de": "Anfangsdatum"
    },
    "START_HOUR": {
        "pl": "Godzina rozpoczęcia",
        "en": "Start hour",
        "de": "Die Anfangszeit"
    },
    "START_SHOULD_BE_AFTER_END": {
        "pl": "Data zakończenia powinna być większa niż rozpoczęcia",
        "en": "End date should be bigger than start",
        "de": "Das Enddatum sollte größer sein als das Anfangsdatum"
    },
    "STATUS": {
        "pl": "Status: ",
        "en": "Status: ",
        "de": "Status: "
    },
    "STOP_BEING_INSTR": {
        "pl": "Przestań być instruktorem",
        "en": "Stop being instructor",
        "de": "Beende deine Tätigkeit des Trainers"
    },
    "SUB_NAME": {
        "pl": "Karnet",
        "en": "Carnet",
        "de": "Abonnement"
    },
    "SUCCESS": {
        "pl": "Sukces",
        "en": "Success",
        "de": "Erfolg"
    },
    "TAGS": {
        "pl": "Tagi",
        "en": "Tags",
        "de": "Tags"
    },
    "TEMPORARILY_DEACTIVATE_INSTR": {
        "pl": "Tymczasowo deaktywuj swoje konto instruktora",
        "en": "Temporarily deactivate Your instructor account",
        "de": "Deaktiviere vorübergehend dein Trainerkonto"
    },
    "TEMP_INACTIVE": {
        "pl": "Twoje konto instruktora zostało tymczasowo wyłączone",
        "en": "Your account is temporarily disabled",
        "de": "Dein Trainerkonto wurde vorübergehen deaktiviert"
    },
    "TERM": {
        "pl": "Termin",
        "en": "Term",
        "de": "der Termin"
    },
    "TERMS_FOR": {
        "pl": "terminy na: ",
        "en": "terms for: ",
        "de": "Termine für "
    },
    "TERMS_OF_SERVICE": {
        "pl": "Warunki korzystania z usługi",
        "en": "Terms of service",
        "de": "Nutzungsbedingungen"
    },
    "TEXT_SECTIONS": {
        "pl": "Twoje info",
        "en": "Text sections",
        "de": "Deine Informationen"
    },
    "TITLE": {
        "pl": "Tytuł",
        "en": "Title",
        "de": "Titel"
    },
    "TO": {
        "pl": "Do",
        "en": "To",
        "de": "Bis"
    },
    "TODAY": {
        "pl": "dzisiaj ",
        "en": "today ",
        "de": "heute "
    },
    "TODAY_AT": {
        "pl": "dzisiaj o ",
        "en": "today at ",
        "de": "heute um "
    },
    "TOKEN_HAS_EXPIRED": {
        "pl": "Token wygasł - spróbuj się zarejestrować jeszcze raz",
        "en": "Token has expired - please try to create account again",
        "de": "Der Token ist abgelaufen - bitte versuch dein Konto erneut zu erstellen"
    },
    "TOMORROW": {
        "pl": "jutro ",
        "en": "tomorrow ",
        "de": "morgen"
    },
    "TOMORROW_AT": {
        "pl": "jutro o ",
        "en": "tomorrow at ",
        "de": "morgen um "
    },
    "TOTAL_COST": {
        "pl": "Koszt całkowity",
        "en": "Total cost",
        "de": "Gesamtkosten"
    },
    "TO_CONFIRM_YOU_MUST_BE_LOGGED_IN": {
        "pl": "Pamiętaj, by móc potwierdzić rezerwację musisz być zalogowany na telefonie na swoje konto instruktora",
        "en": "Remember to confirm reservation you must be logged in to your instructor account",
        "de": "Beachte, um die Reservierung zu bestätigen, musst du auf deinem Konto des Trainers angemeldet sein"
    },
    "TO_DEACTIVATE_INSTR_ACC_": {
        "pl": "By móc usunąć swoje konto instruktora musisz zamknąć wszystkie aktywne rezerwacje: anulować je, albo je wykonać i potem usunąć konto",
        "en": "To be able to delete Your instructor account you must close all active bookings: cancel them or finish them, and then delete the account",
        "de": "Um dein Trainerkonto löschen zu können, musst du alle aktiven Reservierungen schließen: stornieren oder beenden, und dann das Konto löschen"
    },
    "TO_FINISH_REGISTER": {
        "pl": "by zakończyć proces rejestracji",
        "en": "to finish registration",
        "de": "um die Registrierung abzuschliessen"
    },
    "TO_REJECT_WITHOUT_CONSEC": {
        "pl": "by wciąż odrzucić rezerwacje bez poniesienia konsekwencji",
        "en": "To still reject reservations without incurring any consequences",
        "de": "Um die Reservierung ohne Folgen abzulehnen"
    },
    "TRAINING": {
        "pl": "Trening",
        "en": "Training",
        "de": "Training"
    },
    "TRAININGS": {
        "pl": "Treningi",
        "en": "Trainings",
        "de": "Trainings"
    },
    "TRAININGS_DESC": {
        "pl": "Przy pomocy naszego edytora, szybko stworzysz swoją ofertę treningową.",
        "en": "Using our editor create Your offer in minutes!",
        "de": "Mithilfe unseres Editors wirst du schnell dein eigenes Angebott erstellen können"
    },
    "TRAINING_GONE": {
        "pl": "Rezerwacja jest aktualna, ale trening już nie istnieje, lub został przesunięty",
        "en": "Reservation is still valid, but training no longer exists, or has been moved",
        "de": "Die Buchung ist noch gültig, aber das Training existiert nicht mehr oder wurde verschoben"
    },
    "TRAINING_HAS_BEEN_REMOVED": {
        "pl": "Trening został usunięty",
        "en": "Training has been removed",
        "de": "das Training wurde entfernt"
    },
    "TRAINING_INFO": {
        "pl": "Dane i edycja treningu",
        "en": "Training info and management",
        "de": "Daten und Verwaltung des Trainings"
    },
    "TRAINING_NO": {
        "pl": "zajęcia nr ",
        "en": "Training No. ",
        "de": "Training Nummer "
    },
    "TRAINING_NOT_FOUND": {
        "pl": "Nie znaleźliśmy podanego treningu",
        "en": "Training not found",
        "de": "Das Training wurde nicht gefunden"
    },
    "TRAINING_PHOTOS": {
        "en": "Training photos",
        "pl": "Zdjęcia treningu",
        "de": "Fotos vom Training"
    },
    "TRAINING_PRICE": {
        "pl": "Cena treningu",
        "en": "Training price",
        "de": "Der Preis des Trainings"
    },
    "TRANSACTION_COST": {
        "pl": "Koszt transakcji",
        "en": "Transaction cost",
        "de": "Kosten der Transaktion"
    },
    "TRANSFER_IN_6H": {
        "pl": "Przelew już w 6h!",
        "en": "Money transfer in 6h!",
        "de": "Geldüberweisung innerhalb von 6h!"
    },
    "TWOFA_NOT_SUPPORTED": {
        "en": "Server does not support 2fa",
        "pl": "Serwer nie supportuje 2FA",
        "de": "Server unterstützt keine 2fa"
    },
    "TYPE": {
        "pl": "Typ",
        "en": "Type",
        "de": "Typ"
    },
    "T_NAME_REQUIRED": {
        "pl": "Nazwa treningu jest wymagana",
        "en": "Training name is required",
        "de": "Name des Trainings ist erforderlich"
    },
    "UNEXPECTED_ERROR": {
        "pl": "Wystąpił nieoczekiwany błąd - spróbuj jeszcze raz później",
        "en": "Unexpected error occurred - try again later",
        "de": "Unerwarteter Fehler aufgetreten - versuch es später noch einmal"
    },
    "UNKOWN_ERROR_OCCURRED": {
        "pl": "Wystąpił nieznany błąd",
        "en": "An unkown error occurred",
        "de": "ein unbekannter Fehler ist aufgetreten"
    },
    "UNLIMITED": {
        "pl": "Nieograniczona",
        "en": "Unlimited",
        "de": "Unbegrenzt"
    },
    "UNTIL_PROBLEM_SOLVED": {
        "pl": "Do czasu rozwiązania sporu flow zostanie zatrzymany",
        "en": "Until we solve Your problem, flow of this reservation is halted",
        "de": "Bis wir dein Problem gelöst haben, ist die Bearbeitung dieser Reservierung gestoppt"
    },
    "USER_SPACE": {
        "pl": "Strefa użytkownika",
        "en": "User space",
        "de": "User zone"
    },
    "USE_ACCOUNT_DATA": {
        "pl": "Użyj danych z konta",
        "en": "Use account data",
        "de": "Kontodaten verwenden"
    },
    "USE_ANY_NAME_OR": {
        "pl": "Możesz użyć dowolnej nazwy kodu lub",
        "en": "You can use any code name or",
        "de": "Du kannst jeden beliebigen Codenamen benutzen oder"
    },
    "VACATION": {
        "pl": "Wolne",
        "en": "Vacations",
        "de": "Urlaub"
    },
    "VACATION_WHENEVER_YOU_WANT": {
        "pl": "Zrób sobie wolne kiedy chcesz!",
        "en": "Vaction whenever you want!",
        "de": "Urlaub wann immer du willst!"
    },
    "VALID_FROM": {
        "pl": "Ważny od",
        "en": "Valid from",
        "de": "Gültig vom"
    },
    "VALID_TO": {
        "pl": "Ważny to",
        "en": "Valid until",
        "de": "Gültig bis"
    },
    "VALUE_YOUR_TIME": {
        "pl": "Dopasuj pod siebie",
        "en": "Value your time",
        "de": "Passe es an deine Bedürfnisse an"
    },
    "VERIFY_CARNET_LIKE_RSV": {
        "pl": "Przy wejściu na trening możesz zweryfikować karnet tak jak weryfikujesz rezerwacje - poprzez kod QR",
        "en": "When entering training, you can verify the carnet the same way as you verify your reservations - via the QR code",
        "de": "Beim Antreten des Trainings kannst du dein Abonnement verifizieren, auf dieselbe Weise wie du deine Reservierungen verifizierst, über den QR-Code"
    },
    "VER_OF_RSV_QR_CODE": {
        "pl": "Weryfikacja kodu QR na rezerwację",
        "en": "Verification of QR code for reservation",
        "de": "Überprüfung des QR-Codes für die Reservierung"
    },
    "VER_OF_SUB_QR_CODE": {
        "pl": "Weryfikacja kodu QR (karnet)",
        "en": "Verification of the QR code (carnet)",
        "de": "Überprüfung des QR-Codes (Abonnement)"
    },
    "WAITING_FOR_PAYMENT": {
        "pl": "Czekanie na wykonanie płatności",
        "en": "Waiting for payment",
        "de": "Warten auf die Zahlung"
    },
    "WAIT_A_MOMENT": {
        "pl": "Chwilka...",
        "en": "Wait a moment, please...",
        "de": "Einen Moment, bitte..."
    },
    "WARNING": {
        "pl": "Uwaga",
        "en": "Warning",
        "de": "Achtung"
    },
    "WEEKLY": {
        "pl": "Co tydzień",
        "en": "Weekly",
        "de": "Wöchentlich"
    },
    "WE_DONT_KNOW_CARD": {
        "pl": "nie znamy ani nie przechowujemy danych twojej karty w naszym systemie",
        "en": "we do not know or store Your card details in our system",
        "de": "wir kennen und speichern deine Kartendaten nicht in unserem System"
    },
    "WE_ENCOURAGE_EMAIL": {
        "pl": "Zalecamy podanie emaila do kontaktu i notyfikacji",
        "en": "We encourage passing the email to allow contact and notifications",
        "de": "Wir empfehlen die Weitergabe der E-Mail-Adresse, um den Kontakt und Benachrichtigungen zu ermöglichen"
    },
    "WE_ENVY_YOU_SUCCESS": {
        "pl": "Już Ci zazdrościmy efektów! ;)",
        "en": "Enjoy your training!",
        "de": "Wir freuen uns über deine Erfolge! ;)"
    },
    "WE_TAKE_SECURITY": {
        "pl": "Bezpieczeństwo twoich danych płatności traktujemy bardzo poważnie, dlatego ",
        "en": "We take security of Your payment details very seriously, which is why ",
        "de": "Wir nehmen die Sicherheit deiner Zahlungsdaten sehr ernst, deshalb "
    },
    "WE_WILL_CHARGE_AFTER": {
        "pl": "Środki z konta dopiero pobierzemy za 24 godziny",
        "en": "Bank account will be charged after 24 hours",
        "de": "Die Mittel werden vom Konto nach 24 Stunden abgebucht"
    },
    "WE_WILL_RESPOND": {
        "pl": "Odpowiemy na każdego! :)",
        "en": "We respond to all emails!",
        "de": "Wir beantworten jede E-Mail!"
    },
    "WITH": {
        "pl": "z",
        "en": "with"
    }, 
    "WHAT_DO_YOU_WANT_TO_TRAIN": {
        "pl": "Co chcesz trenować?",
        "en": "What do you want to train?",
        "de": "Was willst du trainieren?"
    },
    "WHEN": {
        "pl": "Kiedy: ",
        "en": "When: ",
        "de": "Wann: "
    },
    "WHEN_DO_YOU_WANT_TO_SIGN_UP": {
        "pl": "Na kiedy chcesz się zapisać",
        "en": "What date do you want to sign up for",
        "de": "Wann willst du dich anmelden"
    },
    "WHERE": {
        "pl": "Gdzie?",
        "en": "Where?",
        "de": "Wo?"
    },
    "WHERE_WILL_BE_T": {
        "pl": "Gdzie odbędzie się trening?",
        "en": "Where will training take place?",
        "de": "Wo wird das Training stattfinden?"
    },
    "WHICH_DATES_DO_YOU_PREFER": {
        "pl": "Jakie terminy cię interesują?",
        "en": "Which dates do you prefer?",
        "de": "Bevorzugte Termine"
    },
    "WHICH_DAYS_DO_YOU_PREFER": {
        "pl": "W jakie dni chcesz trenować",
        "en": "When are you free",
        "de": "An welchen Tagen möchtest du trainieren"
    },
    "WILL_BE_VALID_FOR": {
        "pl": "Będzie on ważny przez",
        "en": "It will be valid for",
        "de": "Es ist gültig für"
    },
    "WITHOUT_CONTACT_DATA": {
        "pl": "Bez danych kontaktowych twoi klienci nie będą mieli jak się z tobą porozumieć w razie pytań lub wątpliwości",
        "en": "Without your contact details, customers won't be able to communicate with you in case they have questions or doubts",
        "de": "Ohne deine Kontaktdaten können die Kunden nicht mit dir kommunizieren, wenn sie Fragen oder Zweifel haben"
    },
    "WITHOUT_PAYOUT_DETAILS": {
        "pl": "Bez danych do wypłat nie będziemy ci w stanie przelać wypłat za rezerwacje",
        "en": "Without payment details, we wont be able to transfer payouts for reservations",
        "de": "Ohne Zahlungsinformationen können wir keine Auszahlungen für Reservierungen überweisen"
    },
    "WITH_WHOM": {
        "pl": "Z kim: ",
        "en": "With whom: ",
        "de": "Mit wem: "
    },
    "WRITE_ABOUT_YOURSELF": {
        "pl": "Napisz coś o sobie",
        "en": "Write something about yourself",
        "de": "Schreib etwas über dich selbst"
    },
    "YES": {
        "pl": "tak",
        "en": "yes",
        "de": "ja"
    },
    "YES_CANCEL_RSV": {
        "pl": "Tak, anuluj rezerwację",
        "en": "Yes, cancel reservation",
        "de": "Ja, sorniere die Reservierung"
    },
    "YOUR_CARNETS": {
        "pl": "Twoje karnety",
        "en": "Your carnets",
        "de": "Deine Abonnements"
    },
    "YOUR_EMAIL": {
        "pl": "Twój email",
        "en": "Your email",
        "de": "Deine E-Mail"
    },
    "YOUR_FUNCTIONALITIES": {
        "pl": "Twoje funkcjonalności",
        "en": "Your functionalities",
        "de": "Deine Funktionalitäten"
    },
    "YOUR_GEAR": {
        "en": "Your gear",
        "pl": "Twój sprzęt",
        "de": "deine Ausrüstung"
    },
    "YOUR_INFO": {
        "pl": "Twoje dane",
        "en": "Your data",
        "de": "Deine Daten"
    },
    "YOUR_LOCATION": {
        "pl": "Twoja lokacja",
        "en": "Your location",
        "de": "Dein Standort"
    },
    "YOUR_NAME": {
        "pl": "Twoja nazwa",
        "en": "Your name",
        "de": "Your name"
    },
    "YOUR_QUESTION": {
        "pl": "Twoje pytanie",
        "en": "Your question",
        "de": "Your question"
    },
    "YOUR_REVIEW": {
        "pl": "Twoja opinia",
        "en": "Your review",
        "de": "Deine Bewertung"
    },
    "YOUR_RSV": {
        "pl": "Twoja rezerwacja",
        "en": "Your reservation",
        "de": "Deine Reservierung"
    },
    "YOUR_RSVS": {
        "pl": "Twoje rezerwacje",
        "en": "Your reservations",
        "de": "Deine Reservierungen"
    },
    "YOUR_RSV_TAB": {
        "pl": "zakładce twoich rezerwacji",
        "en": "Your reservation tab",
        "de": "Registerkarte mit deinen Reservierungen"
    },
    "YOUR_VACATION_WILL_START": {
        "pl": "Twoje wolne się zacznie @ ",
        "en": "Your Vacation will start @ ",
        "de": "Dein Urlabub wird beginnen @"
    },
    "YOU_CAN_SHARE_PROFILE_ANYWHERE": {
        "pl": "Możesz podlinkować swój profil gdziekolwiek chcesz",
        "en": "You can share your profile, wherever you want",
        "de": "Du kannst dein Profil verlinken, wohin du willst"
    },
    "YOU_HAVE_DAYS_UNTIL_RSV_FMT": {
        "pl": "Masz %d dni do rezerwacji",
        "en": "You have %d days until reservation",
        "de": "Du hast %d Tage bis zu deiner Reservierung"
    },
    "YOU_MAY_CLOSE_THIS_WINDOW": {
        "pl": "Możesz teraz zamknąć to okno",
        "en": "You may close this tab now",
        "de": "Du kannst das Fester jetzt schliessen"
    },
    "YOU_MAY_DEACITVATE_INSTR_ACC_": {
        "pl": "Możesz permamentnie deaktywować swoje konto instruktora i usunąć wszystkie swoje dane.",
        "en": "You can permanently deactivate Your instructor account and delete all of Your data.",
        "de": "Du kannst dein Konto des Trainers dauerhaft deaktivieren und alle deine Daten löschen"
    },
    "YOU_MUST_BE_LOGGED_IN_TO_BUY_CARNET": {
        "pl": "Musisz być zalogowany by móc kupić karnet",
        "en": "You must be logged in to buy a carnet",
        "de": "Du musst angemeldet sein, um das Abonnement zu kaufen"
    },
    "YOU_WILL_HAVE_TO_COMPLETE_RSVS": {
        "pl": "Wciąż będziesz musiał zamknąć aktywne rezerwacje",
        "en": "You will still have to complete active reservations",
        "de": "Du musst noch die aktiven Reservierungen abschliessen"
    },
    "YOUR_MANAGEMENT_CENTER": {
        "pl": "Centrum zarządzania trenera",
        "en": "Trainer's management center"
    },
    "MISSION": {
        "pl": "W skrócie",
        "en": "In shortcut"
    },
    "MISSION_DESC": {
        "pl": "Jako trener dodawaj swoje usługi za darmo, zarządzaj rezerwacjami i komunikuj się z klientami w prosty i efektywny sposób! Bez ukrytych kosztów, bez żadnych haczyków - tylko transparentna i korzystna współpraca z Veidly!",
        "en": "As a coach, add your services for free, manage your bookings and communicate with your clients in a simple and effective way! No hidden costs, no catches - just transparent and beneficial cooperation with Veidly!"
    },
    "CHOICE": {
        "pl": "Wybór",
        "en": "Choice",
    },
    "CHOICE_DESC": {
        "pl": "Jesteś szefem swojego czasu i decydujesz, jak i gdzie chcesz pracować. Dzięki Veidly możesz swobodnie dodawać swoje treningi na siłowni, stoku narciarskim czy lokalnej pływalni. Załóż konto już dziś i zacznij zarabiać na swojej pasji jako niezależny trener!",
        "en": "You are the boss of your time and decide how and where you want to work out. With Veidly, you are free to add your workouts at the gym, ski slope or local swimming pool. Create an account today and start earning from your passion as an independent trainer!"
    },
    "PARTNERSHIP": {
        "pl": "Misja",
        "en": "Mission"
    },
    "PARTNERSHIP_DESC": {
        "pl": "Jako pasjonaci sportu i zdrowego stylu życia, chcemy, aby każdy miał dostęp do najlepszych trenerów i usług fitness. Veidly powstało właśnie z tego powodu - aby pomóc Ci w osiągnięciu Twoich celów treningowych. Jesteśmy otwarci na dialog i zawsze słuchamy naszych użytkowników, dlatego zachęcamy do dzielenia się z nami swoimi pomysłami i potrzebami. Razem stworzymy jeszcze lepszą platformę, która spełni oczekiwania każdego trenera i ucznia fitness!",
        "en": "YAs passionate about sports and healthy living, we want everyone to have access to the best trainers and fitness services. Veidly was created for this very reason - to help you achieve your workout goals. We are open to dialogue and always listen to our users, so we encourage you to share your ideas and needs with us. Together we will create an even better platform that meets the needs of every fitness trainer and student!"
    },
    "FEATURES": {
        "pl": "Funkcjonalności",
        "en": "Featues"
    },
    "CUSTOMER_IN_CENTER": {
        "pl": "Dodawanie pojedynczych wejściówek na trening personalny, wejście na siłownię lub dowolną atrakcję jeszcze nigdy nie było takie proste! Dołącz do naszej społeczności już teraz i zwiększ swoją widoczność wśród osób zainteresowanych Twoimi usługami.",
        "en": "Adding single passes for personal training, gym entry or any attraction has never been so easy! Join our community now and increase your visibility among those interested in your services."
    },
    "EXPRESS_TRANSFER": {
        "pl": "Ekspresowy przelew",
        "en": "Express payout",
    },
    "EXPRESS_TRANSFER_MARKETING": {
        "pl": "Przelew w 6h po zakończonym treningu.",
        "en": "Payout in 6h after completion of training."
    },
    "FREE_APP": {
        "pl": "Darmowy portal",
        "en": "Free app"
    },
    "FREE_APP_MARKETING": {
        "pl": "Konto, oferta i narzędzia są darmowe, pobieramy jedynie 5% od transakcji jako koszt utrzymania",
        "en": "The account, offer and tools are free, we only charge 5% per transactions as a maintenance cost."
    },
    "FIND_ACTIVITY_FOR_YOU": {
        "pl": "Znajdź aktywność dla siebie",
        "en": "Find an activity for You"
    },
    "MASSAGE_SAUNA_PHYSIOTHERAPY": {
        "pl": "Masaż, sauna, fizjoterapia?",
        "en": "Massage, sauna, physiotherapy?"
    },
    "SINGLE_ENTRY_OR_CARNET": {
        "pl": "Pojedyncza wejściówka czy karnet?",
        "en": "Single ticket or subscription?"
    },
    "FIND_HERE_EVERYTHING": {
        "pl": "U nas znajdziesz wszystko.",
        "en": "With us you will find everything"
    },
    "PAY_FOR_SERVICE": {
        "pl": "Opłać usługę",
        "en": "Pay for the service"
    },
    "EXPRESS_PAYMENTS_MARKETING": {
        "pl": "Ekspresowe i bezpieczne płatności online.",
        "en": "Express and secure online payments."
    },
    "IN_CASE_OF_TROUBLES": {
        "pl": "W razie kłopotów{<br/>}Jesteśmy tu dla Ciebie.",
        "en": "In case of trouble{<br/>}We are here for you."
    },
    "QR_EXPLAIN_MARKETING": {
        "pl": "Otrzymasz kod QR na maila oraz będzie widoczny",
        "en": "You will receive a QR code to your email and it will be visible"
    },
    "QR_EXPLAIN_MARKETING2": {
        "pl": "na portalu, jeśli założysz konto.",
        "en": "On the portal if you create an account."
    },
    "QR_EXPLAIN_MARKETING3": {
        "pl": "Jeśli nie chcesz zostawiać żadnych danych będziesz mógł go pobrać po opłaceniu wybranej usługi.",
        "en": "If you don't want to leave any data you will be able to download it after paying for the selected service."
    },
    "WE_RESPECT_PRIVACY": {
        "pl": "Szanujemy Twoją prywatność.",
        "en": "We respect your privacy."
    },
    "TRENUJ_WYPOCZYWAJ": {
        "pl": "Trenuj, wypoczywaj i ciesz sie życiem.",
        "en": "Train, relax and enjoy life."
    },
    "WELCOME_TO_VEIDLY": {
        "pl": "Witaj w Veidly",
        "en": "Welcome to Veidly"
    },
    "WELCOME_AGAIN": {
        "pl": "Miło Cię znów gościć",
        "en": "Welcome again"
    },
    "WHAT_ARE_BENEFITS": {
        "pl": "Co zyskujesz umawiając trening przez Veidly?",
        "en": "What do you gain by arranging your training through Veidly?"
    },
    "PROF_PARTNERS": {
        "pl": "Szeroki wybór usług",
        "en": "Variety of services"
    },
    "PROFPARTNERS_SUB_TEXT": {
        "pl": "Veidly łączy profesjonalnych trenerów i ośrodki sportowe ze wszystkimi tymi, którzy czują, że nadszedł czas na zmiany i rozwój.",
        "en": "Veidly connects professional trainers and sports centres with all those who feel it is time for change and development."
    },
    "SIMPLICITY": {
        "pl": "Wygodę",
        "en": "Comfort"
    },
    "SIMPLICITY_SUB_TEXT": {
        "pl": "Po akceptacji przez trenera dostajesz kod QR, który trener zeskanuje jako potwierdzenie rezerwacji - dzięki takiemu rozwiązaniu - wszystkie karnety i wejściówki masz w jednym miejscu.",
        "en": "When you pay your entrance fee, you get a QR code that the trainer scans as confirmation of your reservation - thanks to this solution - you have all your passes and entrance fees in one place."
    },
    "TIME": {
        "pl": "Czas",
        "en": "Time",
    },
    "TIME_SUB_TEXT": {
        "pl": "U nas zawsze znasz cenę przed zakupem, znasz wymagania co do treningu, a w razie wątpliwości możesz napisać używając wbudowanego czatu. W razie problemów, Veidly aktywnie pomoże rozwiązać wszelkie spory.",
        "en": "With us you always know the price before you buy, you know your training requirements and if you have any doubts you can ask using the built-in chat. In case of problems, Veidly will proactively help to resolve any disputes."
    },
    "HOW_TO_START_TRAININGS": {
        "pl": "Jak rozpocząć treningi?",
        "en": "How to start training?"
    },
    "PICK_DATE_AND_PLACE": {
        "pl": "Wybierz odpowiednie miejsce i termin",
        "en": "Choose the right place and date"
    },
    "FAST_AND_SAVE": {
        "pl": "bezpiecznie i szybko",
        "en": "fast and save"
    },
    "QR_SHOW_IT_BRO": {
        "pl": "O umówionej porze udaj się na trening i pokaż trenerowi kod QR",
        "en": "Go to training at the appointed time and show the trainer the QR code"
    },
    "CHECK_AMOUNT_OF_BENEFITS": {
        "pl": "Odkryj wszystkie korzyści, jakie możesz uzyskać, decydując się na współpracę z nami",
        "en": "Discover all the benefits you can get by choosing to work with us"
    },
    "CHECK_HOW_EASY_VEIDLY_IS": {
        "pl": "Przekonaj się, jak łatwo i przyjemnie możesz korzystać z Veidly - platformy stworzonej z myślą o Twojej wygodzie!",
        "en": "Discover the benefits of working with us on Veidly - an easy and enjoyable platform designed for your convenience."
    },
    "GIVE_ME_ADDRESS": {
        "pl": "Zacznij wpisywać, aby wyszukać miejsce",
        "en": "Start typing to find place"
    },
    "MARKETING_DESCRIPTION_MAIN_PAGE": {
        "pl":   `Platforma do wyszukiwania treningów i zajęć sportowych Veidly to innowacyjna aplikacja webowa dzięki której odkryjesz nowy wymiar świata sportu i wellness.
        Veidly zapewnia dostęp do kompleksowych usług branży sportowej, terapii masażu, spa i fizjoterapii. 
        
        Co więcej, ta niesamowita platforma treningowa jest dostępna całkowicie za darmo - skorzystaj z niej już dziś
        
        Nasza aplikacja oferuje wsparcie dla osób niepełnosprawnych, intuicyjny i prosty kalendarz oraz wiele narzędzi zarówno dla profesjonalnych trenerów jak i użytkowników. 
        
        Chcesz śledzić swoje postępy? Motywuje Cię wyznaczanie celów treningowych? Dbasz o swoje zdrowie i dobre samopoczucie? A może chcesz zacząć swój pierwszy trening? 
        Veidly ma wszystko, czego potrzebujesz, w jednym miejscu. 
        
        W Veidly wierzymy, że każdy powinien mieć dostęp do usług z zakresu sportu i wellness najwyższej jakości. Dlatego polegamy na wsparciu naszej społeczności poprzez dobrowolne datki, aby utrzymać i ciągle rozwijać naszą platformę. Razem tworzymy idealne miejsce dla miłośników aktywnego stylu życia!`,
        "en":   `Veidly's sports training and activities search platform is an innovative web application through which you will discover a new dimension in the world of sports and wellness.
        Veidly provides access to comprehensive services of the sports industry, massage therapy, spa and physiotherapy. 
        
        What's more, this amazing training platform is available completely free - take advantage of it today
        
        Our app offers support for people with disabilities, an intuitive and simple calendar and many tools for both professional trainers and users. 
        
        Do you want to track your progress? Are you motivated by setting training goals? Do you care about your health and well-being? Or maybe you want to start your first workout? 
        Veidly has everything you need in one place. 
        
        At Veidly, we believe that everyone should have access to top-quality sports and wellness services. That's why we rely on the support of our community through voluntary donations to maintain and continuously develop our platform. Together we are creating the perfect place for active lifestyle lovers!
        `
    },
    "USER_BENEFITS_CONTACT_WITH_INTSRUCTOR": {
        "pl": "Zapisz się na trening i dograj szczegóły z instruktorem na naszym czacie",
        "en": "Sign up for a training session and catch up on the details with the instructor in our chat room."
    },
    "SUPPORT": {
        "pl": "Dział pomocy",
        "en": "Support"
    },
    "SUPPORT_DESC": {
        "pl": "Znajdziesz tu dokumenty, regulaminy oraz wszystkie sposoby kontaktu z nami",
        "en": "You'll find documents, regulations and all the ways to contact us."
    },
    "SUPPORT_VEIDLY": {
        "pl": "Wesprzyj Veidly",
        "en": "Support Veidly"
    },
    "FREE_TO_USE_APP": {
        "pl": "Nasz portal jest darmowy, dzięki Waszym wpłatom, wesprzyj, aby taki pozostał",
        "en": "Our portal is free, thanks to your contributions, support to keep it that way",
    },
    "SEO_FOOTER": {
        "pl": "Zacznij swoją przygodę ze zdrowym stylem życia dzięki platformie treningowej Veidly. Oferujemy usługi z zakresu sportu, wellnes&SPA, masażu oraz fizjoterapii, aby pomóc Ci osiągnąć swoje cele treningowe. Z naszą bogatą bazą trenerów personalnych na Śląsku i intuicyjną wyszukiwarką, łatwo znajdziesz idealny trening dla siebie.",
        "en": "Start your adventure into a healthy lifestyle with Veidly's workout platform. We offer sports, wellness&SPA, massage and physiotherapy services to help you achieve your workout goals. With our extensive database of personal trainers in Silesia and intuitive search engine, it's easy to find the perfect workout for you."
    },
    "SUPPORT": {
        "pl": "Wesprzyj!",
        "en": "Donate!"
    },
    "BROKEN_EMAIL": {
        "pl": "Niepoprawny email!",
        "en": "Invalid email!"
    },
    "THANKS": {
        "pl": "Dziękujemy",
        "en": "Thank You!"
    },
    "DONATION_WORKS": {
        "pl": "Zarejestrowaliśmy wpłatę! Jesteście najlepsi!",
        "en": "Donation registered! You're awesome!"
    }

}

function sortObjectKeys(o) {
    var sorted = {},
    key, a = [];

    for (key in o) {
        if (o.hasOwnProperty(key)) {
                a.push(key);
        }
    }

    a.sort();

    for (key = 0; key < a.length; key++) {
        sorted[a[key]] = o[a[key]];
    }
    return sorted;
}
// keep alphabetical order
// console.log(JSON.stringify(sortObjectKeys(locale2)))
