/** @type {import('stylelint').Config} */
export default {
    "extends": [
        "stylelint-config-standard-scss"
    ],
    "plugins": [
        "stylelint-declaration-strict-value"
    ],
    "rules": {
        "scss/at-rule-no-unknown": [
            true,
            {
                "ignoreAtRules": ["tailwind"]
            }
        ],
        "scale-unlimited/declaration-strict-value": [
            ["/color/", "fill", "stroke"],
            {
                "ignoreValues": ["currentColor", "transparent", "inherit", "none", "initial"],
                "expandShorthand": true,
                "message": "Use variables from colors.scss for colors only"
            }
        ],
        "custom-property-pattern": "^[a-zA-Z0-9-]+$",
        "selector-class-pattern": "^[a-z][a-zA-Z0-9]*$",
        "color-hex-length": "long"
    },
    "overrides": [
        {
            "files": ["src/**/colors.scss"],
            "rules": {
                "scale-unlimited/declaration-strict-value": null
            }
        },
        {
            "files": ["src/**/MdEditor.module.scss", "src/**/MdRenderer.module.scss", "src/**/TargetAudience.module.scss"],
            "rules": {
                "selector-pseudo-class-no-unknown": [
                    true,
                    {
                        "ignorePseudoClasses": ["global"]
                    }
                ],
                "selector-class-pattern": null
            }
        }
    ]
}
