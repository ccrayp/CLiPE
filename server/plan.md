# План создания графического интерфейса для CLiPE

## 1. Что уже есть в проекте

Серверная часть уже предоставляет API для централизованного управления доступом и аудита:

- `POST /auth/login` — вход по `username/password`, возвращает `access_token` и `refresh_token`
- `POST /auth/refresh` — обновление access token по refresh token
- `POST /auth/logout` — выход
- `POST /users/search` — пользователи
- `POST /hosts/search` — хосты
- `POST /services/search` — сервисы
- `POST /rules/search` — правила
- `POST /policies/search` — политики
- `POST /policy-contents/search` — связи политика ↔ сервис ↔ правило
- `POST /requests/search` — журнал запросов
- `POST /decisions/search` — журнал решений
- `POST /aggregator` — служебный API для получения итоговой политики/правила по паре `user_name + service_name`

Фронтенд-заготовка уже есть в папке `/Users/roman/projects/CLiPE/server/panel`, сейчас это стандартный шаблон Vite без прикладного UI.

## 2. Сущности, которые должны быть отражены в UI

### Основные справочники

- `hosts`
  - поля: `host_id`, `ip`
- `users`
  - поля: `user_id`, `user_name`, `uid`, `gid`, `host_id`
- `services`
  - поля: `service_id`, `service_name`

### Управление доступом

- `rules`
  - поля: `rule_id`, `rule_name`, `conditions`, `effect`
  - `conditions` хранится в JSON
- `policies`
  - поля: `policy_id`, `policy_name`, `user_id`, `status`
- `policy_contents`
  - поля: `policy_id`, `service_id`, `rule_id`
  - фактически это матрица привязки сервиса и правила к политике

### Аудит и принятие решений

- `requests`
  - поля: `request_id`, `user_id`, `context`, `timestamp`
  - `context` хранится в JSONB
- `decisions`
  - поля: `decision_id`, `request_id`, `policy_id`, `result`, `timestamp`

### Аутентификация

- `sys_users` напрямую не отдается отдельным CRUD API, но используется для входа в систему
- UI должен работать через JWT:
  - `Authorization: Bearer <access_token>`
  - при истечении токена выполнять refresh
  - при неуспешном refresh переводить пользователя на страницу логина

## 3. Целевая структура фронтенда

Фронтенд предлагается развивать в существующей папке:

- `/Users/roman/projects/CLiPE/server/panel`

Базовый стек:

- `React`
- `React Router`
- `React-Bootstrap`
- `Bootstrap`
- `Axios`

Дополнительно рекомендуется:

- `react-hook-form` или управляемые формы на `Form` из React-Bootstrap
- `yup`/`zod` для валидации форм
- `dayjs` для форматирования дат

## 4. Архитектурный подход

### 4.1. Общая схема приложения

Приложение должно быть разделено на 2 зоны:

1. Публичная зона
   - страница входа
2. Защищенная зона
   - все рабочие страницы после логина

### 4.2. Предлагаемая структура каталогов

```text
panel/src/
  app/
    router.jsx
    providers.jsx
  api/
    http.js
    authApi.js
    usersApi.js
    hostsApi.js
    servicesApi.js
    rulesApi.js
    policiesApi.js
    policyContentsApi.js
    requestsApi.js
    decisionsApi.js
  auth/
    AuthContext.jsx
    ProtectedRoute.jsx
    tokenStorage.js
  components/
    layout/
      AppNavbar.jsx
      Sidebar.jsx
      PageShell.jsx
    common/
      DataTable.jsx
      ConfirmModal.jsx
      SearchPanel.jsx
      EmptyState.jsx
      JsonViewer.jsx
      StatusBadge.jsx
      LoadingOverlay.jsx
    forms/
      UserForm.jsx
      HostForm.jsx
      ServiceForm.jsx
      RuleForm.jsx
      PolicyForm.jsx
      PolicyContentForm.jsx
      RequestForm.jsx
      DecisionForm.jsx
  pages/
    LoginPage.jsx
    DashboardPage.jsx
    users/
    hosts/
    services/
    rules/
    policies/
    policyContents/
    requests/
    decisions/
    accessCheck/
  utils/
    auth.js
    errors.js
    formatters.js
    constants.js
```

### 4.3. Layout

Использовать классическую административную компоновку на React-Bootstrap:

- верхний `Navbar`
- левый `Offcanvas`/боковое меню
- центральная рабочая область `Container` + `Row` + `Col`
- таблицы на `Table`
- формы на `Form`, `FloatingLabel`, `InputGroup`
- действия через `Button`, `Dropdown`, `Modal`, `Alert`, `Toast`

## 5. Страницы, которые нужно реализовать

### 5.1. Страница входа `LoginPage`

Назначение:

- ввод `username` и `password`
- отправка запроса на `POST /auth/login`
- сохранение `access_token` и `refresh_token`
- переход в защищенную часть приложения

Поведение:

- при ошибке входа показывать `Alert`
- при наличии валидной сессии сразу редиректить на dashboard
- добавить кнопку выхода, вызывающую `POST /auth/logout`

### 5.2. Главная страница `DashboardPage`

Назначение:

- стартовая точка после логина
- быстрый обзор системы

Содержимое:

- карточки со счетчиками: users, hosts, services, rules, policies, requests, decisions
- быстрые переходы к разделам
- блок “последние решения” или “последние запросы”

### 5.3. Страница пользователей

Сценарии:

- просмотр списка пользователей
- фильтрация по `user_id`, `user_name`, `uid`, `gid`, `host_id`
- создание пользователя
- редактирование пользователя
- удаление пользователя

UI:

- таблица пользователей
- модальное окно формы
- селект хоста с подгрузкой `hosts`

### 5.4. Страница хостов

Сценарии:

- просмотр списка хостов
- фильтрация по `host_id`, `ip`
- создание
- редактирование
- удаление

UI:

- таблица
- форма с валидацией IP

### 5.5. Страница сервисов

Сценарии:

- просмотр списка сервисов
- фильтрация по `service_id`, `service_name`
- создание
- редактирование
- удаление

### 5.6. Страница правил

Сценарии:

- просмотр списка правил
- фильтрация по `rule_id`, `rule_name`, `effect`
- создание
- редактирование
- удаление

Особенности:

- `conditions` — JSON; нужен удобный редактор
- на первом этапе можно использовать `textarea` с JSON-валидацией
- далее можно улучшить до конструктора условий

UI:

- таблица правил
- бейдж `Allow/Deny` для `effect`
- редактор JSON-условий

### 5.7. Страница политик

Сценарии:

- просмотр списка политик
- фильтрация по `policy_id`, `policy_name`, `user_id`, `status`
- создание
- редактирование
- удаление

Особенности:

- политика привязана к одному пользователю
- нужен селект пользователя
- статус удобно визуализировать `Badge`/`Form.Check`

### 5.8. Страница policy contents

Назначение:

- управление связями между политикой, сервисом и правилом
- это ключевой экран реального управления доступом

Сценарии:

- просмотр списка связей
- фильтрация по `policy_id`, `service_id`, `rule_id`
- создание связи
- редактирование `rule_id` для пары `(policy_id, service_id)`
- удаление

UI:

- таблица связей
- селекты для политики, сервиса и правила
- желательно показывать человекочитаемые названия рядом с ID

### 5.9. Страница журнала запросов

Сценарии:

- просмотр списка запросов
- фильтрация по `request_id`, `user_id`
- просмотр `context`
- при необходимости ручное редактирование и удаление

Особенности:

- `context` — JSONB, нужен читабельный просмотр
- редактирование лучше оставить доступным, но визуально пометить как административное действие

UI:

- таблица
- модальное окно “просмотр JSON”
- форматирование timestamp

### 5.10. Страница журнала решений

Сценарии:

- просмотр списка решений
- фильтрация по `decision_id`, `request_id`, `policy_id`
- просмотр результата `allow/deny`
- при необходимости редактирование и удаление

UI:

- таблица
- цветовая индикация `result`
- ссылки на связанные `request` и `policy`

### 5.11. Страница проверки доступа

Несмотря на то, что `aggregator` — служебный endpoint для `decision_server`, для GUI полезно заложить отдельную страницу-концепт:

- форма `user_name + service_name`
- отображение найденной политики и правила

Но в текущем API этот endpoint требует machine-auth (`DecisionServer`), а не пользовательский JWT. Поэтому:

- в первой версии GUI эту страницу не реализовывать как рабочую
- в плане заложить как будущий экран после появления пользовательского маршрута для ручной проверки доступа

## 6. JWT-аутентификация во фронтенде

### 6.1. Что хранить

- `access_token`
- `refresh_token`
- признак авторизованной сессии
- опционально `username`, если его хотим показывать в шапке

### 6.2. Где хранить

Для первой версии:

- `access_token` и `refresh_token` в `localStorage` или `sessionStorage`

Более безопасная перспектива:

- перевести refresh token в httpOnly cookie на стороне backend

### 6.3. Axios interceptor

Нужен общий HTTP-клиент:

- автоматически добавляет `Authorization: Bearer ...`
- при `401` пытается один раз выполнить `POST /auth/refresh`
- если refresh успешен, повторяет исходный запрос
- если нет, очищает токены и редиректит на `/login`

### 6.4. Protected routes

Все маршруты кроме `/login` должны быть обернуты в `ProtectedRoute`:

- нет токена → редирект на логин
- токен есть → рендер layout и страницы

## 7. Модель маршрутов React Router

```text
/login
/
/users
/hosts
/services
/rules
/policies
/policy-contents
/requests
/decisions
```

Роут `/` должен вести на `DashboardPage`.

## 8. API-слой фронтенда

Для каждой сущности нужен отдельный сервис с методами:

- `search(filters, limit, offset)`
- `create(payload)`
- `update(id | compositeKey, payload)`
- `remove(id | compositeKey)`

Особенности:

- у backend поиск реализован через `POST .../search`, а не `GET`
- пагинация передается через query params `limit` и `offset`
- фильтры передаются в JSON body
- у `policy-contents` составной ключ: `policy_id + service_id`

## 9. Повторно используемые UI-компоненты

Чтобы не писать каждый экран заново, стоит сразу заложить общие блоки:

- `AppNavbar`
- `Sidebar`
- `PageHeader`
- `EntityTable`
- `EntityToolbar`
- `SearchPanel`
- `EntityModal`
- `DeleteConfirmModal`
- `JsonViewer`
- `PaginationBar`
- `StatusBadge`

Это позволит все CRUD-разделы собрать из одинаковых шаблонов.

## 10. Поведение форм по сущностям

### Простые формы

- `hosts`
- `services`

### Формы со связями

- `users` — выбор `host_id`
- `policies` — выбор `user_id`
- `policy_contents` — выбор `policy_id`, `service_id`, `rule_id`
- `decisions` — выбор `request_id`, `policy_id`
- `requests` — выбор `user_id`, JSON `context`

### Формы с JSON

- `rules.conditions`
- `requests.context`

Нужна единая стратегия:

- пользователь вводит JSON в текстовом поле
- перед отправкой выполняется `JSON.parse`
- при ошибке показывается текст валидации

## 11. UX-решения для административной панели

### 11.1. Навигация

Левое меню сгруппировать так:

- Dashboard
- Справочники
  - Users
  - Hosts
  - Services
- Управление доступом
  - Rules
  - Policies
  - Policy Contents
- Аудит
  - Requests
  - Decisions

### 11.2. Таблицы

В каждой таблице предусмотреть:

- строку фильтров
- пагинацию
- кнопки `Создать`, `Редактировать`, `Удалить`
- индикатор загрузки
- сообщение при пустом результате

### 11.3. Визуальные акценты

Через классические компоненты React-Bootstrap:

- `Badge bg="success"` для allow/active
- `Badge bg="secondary"` или `danger` для deny/inactive
- `Card` для summary-блоков
- `Modal` для подтверждения удаления и форм

## 12. Ограничения и нюансы backend, которые надо учесть

### 12.1. Не все create/update доступны обычному JWT одинаково

По коду:

- большая часть админских экранов доступна пользователю с JWT
- `aggregator` доступен только `decision_server`
- создание `requests` и `decisions` тоже ориентировано на machine principal

Следствие для UI:

- страницы `requests` и `decisions` нужны прежде всего как журналы просмотра и администрирования
- акцент в интерфейсе должен быть на чтение, фильтрацию и просмотр связей

### 12.2. Есть расхождения между SQL и Go-моделями

Например:

- в `init.sql` у `policies` и тестовых данных заметны несоответствия
- в коде `policy_contents` хранит связь с rule отдельно

Следствие:

- фронтенд должен опираться на фактическое API, а не только на SQL-дамп
- на этапе реализации полезно вручную проверить реальные ответы API

### 12.3. Некоторые search endpoints фильтруют по zero-values

Например:

- для `rules` фильтр по `effect` может вести себя не совсем нейтрально без отдельной обработки

Следствие:

- во фронтенде для boolean-фильтров нужен режим “не задано / true / false”
- не отправлять поле в body, если фильтр не выбран

## 13. Пошаговый план реализации

### Этап 1. Подготовка фронтенд-основы

- очистить шаблонный Vite UI
- установить `react-router-dom`, `react-bootstrap`, `bootstrap`, `axios`
- подключить bootstrap CSS
- создать базовый layout
- настроить маршрутизацию

### Этап 2. Аутентификация

- реализовать `AuthContext`
- реализовать страницу логина
- настроить хранение токенов
- добавить Axios interceptors
- добавить `ProtectedRoute`
- реализовать logout

### Этап 3. Базовый каркас админки

- верхняя навигация
- боковое меню
- dashboard
- общие компоненты таблиц, модалок и пагинации

### Этап 4. CRUD для справочников

- `hosts`
- `users`
- `services`

Это даст базу для зависимых селектов в остальных формах.

### Этап 5. CRUD для управления доступом

- `rules`
- `policies`
- `policy_contents`

Это главный бизнес-блок системы.

### Этап 6. Аудит

- `requests`
- `decisions`
- просмотр JSON и связанных сущностей

### Этап 7. Полировка UX

- единые уведомления об успехе/ошибке
- скелетоны/спиннеры
- подтверждение удаления
- форматирование дат
- удобное отображение JSON

### Этап 8. Будущие улучшения

- страница ручной проверки доступа
- визуальный конструктор правил вместо raw JSON
- drill-down: переход из policy в связанные contents, из decision в request и policy
- role-based UI, если backend начнет различать админов и операторов

## 14. Минимальный MVP состава страниц

Для первой рабочей версии достаточно:

- Login
- Dashboard
- Users
- Hosts
- Services
- Rules
- Policies
- Policy Contents
- Requests
- Decisions

## 15. Что я рекомендую делать первым

Оптимальная очередность разработки:

1. Поднять инфраструктуру React + Router + React-Bootstrap
2. Сделать JWT login/refresh/logout
3. Собрать общий admin layout
4. Реализовать `Hosts`, `Users`, `Services`
5. Реализовать `Rules`, `Policies`, `Policy Contents`
6. Добавить `Requests` и `Decisions`
7. После этого переходить к улучшению UX и более умному редактору JSON

## 16. Итог

Для этого проекта подходит классическая админ-панель на React + React-Bootstrap с защищенными маршрутами и единым CRUD-паттерном для всех сущностей. Основной фокус UI должен быть на:

- удобной навигации по сущностям
- понятном управлении политиками и правилами
- безопасной JWT-сессии с refresh
- хорошем просмотре журналов и JSON-структур

Такой подход позволит быстро получить рабочий интерфейс без тяжелого кастомного дизайна и без конфликта с текущим устройством backend API.
