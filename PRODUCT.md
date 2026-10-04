# GMailer

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Vue 3 with Vite, Vue Router, Axios, and Pinia; Go HTTP API and PostgreSQL are already present in the repository.

## Users

Inferred from the repository and request: a signed-in owner of an email list who needs a compact workspace for preparing and monitoring mailings.

## Product Purpose

GMailer provides registration and sign-in for a mailing workspace. The current frontend should make authentication usable and provide a clear dashboard shell for the contacts, mailings, and delivery-message entities already represented by the database schema.

## Positioning

Open decision: no differentiated commercial positioning or claims were supplied. The implementation should demonstrate the real workflow without inventing product claims.

## Operating Context

Users enter by registering or signing in, then work from an authenticated desktop-first dashboard that remains usable on mobile. The current API exposes authentication; dashboard business data is therefore clearly labeled demo content until corresponding endpoints exist.

## Capabilities and Constraints

- Registration uses `POST /api/auth/register` with `email` and `password`.
- Sign-in uses `POST /api/auth/login` and returns a JWT plus a user object.
- No authenticated profile, contacts, mailings, or analytics endpoints currently exist.
- The frontend must use Vue Router, Axios, and Pinia.
- Authentication state may be persisted locally until the API provides a refresh-token/session mechanism.

## Brand Commitments

The product name is GMailer. The user explicitly requested a color direction based on Gmail: a light, cool-neutral workspace with Google-blue emphasis and restrained status colors.

## Evidence on Hand

The repository contains Go authentication handlers, user/contact/mailing models and migrations. There are no logos, customer evidence, benchmarks, or production dashboard data; none should be fabricated as product claims.

## Product Principles

- Make the next mailing task obvious at a glance.
- Keep authentication direct, forgiving, and transparent about errors.
- Clearly separate live account data from illustrative dashboard content.
- Preserve a calm, scan-friendly workspace on both desktop and mobile.

## Accessibility & Inclusion

Use semantic controls, visible focus states, keyboard-accessible navigation, sufficient contrast, and reduced-motion support.
