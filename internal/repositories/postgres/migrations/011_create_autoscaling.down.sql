-- Migration: 011_create_autoscaling.down.sql

DROP TABLE IF EXISTS scaling_group_instances;
DROP TABLE IF EXISTS scaling_policies;
DROP TABLE IF EXISTS scaling_groups;
