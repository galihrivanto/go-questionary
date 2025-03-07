from langchain.llms import Ollama
from langchain.chains import LLMChain
from langchain.prompts import PromptTemplate
from abc import ABC, abstractmethod

class BaseAgent(ABC):
    def __init__(self, model_name="llama2"):
        self.llm = Ollama(model=model_name)
        self.chain = self._create_chain()
    
    @abstractmethod
    def get_prompt_template(self) -> str:
        pass
    
    def _create_chain(self) -> LLMChain:
        prompt = PromptTemplate(
            input_variables=["token_data"],
            template=self.get_prompt_template()
        )
        return LLMChain(llm=self.llm, prompt=prompt)
    
    @abstractmethod
    async def analyze(self, token_data: dict) -> dict:
        pass        
        