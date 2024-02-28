# L10n, library for application localization, based on json files.

## Supported programming languages.

| Language |
|----------|
| Python   |

## How to use
### Downloading the library.
```bash
wget https://github.com/A1ex3/l10n/releases/download/{version}_{programming language}/{file}
```
#### Example.
```bash
wget https://github.com/A1ex3/l10n/releases/download/0.0.1_python/l10n_0.0.1_python.tar.gz
```

### Unpack the archive.
```bash
tar -xvf l10n_0.0.1_python.tar.gz
```

### Library installation.
```bash
pip install -e l10n 
```

### Creating a configuration file.

#### Example `configuration.yml`.
```yaml
pathToTranslates: app/translates/
pathToOut: app/app_localization.py
defaultTranslateFile: l10n_en.json
className: AppLocalization
```

| Parameter | Description |
|-----------|-------------|
| pathToTranslates | Path to the directory with translations where l10n_language.json files are stored. |
| pathToOut | The path where the file will be generated, you must first create a directory where the file will be saved. |
| defaultTranslateFile | Default translation file where the translation status will always be 100%, no need to specify the full path to the file, just the name. |
| className | The name of the class whose name will be assigned to the main class. |

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

#### Table of optional values to be filled in.

| Key | Description |
|-----|-------------|
| #params | If you need parameters for translation. |
| description | |
| example | |
| defaultValue | |

### Generation.
```bash
python -m l10n.generator --config="app/configuration.yml"
```

| Key | args |
|-----|-------------|
| --configuration | Path to configuration [file](#example-configurationyml) |

### Result `app_localization.py`.
```python
#NOTE THIS IS AN AUTO-GENERATED FILE, DO NOT EDIT IT.

import inspect

class BaseAppLocalization:
    def __init__(self):...
    

    @property
    def helloWorld(self) -> str:
        """        """
        raise NotImplementedError(f"{self.__class__.__name__}.{inspect.currentframe().f_code.co_name} method must be implemented in subclass")
    
    
    def bye(self, value: str = "World") -> str:
        """Saying goodbye to someone\n
        Bye World\n
        Args:
           value (str) Default value: World.
        """
        raise NotImplementedError(f"{self.__class__.__name__}.{inspect.currentframe().f_code.co_name} method must be implemented in subclass")
    
    
    def numberOfUsers(self, values: str, number: int = 0) -> str:
        """No description provided for #numberOfUsers.\n
        Args:
           number (int) Default value: 0.
           values (str) Default value: .
        """
        raise NotImplementedError(f"{self.__class__.__name__}.{inspect.currentframe().f_code.co_name} method must be implemented in subclass")
    
class AppLocalizationEn(BaseAppLocalization):
    def __init__(self):...
    
    @property
    def helloWorld(self) -> str:
        return f"Hello World"
    
    
    def bye(self, value: str = "World") -> str:
        return f"Bye {value}"
    
    
    def numberOfUsers(self, values: str, number: int = 0) -> str:
        return f"Number of users: {number}"
    
class AppLocalizationRu(BaseAppLocalization):
    def __init__(self):...
    
    @property
    def helloWorld(self) -> str:
        return f"Привет мир"
    
    
    def bye(self, value: str = "Мир") -> str:
        return f"Пока {value}"
    
    
    def numberOfUsers(self, values: str, number: int = 0) -> str:
        return f"Number of users: {number}"
    
class AppLocalization:
    def __init__(self, locale: str = "en"):
        self.__locale = locale
    
    @property
    def default_locale(self) -> str:
        return "en"
    
    @property
    def current_locale(self) -> str:
        return self.__locale
    
    @property
    def locales(self) -> list[str]:
        return [
           "en",
           "ru",
        ]
    
    
    def of(self) -> BaseAppLocalization:
        
        if self.__locale == "en":
            return AppLocalizationEn()
                
        elif self.__locale == "ru":
            return AppLocalizationRu()
                
        else:
            raise ValueError(f"No {self.__locale} localization.")

```
