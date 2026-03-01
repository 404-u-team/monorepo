// For more info, see https://github.com/storybookjs/eslint-plugin-storybook#configuration-flat-config-format
import storybook from "eslint-plugin-storybook";

import { defineConfig, globalIgnores } from 'eslint/config'
import js from '@eslint/js'
import globals from 'globals'

import tsEslint, { configs as tsConfigs, plugin as tsPlugin } from 'typescript-eslint'

import react from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'
// eslint-disable-next-line import-x/no-named-as-default
import reactRefresh from 'eslint-plugin-react-refresh'
import { flatConfigs as importXFlatConfigs } from 'eslint-plugin-import-x'
import { createTypeScriptImportResolver } from 'eslint-import-resolver-typescript'
import unicorn from 'eslint-plugin-unicorn'
import unusedImports from 'eslint-plugin-unused-imports'

import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

export default defineConfig([
  globalIgnores(['dist', 'build', 'node_modules', 'src/app/generated/routeTree.gen.ts']),
  js.configs.recommended,
  ...tsEslint.config({
    extends: [
      ...tsConfigs.strictTypeChecked,
      ...tsConfigs.stylisticTypeChecked,
    ],
    files: ['**/*.{ts,tsx}'],
  }),
  importXFlatConfigs.recommended,
  importXFlatConfigs.typescript,
  {
    files: ['**/*.{ts,tsx}'],

    languageOptions: {
      ecmaVersion: 'latest',
      globals: {
        ...globals.browser,
        ...globals.es2020,
      },
      parserOptions: {
        project: ['./tsconfig.app.json', './tsconfig.node.json'],
        tsconfigRootDir: __dirname,
      },
    },

    plugins: {
      '@typescript-eslint': tsPlugin,
      react,
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
      unicorn,
      'unused-imports': unusedImports,
    },

    rules: {
      ...reactHooks.configs.recommended.rules,

      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/no-unsafe-assignment': 'error',
      '@typescript-eslint/no-unsafe-call': 'error',
      '@typescript-eslint/no-unsafe-member-access': 'error',
      '@typescript-eslint/no-unsafe-return': 'error',
      '@typescript-eslint/no-unsafe-argument': 'error',
      '@typescript-eslint/strict-boolean-expressions': 'error',
      '@typescript-eslint/no-floating-promises': 'error',
      '@typescript-eslint/no-misused-promises': 'error',
      '@typescript-eslint/restrict-template-expressions': 'error',
      '@typescript-eslint/no-unnecessary-type-assertion': 'error',
      '@typescript-eslint/no-non-null-assertion': 'error',

      '@typescript-eslint/explicit-function-return-type': ['error', { allowExpressions: false }],
      '@typescript-eslint/explicit-module-boundary-types': 'error',

      'import-x/no-cycle': ['error', { maxDepth: Infinity }],
      'import-x/no-duplicates': 'error',
      'import-x/no-self-import': 'error',
      'import-x/no-useless-path-segments': 'error',
      'import-x/no-relative-parent-imports': 'error',
      'import-x/no-restricted-paths': ['error', {
        zones: [
          {
            target: './src/shared',
            from: ['./src/entities', './src/features', './src/widgets', './src/routes', './src/app'],
            message: '[FSD] "shared" cannot import from layers above it.',
          },
          {
            target: './src/entities',
            from: ['./src/features', './src/widgets', './src/routes', './src/app'],
            message: '[FSD] "entities" cannot import from layers above it.',
          },
          {
            target: './src/features',
            from: ['./src/widgets', './src/routes', './src/app'],
            message: '[FSD] "features" cannot import from layers above it.',
          },
          {
            target: './src/widgets',
            from: ['./src/routes', './src/app'],
            message: '[FSD] "widgets" cannot import from layers above it.',
          },
          {
            target: './src/routes',
            from: ['./src/app'],
            message: '[FSD] "routes" cannot import from "app" layer.',
          },
        ],
      }],

      'eqeqeq': ['error', 'always'],
      'no-console': ['error', { allow: ['warn', 'error'] }],
      'no-var': 'error',
      'prefer-const': 'error',
      'no-nested-ternary': 'error',
      'no-multi-assign': 'error',
      'no-return-await': 'error',
      'no-useless-return': 'error',

      '@typescript-eslint/no-unused-vars': 'off',
      'unused-imports/no-unused-imports': 'error',
      'unused-imports/no-unused-vars': [
        'error',
        { vars: 'all', varsIgnorePattern: '^_', argsIgnorePattern: '^_' }
      ],

      'react/jsx-key': 'error',
      'react/void-dom-elements-no-children': 'error',
      'react/jsx-no-useless-fragment': 'error',
      'react/self-closing-comp': 'error',

      'unicorn/prevent-abbreviations': ['error', { replacements: { props: false } }],
      'unicorn/no-array-callback-reference': 'error',
      'unicorn/no-null': 'error',
      'unicorn/prefer-module': 'error',
      'unicorn/explicit-length-check': 'error',
    },

    settings: {
      react: { version: 'detect' },
      'import-x/resolver-next': [
        createTypeScriptImportResolver({ alwaysTryTypes: true }),
      ],
    },
  },
  ...storybook.configs["flat/recommended"]
])