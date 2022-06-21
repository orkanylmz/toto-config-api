create table sku_configs
(
    uuid           varchar(255)               not null
        constraint sku_configs_pk
            primary key,
    package        varchar(155)       not null,
    country_code   varchar(3)         not null,
    percentile_min smallint CHECK ( percentile_min between 0 and 100) default 0 not null,
    percentile_max smallint CHECK ( percentile_max between 0 and 100) default 100,
    sku            varchar(255)       not null
);

create index sku_configs_country_code_index
    on sku_configs (country_code);

create index sku_configs_percentile_min_percentile_max_index
    on sku_configs (percentile_min, percentile_max);

create unique index sku_configs_uuid_uindex
    on sku_configs (uuid);

--
-- INSERT INTO public.sku_configs (uuid, package, country_code, percentile_min, percentile_max, sku)
-- VALUES ('5b4dcb45-32ad-480c-b6fe-91d4c4ee3198
--
-- ', 'com.softinit.iquitos.mainapp', 'US', 0, 25, 'rdm_premium_v3_020_trial_7d_monthly');
--
-- INSERT INTO public.sku_configs (uuid, package, country_code, percentile_min, percentile_max, sku)
-- VALUES ('5cbbba53-43d3-4add-a01b-cb8e29288274
--
-- ', 'com.softinit.iquitos.mainapp', 'US', 25, 50, 'rdm_premium_v3_030_trial_7d_monthly');
--
-- INSERT INTO public.sku_configs (uuid, package, country_code, percentile_min, percentile_max, sku)
-- VALUES ('bd88fd75-d336-46c2-87cb-a66b73ff06c3
--
-- ', 'com.softinit.iquitos.mainapp', 'US', 50, 75, 'rdm_premium_v3_100_trial_7d_yearly');
--
-- INSERT INTO public.sku_configs (uuid, package, country_code, percentile_min, percentile_max, sku)
-- VALUES ('72ee460f-7e01-4b2d-be57-0a826c6ecb29
--
-- ', 'com.softinit.iquitos.mainapp', 'US', 75, 100, 'rdm_premium_v3_150_trial_7d_yearly');
--
-- INSERT INTO public.sku_configs (uuid, package, country_code, percentile_min, percentile_max, sku)
-- VALUES ('38af6ce7-cde4-4a69-a82d-5cdfc7f29487
--
-- ', 'com.softinit.iquitos.mainapp', 'ZZ', 0, 100, 'rdm_premium_v3_050_trial_7d_yearly');

