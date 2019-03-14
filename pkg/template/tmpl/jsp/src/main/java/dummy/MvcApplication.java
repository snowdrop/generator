package {{.PackageName}};

{{$hasAp4k := .HasAp4k}}

{{if $hasAp4k}}
import io.ap4k.kubernetes.annotation.KubernetesApplication;
import io.ap4k.openshift.annotation.OpenshiftApplication;
{{end}}
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.boot.web.support.SpringBootServletInitializer;

@SpringBootApplication{{if $hasAp4k}}
@KubernetesApplication
@OpenshiftApplication{{end}}
public class MvcApplication extends SpringBootServletInitializer {

	@Override
	protected SpringApplicationBuilder configure(SpringApplicationBuilder application) {
		return application.sources(MvcApplication.class);
	}

	public static void main(String[] args) throws Exception {
		SpringApplication.run(MvcApplication.class, args);
	}
}
