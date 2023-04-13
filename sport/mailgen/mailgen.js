const { off } = require('process')
var fs = require('fs')
var Mailgen = require('./mailgen/index')

const templates_directory = "email_templates"
const pl_greeting = "Cześć"
const pl_signature = "Do zobaczenie na platformie"
const pl_outro = [
    'W razie jakichkolwiek problemów, prosimy skontaktowac się z supportem pod adresem: <h3><a href="mailto:support@veidly.com">support@veidly.com</a></h3>',
    'lub napisz do nas na czacie: <h3><a href="https://m.me/veidly"><svg style="color: rgb(36, 192, 132); width: 15px; --darkreader-inline-color:#85c9f0;" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" data-darkreader-inline-color=""><path d="M256.55 8C116.52 8 8 110.34 8 248.57c0 72.3 29.71 134.78 78.07 177.94 8.35 7.51 6.63 11.86 8.05 58.23A19.92 19.92 0 0 0 122 502.31c52.91-23.3 53.59-25.14 62.56-22.7C337.85 521.8 504 423.7 504 248.57 504 110.34 396.59 8 256.55 8zm149.24 185.13l-73 115.57a37.37 37.37 0 0 1-53.91 9.93l-58.08-43.47a15 15 0 0 0-18 0l-78.37 59.44c-10.46 7.93-24.16-4.6-17.11-15.67l73-115.57a37.36 37.36 0 0 1 53.91-9.93l58.06 43.46a15 15 0 0 0 18 0l78.41-59.38c10.44-7.98 24.14 4.54 17.09 15.62z" fill="#24c084" data-darkreader-inline-fill="" style="--darkreader-inline-fill:#0a3b57;"></path></svg>Messenger</a></h3>',
]

const en_greeting = "Hi"
const en_signature = "See You on platform"
const en_outro = [
    'In case of any troubles, please write to: <h3><a href="mailto:support@veidly.com">support@veidly.com</a></h3>',
    'or write to us on chat: <h3><a href="https://m.me/veidly"><svg style="color: rgb(36, 192, 132); width: 15px; --darkreader-inline-color:#85c9f0;" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" data-darkreader-inline-color=""><path d="M256.55 8C116.52 8 8 110.34 8 248.57c0 72.3 29.71 134.78 78.07 177.94 8.35 7.51 6.63 11.86 8.05 58.23A19.92 19.92 0 0 0 122 502.31c52.91-23.3 53.59-25.14 62.56-22.7C337.85 521.8 504 423.7 504 248.57 504 110.34 396.59 8 256.55 8zm149.24 185.13l-73 115.57a37.37 37.37 0 0 1-53.91 9.93l-58.08-43.47a15 15 0 0 0-18 0l-78.37 59.44c-10.46 7.93-24.16-4.6-17.11-15.67l73-115.57a37.36 37.36 0 0 1 53.91-9.93l58.06 43.46a15 15 0 0 0 18 0l78.41-59.38c10.44-7.98 24.14 4.54 17.09 15.62z" fill="#24c084" data-darkreader-inline-fill="" style="--darkreader-inline-fill:#0a3b57;"></path></svg>Messenger</a></h3>',
]
var mailGenerator = new Mailgen({
    theme: 'cerberus',
    product: {
        name: "<a href='https://veidly.com'>Veidly.com</a>",
        link: 'https://veidly.com',
        logo: 'https://veidly.com/logo-light.png'
    },
})

var templates = {
    user: {
        "pl.register.html": {
            body: {
                name: '{{.Name}}',
                intro: 'Witaj w Veidly.com, aby dokończyć rejestrację kliknij w poniższy link.',
                action: {
                    instructions: 'Link aktywacyjny:',
                    button: {
                        color: '#24C084',
                        text: 'Potwierdź konto',
                        link: '{{.Url}}'
                    }
                },
            }
        },
        "pl.forgot_pass.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Dostaliśmy infromację, że poproszono o reset hasła dla tego adresu email.',
                    'Jeśli to nie Ty, zignoruj tą wiadomość.'
                ],
                action: {
                    instructions: 'Link do resetu hasła:',
                    button: {
                        color: '#24C084',
                        text: 'Resetuj hasło',
                        link: '{{.Url}}'
                    }
                },
            }
        },
        "en.register.html": {
            body: {
                name: '{{.Name}}',
                intro: 'Welcome in Veidly.com, to finish registration process, click below link.',
                action: {
                    instructions: 'Activation link:',
                    button: {
                        color: '#24C084',
                        text: 'Confirm registration',
                        link: '{{.Url}}'
                    }
                },
            }
        },
        "en.forgot_pass.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We received information that You forgot your password',
                    'If that\s not You, ignore this email.'
                ],
                action: {
                    instructions: 'Click below to recover password:',
                    button: {
                        color: '#24C084',
                        text: 'Recover password',
                        link: '{{.Url}}'
                    }
                },
            }
        }
    },
    sub: {
        "en.sub_cancel.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Subscription {{.Name}} has been canceled',
                ],
                action: {
                    instructions: ' Click below to see details of this subsctiption:',
                    button: {
                        color: '#24C084',
                        text: 'Your subscription',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "en.sub_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Your subcription {{.SubName}} has been payed off',
                ],
                action: {
                    instructions: 'To see details of this subscription click below',
                    button: {
                        color: '#24C084',
                        text: 'Your subscription',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "en.sub_dispute.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We received a ticket about issues with subscription',
                    '{{.SubName}}',

                ],
                action: {
                    instructions: 'Our team is trying to solve the problem, and we will contact you soon!',
                    button: {
                        color: '#24C084',
                        text: 'Your subscription',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "en.sub_fail_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We couldn\'t finalize your payment for subscription:',
                    '{{.SubName}}',
                    'We will keep retrying for a while but if problem wont be solved then we will cancel this transaction'

                ],
                action: {
                    instructions: 'Click below to see status of your subscription',
                    button: {
                        color: '#24C084',
                        text: 'Your subscription',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "en.sub_fail_payout.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We cannot transfer your payout for subscription:',
                    '{{.SubName}}',
                ],
                action: [
                    {
                        instructions: 'Click below to see status of your subscription',
                        button: {
                            color: '#24C084',
                            text: 'Your subscription',
                            link: '{{.SubUrl}}'
                        },
                    }, {
                        instructions: 'But maybe your payout data is not configured properly?',
                        button: {
                            color: '#24C084',
                            text: 'Check Your configuration',
                            link: '{{.PaymentConfig}}',
                        },
                    },
                ]
            },
        },
        "en.sub_link.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Thanks for creating veidly subscription:',
                    '{{.SubName}}',
                ],
                action: [
                    {
                        instructions: 'Click below to see details of your subscription:',
                        button: {
                            color: '#24C084',
                            text: 'Your subscription',
                            link: '{{.SubUrl}}'
                        },
                    }, {
                        instructions: 'Check Your payment:',
                        button: {
                            color: '#24C084',
                            text: 'Check Your payment',
                            link: '{{.PaymentUrl}}',
                        },
                    },
                ]
            },
        },
        "pl.sub_cancel.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Twój karnet został anulowany:',
                    '{{.SubName}}',

                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale karnetu',
                    button: {
                        color: '#24C084',
                        text: 'Karnet',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "pl.sub_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Twój karnet został opłacony:',
                    '{{.SubName}}',

                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale karnetu:',
                    button: {
                        color: '#24C084',
                        text: 'Karnet',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "pl.sub_dispute.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Otrzymaliśmy zgłoszenie, że coś jest nie tak z twoim karnetem:',
                    '{{.SubName}}',
                    'Nasz zespół już stara się rozwiązać problem, i wkrótce się z tobą skontaktuje'

                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale karnetu:',
                    button: {
                        color: '#24C084',
                        text: 'Karnet',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "pl.sub_fail_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Nie udało się pobrać płatności za karnet:',
                    '{{.SubName}}',
                    'Będziemy próbowali dalej, ale jeżeli problem nie zostanie rozwiązany anulujemy tą transakcje'
                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale karnetu:',
                    button: {
                        color: '#24C084',
                        text: 'Karnet',
                        link: '{{.SubUrl}}'
                    }
                },
            }
        },
        "pl.sub_fail_payout.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Nie potrafimy przelać ci środków na konto za karnet:',
                    '{{.SubName}}',
                ],
                action: [
                    {
                        instructions: 'Kliknij poniżej by zobaczyć detale karnetu',
                        button: {
                            color: '#24C084',
                            text: 'Karnet',
                            link: '{{.SubUrl}}'
                        },
                    }, {
                        instructions: 'Być może nie masz poprawnie skonfigurowanych danych do przelewu?',
                        button: {
                            color: '#24C084',
                            text: 'Skonfiguruj dane do przelewów',
                            link: '{{.PaymentConfig}}',
                        },
                    },
                ]
            },
        },
        "pl.sub_link.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Dziękujemy za kupno karnetu w Veidly:',
                    '{{.SubName}}',
                ],
                action: [
                    {
                        instructions: 'Kliknij poniżej by zobaczyć detale karnetu',
                        button: {
                            color: '#24C084',
                            text: 'Karnet',
                            link: '{{.SubUrl}}'
                        },
                    }, {
                        instructions: 'Kliknij tu by przejść do swojej transakcji',
                        button: {
                            color: '#24C084',
                            text: 'Status płatności',
                            link: '{{.PaymentUrl}}',
                        },
                    },
                ]
            },
        },
    },
    rsv: {
        "en.rsv_cancel.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Your reservation for training:',
                    '{{.Training}}',
                    'has been cancelled',
                ],
                action: {
                    instructions: 'Click below to see status of your reservation:',
                    button: {
                        color: '#24C084',
                        text: 'Your reservation',
                        link: '{{.RsvUrl}}'
                    }
                },
            },
        },
        "en.rsv_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Your training has been paid for!',
                    '{{.Training}}',
                ],
                action: {
                    instructions: 'To see status of your reservation or generate QR code click below:',
                    button: {
                        color: '#24C084',
                        text: 'Your reservation',
                        link: '{{.RsvUrl}}'
                    }
                },
            },
        },
        "en.rsv_dispute.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We received ticket about issues with reservation for training:',
                    '{{.Training}}',
                    'Our team is trying to solve the problem, and we will contact you soon!',
                ],
                action: {
                    instructions: 'Click below to see status of your reservation:',
                    button: {
                        color: '#24C084',
                        text: 'Your reservation',
                        link: '{{.RsvUrl}}'
                    }
                },
            },
        },
        "en.rsv_fail_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We couldnt finalize your payment for training:',
                    '{{.Training}}',
                    'We will keep retrying for a while but if problem wont be solved then we will cancel this reservation.',
                ],
                action: {
                    instructions: 'Click below to see status of your reservation:',
                    button: {
                        color: '#24C084',
                        text: 'Your reservation',
                        link: '{{.RsvUrl}}'
                    }
                },
            },
        },
        "en.rsv_fail_payout.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'We cannot transfer your payout for training:',
                    '{{.Training}}',
                ],
                action: [
                    {
                        instructions: 'Click below to see status of your training',
                        button: {
                            color: '#24C084',
                            text: 'Your training',
                            link: '{{.RsvUrl}}'
                        },
                    }, {
                        instructions: 'But maybe your payout data is not configured properly?',
                        button: {
                            color: '#24C084',
                            text: 'Check Your configuration',
                            link: '{{.PaymentConfig}}',
                        },
                    },
                ]
            },
        },
        "en.rsv_hold.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Your reservation for training:',
                    '{{.Training}}',
                    'For day:',
                    '{{.RsvDate}}',
                    'Has been paid. For now your money is put on hold, we will capture it:',
                    '{{.CaptureDate}}'
                ],
                action: {
                    instructions: 'Click below to see status of your reservation:',
                    button: {
                        color: '#24C084',
                        text: 'Your reservation',
                        link: '{{.RsvUrl}}'
                    }
                },
            },
        },
        "en.rsv_link.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Thanks for making reservation via Veidly for training:',
                    '{{.Training}}',
                    'For day:',
                    '{{.RsvDate}}',
                ],
                action: [
                    {
                        instructions: 'Click below to see details of your reservation:',
                        button: {
                            color: '#24C084',
                            text: 'Your reservation',
                            link: '{{.RsvUrl}}'
                        },
                    }, {
                        instructions: 'Check Your reservation:',
                        button: {
                            color: '#24C084',
                            text: 'Check Your reservation',
                            link: '{{.PaymentUrl}}',
                        },
                    },
                ]
            },
        },
        "pl.rsv_cancel.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Twoja rezerwacja:',
                    '{{.DateStart}}',
                    'Na trening:',
                    '{{.Training}}',
                    'została anulowana.'

                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale tej rezerwacji:',
                    button: {
                        color: '#24C084',
                        text: 'Rezerwacja',
                        link: '{{.RsvUrl}}'
                    }
                },
            }
        },
        "pl.rsv_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Twoja rezerwacja:',
                    '{{.Training}}',
                    'Na dzień:',
                    '{{.TrainingDate}}',
                    'została opłacona.'

                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale tej rezerwacji:',
                    button: {
                        color: '#24C084',
                        text: 'Twoja rezerwacja',
                        link: '{{.RsvUrl}}'
                    }
                },
            }
        },
        "pl.rsv_dispute.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Otrzymaliśmy zgłoszenie, że coś jest nie tak z Twoją rezerwacją:',
                    '{{.Training}}',
                    'Nasz zespół już stara się rozwiązać problem, i wkrótce się z tobą skontaktuje'

                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale rezerwacji:',
                    button: {
                        color: '#24C084',
                        text: 'Twoja rezerwacja',
                        link: '{{.RsvUrl}}'
                    }
                },
            }
        },
        "pl.rsv_fail_capture.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Nie udało się pobrać płatności za rezerwację:',
                    '{{.Training}}',
                    'Będziemy próbowali dalej, ale jeżeli problem nie zostanie rozwiązany anulujemy tą transakcje'
                ],
                action: {
                    instructions: 'Kliknij poniżej by zobaczyć detale rezerwacji:',
                    button: {
                        color: '#24C084',
                        text: 'Twoja rezerwacja',
                        link: '{{.RsvUrl}}'
                    }
                },
            }
        },
        "pl.rsv_fail_payout.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Nie potrafimy przelać ci środków na konto za rezerwację:',
                    '{{.Training}}',
                ],
                action: [
                    {
                        instructions: 'Kliknij poniżej by zobaczyć detale rezerwacji',
                        button: {
                            color: '#24C084',
                            text: 'Karnet',
                            link: '{{.SubUrl}}'
                        },
                    }, {
                        instructions: 'Być może nie masz poprawnie skonfigurowanych danych do przelewu?',
                        button: {
                            color: '#24C084',
                            text: 'Skonfiguruj dane do przelewów',
                            link: '{{.RsvUrl}}',
                        },
                    },
                ]
            },
        },
        "pl.rsv_hold.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Dostaliśmy potwierdzenie o wykonanej transakcji na rezerwację:',
                    '{{.Training}}',
                    'Na dzień:',
                    '{{.RsvDate}}',
                    'Na razie środki są pod blokadą, pobierzemy je dopiero:',
                    '{{.CaptureDate}}'
                ],
                action: {
                    instructions: ' Kliknij poniżej by zobaczyć status swojej rezerwacji:',
                    button: {
                        color: '#24C084',
                        text: 'Twoja rezerwacja',
                        link: '{{.RsvUrl}}'
                    }
                },
            },
        },
        "pl.rsv_link.html": {
            body: {
                name: '{{.Name}}',
                intro: [
                    'Dziękujemy za zrobienie rezerwacji na trening:',
                    '{{.Training}}',
                    'Na dzień:',
                    '{{.RsvDate}}',
                ],
                action: [
                    {
                        instructions: 'Kliknij poniżej by zobaczyć status swojej rezerwacji:',
                        button: {
                            color: '#24C084',
                            text: 'Twoja rezerwacja',
                            link: '{{.RsvUrl}}'
                        },
                    }, {
                        instructions: 'Sprawdź status swojej rezerwacji:',
                        button: {
                            color: '#24C084',
                            text: 'Sprawdź rezerwację',
                            link: '{{.PaymentUrl}}',
                        },
                    },
                ]
            },
        },
    },
    chat: {
        "en.new_chatroom.html": {
            body: {
                name: '{{.Target}}',
                intro: ['You created a new chat channel.', '{{ .ChatName }}'],
                action: {
                    instructions: 'Start Your conversation using below link:',
                    button: {
                        color: '#24C084',
                        text: 'Chat',
                        link: '{{.ChatUrl}}'
                    }
                },
            }
        },
        "en.unread_messages.html": {
            body: {
                name: '{{.Target}}',
                intro: [
                    'You have new unread messages in channel:',
                    '{{ .ChatName }}',
                    'Link will be valid to: {{ .ValidTo }}'
                ],
                action: {
                    instructions: 'Link to chat:',
                    button: {
                        color: '#24C084',
                        text: 'Chat',
                        link: '{{.ChatUrl}}'
                    }
                },
                outro: [
                    'Last messages:',
                    '{{ range $x := .Notifications }}',
                    '<strong>{{ .Author }}</strong> {{ .TimestampStr }}: {{ .MsgSummary }}',
                    '{{ end }}',
                ]
            }
        },
        "pl.new_chatroom.html": {
            body: {
                name: '{{.Target}}',
                intro: ['Utworzyłeś nowy kanał:', '{{ .ChatName }}'],
                action: {
                    instructions: 'Rozpocznij rozmowę używając poniższego linku:',
                    button: {
                        color: '#24C084',
                        text: 'Chat',
                        link: '{{.ChatUrl}}'
                    }
                },
            }
        },
        "pl.unread_messages.html": {
            body: {
                name: '{{.Target}}',
                intro: [
                    'Masz nowe nieprzeczytane wiadomości na kanale:',
                    '{{ .ChatName }}',
                    'Link będzie ważny do: {{ .ValidTo }}'
                ],
                action: {
                    instructions: 'Link do chatu:',
                    button: {
                        color: '#24C084',
                        text: 'Chat',
                        link: '{{.ChatUrl}}'
                    }
                },
                outro: [
                    'Ostatnie wiadomości',
                    '{{ range $x := .Notifications }}',
                    '<strong>{{ .Author }}</strong> {{ .TimestampStr }}: {{ .MsgSummary }}',
                    '{{ end }}',
                ]
            }
        },
    },
    marketing: {
        "pl.zapraszamy_na_portal.html": {
            body: {
                name: 'trenerze',
                intro: [
                    "Mam przyjemność zaprosić Cię do dołączenia do naszej dynamicznej i innowacyjnej platformy - <a href='https://veidly.com'>Veidly.com</a>! Jesteśmy przekonani, że to miejsce, w którym Twoje umiejętności jako trenera, masażysty lub fizjoterapeuty zostaną docenione i przyczynią się do rozwoju zdrowego stylu życia naszych użytkowników.",
                    "Dlaczego warto dołączyć do Veidly.com?",
                    "🚀 Szeroki zasięg: Nasza platforma gromadzi klientów z różnych regionów, zainteresowanych usługami specjalistów takich jak Ty.Dzięki temu zyskasz możliwość współpracy z nowymi klientami i rozwoju swojej działalności.",
                    "🌟 Profesjonalne portfolio: Twórz swoją stronę zawierającą opis Twojej oferty, zdjęcia, certyfikaty, referencje oraz linki do mediów społecznościowych.W ten sposób klienci będą mogli lepiej Cię poznać i zdecydować się na Twoje usługi.",
                    "📅 Elastyczne zarządzanie terminami: Nasz system rezerwacji umożliwia łatwe zarządzanie dostępnością oraz umawianie spotkań z klientami zgodnie z Twoim harmonogramem.",
                    "💬 Komunikacja z klientami: Bezpośredni kontakt z klientami poprzez wbudowany chat ułatwi ustalanie szczegółów współpracy, a nasz system ocen i recenzji pozwoli Ci pozyskiwać nowych zainteresowanych Twoimi usługami.",
                    "🛡 Bezpieczeństwo: Gwarantujemy bezpieczne transakcje oraz dbamy o prywatność Twoich danych, dzięki czemu możesz się skupić na swojej pracy i rozwijaniu kariery.",
                    "🧑‍🦼 Jako jedyni na rynku wprowadzamy wsparcie dla niepełnosprawnych",
                    "💵 Nie potrzebujesz działalności gospodarczej, aby zacząć trenować innych - nasz portal jest robiony z myślą o każdym pasjonacie, także tych, którzy chcą \"dorobić\" na etacie!",
                    "Jak dołączyć do <a href='https://veidly.com'>Veidly.com</a>?",

                    "Zarejestruj się na stronie <a href='https://veidly.com/register'>Veidly.com/register</a>",
                    "Wypełnij swój profil, dodając zdjęcia, certyfikaty oraz opis swojej oferty.",
                    "Ustaw dostępność w swoim kalendarzu.",
                    "Gotowe! Możesz teraz oczekiwać na kontakt od zainteresowanych klientów.",
                    "Nie przegap tej okazji! Dołącz do naszej społeczności już dziś i wykorzystaj pełen potencjał Veidly.com.",
                    "Razem pokażemy, jak ważny jest zdrowy styl życia!",

                ]
            }
        }
    }
}

for (const x of Object.keys(templates)) {
    // concatenate and create directory based on keys in templates object
    const directory = templates_directory + "/" + x
    fs.promises.mkdir(directory, { recursive: true }).catch(console.error)

    for (let index = 0; index < Object.keys(templates[x]).length; index++) {
        // object containing email variables
        const element = Object.values(templates[x])[index];
        // inject proper locale into template, so it's easier to maintain in the future => single source of truth
        if (/^pl/.test(Object.keys(templates[x])[index])) {
            element.body.greeting = pl_greeting
            // if someone already declared outro, just merge locale_outro into that
            if (element.body.outro) {
                let merged = [].concat(element.body.outro, pl_outro);
                element.body.outro = merged
            } else {
                element.body.outro = pl_outro
            }
            element.body.signature = pl_signature
        } else {
            element.body.greeting = en_greeting
            // if someone already declared outro, just merge locale_outro into that
            if (element.body.outro) {
                let merged = [].concat(element.body.outro, en_outro);
                element.body.outro = merged
            } else {
                element.body.outro = en_outro
            }
            element.body.signature = en_signature
        }
        // path to file where we write HTML
        const path = directory + "/" + Object.keys(templates[x])[index]
        // generate() generates object necessary for writeFileSync based on our element
        var emailBody = mailGenerator.generate(element);
        fs.writeFileSync(path, emailBody, 'utf8')
    }
}
