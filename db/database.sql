--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: golang_gin_db; Type: DATABASE; Schema: -; Owner: postgres
--
DROP DATABASE golang_gin_db;

CREATE DATABASE golang_gin_db WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


ALTER DATABASE golang_gin_db OWNER TO postgres;

\connect golang_gin_db

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


CREATE FUNCTION created_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN
	NEW.updated_at = EXTRACT(EPOCH FROM NOW());
	NEW.created_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;

$$;


ALTER FUNCTION public.created_at_column() OWNER TO postgres;

--
-- TOC entry 190 (class 1255 OID 36646)
-- Name: update_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION update_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;

$$;


ALTER FUNCTION public.update_at_column() OWNER TO postgres;


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: article; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--

CREATE TABLE article (
    id integer NOT NULL,
    user_id integer,
    title character varying,
    content text,
    updated_at integer,
    created_at integer
);


ALTER TABLE article OWNER TO postgres;

--
-- Name: article_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE article_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE article_id_seq OWNER TO postgres;

--
-- Name: article_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE article_id_seq OWNED BY article.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--

CREATE TABLE "user" (
    id integer NOT NULL,
    user_id character varying,
    email character varying,
    password character varying,
    name character varying,
    icon character varying,
    joined_at integer,
    updated_at integer,
    created_at integer
);


ALTER TABLE "user" OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_id_seq OWNER TO postgres;


CREATE TABLE "messages" (
    id integer NOT NULL,
    user_id character varying,
    group_id character varying,
    message character varying,
    updated_at integer,
    created_at integer
);

ALTER TABLE "messages" OWNER TO postgres;

CREATE SEQUENCE messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE messages_id_seq OWNER TO postgres;

ALTER SEQUENCE messages_id_seq OWNED BY "messages".id;

ALTER TABLE ONLY "messages" ALTER COLUMN id SET DEFAULT nextval('messages_id_seq'::regclass);

SELECT pg_catalog.setval('messages_id_seq', 1, false);

ALTER TABLE ONLY "messages"
    ADD CONSTRAINT messages_id PRIMARY KEY (id);
CREATE TRIGGER create_messages_created_at BEFORE INSERT ON "messages" FOR EACH ROW EXECUTE PROCEDURE created_at_column();
CREATE TRIGGER update_messages_updated_at BEFORE UPDATE ON "messages" FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE user_id_seq OWNED BY "user".id;


CREATE FUNCTION last_activity_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN
	INSERT INTO
        conversations (user_id, group_id, last_activity)
        VALUES(new.user_id, new.group_id, EXTRACT(EPOCH FROM NOW()))
        ON CONFLICT (user_id, group_id) DO UPDATE 
        SET last_activity = EXTRACT(EPOCH FROM NOW());
        RETURN new;
END;

$$;


ALTER FUNCTION public.last_activity_column() OWNER TO postgres;


CREATE FUNCTION last_activity_column_for_new_row() RETURNS trigger
    LANGUAGE plpgsql
    AS $$

BEGIN
        NEW.last_activity = EXTRACT(EPOCH FROM NOW());
        RETURN NEW;
END;

$$;


ALTER FUNCTION public.last_activity_column_for_new_row() OWNER TO postgres;



CREATE TABLE "conversations" (
    id integer NOT NULL,
    user_id character varying,
    group_id character varying,
    last_activity integer,
    updated_at integer,
    created_at integer
);

ALTER TABLE "conversations" OWNER TO postgres;

CREATE SEQUENCE conversations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE conversations_id_seq OWNER TO postgres;

ALTER SEQUENCE conversations_id_seq OWNED BY "conversations".id;

ALTER TABLE ONLY "conversations" ALTER COLUMN id SET DEFAULT nextval('conversations_id_seq'::regclass);

SELECT pg_catalog.setval('conversations_id_seq', 1, false);

ALTER TABLE ONLY "conversations"
    ADD CONSTRAINT conversations_id PRIMARY KEY (id);
    
CREATE TRIGGER create_conversations_created_at BEFORE INSERT ON "conversations" FOR EACH ROW EXECUTE PROCEDURE created_at_column();

CREATE TRIGGER create_conversations_created_at_to_last_activity BEFORE INSERT ON "conversations" FOR EACH ROW EXECUTE PROCEDURE last_activity_column_for_new_row();

CREATE UNIQUE INDEX last_activity_check_idx on conversations (user_id, group_id);

CREATE TRIGGER update_conversations_updated_at BEFORE UPDATE ON "conversations" FOR EACH ROW EXECUTE PROCEDURE update_at_column();

CREATE TRIGGER update_last_activity AFTER INSERT ON "messages" FOR EACH ROW EXECUTE PROCEDURE last_activity_column();



--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY article ALTER COLUMN id SET DEFAULT nextval('article_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY "user" ALTER COLUMN id SET DEFAULT nextval('user_id_seq'::regclass);


--
-- Data for Name: article; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY article (id, user_id, title, content, updated_at, created_at) FROM stdin;
\.


--
-- Name: article_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('article_id_seq', 1, false);


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY "user" (id, user_id, email, password, name, updated_at, created_at) FROM stdin;
\.


--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('user_id_seq', 1, false);


--
-- Name: article_id; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--

ALTER TABLE ONLY article
    ADD CONSTRAINT article_id PRIMARY KEY (id);


--
-- Name: user_id; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_id PRIMARY KEY (id);


--
-- Name: article_user_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY article
    ADD CONSTRAINT article_user_id FOREIGN KEY (user_id) REFERENCES "user"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 2284 (class 2620 OID 36647)
-- Name: article create_article_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER create_article_created_at BEFORE INSERT ON article FOR EACH ROW EXECUTE PROCEDURE created_at_column();


--
-- TOC entry 2286 (class 2620 OID 36653)
-- Name: user create_user_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER create_user_created_at BEFORE INSERT ON "user" FOR EACH ROW EXECUTE PROCEDURE created_at_column();


--
-- TOC entry 2285 (class 2620 OID 36648)
-- Name: article update_article_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_article_updated_at BEFORE UPDATE ON article FOR EACH ROW EXECUTE PROCEDURE update_at_column();


--
-- TOC entry 2287 (class 2620 OID 36654)
-- Name: user update_user_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE PROCEDURE update_at_column();



--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--
