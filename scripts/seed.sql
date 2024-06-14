BEGIN TRANSACTION;

INSERT OR IGNORE INTO users (id, username, password_hash, email, phone, first_name, last_name)
VALUES (1, 'admin', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'admin@email.com', '123123123', 'Adminsław', 'Userowicz'),
        (2, 'user', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'user@email.com', '678678678', 'Mikołaj', 'Święty'),
        (3, 'user2', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'asd@email.com', '345345890', 'Patryk', 'Pogorzelczyk'),
        (4, 'user3', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'qwe@gmail.com', '101020999', 'Jan', 'Kowalski'),
        (5, 'user4', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'xyz@email.com', '123456789', 'Przemysław', 'Nawrot'),
        (6, 'user5', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'patryk@email.com', '987654321', 'Patryk', 'Nowak'),
        (7, 'user6', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'adam@email.com', '123456789', 'Adam', 'Nowak'),
        (8, 'user7', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'michal@email.com', '987654321', 'Michał', 'Iksinski'),
        (9, 'user8', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKkdhO', 'rysiek@email.com', '123456789', 'Ryszard', 'Kowalski'),
        (10, 'user9', '$2a$14$HQKt9/9GsMV/PXlYrnWmveP6exu6QV8r8MyUhViLL22YPysaKdhO', 'adrian@email.com', '987654321', 'Adrian', 'Nowak');


INSERT OR IGNORE INTO lists (id, owner_id, title, description, event_date)
VALUES (1, 1, 'Impreza urodzinowa', 'Organizuję imprezę z okazji moich urodzin', '2020-01-01');

INSERT OR IGNORE INTO list_items (id, list_id, title, description, price, link, photo)
VALUES (1, 1, 'Ciasto', 'Pyszne ciasto', 200.0, 'https://wiejskiewypieki.pl/pl/products/ciasto-jogurtowe-z-owocami-i-serem-400g-73.html', 'https://assets.tmecosys.com/image/upload/t_web767x639/img/recipe/ras/Assets/14161EA0-3432-4056-9204-8B846A8BAD42/Derivates/9d9ab322-9ac3-4c9d-8113-97f58696edef.jpg'),
        (2, 1, 'Skarpetki do biegania', 'Pakiet kolorowych skarpetek', 50.0, 'https://example.com', ''),
        (3, 1, 'Lampka do rowera', 'Taka mocna, nie byle jaka', 60.0, 'https://pancernik.eu/akcesoria-rowerowe/lampka-do-roweru/lampka-rowerowa-przednia-rockbros-260lm-biale-swiatlo-czarna?cd=17593697947&ad=&kd=', 'https://pancernik.eu/environment/cache/images/500_500_productGfx_371470/rockbros-lampka-rowerowa-przednia-260lm-black-15.webp'),
        (4, 1, 'Bilety na festival', 'Bilety na festival muzyczny', 100.0, 'https://www.ebilet.pl/muzyka/festiwale/before-festival', ''),
        (5, 1, 'Piłka do kopania', 'Piłka nozna oficjalna EURO 2024', 500.0, 'https://www.ambersport.pl/pilka-nozna-adidas-euro24-fussballliebe-pro-iq3682', 'https://www.ambersport.pl/upload/220728-1700068915.jpg'),
        (6, 1, 'Playstation 5', 'Nie wiem czy uda nam się złozyć, wpiszcie ile jesteście w stanie dać i jak dobijemy to olać resztę prezentów', 2500.0, 'https://www.mediaexpert.pl/gaming/playstation-5/konsole-ps5/konsola-sony-playstation-5-digital-slim', 'https://prod-api.mediaexpert.pl/api/images/gallery/thumbnails/images/60/6007768/Konsola-SONY-PlayStation-5-Digital-Slim-skos.jpg'),
        (7, 1, 'Taniec', 'Lekcja tańca', 100.0, 'https://example.com', ''),
        (8, 1, 'Gra planszowa Monopoly', 'Monopoly edycja specjalna', 337.0, 'https://example.com', ''),
        (9, 1, 'Gry', 'Pakiet gier na ps5', 599.0, 'https://example.com', ''),
        (10, 1, 'Niespodzianka', 'Nie mam pojęcia co', 1234.0, 'https://example.com', '');

INSERT OR IGNORE INTO list_assignments (list_id, user_id)
VALUES (1, 2),
        (1, 3),
        (1, 4);

INSERT OR IGNORE INTO list_item_assignments (list_item_id, user_id, amount)
VALUES (1, 4, 200.0),
        (3, 2, 40.0),
        (4, 3, 30.0),
        (5, 4, 70.0),
        (6,3, 300.0),
        (6, 4, 100.0),
        (6, 5, 400.0),
        (6, 6, 50.0),
        (6, 7, 100.0),
        (6, 8, 100.0),
        (6, 9, 100.0),
        (6, 10, 100.0),
        (7, 4, 30.0),
        (7, 3, 40.0),
        (8, 4, 213.0);


INSERT OR IGNORE INTO lists (id, owner_id, title, description, event_date)
VALUES (2, 2, 'Urodziny Bronisława', 'Bronek ma urodzinki i podał mi listę prezentów które chciałby dostać', '2025-01-01');

INSERT OR IGNORE INTO list_items (id, list_id, title, description, price, link, photo)
VALUES (11, 2, 'Czapka', 'Czapka na zimę', 100.0, 'https://kabak.com.pl/pl/czapka-zimowa-z-bawelny-organicznej-jasny-roz', 'https://www.lego.com/cdn/cs/set/assets/blt14d2306147e2a6f3/10333_alt1.png?format=webply&fit=bounds&quality=60&width=1200&height=1200&dpr=2'),
        (12, 2, 'Samochodzik', 'Taki przypominający jego aktualny samochód', 88.0, 'https://www.amazon.pl/Friki-Monkey-Hot-Wheels-Mercedes-Benz/dp/B0C1K72YCS', ''),
        (13, 2, 'Lego', 'Lego z Władców pierścieni, jakieś ogromne', 2000.0, 'https://www.lego.com/pl-pl/product/the-lord-of-the-rings-barad-dur-10333', ''),
        (14, 2, 'Mop', 'Mop dla Jadzi', 10.0, 'https://example.com', ''),
        (15, 2, 'Wiadro', 'Wiadro dla Jadzi', 70.0, 'https://example.com', ''),
        (16, 2, 'Kubek do kawy', 'Kubek na trudne poranki', 7.99, 'https://example.com', ''),
        (17, 2, 'Skuter', 'Skuter', 1000.0, 'https://example.com', ''),
        (18, 2, 'Kolekcjonerska książka', 'Kolekcja Władców Pierścieni', 337.0, 'https://www.empik.com/wladca-pierscieni-tomy-1-3-tolkien-john-ronald-reuel,p1360315368,ksiazka-p', 'https://ecsmedia.pl/cdn-cgi/image/format=webp,width=544,height=544,/c/wladca-pierscieni-tomy-1-3-b-iext151541671.jpg');

INSERT OR IGNORE INTO list_assignments (list_id, user_id)
VALUES (2, 1),
        (2, 3);

INSERT OR IGNORE INTO list_item_assignments (list_item_id, user_id, amount)
VALUES (11, 1, 100.0),
        (12, 3, 50.0),
        (13, 1, 300.0),
        (15, 1, 70.0),
        (16, 3, 7.99);

COMMIT;
