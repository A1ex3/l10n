# L10n, an application localization program based on json files.

## Supported programming languages.

| Language   |
|------------|
| Python     |
| Java       |
| Kotlin     |
| Cpp        |
| JavaScript |
| TypeScript |

## How to use
### Downloading the program.
```bash
wget https://github.com/A1ex3/l10n/releases/download/{version}/{file}
```
#### Example.
```bash
wget https://github.com/A1ex3/l10n/releases/download/v1.0.0/l10n_windows.exe
```

### Creating a configuration file.

#### Example `configuration.yml`.
```yaml
dir: app/translates/
output_localization_file: app/app_localization.py
template: en
class_name: AppLocalization
programming_language: Python
```

| Parameter | Description |
|-----------|-------------|
| dir | Path to the directory with translations where l10n_[language code]-[optional country code].json files are stored. |
| output_localization_file | The path where the file will be generated, you must first create a directory where the file will be saved. |
| template | Default translation where the translation status will always be 100%, no need to specify the full path to the file, just the name. |
| class_name | The name of the class whose name will be assigned to the main class. |
| programming_language | Programming language for which the code will be generated. In programming languages where a package name is required, such as `Java`, `Kotlin`, etc. It can be specified in the value for `programming_language`. For example `Java.dev.example`, then the package name will be `dev.example`, the default package name is `l10n`. |

### Creating a translation file.
#### Example.
```json
{
    "helloWorld": "Hello World",
    "bye": "Bye {value}",
    "#bye": {
        "description": "Saying goodbye to someone",
        "example": "Bye World",
        "variables": {
            "value": {
                "defaultValue": "World",
                "type": "string"
            }
        }
    },
    "numberOfUsers": "Number of users: {number}",
    "#numberOfUsers": {
        "variables": {
            "number": {
                "defaultValue": 0,
                "type": "int"
            },
            "values": {
                "type": "string"
            }
        }
    }
}
```
#### Description.
#### The file contains translations in json format. Define translation and its `optional` parameters (keys whose value starts with `#`, e.g. `#bye`).
```json
{
    "key": "Translate",
    "#key": {}
}
```
#### Values such as, for example, description `optional`, example `optional` can be added internally. The variables to be inserted into the text are defined here in `variables` is optional if no `variables` are required to be inserted into the text.
#### The names of the variables should be different within the same translation. In variables, you must set the variable type, defaultValue `optional` (the defaultValue type must match the specified variable type, otherwise an error will occur).
```json
{
    "helloAny": "Hello {value}",
    "#helloAny": {
        "description": "Greeting anyone",
        "example": "Hello World",
        "variables": {
            "value": {
                "defaultValue": "World",
                "type": "string"
            }
        }
    }
}
```

| Types  |
|--------|
| string |
| int    |
| float  |

#### Table of optional values to be filled in.

| Key | Description |
|-----|-------------|
| #params | If you need parameters for translation. |
| description | |
| example | |
| defaultValue | |

### Generation.
```powershell
.\build\l10n.exe -FILE .\build\config.yml
```

| Key | args |
|-----|-------------|
| -FILE | Path to configuration [file](#example-configurationyml) |

or

### Generation.
```powershell
.\build\l10n.exe -DIR "app/translates/" -TEMPLATE "en" -OUTPUT_LOCALIZATION_FILE "app/app_localization.py" -PROGRAMMING_LANGUAGE "Python" -CLASS_NAME "AppLocalization"
```

| Key |
|-----|
| -DIR |
| -TEMPLATE |
| -OUTPUT_LOCALIZATION_FILE |
| -PROGRAMMING_LANGUAGE |
| -CLASS_NAME |

### Result `app_localization.py`.
```python
class BaseApplocalization:
    def helloWorld(self, ) -> str:
        raise NotImplementedError("Method must be implemented in subclass!")
    def bye(self, value: str) -> str:
        raise NotImplementedError("Method must be implemented in subclass!")
    def numberOfUsers(self, number: int, values: str) -> str:
        raise NotImplementedError("Method must be implemented in subclass!")

class AppLocalizationEnUs(BaseApplocalization):
    def __init__(self) -> None:
        self.LANGUAGE_CODE: str = "en-US"
    def bye(self, value: str) -> str:
        """Description: Saying goodbye to someone
        Example: Bye World
        """
        return f"Bye {value}"
    def helloWorld(self, ) -> str:
        """Description: 
        Example: 
        """
        return f"Hello World"
    def numberOfUsers(self, number: int, values: str) -> str:
        """Description: 
        Example: 
        """
        return f"Number of users: {number}"

class AppLocalizationRu(BaseApplocalization):
    def __init__(self) -> None:
        self.LANGUAGE_CODE: str = "ru"
    def bye(self, value: str) -> str:
        """Description: Saying goodbye to someone
        Example: Bye World
        """
        return f"Bye {value}"
    def helloWorld(self, ) -> str:
        """Description: 
        Example: 
        """
        return f"Привет мир"
    def numberOfUsers(self, number: int, values: str) -> str:
        """Description: 
        Example: 
        """
        return f"Количество пользователей: {number}"

class AppLocalization:
    DEFAULT_LANGUAGE_CODE: str = "en-US"
    current_language_code: str = "en-US"
    languages: dict[str, BaseApplocalization] = {
        "en-US": AppLocalizationEnUs(),
        "ru": AppLocalizationRu(),
    }

    @staticmethod
    def get() -> BaseApplocalization:
        return AppLocalization.languages[AppLocalization.current_language_code]
```

### Result `AppLocalization.java`.
```java
package l10n;

import java.text.MessageFormat;
import java.util.HashMap;
import java.util.Map;

interface BaseApplocalization {
    /**
     * Description: <b>  </b>
     * Example: <b>  </b>
     */
    String numberOfUsers(Integer number, String values);
    /**
     * Description: <b> Saying goodbye to someone </b>
     * Example: <b> Bye World </b>
     */
    String bye(String value);
    /**
     * Description: <b>  </b>
     * Example: <b>  </b>
     */
    String helloWorld();
}

final class AppLocalizationRu implements BaseApplocalization {
    public final String LANGUAGE_CODE = "ru";
    @Override
    public String bye(String value) {
        return MessageFormat.format("Bye {0}", value);
    }
    @Override
    public String helloWorld() {
        return "Привет мир";
    }
    @Override
    public String numberOfUsers(Integer number, String values) {
        return MessageFormat.format("Количество пользователей: {0}", number, values);
    }
}

final class AppLocalizationEnUs implements BaseApplocalization {
    public final String LANGUAGE_CODE = "en-US";
    @Override
    public String bye(String value) {
        return MessageFormat.format("Bye {0}", value);
    }
    @Override
    public String helloWorld() {
        return "Hello World";
    }
    @Override
    public String numberOfUsers(Integer number, String values) {
        return MessageFormat.format("Number of users: {0}", number, values);
    }
}

public final class AppLocalization {
    public static final String DEFAULT_LANGUAGE_CODE = "en-US";
    private static String currentLanguageCode = "en-US";
    public static final Map<String, BaseApplocalization> languages = new HashMap<>();

    static {
        languages.put("ru", new AppLocalizationRu());
        languages.put("en-US", new AppLocalizationEnUs());
    }

    public static BaseApplocalization get() {
        return languages.get(currentLanguageCode);
    }

    public static boolean setCurrentLanguageCode(String languageCode) {
        if (languages.containsKey(languageCode)) {
            currentLanguageCode = languageCode;
            return true;
        } else {
            return false;
        }
    }

    public static String getCurrentLanguageCode() {
        return currentLanguageCode;
    }
}
```

# TODO
- [ ] Add generation of class methods with default parameters.
