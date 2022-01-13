create table IF NOT EXISTS events
(
    id       serial             not null
        constraint events_pkey
            primary key,
    name     varchar(255)       not null,
    date     date default now() not null,
    location varchar(255)
);

create table pictures
(
    id            uuid         not null
        constraint pictures_pkey
            primary key,
    original_path varchar(255) not null,
    preview_path  varchar(255) not null,
    published     timestamp default CURRENT_TIMESTAMP,
    eventid       integer
        constraint pictures_eventid_fkey
            references events
            on delete cascade
);

create table pictures_text_detection
(
    picturesid    uuid
        constraint pictures_text_detection_picturesid_fkey
            references pictures
            on delete cascade,
    detected_text varchar  not null,
    confidence    smallint not null
);

create table users
(
    id        serial             not null
        constraint users_pkey
            primary key,
    username  varchar(255)       not null,
    passwd    text               not null,
    user_role smallint default 1 not null
);

create unique index users_username_uindex
    on users (username);

insert into public.users (id, username, passwd, user_role)
values  (7, 'admin', '$2a$10$lZBJMPMSt3LKV08jWDrh7.KEIxvvZw8efjQX5GENsgrVnPye56wX.', 0);


insert into public.events (id, name, date, location)
values  (20, 'asdfasdfasd', '2021-07-27', 'asfadsf'),
        (21, 'Событие в испаниия', '2021-07-07', 'Локация'),
        (22, 'Событие в испаниия', '2021-07-22', 'Локация'),
        (23, 'Событие в испаниия', '2021-07-09', 'Локация'),
        (24, 'ВЭБ.РФ IRONMAN 70.3 ST. PETERSBURG', '2021-08-01', ' St. Petersburg, Russia  '),
        (26, 'Событие в испаниия', '2021-07-23', 'Локация');

insert into public.pictures (id, original_path, preview_path, published, eventid)
values  ('b41787a9-3667-4cbc-b363-5aa0bde32226', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/87d2fd79-3448-4c22-8fe5-7eaca6156864830a8890-179b-43c7-add1-65a593a2e6fcjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/87d2fd79-3448-4c22-8fe5-7eaca6156864-thumb459986e2-ed2b-412e-bb59-2b0b2f071814jpg', '2021-07-27 16:01:33.003177', 26),
        ('86484717-a4a2-48fd-8230-4274bde43870', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/3b7c52d6-8a8b-439d-91eb-1406b045f086d699905b-c1f4-4b91-b9a3-aa4077853ad7jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/3b7c52d6-8a8b-439d-91eb-1406b045f086-thumb1bc06f9f-aa48-40f0-abd1-782fa606f4b6jpg', '2021-07-27 16:01:33.003177', 26),
        ('96d9af0e-c501-462a-b0b8-e637a7505361', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/71bc3e22-c5d7-4fe9-89d5-a884989d32f677a9cd49-24b8-4497-9c9d-36b625614ac5jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/71bc3e22-c5d7-4fe9-89d5-a884989d32f6-thumb36d06f0c-6f3f-41e2-92f6-041896536113jpg', '2021-07-27 16:01:33.003177', 26),
        ('91134895-0017-4838-8af0-ff91deaa171d', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/a5c0c5b2-1d3e-43b9-b2c6-6fc28fb9960d36c2d5cb-04cd-4932-9295-ebf16c434dc0jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/a5c0c5b2-1d3e-43b9-b2c6-6fc28fb9960d-thumb7d7f9632-c089-4cb7-aab2-7f64cb879f26jpg', '2021-07-27 16:01:33.003177', 26),
        ('09b6cbd2-3ea2-4198-ac80-6968c764b271', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/b839dc76-d9e3-4f78-ba3b-7dd79f24eb59966563f1-fe99-4891-8c34-3e8bdfe2168ejpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/b839dc76-d9e3-4f78-ba3b-7dd79f24eb59-thumb311da702-d806-4206-aa95-c138b26de246jpg', '2021-07-27 16:01:33.003177', 26),
        ('8972c2d3-d40a-465d-806b-4902b3e8b63c', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/a4c730b8-6e7c-4647-a04f-84574b383d98006cfb2a-bced-4dff-80c7-20ce9db6018djpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/26/a4c730b8-6e7c-4647-a04f-84574b383d98-thumb6814e553-fb79-443b-ac97-4d61ca3a937ejpg', '2021-07-27 16:01:33.003177', 26),
        ('0cb84973-6943-4cb5-b798-eb16276efdd3', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/93572390-46a2-4239-8e23-d680cc0d14b0c7710de6-4837-48b3-8781-2e7d0c7a7978jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/93572390-46a2-4239-8e23-d680cc0d14b0-thumbdec547c3-d1fb-4c3b-ac20-7b23ec6a6863jpg', '2021-07-27 05:31:37.282302', 24),
        ('d5559f8f-04d0-4e7c-877f-67cb95f9c86e', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ee97d176-00dc-4419-bfbf-6c3611f73766bad15727-af4e-4895-86fd-c3e1483ad383jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ee97d176-00dc-4419-bfbf-6c3611f73766-thumbbd0e344d-d94c-4346-9ad6-b6272da08719jpg', '2021-07-27 05:31:37.282302', 24),
        ('ff548bb3-7a10-4c8c-97c4-61926174b559', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ad70611a-2142-437c-be45-0e2fda82c7f5fa62bbf8-588a-48a3-a653-49363addec13jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ad70611a-2142-437c-be45-0e2fda82c7f5-thumb0a78c0ae-3a0f-4303-a6ed-6741c5c6e1e7jpg', '2021-07-27 05:31:37.282302', 24),
        ('8046fd74-fe6d-47e0-bc08-24ad75b08e97', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/89ec012c-2202-49ee-b655-c93efb5dd387f5d522a8-1a4c-4ed4-98c3-583c888eb860jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/89ec012c-2202-49ee-b655-c93efb5dd387-thumb4a29bcaa-be7f-4e60-935e-676b945fb20djpg', '2021-07-27 05:31:37.282302', 24),
        ('beea4608-fb8a-47d1-b4f6-3a0dbc5cda47', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/a6cf55a5-3beb-46e7-957a-e2e8a79d99ab3113d395-c826-4255-94f8-eeb4318ed221jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/a6cf55a5-3beb-46e7-957a-e2e8a79d99ab-thumba562084c-0070-48ea-8fbe-f9311b62919cjpg', '2021-07-27 05:31:37.282302', 24),
        ('c419c153-b308-49ec-a774-3699194c33e9', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/d810c192-3fcc-41f2-82f3-6031db5009be0a4c7f0b-58e2-40b3-b014-a8d6ae9a29c9jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/d810c192-3fcc-41f2-82f3-6031db5009be-thumb15b888a9-c009-457c-8c60-8b4c72de610djpg', '2021-07-27 05:31:37.282302', 24),
        ('54c84f55-7d53-4d27-814a-a48cee9e5f9f', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/943b275b-c1eb-4461-a6e8-e8e4018770d1ff631639-46eb-4ac2-8964-1b925410a433jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/943b275b-c1eb-4461-a6e8-e8e4018770d1-thumbbfc14dde-b0c7-4588-9778-6ca4a4855f28jpg', '2021-07-27 05:31:37.282302', 24),
        ('27690c4c-016b-45e0-9df1-174b1d42a9bd', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/42188761-e097-4236-9e64-9833e30da962e465f2ac-0860-4783-8eec-f237fc89669djpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/42188761-e097-4236-9e64-9833e30da962-thumb693f5218-24aa-4343-b18f-19c5eabf12d5jpg', '2021-07-27 05:31:37.282302', 24),
        ('9f8ff147-dd63-4425-bd0b-f31815664ba3', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/527400b4-a193-4ff3-b4b4-fe72a3fea414741ab3a1-6c22-405b-b785-0c2412a22a14jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/527400b4-a193-4ff3-b4b4-fe72a3fea414-thumb8e82346f-5af8-4199-9927-9977f7ac6361jpg', '2021-07-27 05:31:37.282302', 24),
        ('bf853453-f48a-48b1-9c57-fd27ea1beed9', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/6c3b628d-958e-4c4d-a1ab-892d5dd6e170f48aecdd-147d-4dcd-b8c4-e07db208a98bjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/6c3b628d-958e-4c4d-a1ab-892d5dd6e170-thumb1005f944-516c-46be-a0b1-b942b4883f25jpg', '2021-07-27 05:31:37.282302', 24),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/9093c403-22a0-4925-855e-758b5861e9d5810b3550-33c1-4922-9501-7134d43300c9jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/9093c403-22a0-4925-855e-758b5861e9d5-thumba62fe35a-3393-4bb7-8acb-a7b67096940cjpg', '2021-07-27 05:31:37.282302', 24),
        ('0713bbe4-2765-4611-8ae3-8d953f914215', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/da723116-6b53-4340-9ad9-ed6e1922af64f47d45e2-c2bb-43c7-8baf-e58f4caf4455jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/da723116-6b53-4340-9ad9-ed6e1922af64-thumb868e9253-d2ec-4450-bc16-d077cba862d3jpg', '2021-07-27 05:31:37.282302', 24),
        ('655d3903-ebef-4059-b041-e2d66e7cca6d', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/891ccdbd-3f1d-4a0c-8811-36c0647142b627437898-439a-47d6-952c-7ef5d65809d3jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/891ccdbd-3f1d-4a0c-8811-36c0647142b6-thumbdcb60415-23ae-487d-ae4a-0ca6c4953d83jpg', '2021-07-27 05:31:37.282302', 24),
        ('61e0bd33-5062-4ad7-8799-b419ad65e373', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/dbcb4773-e707-4369-a2b2-b5548c43dcba0176eefe-e271-4a9c-a6e8-351f186b9dddjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/dbcb4773-e707-4369-a2b2-b5548c43dcba-thumb4036f90a-86af-40e6-adf4-0fe16aafa0a9jpg', '2021-07-27 05:31:37.282302', 24),
        ('849f5350-e6c5-47ee-b49d-a829bd9b95b8', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/3ef1aacb-6997-4fd0-86d1-8d1205469c1c541f44f0-4a84-4097-967f-4dd8b073bdffjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/3ef1aacb-6997-4fd0-86d1-8d1205469c1c-thumb8e928064-8d60-4897-8e11-2a80e04643a7jpg', '2021-07-27 05:31:37.282302', 24),
        ('978e253f-4e8d-4f54-bf1d-90b6dd8be2b2', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/4e345ebf-27cc-4130-b485-7e881bea8e97df4c2591-5e0b-4254-ae32-efd5e180d117jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/4e345ebf-27cc-4130-b485-7e881bea8e97-thumbb13c1c11-b510-49f3-9a46-bdb1f825004ajpg', '2021-07-27 05:31:37.282302', 24),
        ('b5606da5-69b7-44ac-af00-fc643e47686d', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/8d8218ab-e420-4a13-851d-0b4348036878911fc684-ce3f-44ec-a486-7831f4459b79jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/8d8218ab-e420-4a13-851d-0b4348036878-thumb3c6ad810-0d08-4901-b210-61bd94bcaa2fjpg', '2021-07-27 05:31:37.282302', 24),
        ('d4936a09-af02-43f7-9cff-e0aa2ac3f1a0', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/f986ad43-45df-4460-bcec-854e1d9ff608182fb20b-1fe4-40bf-a594-edbc97a443f3jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/f986ad43-45df-4460-bcec-854e1d9ff608-thumbed2e75f5-05da-45d2-a9b4-7763218bf27djpg', '2021-07-27 05:31:37.282302', 24),
        ('46bd96f7-09b5-4c97-813f-704df87942f2', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1cef9ccc-5486-41d1-b56b-9228e71c73c0266d5f81-6320-4947-b092-9fdcce0f4ba1jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1cef9ccc-5486-41d1-b56b-9228e71c73c0-thumb1a8da46a-ae62-463a-8483-587213996f2ajpg', '2021-07-27 05:31:37.282302', 24),
        ('1672ca56-c24e-4740-9850-ac71a7ff8959', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/24ea0795-9d46-4652-8d58-5ea92ba68a51aa4cd068-6ee2-4266-adeb-200eba35783djpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/24ea0795-9d46-4652-8d58-5ea92ba68a51-thumba188a9da-7a8f-4032-9216-078a34b72910jpg', '2021-07-27 05:31:37.282302', 24),
        ('f6bd0df9-a10c-494f-9d8a-787528be0f62', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/b7d76e91-517c-4c31-ab62-3e2ff7311f05f9f535ca-2c9c-463d-8eae-ddc66623a971jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/b7d76e91-517c-4c31-ab62-3e2ff7311f05-thumbff687c07-4f24-4118-989c-763aa7573eb8jpg', '2021-07-27 05:31:37.282302', 24),
        ('47d4378d-26e4-43f9-a13a-b64a4258e1e5', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/49703f2e-8848-4ec8-9afe-3132bc8a9d540d221260-4010-465b-8fc3-dceba35c4f9djpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/49703f2e-8848-4ec8-9afe-3132bc8a9d54-thumb6befd702-c541-4280-81ff-483f77e9e5c0jpg', '2021-07-27 05:31:37.282302', 24),
        ('8898a95b-5413-4777-9884-29ac8f2d4c1d', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1bdcee1d-dad7-41aa-86e6-72eb785ba9720af2d685-5b8b-4877-bd80-669162187a4cjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1bdcee1d-dad7-41aa-86e6-72eb785ba972-thumb40c8079a-8946-4510-9df3-285e43e0c6a0jpg', '2021-07-27 05:31:37.282302', 24),
        ('f59255ea-a0d5-4c6e-abfc-a6a01a5a06ab', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/726f61d7-7fb8-4ddb-a23c-88b724790852af85b0eb-36f3-48cb-8fce-bb9a078a6af7jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/726f61d7-7fb8-4ddb-a23c-88b724790852-thumb4277e9ef-1dcb-4a74-8bf4-cf84a52bb574jpg', '2021-07-27 05:31:37.282302', 24),
        ('10ce532a-d7fb-46f7-ba27-1225014d7b05', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/efaca5d7-1a1d-4b3c-91c5-3740c9a46d1823b02280-086e-4447-9b6d-6764cfc958acjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/efaca5d7-1a1d-4b3c-91c5-3740c9a46d18-thumb1c1971e8-119c-4ab0-82ff-395a1ba929f9jpg', '2021-07-27 05:31:37.282302', 24),
        ('7c7f8ee0-4fa8-42eb-aa39-bec76596196a', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ea6885c0-38e5-4e1a-991d-5ca3e360e650048c83fb-ea3e-4814-8522-c9b2cf9a67eajpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ea6885c0-38e5-4e1a-991d-5ca3e360e650-thumba98ea135-135c-48ae-9c82-304eb2895080jpg', '2021-07-27 05:31:37.282302', 24),
        ('bd6759bf-14d7-4803-9201-cb4f0a81de5e', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/96b907ac-e285-400e-a98a-a4d448940add2ea26b63-fb48-482a-8596-1092937a3940jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/96b907ac-e285-400e-a98a-a4d448940add-thumb7b992058-72f9-4dee-9a10-13ed2c2d2a5djpg', '2021-07-27 05:31:37.282302', 24),
        ('2872de8e-3829-4276-954b-b0d46b06c0ea', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/eafde733-3592-48f9-b668-1a6326b06c165d723eba-99c9-4407-8e7c-0a0c21c4a8bajpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/eafde733-3592-48f9-b668-1a6326b06c16-thumb78f29238-abeb-416c-a105-15e8198f334fjpg', '2021-07-27 05:31:37.282302', 24),
        ('7f0c0bb2-703e-4f25-a150-2633097d2a32', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/e5107d73-eada-4b97-addf-3aeee21c55580d0bbb81-2f1d-4f47-a310-139fa230a779jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/e5107d73-eada-4b97-addf-3aeee21c5558-thumb993f64d9-9246-4298-8234-c73f6d7da2cejpg', '2021-07-27 05:31:37.282302', 24),
        ('03e47f06-98d7-43b9-91b3-863c3ce0e325', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/e58ea5b9-4b34-4f16-9f73-65a5d581482d1af65b53-fd8a-4e42-a564-ed98a73dc181jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/e58ea5b9-4b34-4f16-9f73-65a5d581482d-thumb375fabf5-837e-4088-a5c2-e922479ff67ajpg', '2021-07-27 05:31:37.282302', 24),
        ('ccabc9f4-846f-4e03-9a4b-c83ec45ca7aa', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/56e29790-e512-4ca8-875d-f95c1816c9fc4a993741-2c33-4cd7-8a9d-aa630e743907jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/56e29790-e512-4ca8-875d-f95c1816c9fc-thumb8c095b7f-88ca-4a8e-8614-306339111a2ajpg', '2021-07-27 05:31:37.282302', 24),
        ('ab1e96e4-72a4-47d0-92b9-ae7d93bd11e4', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ce690798-5d7c-4067-956c-77d280ed076b5ea42330-84b6-431d-9417-f8b941426821jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/ce690798-5d7c-4067-956c-77d280ed076b-thumb23dc65b9-9840-44eb-9a14-482d40e256fdjpg', '2021-07-27 05:31:37.282302', 24),
        ('f9792792-6992-482b-a973-0a7b5202dd51', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1481edf0-ce8b-42b2-b30b-215ce1699d5ba095773e-04a6-4b58-9ce9-400759895517jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1481edf0-ce8b-42b2-b30b-215ce1699d5b-thumb9a1f5aa9-440d-4ee6-956f-440397b04234jpg', '2021-07-27 05:31:37.282302', 24),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8493', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/3729a546-b3ae-4386-8593-8866a9d6a548333526b8-e79e-459b-9ee0-e43beef825e5jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/3729a546-b3ae-4386-8593-8866a9d6a548-thumbbe9c6501-0a13-4014-b8b3-59a8933b09efjpg', '2021-07-27 05:31:37.282302', 24),
        ('51431b70-abc0-4a83-845f-640ae6037503', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/030242c3-5dd6-4223-b4b8-795b707be4da87f36e00-da63-45a8-a8f3-a6ea2ea853ebjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/030242c3-5dd6-4223-b4b8-795b707be4da-thumbd7cad4f7-ad75-4d1d-9e80-7ee399767733jpg', '2021-07-27 05:31:37.282302', 24),
        ('11f400cb-c286-41b1-afa4-db88b8a7622b', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 05:31:37.282302', 24),
        ('0cb84973-6943-4cb5-b798-eb16276ef234', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/93572390-46a2-4239-8e23-d680cc0d14b0c7710de6-4837-48b3-8781-2e7d0c7a7978jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/93572390-46a2-4239-8e23-d680cc0d14b0-thumbdec547c3-d1fb-4c3b-ac20-7b23ec6a6863jpg', '2021-07-27 05:31:37.282302', 20),
        ('f9792792-6992-482b-a973-0a7b5202dd52', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1481edf0-ce8b-42b2-b30b-215ce1699d5ba095773e-04a6-4b58-9ce9-400759895517jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/1481edf0-ce8b-42b2-b30b-215ce1699d5b-thumb9a1f5aa9-440d-4ee6-956f-440397b04234jpg', '2021-07-27 05:31:37.282302', 24),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8492', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/3729a546-b3ae-4386-8593-8866a9d6a548333526b8-e79e-459b-9ee0-e43beef825e5jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/3729a546-b3ae-4386-8593-8866a9d6a548-thumbbe9c6501-0a13-4014-b8b3-59a8933b09efjpg', '2021-07-27 05:31:37.282302', 24),
        ('51431b70-abc0-4a83-845f-640ae6037504', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/030242c3-5dd6-4223-b4b8-795b707be4da87f36e00-da63-45a8-a8f3-a6ea2ea853ebjpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/030242c3-5dd6-4223-b4b8-795b707be4da-thumbd7cad4f7-ad75-4d1d-9e80-7ee399767733jpg', '2021-07-27 05:31:37.282302', 24),
        ('11f400cb-c286-41b1-afa4-db88b8a76223', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 05:31:37.282302', 21),
        ('11f400cb-c286-41b1-afa4-db88b8a76224', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('7f0c0bb2-703e-4f25-a150-2633097d2d33', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8494', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8495', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8496', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8497', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('6491dda8-eed4-11eb-9a03-0242ac130003', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24),
        ('6af6e864-eed4-11eb-9a03-0242ac130003', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f519673da3-7069-4c2d-9c6c-14bd370c59b2jpg', 'https://img-rekongnition-test.s3.eu-central-1.amazonaws.com/24/c6aed2c8-ded3-4485-9149-b926196c23f5-thumbd0f9266b-e5a4-4275-87a2-fea1dd084f8ejpg', '2021-07-27 12:16:10.955833', 24);

insert into public.pictures_text_detection (picturesid, detected_text, confidence)
values  ('0cb84973-6943-4cb5-b798-eb16276efdd3', '2', 93),
        ('0cb84973-6943-4cb5-b798-eb16276efdd3', '128', 100),
        ('d5559f8f-04d0-4e7c-877f-67cb95f9c86e', '2018', 66),
        ('d5559f8f-04d0-4e7c-877f-67cb95f9c86e', '202', 100),
        ('ff548bb3-7a10-4c8c-97c4-61926174b559', '19', 97),
        ('ff548bb3-7a10-4c8c-97c4-61926174b559', '268', 100),
        ('8046fd74-fe6d-47e0-bc08-24ad75b08e97', '1', 98),
        ('8046fd74-fe6d-47e0-bc08-24ad75b08e97', '54', 100),
        ('8046fd74-fe6d-47e0-bc08-24ad75b08e97', '2009', 98),
        ('beea4608-fb8a-47d1-b4f6-3a0dbc5cda47', '0000', 74),
        ('beea4608-fb8a-47d1-b4f6-3a0dbc5cda47', '16', 100),
        ('c419c153-b308-49ec-a774-3699194c33e9', '555', 99),
        ('54c84f55-7d53-4d27-814a-a48cee9e5f9f', '555', 100),
        ('54c84f55-7d53-4d27-814a-a48cee9e5f9f', '866', 100),
        ('54c84f55-7d53-4d27-814a-a48cee9e5f9f', '2019', 78),
        ('27690c4c-016b-45e0-9df1-174b1d42a9bd', '2019', 36),
        ('27690c4c-016b-45e0-9df1-174b1d42a9bd', '550', 76),
        ('27690c4c-016b-45e0-9df1-174b1d42a9bd', '1111', 76),
        ('9f8ff147-dd63-4425-bd0b-f31815664ba3', '342', 100),
        ('bf853453-f48a-48b1-9c57-fd27ea1beed9', '122', 100),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', '67', 68),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', '1.9', 26),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', '0121', 65),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', '550', 99),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', '669', 80),
        ('d301db0c-a66b-415b-9fa1-ea8a9fb0332f', '425', 99),
        ('0713bbe4-2765-4611-8ae3-8d953f914215', '259', 100),
        ('655d3903-ebef-4059-b041-e2d66e7cca6d', '3', 100),
        ('655d3903-ebef-4059-b041-e2d66e7cca6d', '114', 100),
        ('61e0bd33-5062-4ad7-8799-b419ad65e373', '2019', 60),
        ('61e0bd33-5062-4ad7-8799-b419ad65e373', '36', 100),
        ('61e0bd33-5062-4ad7-8799-b419ad65e373', '31', 34),
        ('61e0bd33-5062-4ad7-8799-b419ad65e373', '-29', 34),
        ('849f5350-e6c5-47ee-b49d-a829bd9b95b8', '298', 100),
        ('978e253f-4e8d-4f54-bf1d-90b6dd8be2b2', '287', 100),
        ('978e253f-4e8d-4f54-bf1d-90b6dd8be2b2', '0', 75),
        ('978e253f-4e8d-4f54-bf1d-90b6dd8be2b2', '1015', 43),
        ('b5606da5-69b7-44ac-af00-fc643e47686d', '539', 100),
        ('d4936a09-af02-43f7-9cff-e0aa2ac3f1a0', '111111', 57),
        ('d4936a09-af02-43f7-9cff-e0aa2ac3f1a0', '2019', 98),
        ('d4936a09-af02-43f7-9cff-e0aa2ac3f1a0', '33', 100),
        ('46bd96f7-09b5-4c97-813f-704df87942f2', '287', 100),
        ('1672ca56-c24e-4740-9850-ac71a7ff8959', '314', 100),
        ('f6bd0df9-a10c-494f-9d8a-787528be0f62', '2019', 98),
        ('f6bd0df9-a10c-494f-9d8a-787528be0f62', '2', 99),
        ('f6bd0df9-a10c-494f-9d8a-787528be0f62', '287', 100),
        ('47d4378d-26e4-43f9-a13a-b64a4258e1e5', '114', 100),
        ('8898a95b-5413-4777-9884-29ac8f2d4c1d', '362', 100),
        ('8898a95b-5413-4777-9884-29ac8f2d4c1d', '50', 71),
        ('f59255ea-a0d5-4c6e-abfc-a6a01a5a06ab', '2015', 29),
        ('f59255ea-a0d5-4c6e-abfc-a6a01a5a06ab', '2019', 98),
        ('f59255ea-a0d5-4c6e-abfc-a6a01a5a06ab', '515', 100),
        ('f59255ea-a0d5-4c6e-abfc-a6a01a5a06ab', '243', 100),
        ('10ce532a-d7fb-46f7-ba27-1225014d7b05', '360', 100),
        ('7c7f8ee0-4fa8-42eb-aa39-bec76596196a', '104', 100),
        ('bd6759bf-14d7-4803-9201-cb4f0a81de5e', '207', 100),
        ('2872de8e-3829-4276-954b-b0d46b06c0ea', '550', 100),
        ('7f0c0bb2-703e-4f25-a150-2633097d2a32', '107', 100),
        ('03e47f06-98d7-43b9-91b3-863c3ce0e325', '1247', 49),
        ('03e47f06-98d7-43b9-91b3-863c3ce0e325', '550', 100),
        ('ccabc9f4-846f-4e03-9a4b-c83ec45ca7aa', '317', 64),
        ('ccabc9f4-846f-4e03-9a4b-c83ec45ca7aa', '32.', 86),
        ('ab1e96e4-72a4-47d0-92b9-ae7d93bd11e4', '114', 100),
        ('ab1e96e4-72a4-47d0-92b9-ae7d93bd11e4', '2', 72),
        ('ab1e96e4-72a4-47d0-92b9-ae7d93bd11e4', '3', 100),
        ('f9792792-6992-482b-a973-0a7b5202dd51', '301', 100),
        ('f7e50a7b-aaaa-41c9-87e3-61c03bfc8493', '439', 100),
        ('51431b70-abc0-4a83-845f-640ae6037503', '84', 100),
        ('11f400cb-c286-41b1-afa4-db88b8a7622b', '347', 99),
        ('0cb84973-6943-4cb5-b798-eb16276ef234', '1', 99),
        ('b41787a9-3667-4cbc-b363-5aa0bde32226', '207', 99),
        ('96d9af0e-c501-462a-b0b8-e637a7505361', '458', 100),
        ('91134895-0017-4838-8af0-ff91deaa171d', '1919', 100),
        ('09b6cbd2-3ea2-4198-ac80-6968c764b271', '550', 97),
        ('8972c2d3-d40a-465d-806b-4902b3e8b63c', '287', 96);