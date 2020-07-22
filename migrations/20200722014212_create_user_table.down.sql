BEGIN;

  ALTER TABLE node
  DROP CONSTRAINT node_user_constraint_fk;

  ALTER TABLE node
  DROP COLUMN user_id;

  DROP TABLE IF EXISTS users;
  
  COMMIT;
