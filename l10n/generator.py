import argparse
import json
import logging
import os
from typing import  Any, Optional, Union
import yaml

class Types:
    def __init__(self):
        self.STRING = str
        self.INTEGER = int
        self.FLOAT = float

class L10nParamsVariable:
    def __init__(self) -> None:
        self.__variable_name: str
        self.__default_value: Union[str, None] = None
        self.__type: Union[str, int, float]

    @property
    def variable_name(self) -> str:
        return self.__variable_name

    @variable_name.setter
    def variable_name(self, value: str) -> None:
        self.__variable_name = value

    @property
    def default_value(self) -> Union[str, None]:
        return self.__default_value

    @default_value.setter
    def default_value(self, value: Union[str, None]) -> None:
        self.__default_value = value

    @property
    def type(self) -> Union[str, int, float]:
        return self.__type

    @type.setter
    def type(self, value: Union[str, int, float]) -> None:
        if (
            isinstance(value, str)
            or isinstance(value, int)
            or isinstance(value, float)
        ):
            self.__type = value
        else:
            raise ValueError(f"Incorrect data type: {type(value)}, allowed types: str, int, float")

class L10nParams:
    def __init__(self) -> None:
        self.__description: Union[str, None] = None
        self.__example: Union[str, None] = None
        self.__variables: Union[list[L10nParamsVariable], None] = None
    
    @property
    def description(self) -> Union[str, None]:
        return self.__description

    @description.setter
    def description(self, value: Union[str, None]) -> None:
        self.__description = value
        
    @property
    def example(self) -> Union[str, None]:
        return self.__example

    @example.setter
    def example(self, value: Union[str, None]) -> None:
        self.__example = value

    @property
    def variables(self) -> Union[list[L10nParamsVariable], None]:
        return self.__variables

    @variables.setter
    def variables(self, value: Union[list[L10nParamsVariable], None]) -> None:
        self.__variables = value
        

class L10nObject:
    def __init__(self) -> None:
        self.__value: str
        self.__text: str
        self.__params: Union[L10nParams, None] = None
    
    @property
    def value(self) -> str:
        return self.__value

    @value.setter
    def value(self, value: str) -> None:
        self.__value = value

    @property
    def text(self) -> str:
        return self.__text

    @text.setter
    def text(self, value: str) -> None:
        self.__text = value

    @property
    def params(self) -> Union[L10nParams, None]:
        return self.__params

    @params.setter
    def params(self, value: Union[L10nParams, None]) -> None:
        self.__params = value

class L10nNode(object):
    def __init__(
        self,
        language_code: str,
        translate: list[L10nObject]
    ):
        self.language_code: str = language_code
        self.translate: list[L10nObject] = translate
        self.next = None

class L10nNodeOperations:
    def __init__(self):
        self.head: L10nNode = None

    def append(
        self,
        language_code: str,
        translate: list[L10nObject]
    ) -> None:
        """
        Appends a new node with the given data to the end of the linked list.
        """
        new_node: L10nNode = L10nNode(
            language_code = language_code,
            translate = translate
        )
        if not self.head:
            self.head = new_node
            return
        last_node: L10nNode = self.head
        while last_node.next:
            last_node = last_node.next
        last_node.next = new_node

    def prepend(
        self,
        language_code: str,
        translate: list[L10nObject]
    ) -> None:
        """
        Prepends a new node with the given data to the beginning of the linked list.
        """
        new_node: L10nNode = L10nNode(
            language_code = language_code,
            translate = translate
        )
        new_node.next = self.head
        self.head = new_node

    def delete(
        self,
        language_code: str,
    ) -> None:
        """
        Deletes the first occurrence of the node with the given data from the linked list.
        """
        current_node: L10nNode = self.head
        if current_node and current_node.language_code == language_code:
            self.head = current_node.next
            current_node = None
            return
        prev_node: L10nNode = None
        while current_node and current_node.language_code != language_code:
            prev_node = current_node
            current_node = current_node.next
        if current_node is None:
            return
        prev_node.next = current_node.next
        current_node = None

    def get_value_at_index(self, index: int) -> Union[L10nNode, None]:
        """
        Returns the value of the node at the given index in the linked list.
        Returns None if the index is out of range.
        """
        current_node: L10nNode = self.head
        count: int = 0
        while current_node:
            if count == index - 1:
                return current_node
            count += 1
            current_node = current_node.next
        return None

    def get_index_of_value(
        self,
        language_code: str,
    ) -> Optional[int]:
        """
        Returns the index of the first occurrence of the given value in the linked list.
        Returns None if the value is not found.
        """
        current_node: L10nNode = self.head
        index: int = 0
        while current_node:
            if current_node.language_code == language_code:
                return index
            index += 1
            current_node = current_node.next
        return None

    @property
    def size(self) -> int:
        """
        Returns the number of elements in the linked list.
        """
        current_node: L10nNode = self.head
        count: int = 0
        while current_node:
            count += 1
            current_node = current_node.next
        return count

    def is_empty(self) -> bool:
        """
        Checks if the linked list is empty.
        """
        return self.head is None

    def print_list(self) -> None:
        """
        Prints all the elements in the linked list.
        """
        current_node: L10nNode = self.head
        while current_node:
            print(current_node.language_code, end=' ')
            current_node = current_node.next
        print()

class Configuration:
    def __init__(self):
        self.__path_to_translates: str
        self.__path_to_out: str
        self.__default_translate_file: str
        self.__class_name: str
    
    @property
    def path_to_translates(self) -> str:
        return self.__path_to_translates

    @path_to_translates.setter
    def path_to_translates(self, value: str) -> None:
        self.__path_to_translates = value

    @property
    def path_to_out(self) -> str:
        return self.__path_to_out

    @path_to_out.setter
    def path_to_out(self, value: str) -> None:
        self.__path_to_out = value

    @property
    def default_translate_file(self) -> str:
        return self.__default_translate_file

    @default_translate_file.setter
    def default_translate_file(self, value: str) -> None:
        self.__default_translate_file = value

    @property
    def class_name(self) -> str:
        return self.__class_name

    @class_name.setter
    def class_name(self, value: str) -> None:
        self.__class_name = value

class Generator:
    """
    ### This module is a generator of a python file that will contain translations. \
    The translations themselves are contained in l10n_languageCode.json files (e.g. l10n_en.json).
    #### File contents:
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
                }
            }
        }
    }
    ```
    """
    def __init__(
        self,
        path_to_config_file: str,
    ) -> None:
        self.path_to_config_file: str = path_to_config_file
        self.configuration: Configuration = Configuration()
        self.__l10n_node_operations = L10nNodeOperations()
        self.__locales: list[str] = []

    def parse_config(self) -> None:
        with open(self.path_to_config_file, "r") as file:
            data = yaml.safe_load(file)
        
        self.configuration.path_to_translates = data['pathToTranslates']
        self.configuration.path_to_out = data['pathToOut']
        self.configuration.default_translate_file = data['defaultTranslateFile']
        self.configuration.class_name = data['className']
    
    def __get_file_names(self) -> list[str]:
        """
        Returns an array with the names of files contained in the folder.\n
        Returns: ["l10n_en.json","l10n_ru.json",...]
        """
        file_names_with_extension = []
        with os.scandir(self.configuration.path_to_translates) as entries:
            for entry in entries:
                if entry.is_file() and entry.name.startswith('l10n_') and entry.name.endswith('.json'):
                    file_names_with_extension.append(entry.name)
        
        return file_names_with_extension

    def __convert_type(self, input_type: str) -> type:
        types = Types()
        
        if input_type == 'string':
            return types.STRING
        elif input_type == 'float':
            return types.FLOAT
        elif input_type == 'int':
            return types.INTEGER
        else:
            raise ValueError(f"Incorrect data type: {input_type}")
        
    def __check_type(self, value: Any, expected_type: type) -> bool:
        if isinstance(value, expected_type):
            return True
        else:
            return False

    def unmarshal(self):
        files: list[str] = self.__get_file_names()

        for i in files:
            l10n_object_list: list[L10nObject] = []  # Create a new list for each file

            with open(self.configuration.path_to_translates + i, 'r', encoding='utf-8') as file:
                data: dict[str, Union[str, dict, list]] = json.load(file)
            
            language_code = i.replace('l10n_', '').replace('.json', '')
            self.__locales.append(language_code)

            keys_list: list[str] = list(data.keys())
            
            for key in keys_list:
                l10n_object = L10nObject()
                l10n_object.value = key
                l10n_object.text = key

                if not key.startswith('#'):
                    l10n_object.text = data[key]
                else:
                    continue  # Skip the metadata for now

                if f"#{key}" in keys_list:  # Check if there is metadata for this key
                    l10n_params = L10nParams()
                    meta_key = f"#{key}"
                    
                    if 'description' in data[meta_key]:
                        l10n_params.description = data[meta_key]['description']
                    if 'example' in data[meta_key]:
                        l10n_params.example = data[meta_key]['example']
                    if 'variables' in data[meta_key]:
                        variables_map = data[meta_key]['variables']
                        variables: list[L10nParamsVariable] = []

                        for var_name, var_data in variables_map.items():
                            l10n_params_variable = L10nParamsVariable()
                            l10n_params_variable.variable_name = var_name
                            typeVar = self.__convert_type(var_data['type'])
                            l10n_params_variable.type = typeVar.__name__
                            if 'defaultValue' in var_data:
                                if self.__check_type(var_data['defaultValue'], typeVar):
                                    l10n_params_variable.default_value = var_data['defaultValue']
                                else:
                                    logging.error(f"Default type does not match the specified variable type, default variable type: {type(var_data['defaultValue'])}, expected type: {l10n_params_variable.type}")
                                    return
                                
                            variables.append(l10n_params_variable)
                            l10n_params.variables = variables

                        l10n_object.params = l10n_params
                l10n_object_list.append(l10n_object)

            self.__l10n_node_operations.append(
                language_code = language_code,
                translate = l10n_object_list
            )

    def generate(
        self,
        template_base_class: str,
        template_extend_class: str,
        template_property: str
    ):
        # Extracting default language code from configuration file
        default_language_code: str = self.configuration.default_translate_file.replace(
            'l10n_', ''
        ).replace(
            '.json', ''
        )

        default_language_node: L10nNode

        for i in range(1, self.__l10n_node_operations.size + 1):
            current_node = self.__l10n_node_operations.get_value_at_index(i)
            if current_node.language_code == default_language_code:
                default_language_node = current_node
                break

        BASE_CLASS = f"Base{self.configuration.class_name}"  # Generating base class name

        # Generation of implementation classes.
        extend_class: str = ""  # Initializing extend class string

        # Iterate over all languages
        for i in range(1, self.__l10n_node_operations.size + 1):
            language_code = self.__l10n_node_operations.get_value_at_index(i).language_code.capitalize()
            current_node = self.__l10n_node_operations.get_value_at_index(i)

            # Generating extended class using provided template and language code
            extended_class = template_extend_class.replace(
                '{ClassNameExtend}',
                f"{self.configuration.class_name}{language_code}"
            ).replace(
                '{ExtendClass}',
                BASE_CLASS
            ).replace(
                '{bodyExtend}',
                '...'
            ).replace(
                '{ClassConstructorArgs}',
                ''
            )

            properties_class: str = ''

            # Iterate over the translations of default language
            for default_val in default_language_node.translate:
                # Check if the property exists in the current language
                current_property = next((prop for prop in current_node.translate if prop.value == default_val.value), None)
                
                # If the property does not exist in the current language, use the default language property
                if not current_property:
                    val = default_val
                else:
                    val = current_property

                _property = template_property.replace(
                    '{PropertyName}',
                    val.value
                ).replace(
                    '{PropertyType}',
                    type(val.text).__name__
                ).replace(
                    '{PropertyValue}',
                    f'return f"{val.text}"'
                )

                # Checking if property has parameters
                if val.params is None or val.params.variables is None:
                    _property = _property.replace(
                        '{PropertyArgs}',
                        ''
                    ).replace(
                        '{PropertyIs}',
                        '@property'
                    )
                else:
                    property_args = ""
                    default_args = ""
                    for variable in val.params.variables:
                        # Adding variable and its type to property
                        if variable.default_value is None:
                            default_args += f", {variable.variable_name}: {variable.type}"
                        else:
                            property_args += f", {variable.variable_name}: {variable.type}"
                            # Checking if variable has a default value
                            if isinstance(variable.default_value, str):
                                property_args += f' = "{variable.default_value}"'
                            else:
                                property_args += f' = {variable.default_value}'
                    # Concatenate default_args before property_args
                    property_args = default_args + property_args
                    _property = _property.replace('{PropertyArgs}', property_args).replace('{PropertyIs}', '')

                properties_class += _property  # Appending property to properties string

            extend_class += extended_class + properties_class + "\n"  # Appending extended class to the extend class string

        # Generating base class using provided template and base class name
        base_class = template_base_class.replace(
            '{ClassNameBase}', BASE_CLASS
        ).replace(
            '{bodyBase}',
            '...'
        ).replace(
                '{ClassConstructorArgs}',
                ''
            )

        properties_base_class: str = ""  # Initializing properties string for base class
        for val in default_language_node.translate:
            _property: str = template_property.replace(
                '{PropertyName}',
                val.value
            ).replace(
                '{PropertyValue}',
                '{Comments}'
            )

            property_description: str = ''

            if val.params is not None:
                if val.params.description is not None:
                    property_description += val.params.description + '\\n'
                else:
                    property_description += f'No description provided for #{val.value}.\\n'

                property_description += "\n"

                if val.params.example is not None:
                    property_description += f'        {val.params.example}\\n\n'
                else:
                    property_description += ''

                if val.params.variables is not None:
                    __variables_description: str = ''
                    __variables_description += "        Args:\n"
                    for vars in val.params.variables:
                        __variables_description += f'           {vars.variable_name} ({vars.type}) Default value: {vars.default_value if vars.default_value is not None else ""}.\n'
                    property_description += __variables_description
            _property = _property.replace(
                '{Comments}',
                f'"""{property_description}        """\n        {{PropertyValue}}'
            )
            _property = _property.replace(
                '{PropertyValue}',
                'raise NotImplementedError(f"{self.__class__.__name__}.{inspect.currentframe().f_code.co_name} method must be implemented in subclass")'
            ).replace(
                '{PropertyType}',
                type(val.text).__name__
            )

            # Checking if property has parameters
            if val.params is None or val.params.variables is None:
                _property = _property.replace(
                    '{PropertyArgs}',
                    ''
                ).replace(
                    '{PropertyIs}',
                    '@property'
                )
            else:
                property_args = ""
                default_args = ""
                for variable in val.params.variables:
                    # Adding variable and its type to property
                    if variable.default_value is None:
                        default_args += f", {variable.variable_name}: {variable.type}"
                    else:
                        property_args += f", {variable.variable_name}: {variable.type}"
                        # Checking if variable has a default value
                        if isinstance(variable.default_value, str):
                            property_args += f' = "{variable.default_value}"'
                        else:
                            property_args += f' = {variable.default_value}'
                # Concatenate default_args before property_args
                property_args = default_args + property_args
                _property = _property.replace('{PropertyArgs}', property_args).replace('{PropertyIs}', '')

            properties_base_class += _property  # Appending property to properties string

        # Generating main class using provided template and class name
        main_class = template_base_class.replace(
            '{ClassNameBase}',
            self.configuration.class_name
        ).replace(
            '{bodyBase}',
            '\n        self.__locale = locale'
        ).replace(
            '{ClassConstructorArgs}',
            f', locale: str = "{default_language_code}"'
        )

        main_class += template_property.replace(
            '{PropertyIs}',
            '@property'
        ).replace(
            '{PropertyName}',
            'default_locale'
        ).replace(
            '{PropertyArgs}',
            ''
        ).replace(
            '{PropertyType}',
            'str'
        ).replace(
            '{PropertyValue}',
            f'return "{default_language_code}"'
        )

        main_class += template_property.replace(
            '{PropertyIs}',
            '@property'
        ).replace(
            '{PropertyName}',
            'current_locale'
        ).replace(
            '{PropertyArgs}',
            ''
        ).replace(
            '{PropertyType}',
            'str'
        ).replace(
            '{PropertyValue}',
            f'return self.__locale'
        )

        supported_locales_str_list: str = ''
        supported_locales_str_list += '[\n'
        for i in range(1, self.__l10n_node_operations.size + 1):
            current = self.__l10n_node_operations.get_value_at_index(i)
            if i == 1:
                supported_locales_str_list += f'           "{current.language_code}",\n'
            else:
                supported_locales_str_list += f'           "{current.language_code}",\n'
        supported_locales_str_list += '        ]'

        main_class += template_property.replace(
            '{PropertyIs}',
            '@property'
        ).replace(
            '{PropertyName}',
            'locales'
        ).replace(
            '{PropertyArgs}',
            ''
        ).replace(
            '{PropertyType}',
            'list[str]'
        ).replace(
            '{PropertyValue}',
            f'return {supported_locales_str_list}'
        )

        main_class_method_of = template_property.replace(
            '{PropertyIs}',
            ''
        ).replace(
            '{PropertyName}',
            'of'
        ).replace(
            '{PropertyArgs}',
            ''
        ).replace(
            '{PropertyType}',
            BASE_CLASS
        ).replace(
            '{PropertyValue}',
            f''
        )
        
        method_of_if_else_construct: str = ''
        for i in range(1, self.__l10n_node_operations.size + 1):
            current = self.__l10n_node_operations.get_value_at_index(i)
            if i == 1:
                method_of_if_else_construct += f"""    if self.__locale == "{current.language_code}":
            return {self.configuration.class_name}{current.language_code.capitalize()}()
                """
            else:
                method_of_if_else_construct += f"""\n        elif self.__locale == "{current.language_code}":
            return {self.configuration.class_name}{current.language_code.capitalize()}()
                """
        
        method_of_if_else_construct += """
        else:
            raise ValueError(f"No {self.__locale} localization.")
        """

        main_class_method_of += method_of_if_else_construct

        main_class += main_class_method_of

        imports = (
            "import inspect"  # Importing inspect module
        )

        note: str = "#NOTE THIS IS AN AUTO-GENERATED FILE, DO NOT EDIT IT."

        # Combining all generated components to form the final result
        result = (
            note
            + "\n\n"
            + imports
            + "\n\n"
            + base_class
            + "\n"
            + properties_base_class
            + "\n"
            + extend_class
            + main_class
        )

        # Writing the result to an output file
        with open(self.configuration.path_to_out, 'w', encoding='utf-8') as file:
            file.write(
                result
            )


parser = argparse.ArgumentParser()
parser.add_argument('--config', type=str, help="Path to configuration file")
args = parser.parse_args()

if __name__ == "__main__":
    if not args.config:
        logging.warning("Error: no path to config file provided")
        exit(0)

    template_base_class = """class {ClassNameBase}:
    def __init__(self{ClassConstructorArgs}):{bodyBase}
    """

    template_extend_class = """class {ClassNameExtend}({ExtendClass}):
    def __init__(self{ClassConstructorArgs}):{bodyExtend}
    """

    template_property = """
    {PropertyIs}
    def {PropertyName}(self{PropertyArgs}) -> {PropertyType}:
        {PropertyValue}
    """

    obj = Generator(args.config,)
    obj.parse_config()
    obj.unmarshal()
    obj.generate(
        template_base_class,
        template_extend_class,
        template_property
    )